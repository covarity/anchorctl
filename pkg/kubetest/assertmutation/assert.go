package assertmutation

import (
	"anchorctl/pkg/kubetest/assertjsonpath"
	"anchorctl/pkg/kubetest/common"
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/kubernetes"
)

var log *logging.Logger

func Parse(clientset *kubernetes.Clientset, logger *logging.Logger, input interface{}) (common.KubeTester, error) {
	var execTest *mutationTest
	log = logger
	err := mapstructure.Decode(input, &execTest)
	if err != nil {
		return nil, err
	}
	execTest.client = clientset
	return execTest, nil
}

type mutationTest struct {
	JSONPath string
	Value    string
	client   *kubernetes.Clientset
}

func (am *mutationTest) Valid(res *resource.Resource) (bool, error) {
	return common.ValidTest("MutationTest", map[string]string{
		"JSONPath":        am.JSONPath,
		"JSONValue":       am.Value,
		"Manifest.Path":   res.Manifest.Path,
		"Manifest.Action": res.Manifest.Action,
	})
}

func (am *mutationTest) Test(res *resource.Resource) *resultaggregator.TestRun {
	testRun := resultaggregator.NewTestRun("Mutation Desc", am.JSONPath+" == "+am.Value)

	objectMetadata, err := res.Manifest.Apply(false)
	if err != nil {
		testRun.Invalid = true
		return testRun
	}

	jsonTestResource := &resource.Resource{
		ObjectRef: *objectMetadata,
	}

	jpTest := assertjsonpath.New(am.client, am.JSONPath, am.Value)

	jsonPathTest := jpTest.Test(jsonTestResource)

	testRun.Invalid, testRun.Passed = jsonPathTest.Invalid, jsonPathTest.Passed
	return testRun
}
