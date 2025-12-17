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

// Header represents parts of a commit message header.
type Header struct {
	// Scope denotes where the changes were made.
	Scope string
	// Verb denotes the action of the commit (what happened).
	Verb string
	// Description contains the message (<verb> <content>) of the header.
	Description string
}

func (h Header) String() string {
	if len(h.Scope) == 0 {
		return h.Description
	}
	return fmt.Sprintf("%s: %s", h.Scope, h.Description)
}

type HeaderValidator struct {
	header *regexp.Regexp
	// Accepted words for scope.
	scopes []string
	// Accepted words for verbs.
	verbs []string
	// minimum and maximum.
	headerLength, contentLength [2]int

	// The content of the header in an easy accessible format.
	Header *Header
}

// Help returns a string that describes what the header is expected
// to look like.
func (validator HeaderValidator) Help() string {
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

func NewHeaderValidator(scopes, verbs []string, headerLength, contentLength [2]int) (h *HeaderValidator) {
	h = &HeaderValidator{}

	h.header = header
	h.scopes = scopes
	h.verbs = verbs
	if headerLength[0] == 0 && headerLength[1] == 0 {
		// no particular point in setting a minimum to anything positive
		// here, the other checks will require a certain minimum anyway
		headerLength[0], headerLength[1] = 0, 50
	}
	if contentLength[0] == 0 && contentLength[1] == 0 {
		contentLength[0], contentLength[1] = 3, 50
	}
	h.headerLength, h.contentLength = headerLength, contentLength

	return h
}

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

// ValidateString Validates the header of a commit message, returning
// the different parts of the header, and a non-nil error if validation
// failed at some point.
func (validator *HeaderValidator) ValidateString(possibleHeader string) (e error) {

	// e = InvalidError{}
	validator.Header = &Header{}
	var errs []error

	var matches = validator.header.FindStringSubmatch(possibleHeader)
	if matches == nil {
		return InvalidError{}
	}

	// TODO: move these to separate validation functions
	if len(possibleHeader) < validator.headerLength[0] || len(possibleHeader) > validator.headerLength[1] {
		errs = append(errs, InvalidLengthError{ExpectedMin: validator.headerLength[0], ExpectedMax: validator.headerLength[1], Received: len(possibleHeader)})
	}

	var scope = matches[validator.header.SubexpIndex("scope")]
	if !slices.Contains(validator.scopes, scope) {
		errs = append(errs, InvalidScopeError{Expected: validator.scopes, Received: scope})
	} else {
		validator.Header.Scope = scope
	}

	var desc = matches[validator.header.SubexpIndex("description")]
	var verb, content, found = strings.Cut(desc, " ")
	if !found ||
		!slices.Contains(validator.verbs, verb) ||
		(len(content) < validator.contentLength[0] || len(content) > validator.contentLength[1]) {
		errs = append(errs,
			InvalidDescriptionError{
				Received:   desc,
				contentMin: validator.contentLength[0],
				contentMax: validator.contentLength[1],
				Verbs:      validator.verbs})
	} else {
		validator.Header.Description = desc
		validator.Header.Verb = verb
	}

	if len(errs) == 0 {
		e = nil
	} else {
		e = errors.Join(errs...)
	}
	return e
}
