package kubetest

import (
	"anchorctl/pkg/resultaggregator"
	"strings"
)

func (aj *jsonTest) test(res *resource) (*resultaggregator.TestRun) {
	testRun := &resultaggregator.TestRun{
		TestType: "JSON Path Desc",
		Desc:     aj.JSONPath + " == " + aj.Value,
		Passed:   false,
		Invalid:  false,
	}
	objIterms, err := res.ObjectRef.getObject(aj.client)

	if err != nil {
		testRun.Invalid = true
		return testRun
	}

	passed, err :=  assertJSONPath(objIterms, aj.JSONPath, aj.Value)
	if err != nil {
		testRun.Invalid = true
		return testRun
	}
	testRun.Passed = passed
	return testRun
}

func (av *validationTest) test(res *resource) (*resultaggregator.TestRun) {
	testRun := &resultaggregator.TestRun{
		TestType: "Validation Desc",
		Desc:     "Contains error: " + av.ContainsResponse,
		Passed:   false,
		Invalid:  false,
	}

	_, err := res.Manifest.apply(true)
	if err != nil && strings.Contains(err.Error(), av.ContainsResponse) {
		log.InfoWithFields(map[string]interface{}{
			"test":             "AssertValidation",
			"containsResponse": av.ContainsResponse,
			"status":           "PASSED",
		}, "AssertValidation throws the expected error.")
		testRun.Passed = true
		return testRun
	}

	log.WarnWithFields(map[string]interface{}{
		"test":     "AssertValidation",
		"expected": av.ContainsResponse,
		"got":      "Error: " + err.Error(),
		"status":   "FAILED",
	}, "AssertValidation Failed")

	return testRun
}

func (am *mutationTest) test(res *resource) (*resultaggregator.TestRun) {
	testRun := &resultaggregator.TestRun{
		TestType: "Mutation Desc",
		Desc:     am.JSONPath + " == " + am.Value,
		Passed:   false,
		Invalid:  false,
	}

	if valid := res.Manifest.valid(); !valid {
		testRun.Invalid = true
		return testRun
	}

	objectMetadata, err := res.Manifest.apply(false)
	if err != nil {
		testRun.Invalid = true
		return testRun
	}

	jsonTestResource := &resource{
		ObjectRef: *objectMetadata,
	}

	jpTest := &jsonTest{
		JSONPath: am.JSONPath,
		Value:    am.Value,
		client:   am.client,
	}

	jsonPathTest := jpTest.test(jsonTestResource)

	testRun.Invalid = jsonPathTest.Invalid
	testRun.Passed = jsonPathTest.Passed
	return testRun
}
