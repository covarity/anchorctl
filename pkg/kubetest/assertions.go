package kubetest

import (
	"bytes"
	"fmt"

	"k8s.io/client-go/util/jsonpath"
)

func (aj *jsonTest) test(res resource) (bool, error) {

	obj, err := res.ObjectRef.getObject(aj.client)
	if err != nil {
		return false, err
	}

	return assertJSONPath(obj, aj.JsonPath, aj.Value)
}

func assertJSONPath(obj interface{}, path, value string) (bool, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")
	passed := true
	if err != nil {
		log.Error(err, "Cannot parse JSONPath")
		return false, err
	}

	buf := new(bytes.Buffer)

	objects := getSlice(obj)

	for _, i := range objects {
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

func (av *validationTest) test(res resource) (bool, error) {

	_, _, err := res.Manifest.apply(av.client)
	if err != nil && err.Error() == av.ExpectedError {
		log.InfoWithFields(map[string]interface{}{
			"test":          "AssertValidation",
			"expectedError": av.ExpectedError,
			"status":        "PASSED",
		}, "AssertValidation throws the expected error.")
		return true, nil
	}

	log.WarnWithFields(map[string]interface{}{
		"test":     "AssertValidation",
		"expected": av.ExpectedError,
		"got":      "Error: " + err.Error(),
		"status":   "FAILED",
	}, "AssertValidation Failed")

	return false, nil

}

func (am *mutationTest) test(res resource) (bool, error) {

	if valid := res.Manifest.valid(); valid != true {
		return false, fmt.Errorf("Invalid Manifest to apply")
	}

	_, obj, err := res.Manifest.apply(am.client)
	if err != nil {
		return false, err
	}

	return assertJSONPath(obj, am.JsonPath, am.Value)
}
