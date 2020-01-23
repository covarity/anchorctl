package kubetest

import (
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

var log *logging.Logger
var testFilePath string

// Assert function contains the logic to execute kubetests
func (kt *KubeTest) Assert() *resultaggregator.TestResult {

	testFilePath = kt.Opts.TestFilepath
	log = kt.Opts.Logger

	client, err := getKubeClient(kt.Opts.Incluster, kt.Opts.Kubeconfig)
	if err != nil {
		log.Fatal(err, "Unable to get kubernetes client")
	}

	log.Info("Status", "Post-Start", "Starting post start lifecycle")

	executeLifecycle(kt.Lifecycle.PostStart)

	log.Info("Status", "Tests", "Starting kube tests")

	results := runTests(client, kt)

	log.Info("Status", "Pre-Stop", "Starting pre stop lifecycle")

	executeLifecycle(kt.Lifecycle.PreStop)

	log.Info("Status", "Finished", "Finished Running Kube Tests")
	return results

}

func runTests(client *kubernetes.Clientset, kubeTest *KubeTest) *resultaggregator.TestResult {
	res := resultaggregator.NewTestResult(len(kubeTest.Tests))

	for i, test := range kubeTest.Tests {
		switch test.Type {

		case "AssertJSONPath":
			var jsonTestObj *jsonTest

			err := mapstructure.Decode(test.Spec, &jsonTestObj)
			if err != nil {
				res.AddInvalidTestRun("AssertJSONPath", err)
				log.Warn("Error", err.Error(), "Decoding AssertJSONPath returned error")
				continue
			}
			jsonTestObj.client = client
			res.AddRun(jsonTestObj.test(&test.Resource))

		case "AssertValidation":
			var validationTest *validationTest

			err := mapstructure.Decode(test.Spec, &validationTest)
			if err != nil {
				res.AddInvalidTestRun("AssertValidation", err)
				log.Warn("Error", err.Error(), "Decoding AssertValidation returned error")
				continue
			}
			validationTest.client = client
			res.AddRun(validationTest.test(&test.Resource))

		case "AssertMutation":
			var mutationTest *mutationTest

			err := mapstructure.Decode(test.Spec, &mutationTest)
			if err != nil {
				res.AddInvalidTestRun("AssertMutation", err)
				log.Warn("Error", err.Error(), "Decoding AssertMutation returned error")
				continue
			}
			mutationTest.client = client
			res.AddRun(mutationTest.test(&test.Resource))

		default:
			res.AddInvalidTestRun("Invalid", fmt.Errorf("Unknown Test type"))
			log.Warn("index", strconv.Itoa(i), "Unknown test type")
		}
	}

	return res
}
