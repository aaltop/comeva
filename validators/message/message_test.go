package message

import (
	"comeva/validators"
	"comeva/validators/trailer"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func FixtureValidator() *MessageValidator {
	return NewDefaultMessageValidator()
}

func ValidHeader() []string {
	return []string{
		"feat(main)!: Add feature",
	}
}

func ValidBody() []string {
	return []string{
		"This body gives some further details about the commit,\nand has a number of lines.",
	}
}

func ValidTrailerBlock() []string {
	return []string{
		"BREAKING-CHANGE: This commit introduces a breaking change,\n  which needs a long description to explain.",
	}
}

func CreateMessage(header, body, trailer string) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, body, trailer)
}

func ValidMessage() []string {
	return []string{
		CreateMessage(ValidHeader()[0], ValidBody()[0], ValidTrailerBlock()[0]),
	}
}

// The Validate method works.
func TestValidate(t *testing.T) {
	var validator validators.ReaderValidator = FixtureValidator()

	var message = ValidMessage()[0]
	if e := validator.Validate(strings.NewReader(message)); e != nil {
		t.Errorf("Unexpected error: %v\nMessage:\n%s\n", e, message)
	}
}

// Each part (header, body, trailer) is picked out successfully.
func TestFindParts(t *testing.T) {
	var validator = FixtureValidator()

	var header, body, trailer string = ValidHeader()[0], ValidBody()[0], ValidTrailerBlock()[0]
	var message string = CreateMessage(header, body, trailer)

	if e := validator.ValidateString(message); e != nil {
		t.Errorf("Unexpected error: %v\nMessage:\n%s\n", e, message)
	}

	var newExpectedErrorString = func(typ, expected, received string) string {
		return fmt.Sprintf("%s mismatch:\nExpected:\n%s\nReceived:\n%s\n", typ, expected, received)
	}

	var receivedHeader string = validator.HeaderValidator.Header.String()
	if header != receivedHeader {
		t.Error(newExpectedErrorString("Header", header, receivedHeader))
	}

	var receivedBody string = validator.BodyValidator.Body
	if body != receivedBody {
		t.Error(newExpectedErrorString("Body", body, receivedHeader))
	}

	var receivedTrailer string = ""
	for i, tr := range validator.TrailerValidator.Trailers {
		if i == 0 {
			receivedTrailer += tr.String()
		} else {
			receivedTrailer += "\n" + tr.String()
		}
	}
	if trailer != receivedTrailer {
		t.Error(newExpectedErrorString("Trailer", trailer, receivedTrailer))
	}

	if !validator.FoundHeader || !validator.FoundBody || !validator.FoundTrailer {
		t.Errorf("All Found<part> should be true, found:\nHeader: %t\nBody: %t\nTrailer: %t",
			validator.FoundHeader,
			validator.FoundBody,
			validator.FoundTrailer)
	}
}

// Message having only a header and trailer are still validated correctly.
func TestOnlyHeaderAndTrailer(t *testing.T) {
	var validator = FixtureValidator()
	var header, trailer string = ValidHeader()[0], ValidTrailerBlock()[0]
	var message string = fmt.Sprintf("%s\n\n%s", header, trailer)

	if e := validator.ValidateString(message); e != nil {
		t.Errorf("Unexpected error: %v\nMessage:\n%s\n", e, message)
	}

	if len(validator.HeaderValidator.Header.Type) == 0 {
		t.Errorf("Header was not parsed properly")
	}

	if len(validator.TrailerValidator.Trailers) == 0 {
		t.Errorf("Trailer block was not parsed properly")
	}
}

// If a breaking change exclamation mark is not accompanied by a BREAKING-CHANGE key (trailer)
// and vice versa, a validation error is returned.
func TestBreakingChangesTogether(t *testing.T) {
	var validator = FixtureValidator()

	var validMessage = "feat!: Add feature\n\nBREAKING-CHANGE: broke thing"
	if e := validator.ValidateString(validMessage); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	var exclamationMissingMessage = "feat: Add feature\n\nBREAKING-CHANGE: broke thing"
	var e error
	if e = validator.ValidateString(exclamationMissingMessage); e == nil {
		t.Errorf("Invalid message\n%s\nwas found to be valid (should have breaking change exclamation in header)", exclamationMissingMessage)
	}
	if !errors.As(e, &BreakingChangeError{}) {
		t.Errorf("Expected a BreakingChangeError")
	}

	var trailerMissingMessage = "feat!: Add feature whut"
	if e = validator.ValidateString(trailerMissingMessage); e == nil {
		t.Errorf(
			"Invalid message\n%s\nwas found to be valid (should have BREAKING-CHANGE key in trailer block, had %v)",
			trailerMissingMessage,
			validator.TrailerValidator.Trailers)
	}
	if !errors.As(e, &BreakingChangeError{}) {
		t.Errorf("Expected a BreakingChangeError")
	}
}

// It is checked that any required keys are present (mostly for the situation
// where a trailer is not present at all, but required keys are specified).
func TestTrailerChecksKeys(t *testing.T) {
	var validator = FixtureValidator()

	var requiredKeys trailer.KeyMap = make(trailer.KeyMap)
	requiredKeys.Set("Effect", "")
	validator.TrailerValidator, _ = trailer.NewTrailerValidator(
		requiredKeys, make(trailer.KeyMap), 2, [2]uint{0, 80},
	)

	var header, trailer string = "feat: Add some thing", "Effect: behavioural"
	var messageNoTrailer = fmt.Sprintf("%s\n", header)
	var messageNoTrailer2 = fmt.Sprintf("%s\n\n%s", header, "This is a body text\nwith some text")
	var messageInvalidTrailer = fmt.Sprintf("%s\n\n%s", header, "Some-Key: some info")
	var message = fmt.Sprintf("%s\n\n%s", header, trailer)

	var e error
	if e = validator.ValidateString(message); e != nil {
		t.Errorf("Valid message\n%s\nfound to be invalid: %v", message, e)
	}

	var invalidString = "Invalid message (missing required trailer key 'Effect')\n%s\nfound to be valid"

	if e = validator.ValidateString(messageInvalidTrailer); e == nil {
		t.Errorf(invalidString, messageInvalidTrailer)
	}

	if e = validator.ValidateString(messageNoTrailer); e == nil {
		t.Errorf(invalidString, messageNoTrailer)
	}

	if e = validator.ValidateString(messageNoTrailer2); e == nil {
		t.Errorf(invalidString, messageNoTrailer2)
	}
}

// An error should be returned when there is no empty line between
// the body and trailer
func TestNoGapBetweenBodyAndTrailer(t *testing.T) {
	var validator = FixtureValidator()

	var message = fmt.Sprintf("%v\n\n%v\n%v", "feat: Add commit", "A body\nwith a few lines", "Key: value")

	var e error
	if e = validator.ValidateString(message); e == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(e.Error(), "before trailer block") {
		t.Errorf("Unexpected error message: %v\n", e.Error())
	}
}
