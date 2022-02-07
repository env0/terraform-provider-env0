package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const TESTS_FOLDER = "tests/integration"

func main() {
	err := compileProvider()
	if err != nil {
		log.Fatalln("Couldn't compile go")
		return
	}
	printTerraformVersion()
	makeSureRunningFromProjectRoot()
	testNames := testNamesFromCommandLineArguments()
	log.Println(len(testNames), " tests to run")
	buildFakeTerraformRegistry()
	destroyMode := os.Getenv("DESTROY_MODE")
	for _, testName := range testNames {
		if destroyMode == "DESTROY_ONLY" {
			terraformDestory(testName)
		} else {
			success, err := runTest(testName, destroyMode != "NO_DESTROY")
			if !success {
				log.Fatalln("Halting due to test failure:", err)
			}
		}
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
	terraformCommand(testName, "fmt")
	if destroy {
		defer terraformDestory(testName)
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
	expectedBytes, err := ioutil.ReadFile(path.Join(TESTS_FOLDER, testName, "expected_outputs.json"))
	if err != nil {
		log.Println("Test folder for ", testName, " does not contain expected_outputs.json", err)
		return nil, err
	}
	return bytesOfJsonToStringMap(expectedBytes)
}

func bytesOfJsonToStringMap(data []byte) (map[string]string, error) {
	var stringMapUncasted map[string]interface{}
	err := json.Unmarshal(data, &stringMapUncasted)
	if err != nil {
		log.Println("Unable to parse json:", err)
		log.Println("** JSON Input **")
		log.Println(string(data[:]))
		log.Println("******")
		return nil, err
	}
	result := map[string]string{}
	for key, valueUncasted := range stringMapUncasted {
		switch value := valueUncasted.(type) {
		case string:
			result[key] = value
		case map[string]interface{}:
			result[key] = value["value"].(string)
		}
	}
	return result, nil
}

func terraformDestory(testName string) {
	log.Println("Running destroy to clean up in", testName)
	destroy := exec.Command("terraform", "destroy", "-auto-approve", "-var", "second_run=0")
	destroy.Dir = TESTS_FOLDER + "/" + testName
	destroy.CombinedOutput()
	log.Println("Done running terraform destroy in", testName)
}

func terraformCommand(testName string, arg ...string) ([]byte, error) {
	cmd := exec.Command("terraform", arg...)
	cmd.Dir = TESTS_FOLDER + "/" + testName
	log.Println("Running terraform ", arg, " in ", testName)
	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		log.Println("error running terraform ", arg, " in ", testName, " error: ", err, " output: ", output)
	} else {
		log.Println("Completed successfully terraform ", arg, " in ", testName)
	}
	return outputBytes, err
}

func printTerraformVersion() {
	versionString, err := exec.Command("terraform", "version").Output()
	if err != nil {
		log.Fatalln("Unable to invoke terraform. Perhaps it's not in PATH?", err)
	}
	log.Println("Terraform version: ", string(versionString))
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
			if strings.HasSuffix(testName, "/") {
				testName = testName[:len(testName)-1]
			}
			testNames = append(testNames, testName)
		}
	} else {
		allFilesUnderTests, err := ioutil.ReadDir(TESTS_FOLDER)
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
	data, err := ioutil.ReadFile("terraform-provider-env0.exe")
	if err != nil {
		log.Fatalln("Unable to read provider binary: did you build it?", err)
	}
	err = ioutil.WriteFile(registry_dir+"/terraform-provider-env0", data, 0755)
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
	err = ioutil.WriteFile("tests/terraform.rc", []byte(terraformRc), 0644)
	if err != nil {
		log.Fatalln("Unable to write: ", err)
	}

	err = os.Setenv("TF_CLI_CONFIG_FILE", path.Join(cwd, "tests", "terraform.rc"))
	if err != nil {
		log.Fatalln("Unable to set env:", err)
	}
}
