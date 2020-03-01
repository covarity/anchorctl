package assertjsonpath

import (
	"anchorctl/pkg/kubetest/common"
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"bytes"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/jsonpath"
)

var log *logging.Logger

func Parse(clientset *kubernetes.Clientset, logger *logging.Logger, input interface{}) (common.KubeTester, error) {
	var jsonTestObj *jsonTest
	log = logger
	err := mapstructure.Decode(input, &jsonTestObj)
	if err != nil {
		return nil, err
	}
	jsonTestObj.client = clientset
	return jsonTestObj, nil
}

type jsonTest struct {
	JSONPath string
	Value    string
	client   *kubernetes.Clientset
}

func New(clientset *kubernetes.Clientset, path, value string) *jsonTest {
	return &jsonTest{
		JSONPath: path,
		Value:    value,
		client:   clientset,
	}
}

func (aj *jsonTest) Valid(res *resource.Resource) (bool, error) {
	return common.ValidTest("JSONTest", map[string]string{
		"JSONValue":                aj.Value,
		"JSONPath":                 aj.JSONPath,
		"ObjectRef.Spec.Kind":      res.ObjectRef.Spec.Kind,
		"ObjectRef.Spec.Namespace": res.ObjectRef.Spec.Namespace,
	})
}

func (aj *jsonTest) Test(res *resource.Resource) *resultaggregator.TestRun {
	testRun := resultaggregator.NewTestRun("JSON Path Desc", aj.JSONPath+" == "+aj.Value)
	objItems, err := res.ObjectRef.GetObject(aj.client)

	if err != nil {
		testRun.Invalid = true
		testRun.Details = err.Error()
		return testRun
	}

	passed, err := assertJSONPath(objItems, aj.JSONPath, aj.Value)
	if err != nil {
		testRun.Invalid = true
		return testRun
	}
	testRun.Passed = passed
	return testRun
}

func assertJSONPath(objs []runtime.Object, path, value string) (bool, error) {
	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	passed := true

	err := jp.Parse("{" + path + "}")
	if err != nil {
		log.Error(err, "Cannot parse JSONPath")
		return false, err
	}

	buf := new(bytes.Buffer)

	for _, i := range objs {
		err = jp.Execute(buf, i)
		if err != nil {
			log.Error(err, "Cannot execute JSONPath")
			passed = false
			break
		} else if buf.String() != value {
			log.WarnWithFields(map[string]interface{}{
				"jsonpath": path,
				"expected": value,
				"got":      buf.String(),
				"status":   "FAILED",
			}, "Failed asserting jsonpath on obj")
			passed = false
			break
		}

		buf.Reset()
	}

	log.InfoWithFields(map[string]interface{}{
		"test":   "AssertJSONPath",
		"path":   path,
		"status": "PASSED",
	}, "JSON Path matches expected value.")

	return passed, err
}
