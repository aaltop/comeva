// package body contains functionalities for validating the body
// of a git commit message.
package body

import (
	"bufio"
	"comeva/utils"
	"comeva/validators"
	"errors"
	"io"
	"strings"
)

// BodyValidator validates a commit message's body.
type BodyValidator struct {
	// lineLength describes the lower and upper bound of a line's length.
	lineLength utils.Bounds[uint]

	// Body contains the body.
	Body string
}

// Reset resets any content set during validation.
func (validator *BodyValidator) Reset() {
	validator.Body = ""
}

// NewDefaultBodyValidator creates the base BodyValidator.
func NewDefaultBodyValidator() (validator *BodyValidator) {
	return &BodyValidator{}
}

// NewBodyValidator returns a new BodyValidator, returning a non-nil error
// if the creation is unsuccessful.
//
// lineLength should contain a lower
// and upper bound for a single line's length, where [0,0] means
// that all values for line lengths are accepted.
func NewBodyValidator(lineLength [2]uint) (validator *BodyValidator, e error) {
	e = nil
	bounds, e := utils.NewBounds(lineLength[0], lineLength[1], false, false)
	return &BodyValidator{lineLength: bounds}, e
}

// NewBodyValidatorWithDefaults returns a [BodyValidator] with default values.
func NewBodyValidatorWithDefaults() (validator *BodyValidator) {
	var e error
	validator, e = NewBodyValidator([2]uint{0, 80})
	if e != nil {
		panic(e)
	}
	return
}

// Equal reports whether the two BodyValidators are equal.
func (validator *BodyValidator) Equal(other *BodyValidator) bool {
	if validator == nil || other == nil {
		return validator == other
	}

	return validator.lineLength.Equal(other.lineLength)
}

// ValidateLine returns a non-nil error if the passed line
// does not fulfill the requirements of a commit message's body's line.
func (validator *BodyValidator) ValidateLine(line string, lineNum uint) (e error) {

	// for both at zero, don't check length
	var lower, upper uint = validator.lineLength.Lower, validator.lineLength.Upper
	if lower == 0 && upper == 0 {
		return
	}

	var lenLine uint = uint(len(line))
	if !validator.lineLength.Contains(lenLine) {
		var validatorError = validators.ValidatorError{
			Line: lineNum, MessagePart: validators.MessageParts.Body}
		e = validators.InvalidLineLengthError{
			ValidatorError: validatorError, Expected: validator.lineLength, Received: lenLine}
	}
	return
}

func (validator *BodyValidator) Validate(reader io.Reader) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

// ValidateString validates a message body block.
func (validator *BodyValidator) ValidateString(possibleBody string) (e error) {
	return validator.ValidateStringWithLine(possibleBody, 1)
}

// ValidateStringWithLine validates a message body block.
//
// `startLine` > 0 specifies the line
// on which the string starts in the original message, assuming that the trailer
// is a part of a longer message. This is currently only relevant for accurate
// reporting of the line number on which a validation error occurs.
func (validator *BodyValidator) ValidateStringWithLine(possibleBody string, startLine uint) (e error) {
	return validator.ValidateScannerWithLine(bufio.NewScanner(strings.NewReader(possibleBody)), startLine-1)
}

// ValidateScanner validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner.
func (validator *BodyValidator) ValidateScanner(scanner *bufio.Scanner) (e error) {
	return validator.ValidateScannerWithLine(scanner, 0)
}

// ValidateScanner validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner. `timesScanned` specifies
// the number of times .Scan() has been called on `scanner`, representing the
// line at which the the scanner is.
func (validator *BodyValidator) ValidateScannerWithLine(scanner *bufio.Scanner, timesScanned uint) (e error) {

	validator.Reset()

	var line string
	e = nil
	var errs []error
	var scner = utils.CountingScanner{Scanner: scanner}
	scner.TimesScanned = timesScanned

	// empty body is fine
	if !scner.Scan() {
		return nil
	}
	line = scner.Text()

	// The body should not start on an empty line. It is expected that
	// the header and body are separated by exactly one blank line,
	// and this line is not part of either. After this, there should
	// be actual text.
	if len(strings.TrimSpace(line)) == 0 {
		errs = append(errs, validators.InvalidLineError{
			Reason: "First line of body should not be empty",
			Line:   scner.TimesScanned})
	}

	// Not putting that many constraints on the body, in my opinion it
	// should be allowed to be fairly free-form. Conventional commits
	// says (v1.0.0 Spec 7.) "A commit body is free-form and MAY consist
	// of any number of newline separated paragraphs." This partly agrees,
	// but "any number of newline separated paragraphs" is more contentious,
	// as there are no specifics for "paragraph" nor by how many newlines
	// they should be separated. If paragraph is only defined by "text
	// that is separated from other text by one or more newlines", this
	// makes the idea of paragraph pointless, and if paragraphs are more specifically
	// defined as "one continuous block of text more than one line long which
	// concerns one topic" as they might normally be, then this again -- in my
	// opinion -- limits the contents quite heavily.
	//
	// Regardless, I personally don't care how you format your body apart from
	// the length of the lines.
	// The main thing about a commit message should be that it gives a good
	// idea of what changed and why the change was made (with further details
	// available by looking at the diffs), and that this is done
	// in a clear, not terribly verbose way.
	var body []string
	for {

		// processing the first line requires the odd looping here
		body = append(body, line)
		if e := validator.ValidateLine(line, scner.TimesScanned); e != nil {
			errs = append(errs, e)
		}

		if !scner.Scan() {
			break
		} else {
			line = scner.Text()
		}
	}
	validator.Body = strings.Join(body, "\n")

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}

	return e
}
