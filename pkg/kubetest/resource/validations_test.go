package resource

//func TestValidateTestField(t *testing.T) {
//
//	log.SetVerbosity(0)
//
//	// Setup
//	successAssertJSONPath := makeTestObj("AssertJSONPath", map[string]interface{}{
//		"jsonPath": "Test.json.path",
//		"value":    "Test.json.path.value",
//	})
//	failAssertJSONPath := makeTestObj("AssertJSONPath", map[string]interface{}{
//		"jonPath": "Test.json.path",
//		"val":     "Test.json.path.value",
//	})
//
//	successAssertValidation := makeTestObj("AssertValidation", map[string]interface{}{
//		"containsResponse": "Some expected error",
//	})
//	failAssertValidation := makeTestObj("AssertValidation", map[string]interface{}{
//		"Response": "Some expected error",
//	})
//
//	successAssertMutation := makeTestObj("AssertMutation", map[string]interface{}{
//		"jsonPath": "Test.json.path",
//		"value":    "Test.json.path.value",
//	})
//	failAssertMutation := makeTestObj("AssertMutation", map[string]interface{}{
//		"jonPath": "Test.json.path",
//		"val":     "Test.json.path.value",
//	})
//
//	emptyTestType := makeTestObj("", map[string]interface{}{})
//
//	tables := []struct {
//		message string
//		obj     *test
//		result  bool
//	}{
//		{"Successfully validate AssertJSONPath", successAssertJSONPath, true},
//		{"Fail validate AssertJSONPath", failAssertJSONPath, false},
//		{"Successfully validate AssertValidation", successAssertValidation, true},
//		{"Fail validate AssertValidation", failAssertValidation, false},
//		{"Successfully validate AssertMutation", successAssertMutation, true},
//		{"Fail validate AssertMutation", failAssertMutation, false},
//		{"Fail empty Test type", emptyTestType, false},
//	}
//
//	for i, table := range tables {
//		result, _ := validateTestField(i, *table.obj)
//		assert.Equal(t, table.result, result, table.message)
//	}
//}
//
//func makeTestObj(testType string, spec map[string]interface{}) *test {
//	return &test{
//		Resource: Resource{},
//		Type:     testType,
//		Spec:     spec,
//	}
//}
