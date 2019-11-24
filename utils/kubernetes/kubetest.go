package kubernetes

import (
	"fmt"
	"github.com/anchorageio/anchorctl/utils/logging"
	log "github.com/sirupsen/logrus"
)

func Assert(logging *logging.Logger, threshold float64, incluster bool, kubeconfig, testfile string) error {
	client, err := getKubeClient(incluster, kubeconfig)
	logger := logging.GetLogger()

	if err != nil {
		logger.Fatal("Unable to get kubernetes client")
	}

	fmt.Println(logger.GetLevel())

	kubeTest, err := decodeTestFile(testfile)
	if err != nil {
		logger.Fatal("Unable to decode test file")
	}

	object, objectGetErr := getObject(client, kubeTest)
	if err != nil {
		logger.Error("Unable to decode test file")
	}

	var passed, failed, invalid int
	total := len(kubeTest.Tests)

	for i, test := range kubeTest.Tests {

		if !validateTestField(logging, i, test) {
			invalid++
			continue
		}

		switch test["type"]{

		case "AssertJSONPath":
			if result, err := assertJsonpath(logging, object, test["jsonPath"], test["value"]); err != nil || objectGetErr != nil {
				invalid++
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertJSONPath", "jsonPath": test["jsonPath"],}).Error("Invalid test")
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
			if result, err := assertMutation(client, logging, test["action"], test["filePath"], test["jsonPath"], test["value"]); err != nil {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertMutation", "jsonPath": test["jsonPath"],}).Error("Invalid test")
				invalid++
			} else if result == true {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertMutation", "jsonPath": test["jsonPath"],}).Info("Passed test")
				passed++
			} else {
				logger.WithFields(log.Fields{ "number": i, "testType": "AssertMutation", "jsonPath": test["jsonPath"],}).Warn("Failed test")
				failed++
			}

		default:
			logger.WithFields(log.Fields{ "number": i, "testType": test["type"],}).Error("Invalid test")
			invalid++
		}
	}

	fmt.Printf("Total tests: %d \n", total)
	fmt.Printf("Passed tests: %d \n", passed)
	fmt.Printf("Failed tests: %d \n", failed)
	fmt.Printf("Invalid tests: %d \n", invalid)

	successRatio := ((float64(passed) / float64(total)) * 100)

	fmt.Printf("successRatio: %.2f \n", successRatio)

	if successRatio < threshold {
		return fmt.Errorf("Below threshold: expected: %.2f, actual %.2f", threshold, successRatio)
	}

	return nil
}

func validateTestField(logging *logging.Logger, index int, test map[string]string) bool {

	requiredField := map[string][]string{
		"AssertJSONPath" : {"jsonPath", "value"},
		"AssertValidation": {"action", "filePath", "expectedError"},
		"AssertMutation": {"action", "filePath", "jsonPath", "value"},
	}

	for _, i := range requiredField[test["type"]] {
		if _, ok := test[i]; !ok {
			logging.Log.WithFields(log.Fields{
				"number": index,
				"testType": test["type"],
				"requiredField": i,
			}).Warn("Does not contain required field.")
			return false
		}
	}

	return true
}