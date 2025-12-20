package trailer

import (
	"bufio"
	"comeva/validators"
	"fmt"
	"strings"
	"testing"
)

func FixtureValidator() *TrailerValidator {
	var validator, e = NewTrailerValidator(make(KeyMap), make(KeyMap), 2, [2]uint{0, 80})
	if e != nil {
		panic(e)
	}
	return validator
}

// ValidKey returns formatting-wise valid keys. Also see [InvalidKey].
func ValidKey() []string {
	return []string{
		"BREAKING-CHANGE",
		"Dummy-Key",
	}
}

// ValidKeyValue returns valid key-value pairs. These do NOT have continuation
// lines, by the definition used here; see [ValidTrailer] for those.
func ValidKeyValue() []string {
	return []string{
		ValidKey()[0] + ": changed a thing breakingly",
		ValidKey()[1] + ": I'm just an example",
	}
}

// ValidTrailer returns valid git trailer examples.
func ValidTrailer() []string {
	var validKeyValue = ValidKeyValue()
	for i := range validKeyValue {
		validKeyValue[i] += "\n  Here's a continuation line"
	}
	return validKeyValue
}

func ValidTrailerBlock() []string {
	var validTrailer = ValidTrailer()
	var joinedTrailer = strings.Join(validTrailer, "\n")
	return append(ValidTrailer(), joinedTrailer)
}

// InvalidKey returns formatting-wise invalid keys. Also see [ValidKey].
func InvalidKey() []string {
	return []string{
		"BREAKING CHANGE",
		"Dummy_Key",
	}
}

func InvalidKeyValue() []string {

	return []string{
		"BREAKING CHANGE: key's got a space, no good",
		ValidKey()[0] + ":no space after the colon, no good",
		ValidKey()[0] + ": otherwise fine, but the content in the trailers should also be limited in the number of columns on one line",
	}
}

func InvalidTrailer() []string {
	return []string{
		ValidKey()[0] + ": the value\n  Should have\n a consistent indentation",
		ValidKey()[0] + ": the value\n   Should have\n   the specified indentation",
	}
}

func InvalidTrailerBlock() []string {
	return append(InvalidTrailer(),
		fmt.Sprintf("%s: this is fine\n\n%s: but there shouldn't be gaps between trailers", ValidKey()[0], ValidKey()[1]))
}

// RunTestValues runs validation function `f` against valid and invalid sets of
// values `valid` and `invalid`.
func RunTestValues(f func(tr string) error, valid, invalid []string, t *testing.T) {
	for i, trailer := range valid {
		t.Run(fmt.Sprintf("Valid %d", i), func(t *testing.T) {
			if e := f(trailer); e != nil {
				t.Errorf("Unexpected error: %v", e)
			}
		})
	}

	for i, trailer := range invalid {
		t.Run(fmt.Sprintf("Invalid %d", i), func(t *testing.T) {
			if e := f(trailer); e == nil {
				t.Errorf("Invalid value '%s' was found to be valid", trailer)
			}
		})
	}
}

func RunTestTrailerBlock(f func(tr string) error, t *testing.T) {
	RunTestValues(f, ValidTrailerBlock(), InvalidTrailerBlock(), t)
}

func RunTestTrailer(f func(tr string) error, t *testing.T) {
	RunTestValues(f, ValidTrailer(), InvalidTrailer(), t)
}

// The Validate method works.
// func TestValidate(t *testing.T) {
// 	var validator validators.ReaderValidator = FixtureValidator()

// 	var validation = func(tr string) error {
// 		var reader = strings.NewReader(tr)
// 		return validator.Validate(reader)
// 	}

// 	RunTestTrailerBlock(validation, t)

// }

func TestValidateString(t *testing.T) {
	var validator validators.StringValidator = FixtureValidator()

	var validation = func(tr string) error {
		return validator.ValidateString(tr)
	}

	RunTestTrailerBlock(validation, t)
}

func TestValidateScanner(t *testing.T) {
	var validator interface {
		ValidateScanner(scn *bufio.Scanner) error
	} = FixtureValidator()

	var validation = func(tr string) error {
		var scanner = bufio.NewScanner(strings.NewReader(tr))
		return validator.ValidateScanner(scanner)
	}

	RunTestTrailerBlock(validation, t)
}

// Key-value pairs are validated correctly. Should NOT consider trailers with
// continuation lines as valid, by the definition used here.
func TestValidateKeyValue(t *testing.T) {
	var validator = FixtureValidator()

	var validation = func(tr string) error {
		var _, e = validator.ValidateKeyValue(tr, 1)
		return e
	}

	RunTestValues(validation, ValidKeyValue(), InvalidKeyValue(), t)
}

// Keys' formatting is validated correctly.
func TestValidateKeyFormat(t *testing.T) {
	var validator = FixtureValidator()

	var validation = func(key string) error {
		return validator.ValidateKey(key, 1)
	}

	RunTestValues(validation, ValidKey(), InvalidKey(), t)

}

