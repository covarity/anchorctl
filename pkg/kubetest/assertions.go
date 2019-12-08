package kubetest

import (
	"bytes"
	"k8s.io/client-go/util/jsonpath"
)

func (aj *jsonTest) test(ref objectRef) (bool, error) {

	obj, err := getObject(aj.client, &ref)
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
			}, "Failed asserting jsonpath on obj")
			passed = false
			break
		}
		buf.Reset()
	}

	return passed, err

}

func (av *validationTest) test(ref objectRef) (bool, error) {
	_, _, err := applyAction(av.client, ref.Spec.Path, ref.Spec.Action)

	if err != nil && err.Error() == av.ExpectedError {
		log.Info("Expected", av.ExpectedError, "AssertValidation Passed")
		return true, nil
	}

	log.WarnWithFields(map[string]interface{}{
		"expected": av.ExpectedError,
		"got":      err.Error(),
	}, "AssertValidation Failed")
	return false, nil

}

func (am *mutationTest) test(ref objectRef) (bool, error) {
	_, obj, err := applyAction(am.client, ref.Spec.Path, ref.Spec.Action)
	if err != nil {
		log.Warn("Error", err.Error(), "AssertMutation Failed")
		return false, err
	}

	return assertJSONPath(obj, am.JsonPath, am.Value)
}
