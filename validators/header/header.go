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
)

const typeRegex = `(?<type>.+?)`
// scope is optional, surrounded by parentheses
const scope = `(?:\((?<scope>.+)\))?`
// header can be used to for checking for valid commit message
// headers.
var header = regexp.MustCompile(fmt.Sprintf(`\A%s%s(?<breaking>!)?: (?<description>.+)\z`, typeRegex, scope))

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
	// Breaking denotes whether the header marks the commit as a
	// breaking change
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
	header regexp.Regexp
	// accepted words for type
	types []string
	// Accepted words for scope.
	scopes []string
	// Accepted words for verb.
	verbs []string
	// minimum and maximum.
	headerLength [2]int

	// The content of the header in an easy accessible format.
	Header Header
}


// SetHeaderLength ensures that the header length's min and max are
// valid and sets them in the validator, returning a non-nil error
// if the values are not.
func (validator *HeaderValidator) SetHeaderLength(min, max int) (e error) {
	if min < 0 || max < 0 || min > max {
		return errors.New("min and max should be positive and min =< max")
	}
	validator.headerLength = [2]int{min, max}
	return nil
}

func NewDefaultHeaderValidator() (h *HeaderValidator) {
	h = &HeaderValidator{}
	h.header = *header
	return
}

// NewHeaderValidator returns a new HeaderValidator.
// The arguments types, scopes, and verbs define acceptable values for the
// type, scope, and verb. The headerLength argument sets minimum and maximum length.
// If the headerLength is invalid, a non-nil error is returned.
func NewHeaderValidator(types, scopes, verbs []string, headerLength [2]int) (h *HeaderValidator, e error) {
	h = &HeaderValidator{}

	h.header = *header
	h.types = types
	h.scopes = scopes
	h.verbs = verbs
	e = nil
	e = h.SetHeaderLength(headerLength[0], headerLength[1])

	return h, e
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
		fmt.Sprintf("	minimum and maximum length: [%d, %d]\n", validator.headerLength[0], validator.headerLength[1])
}

type InvalidError struct{}

func (e InvalidError) Error() string {
	return "Invalid header"
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

type InvalidLengthError struct {
	ExpectedMin int
	ExpectedMax int
	Received    int
}

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("Invalid header length %d, should be [%d, %d]", e.Received, e.ExpectedMin, e.ExpectedMax)
}

// Validate the length of the header.
func (validator *HeaderValidator) ValidateHeaderLength(possibleHeader string) (e error) {
	e = nil

	if validator.headerLength[0] == 0 && validator.headerLength[1] == 0 {
		return
	}

	if len(possibleHeader) < validator.headerLength[0] || len(possibleHeader) > validator.headerLength[1] {
		e = InvalidLengthError{ExpectedMin: validator.headerLength[0], ExpectedMax: validator.headerLength[1], Received: len(possibleHeader)}
	}
	return e
}

type InvalidScopeError struct {
	Expected []string
	Received string
}

func (e InvalidScopeError) Error() string {
	return fmt.Sprintf("Invalid scope '%s', should be one of %v", e.Received, e.Expected)
}

// Validate the scope (of Conventional commits syntax).
func (validator *HeaderValidator) ValidateScope(scope string) (e error) {
	if len(validator.scopes) == 0 || slices.Contains(validator.scopes, scope) {
		return nil
	}
	return InvalidScopeError{Expected: validator.scopes, Received: scope}
}

type InvalidTypeError struct {
	Expected []string
	Received string
}

func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("Invalid type '%s', should be one of %v", e.Received, e.Expected)
}

// Validate the type (of Conventional commits syntax).
func (validator *HeaderValidator) ValidateType(typ string) (e error) {
	if len(validator.types) == 0 || slices.Contains(validator.types, typ) {
		return nil
	}
	return InvalidTypeError{Expected: validator.types, Received: typ}
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

type InvalidDescriptionError struct {
	Verbs      []string
	// Received is the received description
	Received   string
}

func (e InvalidDescriptionError) Error() string {
	return fmt.Sprintf("Invalid description '%s', ", e.Received) +
		fmt.Sprintf("should be '<verb> <content>', where <verb> is one of %v ", e.Verbs)
}

// Validate the description of a commit message.
func (validator *HeaderValidator) ValidateDescription(description string) (desc Description, e error) {
	e = nil
	desc, e = validator.processDescription(description)
	if e != nil ||
		validator.ValidateVerb(desc.Verb) != nil  {
		e = InvalidDescriptionError{
				Received:   description,
				Verbs:      validator.verbs}
	}
	return desc, e
}

type InvalidVerbError struct {
	Expected []string
	Received string
}

func (e InvalidVerbError) Error() string {
	return fmt.Sprintf("Invalid verb '%s', should be one of %v", e.Received, e.Expected)
}

// Validate the verb of a description of a commit message.
func (validator *HeaderValidator) ValidateVerb(verb string) (e error) {
	if len(validator.verbs) == 0 || slices.Contains(validator.verbs, verb) {
		return nil
	}
	return InvalidVerbError{Expected: validator.verbs, Received: verb}
}

// ValidateString validates the header of a commit message, returning
// a non-nil error if validation
// failed at some point. This sets the Header of the validator.
func (validator *HeaderValidator) ValidateString(possibleHeader string) (e error) {

	// e = InvalidError{}
	var errs []error

	var matches = validator.header.FindStringSubmatch(possibleHeader)
	if matches == nil {
		return InvalidError{}
	}

	if err := validator.ValidateHeaderLength(possibleHeader); err != nil {
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
	if err := validator.ValidateScope(scope); scope != "" && err != nil {
		errs = append(errs, err)
	}
	validator.Header.Scope = scope

	var desc = matches[validator.header.SubexpIndex("description")]
	var description Description
	var err error
	description, err = validator.ValidateDescription(desc)
	if err != nil {
		errs = append(errs, err)
	}
	validator.Header.Description = description

	if len(errs) == 0 {
		e = nil
	} else {
		e = errors.Join(errs...)
	}
	return e
}
