// Package header contains functionalities for validating the header
// of a git commit message.
package header

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"comeva/internal/utils"
	"comeva/validators"
)

const typeRegex = `(?<type>.+?)`

// scope is optional, surrounded by parentheses
const scope = `(?:\((?<scope>.+)\))?`

// headerRegex can be used to for checking for valid commit message
// headers.
var headerRegex = regexp.MustCompile(fmt.Sprintf(`\A%s%s(?<breaking>!)?: (?<description>.+)\z`, typeRegex, scope))

// Description represents the description part of a commit message
type Description struct {
	// Verb describes the action taken by a commit. May be an empty
	// string if a verb was not parsed.
	Verb string
	// Content is the rest of the description.
	Content string
}

func (desc Description) String() string {
	if desc.Verb == "" {
		return desc.Content
	} else {
		return fmt.Sprintf("%s %s", desc.Verb, desc.Content)
	}
}

// Header represents parts of a commit message header.
type Header struct {
	// Type denotes the type of the commit
	Type string
	// Scope denotes where the changes were made.
	Scope string
	// Breaking reports whether the header marks the commit as a breaking change.
	Breaking bool
	// Description contains the message (<verb> <content>) of the header.
	Description Description
}

func (h Header) String() string {
	var breaking = ""
	if h.Breaking {
		breaking = "!"
	}
	var scope = ""
	if len(h.Scope) != 0 {
		scope = fmt.Sprintf("(%s)", h.Scope)
	}
	return fmt.Sprintf("%s%s%s: %s", h.Type, scope, breaking, h.Description.String())
}

// HeaderValidator validates a commit message's header.
type HeaderValidator struct {
	// Regexp used to pick out parts of the header.
	header *regexp.Regexp
	// Accepted words for type.
	types []string
	// Accepted words for scope.
	scopes []string
	// Accepted words for verb.
	verbs []string
	// Minimum and maximum.
	lineLength utils.Bounds[uint]

	Errors []error

	// The content of the header in an easy accessible format.
	Header Header
}

// Equal reports whether the two HeaderValidators are equal.
func (validator *HeaderValidator) Equal(other *HeaderValidator) bool {
	if validator == nil || other == nil {
		return validator == other
	}

	// NOTE: need to ensure that the two (for all the slices here) are sorted. The validator itself
	// does not require sort order for equality in behaviour, but this comparison
	// does.
	return slices.Compare(validator.types, other.types) == 0 &&
		slices.Compare(validator.scopes, other.scopes) == 0 &&
		slices.Compare(validator.verbs, other.verbs) == 0 &&
		validator.lineLength == other.lineLength &&
		validator.header.String() == other.header.String()
}

// SetLineLength ensures that the line length's min and max are
// valid and sets them in the validator, returning a non-nil error
// if the values are not.
func (validator *HeaderValidator) SetLineLength(min, max uint) (e error) {
	bounds, e := utils.NewBounds(min, max, false, false)
	if e == nil {
		validator.lineLength = bounds
	}
	return e
}

// Reset resets any content set during validation.
func (validator *HeaderValidator) Reset() {
	validator.Header = Header{}
	validator.Errors = make([]error, 0)
}

// NewDefaultHeaderValidator creates the base [HeaderValidator].
func NewDefaultHeaderValidator() (h *HeaderValidator) {
	h = &HeaderValidator{}
	h.header = headerRegex
	return
}

// NewHeaderValidator returns a new [HeaderValidator].
// The arguments types, scopes, and verbs define acceptable values for the
// type, scope, and verb. The lineLength argument sets minimum and maximum length.
// If the lineLength is invalid, a non-nil error is returned.
func NewHeaderValidator(types, scopes, verbs []string, lineLength [2]uint) (h *HeaderValidator, e error) {
	h = &HeaderValidator{}

	h.header = headerRegex
	slices.Sort(types)
	h.types = types
	slices.Sort(scopes)
	h.scopes = scopes
	slices.Sort(verbs)
	h.verbs = verbs
	e = nil
	e = h.SetLineLength(lineLength[0], lineLength[1])

	return h, e
}

// NewHeaderValidatorWithDefaults returns a [HeaderValidator] with default values.
func NewHeaderValidatorWithDefaults() (validator *HeaderValidator) {
	var e error
	validator, e = NewHeaderValidator(
		[]string{"feat", "fix"},
		[]string{},
		[]string{"Add", "Change", "Remove", "Update", "Fix", "Make"},
		[2]uint{0, 80},
	)
	if e != nil {
		panic(e)
	}
	return
}

// Help returns a string that describes what the header is expected
// to look like.
func (validator *HeaderValidator) Help() string {

	var verb = strings.Join(validator.verbs, "|")
	var notConstrained = "any (not constrained)"
	if len(validator.verbs) == 0 {
		verb = notConstrained
	}

	var scope = strings.Join(validator.scopes, "|")
	if len(validator.scopes) == 0 {
		scope = notConstrained
	}

	var typ = strings.Join(validator.types, "|")
	if len(validator.types) == 0 {
		typ = notConstrained
	}

	return "Header should be of the format <type>[(<scope>)][!]: <description>\n" +
		fmt.Sprintf("	<type>: (%s)\n", typ) +
		fmt.Sprintf("	<scope>: (%s)\n", scope) +
		"	<description>: <verb> <content>\n" +
		fmt.Sprintf("	<verb>: (%s)\n", verb) +
		fmt.Sprintln("	<content>: Any content") +
		fmt.Sprintf("	minimum and maximum length: %v\n", validator.lineLength)
}