// Optional keys are validated correctly. If a key is specified as optional,
// it does not have to be in the trailer, but any key that is in the trailer should be in either the
// required or optional keys; if both are empty, any key is valid (assuming correct
// formatting). See also [TestRequiredKey].
func TestOptionalKey(t *testing.T) {
	var validator = FixtureValidator()

	var validKey = ValidKey()

	var key = validKey[0]
	if e := validator.ValidateKey(key, 1); e != nil {
		t.Errorf("Unexpected error for valid key %s: %v", key, e)
	}

	var optionalKey = validKey[1]
	validator.optionalKeys.Set(optionalKey, "")
	if e := validator.ValidateKey(key, 1); e == nil {
		t.Errorf("Invalid key '%s' was found to be valid", key)
	}
	if e := validator.ValidateKey(optionalKey, 1); e != nil {
		t.Errorf("Unexpected error for valid key %s: %v", optionalKey, e)
	}

	// does not error if a key set as optional is not present in trailer
	validator.optionalKeys.Set(key, "")
	if e := validator.ValidateString(fmt.Sprintf("%s: dummy key", key)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

}

// Required keys are validated correctly. If a key is specified as required,
// it has to be in the trailer. See also [TestOptionalKey].
func TestRequiredKey(t *testing.T) {
	var validator = FixtureValidator()
	var validKey = ValidKey()

	var key = validKey[0]
	if e := validator.ValidateKey(key, 1); e != nil {
		t.Errorf("Unexpected error for valid key %s: %v", key, e)
	}

	var requiredKey = validKey[1]
	validator.requiredKeys.Set(requiredKey, "")
	// because one key is not specified as required and there are no
	// optional keys, while the other is required, both should be valid
	for _, k := range []string{key, requiredKey} {
		if e := validator.ValidateKey(k, 1); e != nil {
			t.Errorf("Unexpected error for valid key %s: %v", k, e)
		}
	}

	// required key HAS to be in the trailer
	if e := validator.ValidateString(fmt.Sprintf("%s: non-required key", key)); e == nil {
		t.Errorf("Expected error for missing required key %s", requiredKey)
	}

}

// Value continuations are validated correctly.
func TestValidateValueContinuation(t *testing.T) {
	var validator = FixtureValidator()

	var cont = "that has a continuation"

	var whitespaces = []struct {
		WhiteSpace  string
		ExpectError bool
	}{
		{
			"",
			true,
		},
		{
			" ",
			true,
		},
		{
			"  ",
			false,
		},
		{
			"   ",
			true,
		},
		{
			"\t",
			true,
		},
		{
			"\n",
			true,
		},
	}

	for i, whitespaceAndExpected := range whitespaces {
		t.Run(fmt.Sprintf("whitespace_only_%d", i), func(t *testing.T) {
			if e := validator.ValidateValueContinuation(whitespaceAndExpected.WhiteSpace, 1); e == nil {
				t.Errorf("Continuation should not be solely whitespace")
			}
		})
		t.Run(fmt.Sprintf("indent_%d", i), func(t *testing.T) {
			var continuation = fmt.Sprintf("%s%s", whitespaceAndExpected.WhiteSpace, cont)
			if e := validator.ValidateValueContinuation(continuation, 1); (e == nil) == whitespaceAndExpected.ExpectError {
				t.Errorf("Expected error: %t Got error: %v", whitespaceAndExpected.ExpectError, e)
			}
		})
	}

	if e := validator.ValidateValueContinuation("  Correct indentation, but again, the line should be limited in the number of columns", 1); e == nil {
		t.Error("Invalid line with too many characters was found to be valid")
	}

}

// Successful validation adds the trailers to the list in the validator.
func TestValidationAddsToTrailers(t *testing.T) {
	var validator = FixtureValidator()

	var testLength = func(trailers []Trailer, expectedLength int) {
		if len(trailers) != expectedLength {
			t.Errorf("Expected trailers to be of length %d, was %d", expectedLength, len(trailers))
		}
	}

	testLength(validator.Trailers, 0)
	for i, trailer := range ValidTrailer() {
		validator.ValidateString(trailer)
		testLength(validator.Trailers, i+1)
	}
}

// The value that is parsed matches the original value.
func TestValueParsedWithoutModification(t *testing.T) {
	var validator = FixtureValidator()

	var expected = "A trailer\n  with a continuation line,\n  two in fact."
	var trailerBlock = ValidKey()[0] + ": " + expected

	var e = validator.ValidateString(trailerBlock)
	if e != nil {
		t.Fatalf("Unexpected error: %v", e)
	}

	var received = validator.Trailers[0].Value
	if received != expected {
		t.Errorf("Parsed value did not match original\nParsed:\n%s\nOriginal:\n%s", received, expected)
	}
}
