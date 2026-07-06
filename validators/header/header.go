// Package header contains functionalities for validating the header
// of a git commit message.
package header

import (
	baseErrors "errors"
	"fmt"
	"io"
	baseRegexp "regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/aaltop/comeva/internal/errors"
	"github.com/aaltop/comeva/internal/math"
	"github.com/aaltop/comeva/internal/regexp"
	"github.com/aaltop/comeva/validators"
)

var regexGroups = struct {
	Type, Scope, Breaking, ColonSpace, Description string
}{
	Type:        "type",
	Scope:       "scope",
	Breaking:    "breaking",
	ColonSpace:  "colon_space",
	Description: "description",
}

func createHeaderRegex() *baseRegexp.Regexp {

	var parts = struct {
		Type, Scope, Breaking, ColonSpace, Description string
	}{
		Type:        fmt.Sprintf(`(?<%s>[^(!:]+)`, regexGroups.Type),
		Scope:       fmt.Sprintf(`(?:\((?<%s>.+?)\))`, regexGroups.Scope),
		Breaking:    fmt.Sprintf(`(?<%s>!)`, regexGroups.Breaking),
		ColonSpace:  fmt.Sprintf(`(?<%s>: )`, regexGroups.ColonSpace),
		Description: fmt.Sprintf(`(?<%s>.+)`, regexGroups.Description),
	}

	var tmpl = errors.Panic2(template.New("").Parse(`\A{{.Type}}?{{.Scope}}?{{.Breaking}}?{{.ColonSpace}}?{{.Description}}?\z`))
	var builder = &strings.Builder{}
	errors.Panic(tmpl.Execute(builder, &parts))
	return baseRegexp.MustCompile(builder.String())

}

// headerRegex can be used to for checking for valid commit message
// headers.
// var headerRegex = regexp.MustCompile(fmt.Sprintf(`\A%s?%s?(?<breaking>!)?(?<colon_space>: )?(?<description>.+)?\z`, typeRegex, scope))
var headerRegex = createHeaderRegex()

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
	header *baseRegexp.Regexp
	// Accepted words for type.
	types []string
	// Accepted words for scope.
	scopes []string
	// Accepted words for verb.
	verbs []string
	// Minimum and maximum.
	lineLength math.Bounds[uint]

	Errors []validators.ValidatorErrorChild

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
	bounds, e := math.NewBounds(min, max, false, false)
	if e == nil {
		validator.lineLength = bounds
	}
	return e
}

// Reset resets any content set during validation.
func (validator *HeaderValidator) Reset() {
	validator.Header = Header{}
	validator.Errors = make([]validators.ValidatorErrorChild, 0)
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

// Validate validates the header of a commit message, returning
// a non-nil error if validation
// failed at some point. This sets the Header of the validator.
func (validator *HeaderValidator) Validate(reader io.Reader) (errs []validators.ValidatorErrorChild) {
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
		errs = append(errs, validators.ToValidatorErrorChild(readerError))
		return
	}

	return validator.ValidateString(stringBuilder.String())
}

// Validate the length of the header.
func (validator *HeaderValidator) ValidateLength(possibleHeader string) (errs []validators.ValidatorErrorChild) {

	var lower, upper = int(validator.lineLength.Lower), int(validator.lineLength.Upper)
	if lower == 0 && upper == 0 {
		return
	}

	// not exactly sure why len() returns an int in the first place?
	if !validator.lineLength.Contains(uint(len(possibleHeader))) {
		var validatorError = validators.ValidatorError{
			Line: 1, MessagePart: validators.MessageParts.Header}
		return append(errs, validators.InvalidLineLengthError{
			Expected:       validator.lineLength,
			Received:       uint(len(possibleHeader)),
			ValidatorError: validatorError,
		})
	}
	return
}

