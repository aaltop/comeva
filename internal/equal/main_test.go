package equal

import "testing"

// NilEqual correctly tests for nil values.
func TestNilEqual(t *testing.T) {
	var value, equal, notEqual string
	var nilStr *string
	value = "a string"
	equal = "a string"
	notEqual = "some other string"

	if !NilEqual[string](&value, &equal) {
		t.Errorf("Equal values found to be inequal: %v\n%v\n", value, equal)
	}

	if NilEqual[string](&value, &notEqual) {
		t.Errorf("Inequal values found to be equal: %v\n%v\n", value, notEqual)
	}

	if NilEqual[string](&value, nilStr) {
		t.Errorf("Non-nil found equal to nil: %v,\n%v\n", value, nilStr)
	}

}
