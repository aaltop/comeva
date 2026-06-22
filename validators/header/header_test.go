package header

import (
	"comeva/internal/regexp"
	testingUtils "comeva/internal/testing"
	"comeva/validators"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func FixtureValidator() *HeaderValidator {
	var h, e = NewHeaderValidator([]string{"feat"}, []string{"frontend"}, []string{"Add"}, [2]uint{0, 80})
	if e != nil {
		panic(true)
	}
	return h
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

func toSubMatch(value string) regexp.SubMatch {
	return regexp.SubMatch{Match: value}
}

// The Validate method works.
func TestHeaderValidatorValidate(t *testing.T) {
	for _, header := range ValidHeader() {
		t.Run(header, func(t *testing.T) {
			var headerValidator validators.ReaderValidator = FixtureValidator()
			var headerReader = strings.NewReader(header)

			var err []validators.ValidatorErrorChild = headerValidator.Validate(headerReader)
			if len(err) != 0 {
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
			var err = validator.ValidateString(header)
			if len(err) != 0 {

				t.Errorf("Error matching: %v\nHeader: %v\n", err, header)
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

// ValidateLength validates correctly.
func TestValidateLength(t *testing.T) {
	var validator = FixtureValidator()
	var max = validator.lineLength.Upper
	var min = validator.lineLength.Lower
	var header = make([]byte, max+1)
	for i := range header {
		header[i] = 'a'
	}

	var valid, invalid = string(header[:max]), string(header)
	if validator.ValidateLength(valid) != nil {
		t.Errorf("Valid length %d was found to be invalid with a [min, max] of [%d, %d]", len(valid), min, max)
	}

	if validator.ValidateLength(invalid) == nil {
		t.Errorf("Invalid length %d was found to be valid with a [min, max] of [%d, %d]", len(invalid), min, max)
	}

	var emptyValidator = NewDefaultHeaderValidator()
	for _, v := range []string{valid, invalid} {
		if emptyValidator.ValidateLength(v) != nil {
			t.Errorf("length %d was found to be invalid where all lengths should be valid", len(v))
		}
	}

}

// ValidateScope validates correctly.
func TestValidateScope(t *testing.T) {
	var valids = FixtureValidTypeScopeVerb()
	var valid, invalid = toSubMatch(valids[1]), toSubMatch("backend")

	var validator = FixtureValidator()
	if validator.validateScope(valid) != nil {
		t.Errorf("Valid scope %s was found to be invalid for scopes %v", valid.Match, validator.scopes)
	}

	if validator.validateScope(invalid) == nil {
		t.Errorf("Invalid scope %s was found to be valid for scopes %v", invalid.Match, validator.scopes)
	}

	var emptyValidator = NewDefaultHeaderValidator()
	for _, v := range []regexp.SubMatch{valid, invalid} {
		if emptyValidator.validateType(v) != nil {
			t.Errorf("Scope '%s' was found to be invalid where all scopes should be valid", v.Match)
		}
	}

}

// ValidateType validates correctly.
func TestValidateType(t *testing.T) {
	var valids = FixtureValidTypeScopeVerb()
	var valid, invalid = toSubMatch(valids[0]), toSubMatch("imnotype")

	var validator = FixtureValidator()
	if validator.validateType(valid) != nil {
		t.Errorf("Valid type %s was found to be invalid for types %v", valid.Match, validator.types)
	}

	if validator.validateType(invalid) == nil {
		t.Errorf("Invalid type %s was found to be valid for types %v", valid.Match, validator.types)
	}

	var emptyValidator = NewDefaultHeaderValidator()
	for _, v := range []regexp.SubMatch{valid, invalid} {
		if emptyValidator.validateType(v) != nil {
			t.Errorf("Type '%s' was found to be invalid where all types should be valid", v.Match)
		}
	}
}

// ValidateDescription validates correctly.
func TestValidateDescription(t *testing.T) {
	var validator = FixtureValidator()
	var valid = toSubMatch("Add new feature")
	var invalidNoVerb = toSubMatch("No verb starting this description")
	var err []validators.ValidatorErrorChild
	if _, err = validator.validateDescription(valid); len(err) != 0 {
		t.Errorf("Valid description '%s' was found to be invalid", valid.Match)
	}
	var invalidFormat = "Invalid description '%s' was found to be valid"
	if _, err = validator.validateDescription(invalidNoVerb); len(err) == 0 {
		t.Errorf(invalidFormat, invalidNoVerb)
	}
}

// ValidateVerb validates correctly
func TestValidateVerb(t *testing.T) {
	var validator = FixtureValidator()
	var valids = FixtureValidTypeScopeVerb()
	var valid, invalid = valids[2], "Triangulate"

	if validator.validateVerb(valid) != nil {
		t.Errorf("Valid verb '%s' was found to be invalid", valid)
	}

	if validator.validateVerb(invalid) == nil {
		t.Errorf("Invalid verb '%s' was not found to be invalid", invalid)
	}

	var emptyValidator = NewDefaultHeaderValidator()
	for _, v := range []string{valid, invalid} {
		if emptyValidator.validateVerb(v) != nil {
			t.Errorf("Verb '%s' was found to be invalid where all verbs should be valid", v)
		}
	}
}

// The count of errors received from the validator and the errors set in the
// validator match.
func TestErrorCountMatches(t *testing.T) {

	for i, trailer := range InvalidFormattingHeader() {
		t.Run(fmt.Sprintf("test %d", i+1), func(t *testing.T) {
			var validator = FixtureValidator()

			var returnedErrors = validator.ValidateString(trailer)
			var returnedErrorsCount = len(returnedErrors)

			if returnedErrorsCount == 0 {
				t.Fatal("Should return at least one error")
			}

			var containedErrorsCount = len(validator.Errors)
			if containedErrorsCount != returnedErrorsCount {
				t.Error(testingUtils.ValueMismatch("error count", returnedErrorsCount, containedErrorsCount))
			}
		})
	}
}
