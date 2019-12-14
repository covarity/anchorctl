package kubetest

import (
	"anchorctl/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/kubernetes"
)

var requiredField = map[string][]string{
	"AssertJSONPath":   {"jsonPath", "value"},
	"AssertValidation": {"action", "filePath", "expectedError"},
	"AssertMutation":   {"action", "filePath", "jsonPath", "value"},
}

var log *logging.Logger
var testFilePath string

// Assert function contains the logic to execute kubetests
func Assert(logger *logging.Logger, threshold float64, incluster bool, kubeconfig, testfile string) {

	log = logger

	client, err := getKubeClient(incluster, kubeconfig)
	if err != nil {
		log.Fatal(err, "Unable to get kubernetes client")
	}
	testFilePath = testfile
	kubeTest, err := decodeTestFile(client, testFilePath)
	if err != nil {
		log.Fatal(err, "Unable to decode test file")
	}
	kubeTest.testFilePath = testfile

	executeLifecycle(kubeTest.Spec.Lifecycle.PostStart, client)

	results := runTests(client, kubeTest)
	results.threshold = threshold
	results.total = len(kubeTest.Spec.Tests)
	results.print()

	executeLifecycle(kubeTest.Spec.Lifecycle.PreStop, client)

	results.checkThresholdPass()

	log.Info("Passed", "true", "Passed Test")
}

func runTests(client *kubernetes.Clientset, kubeTest *kubeTest) *testResult {
	res := testResult{}
	for i, test := range kubeTest.Spec.Tests {
		switch test.Type {

		case "AssertJSONPath":
			var jsonTestObj *jsonTest
			res.addResultToRow(i, "AssertJSONPath")
			err := mapstructure.Decode(test.Spec, &jsonTestObj)
			res.addResultToRow(i, "JSONPath: "+jsonTestObj.JSONPath+" Value: "+jsonTestObj.Value)
			if err != nil {
				res.invalid++
				res.addResultToRow(i, "❌")
				log.Warn("Error", err.Error(), "Decoding AssertJSONPath returned error")
				continue
			}
			jsonTestObj.client = client
			runKubeTester(jsonTestObj, &test, &res, i)

		case "AssertValidation":
			var validationTest *validationTest
			res.addResultToRow(i, "AssertValidation")
			err := mapstructure.Decode(test.Spec, &validationTest)
			res.addResultToRow(i, "Expected Response: "+validationTest.ContainsResponse[:25]+"...")
			if err != nil {
				res.invalid++
				res.addResultToRow(i, "❌")
				log.Warn("Error", err.Error(), "Decoding AssertValidation returned error")
				continue
			}
			validationTest.client = client
			runKubeTester(validationTest, &test, &res, i)

		case "AssertMutation":
			var mutationTest *mutationTest
			res.addResultToRow(i, "AssertMutation")
			err := mapstructure.Decode(test.Spec, &mutationTest)
			res.addResultToRow(i, "JSONPath: "+mutationTest.JSONPath+" Value "+mutationTest.Value)
			if err != nil {
				res.invalid++
				res.addResultToRow(i, "❌")
				log.Warn("Error", err.Error(), "Decoding AssertMutation returned error")
				continue
			}
			mutationTest.client = client
			runKubeTester(mutationTest, &test, &res, i)
		}
	}
	return &res
}

func runKubeTester(kubetester kubeTester, test *test, res *testResult, i int) {
	if result, err := kubetester.test(&test.Resource); err != nil {
		res.addResultToRow(i, "❌")
		res.invalid++
	} else if result == false {
		res.addResultToRow(i, "❌")
		res.failed++
	} else {
		res.addResultToRow(i, "✅")
		res.passed++
	}
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
