package kubernetes

import (
	"fmt"
	"github.com/anchorageio/anchorctl/utils/logging"

	"github.com/olekukonko/tablewriter"
	"os"
)

var requiredField = map[string][]string{
	"AssertJSONPath":   {"jsonPath", "value"},
	"AssertValidation": {"action", "filePath", "expectedError"},
	"AssertMutation":   {"action", "filePath", "jsonPath", "value"},
}

var log *logging.Logger

func Assert(logger *logging.Logger, threshold float64, incluster bool, kubeconfig, testfile string) {
	log = logger

	client, err := getKubeClient(incluster, kubeconfig)
	if err != nil {
		log.Fatal(err, "Unable to get kubernetes client")
	}

	fmt.Println(logger.GetLevel())

	kubeTest, err := decodeTestFile(testfile)
	if err != nil {
		log.Fatal(err, "Unable to decode test file")
	}

	object, objectGetErr := getObject(client, kubeTest)
	if err != nil {
		log.Error(err, "Unable to get target kube object")
	}

	var passed, failed, invalid int
	total := len(kubeTest.Tests)

	for i, test := range kubeTest.Tests {

		if !validateTestField(requiredField, i, test) {
			invalid++
			continue
		}

		switch test["type"] {

		case "AssertJSONPath":
			if result, err := assertJsonpath(object, test["jsonPath"], test["value"]); err != nil || objectGetErr != nil {
				invalid++
			} else if result == true {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertJSONPath", "jsonPath": test["jsonPath"],}).Info("Passed test")
				passed++
			} else {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertJSONPath", "jsonPath": test["jsonPath"],}).Warn("Failed test")
				failed++
			}

		case "AssertValidation":
			if result := assertValidation(client, test["action"], test["filePath"], test["expectedError"]); result == true {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertValidation", "filePath": test["filePath"],}).Info("Passed test")
				passed++
			} else {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertValidation", "filePath": test["filePath"],}).Warn("Failed test")
				failed++
			}

		case "AssertMutation":
			if result, err := assertMutation(client, test["action"], test["filePath"], test["jsonPath"], test["value"]); err != nil {
				invalid++
			} else if result == true {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertMutation", "jsonPath": test["jsonPath"],}).Info("Passed test")
				passed++
			} else {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertMutation", "jsonPath": test["jsonPath"],}).Warn("Failed test")
				failed++
			}

		default:
			log.Warn("testType", test["type"], "Invalid test")
			invalid++
		}
	}

	fmt.Println()
	successRatio := ((float64(passed) / float64(total)) * 100)

	data := [][]string{
		{"Total", fmt.Sprintf("%d", total)},
		{"Passed", fmt.Sprintf("%d", passed)},
		{"Failed", fmt.Sprintf("%d", failed)},
		{"Invalid", fmt.Sprintf("%d", invalid)},
		{"Expected Coverage", fmt.Sprintf("%.2f", threshold)},
		{"Actual Coverage", fmt.Sprintf("%.2f", successRatio)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tests", "Number"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render() // Send output

	fmt.Println()

	if successRatio < threshold {
		log.Fatal(fmt.Errorf("Expected %.2f, Got %.2f", threshold, successRatio), "Failed Test Threshold")
	}

	log.Info("Passed", "true", "Passed Test")
}

func validateTestField(requiredField map[string][]string, index int, test map[string]string) bool {

	for _, i := range requiredField[test["type"]] {
		if _, ok := test[i]; !ok {
			log.WarnWithFields(map[string]interface{}{
				"number":        index,
				"testType":      test["type"],
				"requiredField": i,
			}, "Does not contain required field.")
			return false
		}
	}

	return true
}
