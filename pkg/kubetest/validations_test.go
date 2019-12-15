package kubetest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateTestField(t *testing.T) {

	log.SetVerbosity(0)

	// Setup
	successAssertJSONPath := &test{
		Resource: resource{},
		Type:     "AssertJSONPath",
		Spec: map[string]interface{}{
			"jsonPath": "test.json.path",
			"value": "test.json.path.value",
		},
	}

	failAssertJSONPath := &test{
		Resource: resource{},
		Type:     "AssertJSONPath",
		Spec: map[string]interface{}{
			"jonPath": "test.json.path",
			"val": "test.json.path.value",
		},
	}

	// Setup
	successAssertValidation := &test{
		Resource: resource{},
		Type:     "AssertValidation",
		Spec: map[string]interface{}{
			"containsResponse": "Some expected error",
		},
	}

	failAssertValidation := &test{
		Resource: resource{},
		Type:     "AssertValidation",
		Spec: map[string]interface{}{
			"Response": "Some expected error",
		},
	}

	// Setup
	successAssertMutation := &test{
		Resource: resource{},
		Type:     "AssertMutation",
		Spec: map[string]interface{}{
			"jsonPath": "test.json.path",
			"value": "test.json.path.value",
		},
	}

	failAssertMutation := &test{
		Resource: resource{},
		Type:     "AssertMutation",
		Spec: map[string]interface{}{
			"jonPath": "test.json.path",
			"val": "test.json.path.value",
		},
	}


	tables := []struct {
		message string
		obj  *test
		result bool
	}{
		{"Successfully validate AssertJSONPath", successAssertJSONPath, true},
		{"Fail validate AssertJSONPath", failAssertJSONPath, false},
		{"Successfully validate AssertValidation", successAssertValidation, true},
		{"Fail validate AssertValidation", failAssertValidation, false},
		{"Successfully validate AssertMutation", successAssertMutation, true},
		{"Fail validate AssertMutation", failAssertMutation, false},
	}

	for i, table := range tables {
		result, _ := validateTestField(i, *table.obj)
		assert.Equal(t, table.result, result, table.message)
	}
}
