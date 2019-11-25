package kubernetes

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Assert(cmd *cobra.Command, threshold float64, kubeconfig, testfile string) error {
	client, err := getKubeClient(false, kubeconfig)

	if err != nil {
		return fmt.Errorf("Could not get client", err.Error())
	}

	kubeTest, err := decodeTestFile(testfile)
	if err != nil {
		return fmt.Errorf("Could not decode test file", err.Error())
	}

	object, err := getObject(client, kubeTest)
	if err != nil {
		return fmt.Errorf("Failed getting object", err)
	}

	var passed, failed, invalid int
	total := len(kubeTest.Tests)

	for _, test := range kubeTest.Tests {

		if !validateTestField(cmd, test) {
			invalid++
			continue
		}

		switch test["type"]{

		case "AssertJSONPath":
			if result, err := assertJsonpath(cmd, object, test["jsonPath"], test["value"]); err != nil {
				invalid++
				cmd.PrintErrln("AssertJsonPath Failed", err)
			} else if result == true {
				passed++
			} else {
				failed++
			}

		case "AssertValidation":
			if result := assertValidation(client, test["action"], test["filePath"], test["expectedError"]); result == true {
				passed++
			} else {
				failed++
			}

		case "AssertMutation":
			if result, err := assertMutation(client, cmd, test["action"], test["filePath"], test["jsonPath"], test["value"]); err != nil {
				invalid++
				cmd.PrintErrln("AssertMutation Failed", err)
			} else if result == true {
				passed++
			} else {
				failed++
			}

		default:
			cmd.Println(test["type"] + " is not a valid test type")
			invalid++
		}
	}

	cmd.Printf("Total tests: %d \n", total)
	cmd.Printf("Passed tests: %d \n", passed)
	cmd.Printf("Failed tests: %d \n", failed)
	cmd.Printf("Invalid tests: %d \n", invalid)

	successRatio := ((float64(passed) / float64(total)) * 100)

	cmd.Printf("successRatio: %.2f \n", successRatio)

	if successRatio < threshold {
		return fmt.Errorf("Below threshold: expected: %.2f, actual %.2f", threshold, successRatio)
	}

	return nil
}

func validateTestField(cmd *cobra.Command, test map[string]string) bool {

	requiredField := map[string][]string{
		"AssertJSONPath" : {"jsonPath", "value"},
		"AssertValidation": {"action", "filePath", "expectedError"},
		"AssertMutation": {"action", "filePath", "jsonPath", "value"},
	}

	for _, i := range requiredField[test["type"]] {
		if _, ok := test[i]; !ok {
			cmd.Println(test["type"] + " requires `" + i + "` as a key.")
			return false
		}
	}

	return true
}