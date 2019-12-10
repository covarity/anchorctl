package kubetest

import (
	"fmt"
	"strings"
)

func (aj *jsonTest) test(res *resource) (bool, error) {

	objIterms, err := res.ObjectRef.getObject(aj.client)
	if err != nil {
		return false, err
	}

	return assertJSONPath(objIterms, aj.JSONPath, aj.Value)
}

func (av *validationTest) test(res *resource) (bool, error) {

	_, err := res.Manifest.apply(av.client, true)
	if err != nil && strings.Contains(err.Error(), av.ContainsResponse) {
		log.InfoWithFields(map[string]interface{}{
			"test":             "AssertValidation",
			"containsResponse": av.ContainsResponse,
			"status":           "PASSED",
		}, "AssertValidation throws the expected error.")
		return true, nil
	}

	log.WarnWithFields(map[string]interface{}{
		"test":     "AssertValidation",
		"expected": av.ContainsResponse,
		"got":      "Error: " + err.Error(),
		"status":   "FAILED",
	}, "AssertValidation Failed")

	return false, nil
}

func (am *mutationTest) test(res *resource) (bool, error) {

	if valid := res.Manifest.valid(); valid != true {
		return false, fmt.Errorf("Invalid Manifest to apply")
	}

	objectMetadata, err := res.Manifest.apply(am.client, false)
	if err != nil {
		return false, err
	}

	jsonTestResource := &resource{
		ObjectRef: *objectMetadata,
	}

	jpTest := &jsonTest{
		JSONPath: am.JSONPath,
		Value:    am.Value,
		client:   am.client,
	}

	return jpTest.test(jsonTestResource)
}
