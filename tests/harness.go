package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

const TESTS_FOLDER = "tests/integration"

const initMaxAttempts = 3

type testResult struct {
	name         string
	passed       bool
	err          error
	environments []string
}

func main() {
	if err := compileProvider(); err != nil {
		log.Fatalf("failed to compile go: %v", err)
	}

	printTerraformVersion()
	makeSureRunningFromProjectRoot()

	testNames := testNamesFromCommandLineArguments()

	log.Println(len(testNames), " tests to run")

	buildFakeTerraformRegistry()

	destroyMode := os.Getenv("DESTROY_MODE")

	var wg sync.WaitGroup

	results := make([]testResult, len(testNames))

	for i, testName := range testNames {
		wg.Add(1)

		go func(i int, testName string) {
			defer wg.Done()

			if destroyMode == "DESTROY_ONLY" {
				terraformDestroy(testName)
				results[i] = testResult{name: testName, passed: true}

				return
			}

			success, environments, err := runTest(testName, destroyMode != "NO_DESTROY")
			if !success {
				log.Printf("Test '%s' failed: %v\n", testName, err)
			}

			results[i] = testResult{name: testName, passed: success, err: err, environments: environments}
		}(i, testName)
	}

	wg.Wait()

	if destroyMode == "DESTROY_ONLY" {
		return
	}

	failed := printSummary(results)

	writeGithubStepSummary(results)

	if failed > 0 {
		os.Exit(1)
	}
}

func printSummary(results []testResult) int {
	failed := 0

	log.Println("==================== Integration tests summary ====================")

	for _, result := range results {
		if result.passed {
			log.Printf("PASS  %s", result.name)
		} else {
			failed++

			log.Printf("FAIL  %s", result.name)
		}
	}

	for _, result := range results {
		if result.passed {
			continue
		}

		log.Printf("-------------------- failure details: %s --------------------", result.name)

		for _, environment := range result.environments {
			log.Println(environment)
		}

		for _, line := range errorHighlights(result.err) {
			log.Println(line)
		}
	}

	log.Printf("==================== %d/%d tests passed ====================", len(results)-failed, len(results))

	return failed
}

// errorHighlights extracts the interesting lines from a test failure: terraform
// diagnostic blocks (╷...╵) and any other lines mentioning an error, skipping the
// verbose provider logs.
func errorHighlights(err error) []string {
	if err == nil {
		return []string{"test failed without error details"}
	}

	var (
		highlights   []string
		inDiagnostic bool
	)

	for _, line := range strings.Split(err.Error(), "\n") {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "╷"):
			inDiagnostic = true

			highlights = append(highlights, line)
		case strings.HasPrefix(trimmed, "╵"):
			inDiagnostic = false

			highlights = append(highlights, line)
		case inDiagnostic:
			highlights = append(highlights, line)
		case strings.Contains(trimmed, "Error") && !strings.Contains(trimmed, "[INFO]"):
			highlights = append(highlights, line)
		}

		if len(highlights) >= 100 {
			highlights = append(highlights, "... (truncated)")

			break
		}
	}

	if len(highlights) == 0 {
		lines := strings.Split(err.Error(), "\n")
		if len(lines) > 20 {
			lines = lines[len(lines)-20:]
		}

		return lines
	}

	return highlights
}

