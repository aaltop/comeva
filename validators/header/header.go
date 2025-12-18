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

// Header can be used to for checking for valid commit message
// headers.
// var header = regexp.MustCompile(`\A(?:(?<location>frontend|backend): )?(?<message>(?<verb>Add|Remove|Fix) .{3,50})\z`)
var header = regexp.MustCompile(`\A(?:(?<scope>.*): )?(?<description>.*)\z`)

// Description represents the description part of a commit message
type Description struct {
	// Verb describes the action taken by a commit.
	Verb string
	// Content is the rest of the description.
	Content string
}

func (desc Description) String() string {
	return fmt.Sprintf("%s %s", desc.Verb, desc.Content)
}

// Header represents parts of a commit message header.
type Header struct {
	// Scope denotes where the changes were made.
	Scope string
	// Verb denotes the action of the commit (what happened).
	Verb string
	// Description contains the message (<verb> <content>) of the header.
	Description Description
}

func (h Header) String() string {
	if len(h.Scope) == 0 {
		return h.Description.String()
	}
	return fmt.Sprintf("%s: %s", h.Scope, h.Description.String())
}

type HeaderValidator struct {
	// Regexp used to pick out parts of the header.
	header *regexp.Regexp
	// Accepted words for scope.
	scopes []string
	// Accepted words for verbs.
	verbs []string
	// minimum and maximum.
	headerLength, contentLength [2]int

	// The content of the header in an easy accessible format.
	Header Header
}

// Help returns a string that describes what the header is expected
// to look like.
func (validator *HeaderValidator) Help() string {
	return "Header should be of the format [<location>: ]<description>" +
		fmt.Sprintf("	<location>: (%s)\n", strings.Join(validator.scopes, "|")) +
		"	<description>: <verb> <content>\n" +
		fmt.Sprintf("	<verb>: (%s)\n", strings.Join(validator.verbs, "|")) +
		fmt.Sprintf("	<content>: Any content, length [%d, %d]", validator.contentLength[0], validator.contentLength[1])
}

type InvalidError struct{}

func (e InvalidError) Error() string {
	return "Invalid header"
}

type InvalidVerbError struct {
	Expected []string
	Received string
}

func (e InvalidVerbError) Error() string {
	return ""
}

type InvalidLengthError struct {
	ExpectedMin int
	ExpectedMax int
	Received    int
}

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("Invalid header length %d, should be [%d, %d]", e.Received, e.ExpectedMin, e.ExpectedMax)
}

type InvalidScopeError struct {
	Expected []string
	Received string
}

func (e InvalidScopeError) Error() string {
	return fmt.Sprintf("Invalid scope '%s', should be one of %v", e.Received, e.Expected)
}

type InvalidDescriptionError struct {
	Verbs      []string
	contentMin int
	contentMax int
	Received   string
}

func (e InvalidDescriptionError) Error() string {
	return fmt.Sprintf("Invalid description '%s', ", e.Received) +
		fmt.Sprintf("should be '<verb> <content>', where <verb> is one of %v ", e.Verbs) +
		fmt.Sprintf("and <content> of length [%d, %d]", e.contentMin, e.contentMax)
}

func NewDefaultHeaderValidator() (h *HeaderValidator) {
	// Have default config file instead
	return NewHeaderValidator([]string{""}, []string{"Add", "Remove", "Fix"}, [2]int{}, [2]int{})
}

// NewHeaderValidator returns a new HeaderValidator, and should be called
// to make one. Scopes and verbs define acceptable values for the scope
// and verb. The *length arguments set minimum and maximum lengths,
// defaulting to [0, 80] when zero values. 
func NewHeaderValidator(scopes, verbs []string, headerLength, contentLength [2]int) (h *HeaderValidator) {
	h = &HeaderValidator{}

	h.header = header
	h.scopes = scopes
	h.verbs = verbs
	if headerLength[0] == 0 && headerLength[1] == 0 {
		// no particular point in setting a minimum to anything positive
		// here, the other checks will require a certain minimum anyway
		headerLength[0], headerLength[1] = 0, 80
	}
	if contentLength[0] == 0 && contentLength[1] == 0 {
		contentLength[0], contentLength[1] = 3, 80
	}

	// TODO: validate these first for positive
	h.headerLength, h.contentLength = headerLength, contentLength

	return h
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
func (validator *HeaderValidator) ValidateHeaderLength(possibleHeader string) (e error) {
	e = nil
	if len(possibleHeader) < validator.headerLength[0] || len(possibleHeader) > validator.headerLength[1] {
		e = InvalidLengthError{ExpectedMin: validator.headerLength[0], ExpectedMax: validator.headerLength[1], Received: len(possibleHeader)}
	}
	return e
}

// Validate the scope (of Conventional commits syntax).
func (validator *HeaderValidator) ValidateScope(scope string) (e error) {
	e = nil
	if !slices.Contains(validator.scopes, scope) {
		e = InvalidScopeError{Expected: validator.scopes, Received: scope}
	}
	return e
}


// Attempt extraction of constituent parts from the description of
// a commit message. No validation is performed.
func (validator *HeaderValidator) ProcessDescription(description string) (desc Description, e error) {
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
	desc, e = validator.ProcessDescription(description)
	if e != nil ||
		!slices.Contains(validator.verbs, desc.Verb) ||
		(len(desc.Content) < validator.contentLength[0] || len(desc.Content) > validator.contentLength[1]) {
		e = InvalidDescriptionError{
				Received:   description,
				contentMin: validator.contentLength[0],
				contentMax: validator.contentLength[1],
				Verbs:      validator.verbs}
	}
	return desc, e
}

// ValidateString Validates the header of a commit message, returning
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

	var scope = matches[validator.header.SubexpIndex("scope")]
	if err := validator.ValidateScope(scope); err != nil {
		errs = append(errs, err)
	} else {
		validator.Header.Scope = scope
	}

	var desc = matches[validator.header.SubexpIndex("description")]
	if description, err := validator.ValidateDescription(desc); err != nil {
		errs = append(errs, err)
	} else {
		validator.Header.Description = description
	}

	if len(errs) == 0 {
		e = nil
	} else {
		e = errors.Join(errs...)
	}
	return e
}