func (validator *HeaderValidator) ValidatedContent() *Header {
	return &validator.Header
}

// Validate validates the header of a commit message, returning
// a non-nil error if validation
// failed at some point. This sets the Header of the validator.
func (validator *HeaderValidator) Validate(reader io.Reader) (e error) {
	var stringBuilder = strings.Builder{}
	var buffer = make([]byte, 256)
	var readerError error
	var bytesRead int
	for {
		bytesRead, readerError = reader.Read(buffer)

		if bytesRead > 0 {
			stringBuilder.Grow(bytesRead)
			stringBuilder.Write(buffer[:bytesRead])
		}

		// bytes should be handled BEFORE handling the error, according
		// to the documentation of io.Reader
		if readerError != nil {
			break
		}
	}

	if !(readerError == io.EOF) {
		return readerError
	}

	return validator.ValidateString(stringBuilder.String())
}

// Validate the length of the header.
func (validator *HeaderValidator) ValidateLength(possibleHeader string) (e error) {
	e = nil

	var lower, upper = int(validator.lineLength.Lower), int(validator.lineLength.Upper)
	if lower == 0 && upper == 0 {
		return
	}

	// not exactly sure why len() returns an int in the first place?
	if !validator.lineLength.Contains(uint(len(possibleHeader))) {
		var validatorError = validators.ValidatorError{
			Line: 1, MessagePart: validators.MessageParts.Header}
		e = validators.InvalidLineLengthError{
			Expected:       validator.lineLength,
			Received:       uint(len(possibleHeader)),
			ValidatorError: validatorError}
	}
	return e
}

// Validate the scope (of Conventional commits syntax).
func (validator *HeaderValidator) ValidateScope(scope string) (e error) {
	if len(validator.scopes) == 0 || slices.Contains(validator.scopes, scope) {
		return nil
	}
	return InvalidScopeError{
		Expected: validator.scopes,
		Received: scope,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
		},
	}
}

// Validate the type (of Conventional commits syntax).
func (validator *HeaderValidator) ValidateType(typ string) (e error) {
	if len(validator.types) == 0 || slices.Contains(validator.types, typ) {
		return nil
	}
	return InvalidTypeError{
		Expected: validator.types,
		Received: typ,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
		},
	}
}

// Attempt extraction of constituent parts from the description of
// a commit message. No validation is performed.
func (validator *HeaderValidator) processDescription(description string) (desc Description, e error) {
	var verb, content, found = strings.Cut(description, " ")
	if !found {
		return desc, errors.New("no space found in description")
	}

	desc.Verb = verb
	desc.Content = content
	return desc, nil
}

// Validate the description of a commit message.
func (validator *HeaderValidator) ValidateDescription(description string) (desc Description, e error) {
	e = nil
	desc, e = validator.processDescription(description)
	if e != nil ||
		validator.ValidateVerb(desc.Verb) != nil {
		e = InvalidDescriptionError{
			Received: description,
			Verbs:    validator.verbs,
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Header,
				Line:        1,
			},
		}
	}
	return desc, e
}

// Validate the verb of a description of a commit message.
func (validator *HeaderValidator) ValidateVerb(verb string) (e error) {
	if len(validator.verbs) == 0 || slices.Contains(validator.verbs, verb) {
		return nil
	}
	return InvalidVerbError{
		Expected: validator.verbs,
		Received: verb,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
		},
	}
}

// ValidateString validates the header of a commit message, returning
// a non-nil error if validation
// failed at some point. This sets the Header of the validator.
func (validator *HeaderValidator) ValidateString(possibleHeader string) (e error) {

	validator.Reset()
	var errs []error

	defer func() {
		validator.Errors = errs
		e = errors.Join(errs...)
	}()

	var matches = validator.header.FindStringSubmatch(possibleHeader)
	if matches == nil {

		e = InvalidError{
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Header,
				Line:        1,
			},
		}
		errs = append(errs, e)
		return
	}

	if err := validator.ValidateLength(possibleHeader); err != nil {
		errs = append(errs, err)
	}

	var typ = matches[validator.header.SubexpIndex("type")]
	if err := validator.ValidateType(typ); err != nil {
		errs = append(errs, err)
	}
	validator.Header.Type = typ

	var breaking = matches[validator.header.SubexpIndex("breaking")]
	validator.Header.Breaking = breaking == "!"

	var scope = matches[validator.header.SubexpIndex("scope")]
	// scope is assumed to be at least one character, so empty scopes
	// mean that the content was matched correctly but that the scope group
	// did not exist, which is fine
	if e = validator.ValidateScope(scope); scope != "" && e != nil {
		errs = append(errs, e)
	}
	validator.Header.Scope = scope

	var desc = matches[validator.header.SubexpIndex("description")]
	var description Description
	description, e = validator.ValidateDescription(desc)
	if e != nil {
		errs = append(errs, e)
	}
	validator.Header.Description = description

	return
}