// writeGithubStepSummary appends a markdown summary of the test results to the
// GitHub Actions step summary (no-op outside of GitHub Actions).
func writeGithubStepSummary(results []testResult) {
	summaryPath := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryPath == "" {
		return
	}

	var sb strings.Builder

	sb.WriteString("## Integration tests summary\n\n")
	sb.WriteString("| Test | Result |\n|---|---|\n")

	failuresExist := false

	for _, result := range results {
		status := "✅ Pass"
		if !result.passed {
			status = "❌ Fail"
			failuresExist = true
		}

		sb.WriteString(fmt.Sprintf("| %s | %s |\n", result.name, status))
	}

	if failuresExist {
		sb.WriteString("\n### Failures\n")

		for _, result := range results {
			if result.passed {
				continue
			}

			sb.WriteString(fmt.Sprintf("\n#### ❌ %s\n\n", result.name))

			for _, environment := range result.environments {
				sb.WriteString(fmt.Sprintf("- %s\n", environment))
			}

			sb.WriteString(fmt.Sprintf("\n```\n%s\n```\n", strings.Join(errorHighlights(result.err), "\n")))
		}
	}

	file, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("WARNING: unable to open GITHUB_STEP_SUMMARY file:", err)

		return
	}
	defer file.Close()

	if _, err := file.WriteString(sb.String()); err != nil {
		log.Println("WARNING: unable to write GITHUB_STEP_SUMMARY:", err)
	}
}
func compileProvider() error {
	cmd := exec.Command("go", "build")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
func runTest(testName string, destroy bool) (passed bool, environments []string, err error) {
	testDir := TESTS_FOLDER + "/" + testName
	toDelete := []string{
		".terraform",
		".terraform.lock.hcl",
		"terraform.tfstate",
		"terraform.tfstate.backup",
		"terraform.rc",
	}

	for _, oneToDelete := range toDelete {
		os.RemoveAll(path.Join(testDir, oneToDelete))
	}

	log.Println("Running test ", testName)

	_, err = terraformInit(testName)
	if err != nil {
		return false, nil, err
	}

	_, _ = terraformCommand(testName, "fmt")

	if destroy {
		defer terraformDestroy(testName)
	}

	// Registered after the destroy cleanup so it runs before it (defers are LIFO):
	// on failure, capture the env0 environments still in the state before they are cleaned up.
	defer func() {
		if !passed {
			environments = environmentsInState(testName)
		}
	}()

	_, err = terraformCommand(testName, "apply", "-auto-approve", "-var", "second_run=0")
	if err != nil {
		return false, nil, err
	}

	_, err = terraformCommand(testName, "apply", "-auto-approve", "-var", "second_run=1")
	if err != nil {
		return false, nil, err
	}

	expectedOutputs, err := readExpectedOutputs(testName)
	if err != nil {
		return false, nil, err
	}

	outputsBytes, err := terraformCommand(testName, "output", "-json")
	if err != nil {
		return false, nil, err
	}

	outputs, err := bytesOfJsonToStringMap(outputsBytes)
	if err != nil {
		return false, nil, err
	}

	for key, expectedValue := range expectedOutputs {
		value, ok := outputs[key]
		if !ok {
			log.Println("Error: Expected terraform output ", key, " but no such output was created")

			return false, nil, fmt.Errorf("expected terraform output '%s' but no such output was created", key)
		}

		if value != expectedValue {
			log.Printf("Error: Expected output of '%s' to be '%s' but found '%s'\n", key, expectedValue, value)

			return false, nil, fmt.Errorf("expected output of '%s' to be '%s' but found '%s'", key, expectedValue, value)
		}

		log.Printf("Verified expected '%s'='%s' in %s", key, value, testName)
	}

	if destroy {
		_, err = terraformCommand(testName, "destroy", "-auto-approve", "-var", "second_run=0")
		if err != nil {
			return false, nil, err
		}
	}

	log.Println("Successfully finished running test ", testName)

	return true, nil, nil
}

type stateModule struct {
	Resources []struct {
		Type   string `json:"type"`
		Values struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			ProjectId string `json:"project_id"`
		} `json:"values"`
	} `json:"resources"`
	ChildModules []stateModule `json:"child_modules"`
}

// environmentsInState lists the env0 environments currently in the test's terraform state,
// each with a direct link to the environment in the env0 UI - to make it easy to find the
// environment a failed test was running against.
func environmentsInState(testName string) []string {
	stateBytes, err := terraformCommand(testName, "show", "-json")
	if err != nil {
		return nil
	}

	var state struct {
		Values struct {
			RootModule stateModule `json:"root_module"`
		} `json:"values"`
	}

	if err := json.Unmarshal(stateBytes, &state); err != nil {
		return nil
	}

	return collectEnvironments(state.Values.RootModule)
}

func collectEnvironments(module stateModule) []string {
	var environments []string

	for _, resource := range module.Resources {
		if resource.Type != "env0_environment" {
			continue
		}

		environment := fmt.Sprintf("environment '%s' (id: %s)", resource.Values.Name, resource.Values.Id)
		if environmentUrl := environmentUiUrl(resource.Values.ProjectId, resource.Values.Id); environmentUrl != "" {
			environment += " " + environmentUrl
		}

		environments = append(environments, environment)
	}

	for _, child := range module.ChildModules {
		environments = append(environments, collectEnvironments(child)...)
	}

	return environments
}

// environmentUiUrl returns a link to the environment in the env0 UI, derived from the
// API endpoint the tests run against (e.g. https://api.env0.com -> https://app.env0.com).
func environmentUiUrl(projectId string, environmentId string) string {
	apiEndpoint := os.Getenv("ENV0_API_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = "https://api.env0.com"
	}

	parsed, err := url.Parse(apiEndpoint)
	if err != nil || parsed.Host == "" {
		return ""
	}

	var host string

	switch {
	case strings.HasPrefix(parsed.Host, "api-"):
		host = strings.TrimPrefix(parsed.Host, "api-")
	case strings.HasPrefix(parsed.Host, "api."):
		host = "app." + strings.TrimPrefix(parsed.Host, "api.")
	default:
		return ""
	}

	return fmt.Sprintf("https://%s/p/%s/environments/%s", host, projectId, environmentId)
}

func readExpectedOutputs(testName string) (map[string]string, error) {
	expectedBytes, err := os.ReadFile(path.Join(TESTS_FOLDER, testName, "expected_outputs.json"))
	if err != nil {
		log.Println("Test folder for ", testName, " does not contain expected_outputs.json", err)

		return nil, err
	}

	return bytesOfJsonToStringMap(expectedBytes)
}

