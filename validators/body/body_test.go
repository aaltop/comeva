package body

import (
	testingUtils "comeva/internal/testing"
	"comeva/internal/utils"
	"comeva/validators"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func FixtureValidator() *BodyValidator {
	var validator, e = NewBodyValidator([2]uint{0, 80})
	if e != nil {
		panic(e)
	}
	return validator
}

type BodyConstructor struct {
}

// Line returns a string which represents a line of the given length
// in the body of a commit message.
func (constructor BodyConstructor) Line(lineLength uint) string {
	var word = "word"
	var line []string
	for uint(len(line)) < lineLength {
		line = append(line, word)
	}
	return strings.Join(line, " ")[:lineLength]
}

// Paragraph returns a string which represents a paragraph of a commit message
// with lines of the given length.
func (constructor BodyConstructor) Paragraph(lineLength []uint) string {
	var lines = make([]string, len(lineLength))
	for i, length := range lineLength {
		lines[i] = constructor.Line(length)
	}
	return strings.Join(lines, "\n")
}

// Body returns a string which represents the body of a commit message.
// Each of the slices in the argument represents one paragraph, and
// the values of that slice represent the length of lines in that
// paragraph.
func (constructor BodyConstructor) Body(lineLength [][]uint) string {
	var paragraphs = make([]string, len(lineLength))
	for i, lengths := range lineLength {
		paragraphs[i] = constructor.Paragraph(lengths)
	}
	return strings.Join(paragraphs, "\n")
}

func ValidBody() []string {
	var constructor = &BodyConstructor{}
	var par = []uint{77, 78, 79, 55}
	return []string{constructor.Body([][]uint{par}), ""}
}

func InvalidBody() []string {
	return []string{"\n" + ValidBody()[0]}
}

// The Validate method works.
func TestValidate(t *testing.T) {
	var validator validators.ReaderValidator = FixtureValidator()

	for i, body := range ValidBody() {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var reader = strings.NewReader(body)
			if e := validator.Validate(reader); e != nil {
				t.Errorf("%v", e)
			}
		})
	}

}

func TestValidateString(t *testing.T) {
	var validator = FixtureValidator()

	for i, body := range ValidBody() {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if e := validator.ValidateString(body); e != nil {
				t.Errorf("%v", e)
			}
		})
	}
}

func TestInvalidBody(t *testing.T) {
	var validator = FixtureValidator()

	for i, body := range InvalidBody() {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if e := validator.ValidateString(body); e == nil {
				t.Error("Invalid body was found to be valid")
			}
		})
	}
}

func TestValidateLine(t *testing.T) {
	var validator = FixtureValidator()

	var constructor = &BodyConstructor{}

	if e := validator.ValidateLine(constructor.Line(80), 1); e != nil {
		t.Error("Valid line was found to be invalid")
	}

	if e := validator.ValidateLine(constructor.Line(81), 1); e == nil {
		t.Error("Invalid length line was found to be valid")
	}

	// Bounds with 0,0 should mean no length check
	validator.lineLength = utils.Bounds[uint]{}
	if e := validator.ValidateLine(constructor.Line(81), 1); e != nil {
		t.Errorf("Line of length 81 was found to be invalid where all lengths should be valid: %v", e)
	}
}

// The validated body is correctly set in the validator.
func TestSetBodySetInValidator(t *testing.T) {
	var validator = FixtureValidator()

	var bodyString = ValidBody()[0]
	var body []string = strings.Split(bodyString, "\n")
	if e := validator.ValidateString(bodyString); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	if slices.Compare(body, validator.Body) != 0 {
		t.Errorf("Body mismatch:\nExpected:\n%s\nReceived:\n%s\n", body, validator.Body)
	}
}

// The count of errors received from the validator and the errors set in the
// validator match.
func TestErrorCountMatches(t *testing.T) {

	for i, trailer := range InvalidBody() {
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