// Validate the scope (of Conventional commits syntax).
func (validator *HeaderValidator) validateScope(scope regexp.SubMatch) (errs []validators.ValidatorErrorChild) {
	if len(validator.scopes) == 0 || slices.Contains(validator.scopes, scope.Match) {
		return
	}
	return append(errs, InvalidScopeError{
		Expected: validator.scopes,
		Received: scope.Match,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
			Cols:        scope.Cols,
		},
	})
}

// Validate the type (of Conventional commits syntax).
func (validator *HeaderValidator) validateType(typ regexp.SubMatch) (errs []validators.ValidatorErrorChild) {
	if len(validator.types) == 0 || slices.Contains(validator.types, typ.Match) {
		return
	}
	return append(errs, InvalidTypeError{
		Expected: validator.types,
		Received: typ.Match,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
			Cols:        typ.Cols,
		},
	})
}

// Attempt extraction of constituent parts from the description of
// a commit message. No validation is performed.
func (validator *HeaderValidator) processDescription(description string) (desc Description, e error) {
	var verb, content, found = strings.Cut(description, " ")
	if !found {
		return desc, baseErrors.New("no space found in description")
	}

	desc.Verb = verb
	desc.Content = content
	return desc, nil
}

// Validate the description of a commit message.
func (validator *HeaderValidator) validateDescription(description regexp.SubMatch) (desc Description, errs []validators.ValidatorErrorChild) {
	var e error
	desc, e = validator.processDescription(description.Match)
	if e != nil ||
		validator.validateVerb(desc.Verb) != nil {
		errs = append(errs, InvalidDescriptionError{
			Received: description.Match,
			Verbs:    validator.verbs,
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Header,
				Line:        1,
				Cols:        description.Cols,
			},
		})
	}
	return
}

// Validate the verb of a description of a commit message.
func (validator *HeaderValidator) validateVerb(verb string) (errs []validators.ValidatorErrorChild) {
	if len(validator.verbs) == 0 || slices.Contains(validator.verbs, verb) {
		return
	}
	return append(errs, InvalidVerbError{
		Expected: validator.verbs,
		Received: verb,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Header,
			Line:        1,
		},
	})
}

// ValidateString validates the header of a commit message, returning
// a non-nil error if validation
// failed at some point. This sets the Header of the validator.
func (validator *HeaderValidator) ValidateString(possibleHeader string) (errs []validators.ValidatorErrorChild) {

	var ve validators.ValidatorErrorChild
	var tempErrs []validators.ValidatorErrorChild

	validator.Reset()

	defer func() {
		validator.Errors = errs
	}()

	var matches = validator.header.FindStringSubmatch(possibleHeader)
	if matches == nil {

		ve = InvalidError{
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Header,
				Line:        1,
			},
		}
		errs = append(errs, ve)
		return
	}

	var subMatches = regexp.GetSubMatches(validator.header, possibleHeader)

	errs = append(errs, validator.ValidateLength(possibleHeader)...)

	var typ = subMatches[regexGroups.Type]
	errs = append(errs, validator.validateType(typ)...)
	validator.Header.Type = typ.Match

	var breaking = subMatches[regexGroups.Breaking].Match
	validator.Header.Breaking = breaking == "!"

	var scope = subMatches[regexGroups.Scope]
	// scope is assumed to be at least one character, so empty scopes
	// mean that the content was matched correctly but that the scope group
	// did not exist, which is fine
	if scope.Match != "" {
		errs = append(errs, validator.validateScope(scope)...)
	}
	validator.Header.Scope = scope.Match

	var colonSpace = subMatches[regexGroups.ColonSpace].Match
	if colonSpace != ": " {
		errs = append(errs, ColonSpaceError{
			ValidatorError: validators.ValidatorError{
				MessagePart: "Header",
				Line:        1,
			},
		})
	}

	var desc = subMatches[regexGroups.Description]
	var description Description
	description, tempErrs = validator.validateDescription(desc)
	errs = append(errs, tempErrs...)
	validator.Header.Description = description

	return
}