func bytesOfJsonToStringMap(data []byte) (map[string]string, error) {
	var stringMapUncasted map[string]any

	err := json.Unmarshal(data, &stringMapUncasted)
	if err != nil {
		log.Println("Unable to parse json:", err)
		log.Println("** JSON Input **")
		log.Println(string(data))
		log.Println("******")

		return nil, err
	}

	result := map[string]string{}

	for key, valueUncasted := range stringMapUncasted {
		switch value := valueUncasted.(type) {
		case string:
			result[key] = value
		case map[string]any:
			result[key] = value["value"].(string)
		}
	}

	return result, nil
}

func terraformDestroy(testName string) {
	log.Println("Running destroy to clean up in", testName)

	destroy := exec.Command("tofu", "destroy", "-auto-approve", "-var", "second_run=0")
	destroy.Env = os.Environ()
	destroy.Dir = TESTS_FOLDER + "/" + testName

	if err := destroy.Run(); err != nil {
		log.Println("WARNING: error running tofu destroy")
	}

	log.Println("Done running tofu destroy in", testName)
}

// terraformInit retries `tofu init` because provider downloads from the public
// registry occasionally time out with transient network errors.
func terraformInit(testName string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= initMaxAttempts; attempt++ {
		output, err := terraformCommand(testName, "init")
		if err == nil {
			return output, nil
		}

		lastErr = err

		log.Printf("tofu init failed in %s (attempt %d/%d), reason: %v", testName, attempt, initMaxAttempts, err)

		if attempt < initMaxAttempts {
			backoff := time.Duration(attempt*5) * time.Second

			log.Printf("retrying tofu init in %s (%s)", testName, backoff)
			time.Sleep(backoff)
		}
	}

	return nil, lastErr
}

func terraformCommand(testName string, arg ...string) ([]byte, error) {
	cmd := exec.Command("tofu", arg...)
	cmd.Dir = TESTS_FOLDER + "/" + testName
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "INTEGRATION_TESTS=1")
	cmd.Env = append(cmd.Env, "TF_LOG_PROVIDER=info")

	var output, errOutput bytes.Buffer

	cmd.Stderr = &errOutput
	cmd.Stdout = &output

	log.Println("Running tofu ", arg, " in ", testName)

	err := cmd.Run()

	log.Println(errOutput.String())

	if err != nil {
		log.Println("error running tofu ", arg, " in ", testName, " error: ", err)

		err = fmt.Errorf("'tofu %s' failed in %s: %w\n%s", strings.Join(arg, " "), testName, err, errOutput.String())
	} else {
		log.Println("Completed successfully tofu", arg, "in", testName)
	}

	return output.Bytes(), err
}

func printTerraformVersion() {
	versionString, err := exec.Command("tofu", "version").Output()
	if err != nil {
		log.Fatalln("Unable to invoke tofu. Perhaps it's not in PATH?", err)
	}

	log.Println("tofu version: ", string(versionString))
}

func makeSureRunningFromProjectRoot() {
	if _, err := os.Stat("tests"); err != nil {
		log.Fatalln("Please `go run` from root folder")
	}
}

func testNamesFromCommandLineArguments() []string {
	testNames := []string{}

	if len(os.Args) > 1 {
		for _, testName := range os.Args[1:] {
			if strings.HasPrefix(testName, TESTS_FOLDER+"/") {
				testName = testName[len(TESTS_FOLDER+"/"):]
			}

			testName = strings.TrimSuffix(testName, "/")
			testNames = append(testNames, testName)
		}
	} else {
		allFilesUnderTests, err := os.ReadDir(TESTS_FOLDER)
		if err != nil {
			log.Fatalln("Unable to list 'tests' folder", err)
		}

		for _, file := range allFilesUnderTests {
			if strings.HasPrefix(file.Name(), "0") {
				testNames = append(testNames, file.Name())
			}
		}
	}

	return testNames
}

func buildFakeTerraformRegistry() {
	architecture := runtime.GOOS + "_" + runtime.GOARCH
	registry_dir := "tests/fake_registry/terraform-registry.env0.com/env0/env0/6.6.6/" + architecture

	err := os.MkdirAll(registry_dir, 0755)
	if err != nil {
		log.Fatalln("Unable to create registry folder ", registry_dir, " error: ", err)
	}

	data, err := os.ReadFile("terraform-provider-env0")
	if err != nil {
		log.Fatalln("Unable to read provider binary: did you build it?", err)
	}

	err = os.WriteFile(registry_dir+"/terraform-provider-env0", data, 0755)
	if err != nil {
		log.Fatalln("Unable to write: ", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get current working dir", err)
	}

	terraformRc := fmt.Sprintf(`
provider_installation {
  filesystem_mirror {
    path    = "%s/tests/fake_registry"
    include = ["terraform-registry.env0.com/*/*"]
  }
  direct {
	exclude = ["terraform-registry.env0.com/*/*"]
  }
}`, cwd)

	err = os.WriteFile("tests/terraform.rc", []byte(terraformRc), 0644)
	if err != nil {
		log.Fatalln("Unable to write: ", err)
	}

	err = os.Setenv("TF_CLI_CONFIG_FILE", path.Join(cwd, "tests", "terraform.rc"))
	if err != nil {
		log.Fatalln("Unable to set env:", err)
	}
}
