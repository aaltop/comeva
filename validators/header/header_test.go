package header

import (
	"comeva/validators"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var validator = NewHeaderValidator([]string{"feat"}, []string{"Add"}, [2]int{}, [2]int{})

// The validate method works.
func TestHeaderValidatorValidate(t *testing.T) {
	var headerValidator validators.ReaderValidator = validator
	var headerReader = strings.NewReader("feat: Add search")

	var err error = headerValidator.Validate(headerReader)
	if err != nil {
		t.Errorf("Error matching: %v", err)
	}
}

// TestFullHeader tests whether the "full" header format is correctly
// matched by the header validator. The full format includes a "scope"
// at the start: "<scope>: <message>".
func TestFullHeader(t *testing.T) {
	var headerString = "feat: Add search"

	var err error = validator.ValidateString(headerString)
	if err != nil {
		t.Errorf("Error matching: %v", err)
	}
	fmt.Println(validator.Header)
}

// Invalid scope is correctly rejected by the validator.
func TestFullHeaderInvalidScope(t *testing.T) {
	var headerString = "imnoscope: Add search"

	if validator.ValidateString(headerString) == nil {
		t.Error("Matching did not pick up incorrect scope")
	}
	fmt.Println(validator.Header)
}

// TestUnscopedHeader tests whether the header format without a scope is correctly
// matched by the header validator.
func TestUnscopedHeader(t *testing.T) {
	var err error = validator.ValidateString("Add search")
	if !(errors.As(err, &InvalidScopeError{})) {
		t.Errorf("Missing scope was not found to be invalid: %v", err)
	}

	// actually ensure that it also doesn't include the error when it's
	// not supposed to
	err = validator.ValidateString("feat: Add search")
	if errors.As(err, &InvalidScopeError{}) {
		t.Errorf("Included scope was found to be invalid: %v", err)
	}
	fmt.Println(validator.Header)
}

// ValidateHeaderLength validates correctly.
func TestValidateHeaderLength(t *testing.T) {
	var max = validator.headerLength[1]
	var header = make([]byte, max + 1)
	for i := range header {
		header[i] = 'a'
	}
	var headerString = string(header[:max])
	if validator.ValidateHeaderLength(headerString) != nil {
		t.Errorf("Valid length %d was found to be invalid with a max of %d", len(headerString), max)
	}

	header[len(header) - 1] = 'a'
	headerString = string(header)
	if validator.ValidateHeaderLength(headerString) == nil {
		t.Errorf("Invalid length %d was found to be valid with a max of %d", len(headerString), max)
	}

}

// ValidateScope validates correctly.
func TestValidateScope(t *testing.T) {
	var valid, invalid = "feat", "imnoscope"
	if validator.ValidateScope(valid) != nil {
		t.Errorf("Valid scope %s was found to be invalid for scopes %v", valid, validator.scopes)
	}
	if validator.ValidateScope(invalid) == nil {
		t.Errorf("Invalid scope %s was found to be valid for scopes %v", invalid, validator.scopes)
	}
}

// ValidateDescription validates correctly.
func TestValidateDescription(t *testing.T) {
	var valid = "Add new feature"
	var invalidNoVerb = "No verb starting this description"
	var invalidTooLong = "Add some new features in this wonderfully amazing commit that is too long"
	var err error
	if _, err = validator.ValidateDescription(valid); err != nil {
		t.Errorf("Valid description '%s' was found to be invalid", valid)
	}
	var invalidFormat = "Invalid description '%s' was found to be valid"
	if _, err = validator.ValidateDescription(invalidNoVerb); err == nil {
		t.Errorf(invalidFormat, invalidNoVerb)
	}
	if _, err = validator.ValidateDescription(invalidTooLong); err == nil {
		t.Errorf(invalidFormat, invalidTooLong)
	}
}