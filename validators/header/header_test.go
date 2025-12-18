package header

import (
	"comeva/validators"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func FixtureValidator() *HeaderValidator {
	return NewHeaderValidator([]string{"feat"}, []string{"frontend"}, []string{"Add"}, [2]int{}, [2]int{})
}

// FixtureValidTypeScopeVerb returns a valid set of the mentioned values.
func FixtureValidTypeScopeVerb() [3]string {
	return [3]string{"feat", "frontend", "Add"}
}

type HeaderConstructor struct {
	ValidFormatting bool
}



func (constr HeaderConstructor) WithScope(typ, scope, verb string) []string {

	if constr.ValidFormatting {
		return []string{
			fmt.Sprintf("%s(%s)!: %s search", typ, scope, verb),
			fmt.Sprintf("%s(%s): %s search", typ, scope, verb),
		}
	}
	return []string{
		fmt.Sprintf("(%s)!: %s search", scope, verb),
		fmt.Sprintf("(%s): %s search", scope, verb),
	}

}

func (constr HeaderConstructor) WithoutScope(typ, verb string) []string {
	
	if constr.ValidFormatting {
		return []string{
			fmt.Sprintf("%s!: %s search", typ, verb),
			fmt.Sprintf("%s: %s search", typ, verb),
		}
	}
	return []string{
		fmt.Sprintf("%s()!: %s search", typ, verb),
		fmt.Sprintf("%s(): %s search", typ, verb),
		fmt.Sprintf("%s:%s search", typ, verb),
	}
}

// Header creates header examples.
// The typ argument sets the type, the scope argument the scope,
// and the verb argument the verb. The validity is therefore
// determined by the arguments, the validity of the FixtureHeader, and the used validator.
func (constr HeaderConstructor) Header(typ, scope, verb string) []string {
	return slices.Concat(constr.WithScope(typ, scope, verb), constr.WithoutScope(typ, verb))
}

// ValidHeader creates header examples which are valid.
func ValidHeader() []string {
	var header = HeaderConstructor{ValidFormatting: true}
	var valids = FixtureValidTypeScopeVerb()
	return header.Header(valids[0], valids[1], valids[2])
}

// InvalidFormattingHeader creates header examples whose formatting,
// but not type, scope, or verb is invalid.
func InvalidFormattingHeader() []string {
	var header = HeaderConstructor{ValidFormatting: false}
	var valids = FixtureValidTypeScopeVerb()
	return header.Header(valids[0], valids[1], valids[2])
}

// The Validate method works.
func TestHeaderValidatorValidate(t *testing.T) {
	for _, header := range ValidHeader() {
		t.Run(header, func(t *testing.T) {
			var headerValidator validators.ReaderValidator = FixtureValidator()
			var headerReader = strings.NewReader(header)
		
			var err error = headerValidator.Validate(headerReader)
			if err != nil {
				t.Errorf("Error matching: %v", err)
			}
		})
	}

}

// TestValidHeader tests whether a valid header format is correctly
// matched by the header validator.
func TestValidHeader(t *testing.T) {

	for _, header := range ValidHeader() {
		t.Run(header, func(t *testing.T) {
			var validator = FixtureValidator()
			var err error = validator.ValidateString(header)
			if err != nil {
				t.Errorf("Error matching: %v", err)
			}
		})
	}

}

// TestHeaderStringIsValid tests whether the String() return of a valid header created by
// a validator is the same as the original header value.
func TestHeaderStringIsValid(t *testing.T) {
	for _, header := range ValidHeader() {
		t.Run(header, func(t *testing.T) {
			var validator = FixtureValidator()
			validator.ValidateString(header)
			var got = validator.Header.String()
			if got != header {
				t.Errorf("'%s' is not equal to %s", got, header)
			}
		})
	}
}

func TestInvalidHeader(t *testing.T) {
	for _, header := range InvalidFormattingHeader() {
		t.Run(header, func(t *testing.T) {
			var validator = FixtureValidator()
			if validator.ValidateString(header) == nil {
				t.Errorf("invalid header '%s' was found to be valid", header)
			}
		})
	}
}

// ValidateHeaderLength validates correctly.
func TestValidateHeaderLength(t *testing.T) {
	var validator = FixtureValidator()
	var max = validator.headerLength[1]
	var header = make([]byte, max + 1)
	for i := range header {
		header[i] = 'a'
	}
	var headerString = string(header[:max])
	if validator.ValidateHeaderLength(headerString) != nil {
		t.Errorf("Valid length %d was found to be invalid with a max of %d", len(headerString), max)
	}

	headerString = string(header)
	if validator.ValidateHeaderLength(headerString) == nil {
		t.Errorf("Invalid length %d was found to be valid with a max of %d", len(headerString), max)
	}

}

// ValidateScope validates correctly.
func TestValidateScope(t *testing.T) {
	var valids = FixtureValidTypeScopeVerb()
	var valid, invalid = valids[1], "backend"

	var validator = FixtureValidator()
	if validator.ValidateScope(valid) != nil {
		t.Errorf("Valid scope %s was found to be invalid for scopes %v", valid, validator.scopes)
	}

	if validator.ValidateScope(invalid) == nil {
		t.Errorf("Invalid scope %s was found to be valid for scopes %v", invalid, validator.scopes)
	}

}

// ValidateType validates correctly.
func TestValidateType(t *testing.T) {
	var valids = FixtureValidTypeScopeVerb()
	var valid, invalid = valids[0], "imnotype"

	var validator = FixtureValidator()
	if validator.ValidateType(valid) != nil {
		t.Errorf("Valid type %s was found to be invalid for types %v", valid, validator.types)
	}

	if validator.ValidateType(invalid) == nil {
		t.Errorf("Invalid type %s was found to be valid for types %v", valid, validator.types)
	}
}

// ValidateDescription validates correctly.
func TestValidateDescription(t *testing.T) {
	var validator = FixtureValidator()
	var valid = "Add new feature"
	var invalidNoVerb = "No verb starting this description"
	var invalidTooLong = "Add some new features in this wonderfully amazing commit whose commit message is just too long"
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