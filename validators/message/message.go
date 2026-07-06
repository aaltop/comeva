package message

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/aaltop/comeva/internal/utils"
	"github.com/aaltop/comeva/validators"
	"github.com/aaltop/comeva/validators/body"
	"github.com/aaltop/comeva/validators/header"
	"github.com/aaltop/comeva/validators/trailer"
)

// MessageValidator validates a git commit message.
type MessageValidator struct {
	HeaderValidator  *header.HeaderValidator
	BodyValidator    *body.BodyValidator
	TrailerValidator *trailer.TrailerValidator

	// Errors contains any errors particular to the MessageValidator itself;
	// sub-validator-specific errors are found in the relevant sub-validator.
	Errors []validators.ValidatorErrorChild

	FoundHeader, FoundBody, FoundTrailer bool
}

// Reset resets any content set during validation.
func (validator *MessageValidator) Reset() {
	validator.HeaderValidator.Reset()
	validator.BodyValidator.Reset()
	validator.TrailerValidator.Reset()

	validator.FoundHeader = false
	validator.FoundBody = false
	validator.FoundTrailer = false

	validator.Errors = make([]validators.ValidatorErrorChild, 0)
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

// AllErrors returns all the validation errors encountered by the validator
// and its subvalidators.
func (validator *MessageValidator) AllErrors() (errs []validators.ValidatorErrorChild) {
	errs = append(errs, validator.Errors...)
	errs = append(errs, validator.HeaderValidator.Errors...)
	errs = append(errs, validator.BodyValidator.Errors...)
	errs = append(errs, validator.TrailerValidator.Errors...)
	return
}

func (validator *MessageValidator) Validate(reader io.Reader) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

func (validator *MessageValidator) ValidateString(possibleMessage string) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateScanner(bufio.NewScanner(strings.NewReader(possibleMessage)))
}

// trailerEndRegex matches the end of a trailer block.
//
// "The group must either be at the end of the input or be the last
// non-whitespace lines before a line that starts with ---
// (followed by a space or the end of the line)."
var trailedEndRegex = regexp.MustCompile(`---(\r?\n)?`)

// validateBreakingChange checks whether an exclamation mark (denoting a breaking change)
// in the header is accompanied by a BREAKING-CHANGE trailer key and vice versa.
// To be run after a message has been processed.
func (validator *MessageValidator) validateBreakingChange() (errs []validators.ValidatorErrorChild) {
	// hasBreaking reports whether the trailers have a BREAKING-CHANGE key
	var hasBreaking bool = slices.ContainsFunc(validator.TrailerValidator.Trailers,
		func(tr trailer.Trailer) bool {
			return tr.Key == "BREAKING-CHANGE"
		},
	)

	if validator.HeaderValidator.Header.Breaking {
		if !hasBreaking {
			errs = append(errs, BreakingChangeError{})
		}
	} else {
		if hasBreaking {
			errs = append(errs, BreakingChangeError{})
		}
	}
	return
}

func (validator *MessageValidator) ValidateScanner(scanner *bufio.Scanner) (errs []validators.ValidatorErrorChild) {

	validator.Reset()
	var scner = utils.CountingScanner{Scanner: scanner}
	var ve validators.ValidatorErrorChild
	var tempErrs []validators.ValidatorErrorChild

	if !scner.Scan() {
		ve = UnexpectedEOFError{
			ValidatorError: validators.ValidatorError{
				Line: scner.TimesScanned,
			},
		}
		// MessageValidator-specific errors should be put here as well
		validator.Errors = append(validator.Errors, ve)
		errs = append(errs, ve)
		return
	}

	var possibleHeader string = scner.Text()
	errs = append(errs, validator.HeaderValidator.ValidateString(possibleHeader)...)
	validator.FoundHeader = true

	// if only header found, fine: return
	if !scner.Scan() {

		tempErrs = validator.validateBreakingChange()
		validator.Errors = append(validator.Errors, tempErrs...)
		errs = append(errs, tempErrs...)

		// Needs to be passed to Errors here because just running
		// EnsureRequiredKeys doesn't actually put the errors in the field,
		// but ValidatedContent needs the errors to be there.
		// Maybe have each validator method replace the errors? Not perfect,
		// but generally the actual full validation should be run anyway, and
		// not individual validation methods.
		validator.TrailerValidator.Errors = validator.TrailerValidator.EnsureRequiredKeys()
		errs = append(errs, validator.TrailerValidator.Errors...)

		return
	}

	var line string = scner.Text()
	// line after header should be a newline (\n or \r\n)
	if line != "" {
		ve = validators.InvalidLineError{
			Reason: fmt.Sprintf("Expected empty newline after header, got %#v", line),
			ValidatorError: validators.ValidatorError{
				Line: scner.TimesScanned,
			}}
		validator.Errors = append(validator.Errors, ve)
		errs = append(errs, ve)
	}

	var bodyStart, trailerStart = -1, -1
	var bodyContent, trailerContent []string

	// assigned non-nil if no line precedes the trailer block.
	var noEmptyBeforeTrailerError validators.ValidatorErrorChild
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

	errs = append(errs, validator.BodyValidator.ValidateStringWithLine(strings.Join(bodyContent, "\n"), uint(bodyStart))...)

	if noEmptyBeforeTrailerError != nil {
		validator.Errors = append(validator.Errors, noEmptyBeforeTrailerError)
		errs = append(errs, noEmptyBeforeTrailerError)
	}

	errs = append(errs, validator.TrailerValidator.ValidateStringWithLine(strings.Join(trailerContent, "\n"), uint(trailerStart))...)

	tempErrs = validator.validateBreakingChange()
	validator.Errors = append(validator.Errors, tempErrs...)
	errs = append(errs, tempErrs...)

	return
}
