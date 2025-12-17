package header

import (
	"comeva/validators"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// The validate method works.
func TestHeaderValidatorValidate(t *testing.T) {
	var headerValidator validators.Validator = NewHeaderValidator([]string{"frontend"}, []string{"Add"}, [2]int{}, [2]int{})
	var headerReader = strings.NewReader("frontend: Add search")

	var err = headerValidator.Validate(headerReader)
	if err != nil {
		t.Errorf("Error matching: %v", err)
	}
}

// TestFullHeader tests whether the "full" header format is correctly
// matched by the header validator. The full format includes a "scope"
// at the start: "<scope>: <message>".
func TestFullHeader(t *testing.T) {
	var headerValidator = NewHeaderValidator([]string{"frontend"}, []string{"Add"}, [2]int{}, [2]int{})
	var headerString = "frontend: Add search"

	var err = headerValidator.ValidateString(headerString)
	if err != nil {
		t.Errorf("Error matching: %v", err)
	}
	fmt.Println(headerValidator.Header)
}

func TestFullHeaderInvalidScope(t *testing.T) {
	var headerValidator = NewHeaderValidator([]string{"backend"}, []string{"Add"}, [2]int{}, [2]int{})
	var headerString = "frontend: Add search"

	var err = headerValidator.ValidateString(headerString)
	if err == nil {
		t.Error("Matching did not pick up incorrect scope")
	}
	fmt.Println(headerValidator.header)
}

// TestUnscopedHeader tests whether the header format without a scope is correctly
// matched by the header validator.
func TestUnscopedHeader(t *testing.T) {
	var headerValidator = NewHeaderValidator([]string{"frontend"}, []string{"Add"}, [2]int{}, [2]int{})
	var err = headerValidator.ValidateString("Add search")
	if !(errors.As(err, &InvalidScopeError{})) {
		t.Errorf("Missing scope was not found to be invalid: %v", err)
	}

	// actually ensure that it also doesn't include the error when it's
	// not supposed to
	err = headerValidator.ValidateString("frontend: Add search")
	if errors.As(err, &InvalidScopeError{}) {
		t.Errorf("Included scope was found to be invalid: %v", err)
	}
	fmt.Println(headerValidator.Header)
}
