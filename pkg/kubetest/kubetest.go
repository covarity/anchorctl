package kubetest

import (
	"anchorctl/pkg/kubetest/assertexec"
	"anchorctl/pkg/kubetest/assertjsonpath"
	"anchorctl/pkg/kubetest/assertmutation"
	"anchorctl/pkg/kubetest/assertvalidation"
	"anchorctl/pkg/kubetest/common"
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"fmt"
	"k8s.io/client-go/kubernetes"
)

type TestType string

var (
	AssertJSONPath   TestType = "AssertJSONPath"
	AssertValidation TestType = "AssertValidation"
	AssertMutation   TestType = "AssertMutation"
	AssertExec       TestType = "AssertExec"

	log *logging.Logger

	testImpl = map[TestType]func(clientset *kubernetes.Clientset, log *logging.Logger, input interface{}) (common.KubeTester, error){}
)

func init() {
	testImpl[AssertJSONPath] = assertjsonpath.Parse
	testImpl[AssertExec] = assertexec.Parse
	testImpl[AssertValidation] = assertvalidation.Parse
	testImpl[AssertMutation] = assertmutation.Parse
}

// Assert function contains the logic to execute kubetests
func (kt *KubeTest) Assert() *resultaggregator.TestResult {

	resource.TestFilePath = kt.Opts.TestFilepath
	log = kt.Opts.Logger

	client, err := common.GetKubeClient(kt.Opts.Incluster, kt.Opts.Kubeconfig)
	if err != nil {
		log.Fatal(err, "Unable to get kubernetes client")
	}

	log.Info("state", "Post-Start", "Starting post start lifecycle")

	resource.ExecuteLifecycle(kt.Lifecycle.PostStart)

	log.Info("state", "Tests", "Starting kube tests")

	results := runTests(client, kt)

	log.Info("state", "Pre-Stop", "Starting pre stop lifecycle")

	resource.ExecuteLifecycle(kt.Lifecycle.PreStop)

	log.Info("state", "Finished", "Finished Running Kube Tests")
	return results

}

func runTests(client *kubernetes.Clientset, kubeTest *KubeTest) *resultaggregator.TestResult {
	res := resultaggregator.NewTestResult(len(kubeTest.Tests), log)

	for _, test := range kubeTest.Tests {
		tester, err := getTest(test, client)
		if err != nil {
			res.AddInvalidTestRun(string(test.Type), err)
			continue
		}

		if valid, err := tester.Valid(&test.Resource); !valid {
			res.AddInvalidTestRun(string(test.Type), err)
			continue
		}

		res.AddRun(tester.Test(&test.Resource))
	}

	return res
}

func getTest(test test, client *kubernetes.Clientset) (common.KubeTester, error) {

	if constructTest, ok := testImpl[test.Type]; ok {
		theTest, err := constructTest(client, log, test.Spec)
		if err != nil {
			log.Warn("Type", string(test.Type), "Error parsing test")
			return nil, fmt.Errorf("Type", string(test.Type), "Error parsing test")
		}
		return theTest, nil
	}

	return nil, fmt.Errorf("Unknown Test type")
}
