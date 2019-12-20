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
