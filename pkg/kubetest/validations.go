package kubetest

import "fmt"

func validateTestField(index int, test test) (bool, error) {

	requiredFields := map[string][]string{
		"AssertJSONPath":   {"jsonPath", "value"},
		"AssertValidation": {"containsResponse"},
		"AssertMutation":   {"jsonPath", "value"},
	}

	if test.Type == "" {
		log.WarnWithFields(map[string]interface{}{
			"number":        index,
			"testType":      test.Type,
			"requiredField": "type",
		}, "Does not contain required field.")
		return false, fmt.Errorf("does not contain required field")
	}

	for _, i := range requiredFields[test.Type] {
		if _, ok := test.Spec[i]; !ok {
			log.WarnWithFields(map[string]interface{}{
				"number":        index,
				"testType":      test.Type,
				"requiredField": i,
			}, "Does not contain required field.")
			return false, fmt.Errorf("does not contain required field")
		}
	}

	return true, nil
}

func (ob objectRef) valid() bool {
	if ob.Type == "" || ob.Spec.Kind == "" || ob.Spec.Namespace == "" ||
		ob.Spec.Labels == nil {

		log.WarnWithFields(map[string]interface{}{
			"resource": "objectRef",
			"expected": "Resource ObjectRef type, kind, namespace, label value and label key should be specified",
			"got":      "Type: " + ob.Type + " Kind: " + ob.Spec.Kind + " Namespace: " + ob.Spec.Namespace,
		}, "Failed getting the resource to apply.")

		return false
	}

	return true
}

func (mf manifest) valid() bool {
	if mf.Path == "" || mf.Action == "" {
		log.WarnWithFields(map[string]interface{}{
			"resource": "manifest",
			"expected": "Resource Manifest path and action should be specified",
			"got":      "Path: " + mf.Path + " Action: " + mf.Action,
		}, "Failed getting the resource to apply.")

		return false
	}

	return true
}
