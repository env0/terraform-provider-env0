package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
)

const TESTS_FOLDER = "tests/integration"

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

	for _, testName := range testNames {
		wg.Add(1)

		go func(testName string) {
			if destroyMode == "DESTROY_ONLY" {
				terraformDestroy(testName)
			} else {
				success, err := runTest(testName, destroyMode != "NO_DESTROY")
				if !success {
					log.Fatalf("Halting due to test '%s' failure: %s\n", testName, err)
				}
			}

			wg.Done()
		}(testName)
	}

	wg.Wait()
}
func compileProvider() error {
	cmd := exec.Command("go", "build")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
func runTest(testName string, destroy bool) (bool, error) {
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

	_, err := terraformCommand(testName, "init")
	if err != nil {
		return false, err
	}

	_, _ = terraformCommand(testName, "fmt")

	if destroy {
		defer terraformDestroy(testName)
	}

	_, err = terraformCommand(testName, "apply", "-auto-approve", "-var", "second_run=0")
	if err != nil {
		return false, err
	}

	_, err = terraformCommand(testName, "apply", "-auto-approve", "-var", "second_run=1")
	if err != nil {
		return false, err
	}

	expectedOutputs, err := readExpectedOutputs(testName)
	if err != nil {
		return false, err
	}

	outputsBytes, err := terraformCommand(testName, "output", "-json")
	if err != nil {
		return false, err
	}

	outputs, err := bytesOfJsonToStringMap(outputsBytes)
	if err != nil {
		return false, err
	}

	for key, expectedValue := range expectedOutputs {
		value, ok := outputs[key]
		if !ok {
			log.Println("Error: Expected terraform output ", key, " but no such output was created")

			return false, nil
		}

		if value != expectedValue {
			log.Printf("Error: Expected output of '%s' to be '%s' but found '%s'\n", key, expectedValue, value)

			return false, nil
		}

		log.Printf("Verified expected '%s'='%s' in %s", key, value, testName)
	}

	if destroy {
		_, err = terraformCommand(testName, "destroy", "-auto-approve", "-var", "second_run=0")
		if err != nil {
			return false, err
		}
	}

	log.Println("Successfully finished running test ", testName)

	return true, nil
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
	destroy.Dir = TESTS_FOLDER + "/" + testName

	if err := destroy.Run(); err != nil {
		log.Println("WARNING: error running tofu destroy")
	}

	log.Println("Done running tofu destroy in", testName)
}

func terraformCommand(testName string, arg ...string) ([]byte, error) {
	cmd := exec.Command("tofu", arg...)
	cmd.Dir = TESTS_FOLDER + "/" + testName
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "INTEGRATION_TESTS=1")
	cmd.Env = append(cmd.Env, "TF_LOG_PROVIDER=info")

	var output, errOutput bytes.Buffer

	cmd.Stderr = bufio.NewWriter(&errOutput)
	cmd.Stdout = bufio.NewWriter(&output)

	log.Println("Running tofu ", arg, " in ", testName)

	err := cmd.Run()

	log.Println(errOutput.String())

	if err != nil {
		log.Println("error running tofu ", arg, " in ", testName, " error: ", err)
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
