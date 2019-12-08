package kubetest

import (
	"github.com/covarity/anchorctl/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/kubernetes"
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

	kubeTest, err := decodeTestFile(client, testfile)
	if err != nil {
		log.Fatal(err, "Unable to decode test file")
	}

	executeLifecycle(kubeTest.Spec.Lifecycle.PostStart, client)

	results := runTests(client, kubeTest)
	results.threshold = threshold
	results.total = len(kubeTest.Spec.Tests)
	results.print()

	executeLifecycle(kubeTest.Spec.Lifecycle.PreStop, client)

	log.Info("Passed", "true", "Passed Test")
}

func runTests(client *kubernetes.Clientset, kubeTest *kubeTest) *testResult {
	res := testResult{}
	for _, i := range kubeTest.Spec.Tests {
		switch i.Type {

		case "AssertJSONPath":
			var jsonTestObj *jsonTest
			err := mapstructure.Decode(i.Spec, &jsonTestObj)
			if err != nil {
				res.invalid++
				log.Warn("Error", err.Error(), "Decoding AssertJSONPath returned error")
				continue
			}

			jsonTestObj.client = client
			runKubeTester(jsonTestObj, i, &res)

		case "AssertValidation":
			var validationTest *validationTest
			err := mapstructure.Decode(i.Spec, &validationTest)
			if err != nil {
				res.invalid++
				log.Warn("Error", err.Error(), "Decoding AssertValidation returned error")
				continue
			}

			validationTest.client = client
			runKubeTester(validationTest, i, &res)

		case "AssertMutation":
			var mutationTest *mutationTest
			err := mapstructure.Decode(i.Spec, &mutationTest)
			if err != nil {
				res.invalid++
				log.Warn("Error", err.Error(), "Decoding AssertMutation returned error")
				continue
			}

			mutationTest.client = client
			runKubeTester(mutationTest, i, &res)

		}
	}
	return &res
}

func runKubeTester(kubetester kubeTester, i test, res *testResult) {
	if result, err := kubetester.test(i.Resource); err != nil {
		res.invalid++
	} else if result == false {
		res.failed++
	} else {
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
