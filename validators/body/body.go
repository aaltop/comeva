// package body contains functionalities for validating the body
// of a git commit message.
package body

import (
	"bufio"
	"comeva/utils"
	"comeva/validators"
	"errors"
	"fmt"
	"io"
	"strings"
)

// BodyValidator validates a commit message's body.
type BodyValidator struct {
	// lineLength describes the lower and upper bound of a line's length.
	lineLength utils.Bounds[uint]
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

// ValidateLine returns a non-nil error if the passed line
// does not fulfill the requirements of a commit message's body's line.
func (validator *BodyValidator) ValidateLine(line string, lineNum uint) (e error) {
	e = nil
	var lenLine uint = uint(len(line))
	if !validator.lineLength.Contains(lenLine) {
		e = validators.InvalidLineLengthError{Line: lineNum, Expected: validator.lineLength, Received: lenLine}
	}
	return
}

func (validator *BodyValidator) Validate(reader io.Reader) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

func (validator *BodyValidator) ValidateString(possibleBody string) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(strings.NewReader(possibleBody)))
}

type InvalidLineError struct {
	Reason string
	Line   uint
}

func (e InvalidLineError) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.Reason)
}

func (validator *BodyValidator) ValidateScanner(scanner *bufio.Scanner) (e error) {
	var line string
	e = nil
	var errs []error
	var scner = utils.CountingScanner{Scanner: scanner}

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
		errs = append(errs, InvalidLineError{
			Reason: "First line should not be empty",
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
	for scner.Scan() {
		if e := validator.ValidateLine(scner.Text(), scner.TimesScanned); e != nil {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}

	return e
}
