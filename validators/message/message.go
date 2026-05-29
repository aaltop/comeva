package message

import (
	"bufio"
	"comeva/internal/utils"
	"comeva/validators"
	"comeva/validators/body"
	"comeva/validators/header"
	"comeva/validators/trailer"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
)

// MessageValidator validates a git commit message.
type MessageValidator struct {
	HeaderValidator  *header.HeaderValidator
	BodyValidator    *body.BodyValidator
	TrailerValidator *trailer.TrailerValidator

	FoundHeader, FoundBody, FoundTrailer bool
}

// Reset resets any content set during validation.
func (validator *MessageValidator) Reset() {
	validator.HeaderValidator.Reset()
	validator.BodyValidator.Reset()
	validator.TrailerValidator.Reset()
}

// NewDefaultMessageValidator creates the base MessageValidator.
func NewDefaultMessageValidator() (validator *MessageValidator) {
	validator, _ = NewMessageValidator(
		header.NewDefaultHeaderValidator(),
		body.NewDefaultBodyValidator(),
		trailer.NewDefaultTrailerValidator(),
	)
	return validator
}

func NewMessageValidator(
	headerValidator *header.HeaderValidator,
	bodyValidator *body.BodyValidator,
	trailerValidator *trailer.TrailerValidator,
) (messageValidator *MessageValidator, e error) {
	return &MessageValidator{
		HeaderValidator:  headerValidator,
		BodyValidator:    bodyValidator,
		TrailerValidator: trailerValidator}, nil
}

// NewMessageValidatorWithDefaults returns a [MessageValidator] with default values.
func NewMessageValidatorWithDefaults() (validator *MessageValidator) {
	var e error
	validator, e = NewMessageValidator(
		header.NewHeaderValidatorWithDefaults(),
		body.NewBodyValidatorWithDefaults(),
		trailer.NewTrailerValidatorWithDefaults(),
	)
	if e != nil {
		panic(e)
	}
	return
}

// Equal reports whether the two MessageValidators are equal.
func (validator *MessageValidator) Equal(other *MessageValidator) bool {
	if validator == nil || other == nil {
		return validator == other
	}

	return validator.HeaderValidator.Equal(other.HeaderValidator) &&
		validator.BodyValidator.Equal(other.BodyValidator) &&
		validator.TrailerValidator.Equal(other.TrailerValidator)
}

func (validator *MessageValidator) Validate(reader io.Reader) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

func (validator *MessageValidator) ValidateString(possibleMessage string) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(strings.NewReader(possibleMessage)))
}

// trailerEndRegex matches the end of a trailer block.
//
// "The group must either be at the end of the input or be the last
// non-whitespace lines before a line that starts with ---
// (followed by a space or the end of the line)."
var trailedEndRegex = regexp.MustCompile(`---(\r?\n)?`)

// ValidateBreakingChange checks whether an exclamation mark (denoting a breaking change)
// in the header is accompanied by a BREAKING-CHANGE trailer key and vice versa.
// To be run after a message has been processed.
func (validator *MessageValidator) ValidateBreakingChange() (e error) {
	// hasBreaking reports whether the trailers have a BREAKING-CHANGE key
	var hasBreaking bool = slices.ContainsFunc(validator.TrailerValidator.Trailers,
		func(tr trailer.Trailer) bool {
			return tr.Key == "BREAKING-CHANGE"
		},
	)

	if validator.HeaderValidator.Header.Breaking {
		if !hasBreaking {
			e = BreakingChangeError{}
		}
	} else {
		if hasBreaking {
			e = BreakingChangeError{}
		}
	}
	return
}

func (validator *MessageValidator) ValidateScanner(scanner *bufio.Scanner) (e error) {

	validator.Reset()
	var scner = utils.CountingScanner{Scanner: scanner}
	var errs []error

	if !scner.Scan() {
		return UnexpectedEOFError{
			ValidatorError: validators.ValidatorError{
				Line: scner.TimesScanned,
			},
		}
	}

	var possibleHeader string = scner.Text()
	if e = validator.HeaderValidator.ValidateString(possibleHeader); e != nil {
		errs = append(errs, e)
	}
	validator.FoundHeader = true

	// if only header found, fine: return
	if !scner.Scan() {

		if e = validator.ValidateBreakingChange(); e != nil {
			errs = append(errs, e)
		}

		// Need to check here separately as required keys would not
		// have been checked. Could technically just check whether there
		// are any required keys? Would be faster, and probably less error
		// prone. Does give the proper error message this way, though.
		if e = validator.TrailerValidator.EnsureRequiredKeys(); e != nil {
			errs = append(errs, e)
		}

		if len(errs) > 0 {
			e = errors.Join(errs...)
		}
		return
	}

	var line string = scner.Text()
	// line after header should be a newline (\n or \r\n)
	if line != "" {
		errs = append(errs, validators.InvalidLineError{
			Reason: fmt.Sprintf("Expected empty newline after header, got %#v", line),
			ValidatorError: validators.ValidatorError{
				Line: scner.TimesScanned,
			}})
	}

	var bodyStart, trailerStart = -1, -1
	var bodyContent, trailerContent []string

	// assigned non-nil if no line precedes the trailer block.
	var noEmptyBeforeTrailerError error = nil
	// find assumed body/start of trailer
	for scner.Scan() {
		line = scner.Text()

		if validator.TrailerValidator.ResemblesKeyValue(line) {
			trailerStart = int(scner.TimesScanned)
			trailerContent = append(trailerContent, line)
			validator.FoundTrailer = true
			if len(bodyContent) > 0 {
				if bodyContent[len(bodyContent)-1] != "" {
					noEmptyBeforeTrailerError = validators.InvalidLineError{
						Reason: fmt.Sprintf("Expected empty newline before trailer block, got %#v", line),
						ValidatorError: validators.ValidatorError{
							Line: scner.TimesScanned,
						},
					}
				}
			}
			break
		}

		if bodyStart == -1 {
			bodyStart = int(scner.TimesScanned)
			validator.FoundBody = true
		}

		bodyContent = append(bodyContent, line)

	}

	// find the assumed end of the trailer
	for scner.Scan() {
		line = scner.Text()

		if trailedEndRegex.MatchString(line) {
			break
		}
		trailerContent = append(trailerContent, line)
	}

	if e = validator.BodyValidator.ValidateStringWithLine(strings.Join(bodyContent, "\n"), uint(bodyStart)); e != nil {
		errs = append(errs, e)
	}

	if noEmptyBeforeTrailerError != nil {
		errs = append(errs, noEmptyBeforeTrailerError)
	}

	if e = validator.TrailerValidator.ValidateStringWithLine(strings.Join(trailerContent, "\n"), uint(trailerStart)); e != nil {
		errs = append(errs, e)
	}

	if e = validator.ValidateBreakingChange(); e != nil {
		errs = append(errs, e)
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}

	return
}
