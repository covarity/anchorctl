package assertvalidation

import (
	"anchorctl/pkg/kubetest/common"
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/kubernetes"
	"strings"
)

var log *logging.Logger

type validationTest struct {
	ContainsResponse string `yaml:"containsResponse"`
	client           *kubernetes.Clientset
}

func Parse(clientset *kubernetes.Clientset, logger *logging.Logger, input interface{}) (common.KubeTester, error) {
	var validationTest *validationTest
	log = logger
	err := mapstructure.Decode(input, &validationTest)
	if err != nil {
		return nil, err
	}
	validationTest.client = clientset
	return validationTest, nil
}

func (av *validationTest) Valid(res *resource.Resource) (bool, error) {
	return common.ValidTest("ValidationTest", map[string]string{
		"ContainsResponse": av.ContainsResponse,
		"Manifest.Action":  res.Manifest.Action,
		"Manifest.Path":    res.Manifest.Path,
	})
}

func (av *validationTest) Test(res *resource.Resource) *resultaggregator.TestRun {
	testRun := resultaggregator.NewTestRun("Validation Desc", "Contains error: "+av.ContainsResponse)

	_, err := res.Manifest.Apply(true)
	if err != nil {
		if strings.Contains(err.Error(), av.ContainsResponse) {
			log.InfoWithFields(map[string]interface{}{
				"Test":             "AssertValidation",
				"containsResponse": av.ContainsResponse,
				"status":           "PASSED",
			}, "AssertValidation throws the expected error.")
			testRun.Passed = true
			return testRun
		}

		log.WarnWithFields(map[string]interface{}{
			"Test":     "AssertValidation",
			"expected": av.ContainsResponse,
			"got":      "Error: " + err.Error(),
			"status":   "FAILED",
		}, "AssertValidation Failed")

		testRun.Invalid = true
		testRun.Details = err.Error()
		return testRun
	}

	log.WarnWithFields(map[string]interface{}{
		"Test":     "AssertValidation",
		"expected": av.ContainsResponse,
		"got":      "Did not throw validation error",
		"status":   "FAILED",
	}, "AssertValidation Failed")

	testRun.Details = "Did not throw validation error"

	return testRun
}
