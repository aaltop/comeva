package validators

import (
	"errors"
	"fmt"
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
}

// Help returns a string that describes what the header is expected
// to look like.
func (h HeaderValidator) Help() string {
	return "Header should be of the format [<location>: ]<description>" +
		fmt.Sprintf("	<location>: (%s)\n", strings.Join(h.scopes, "|")) +
		"	<description>: <verb> <content>\n" +
		fmt.Sprintf("	<verb>: (%s)\n", strings.Join(h.verbs, "|")) +
		fmt.Sprintf("	<content>: Any content, length [%d, %d]", h.contentLength[0], h.contentLength[1])
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

// func (h *HeaderValidator) Validate(reader io.Reader) error {
// 	// var matches = header.FindStringSubmatch(possibleHeader)
// 	// if matches == nil {
// 	// 	h.Problems = append(h.Problems, "The header was invalid.")
// 	// 	return h, InvalidError{}
// 	// }
// 	// var loc = matches[header.SubexpIndex("location")]
// 	// if !location.MatchString(loc) {
// 	// 	h.Problems = append(h.Problems, fmt.Sprintf("Invalid location: %s", loc))
// 	// }
// 	// var msg = matches[header.SubexpIndex("message")]
// }

// ValidateString Validates the header of a commit message, returning
// the different parts of the header, and a non-nil error if validation
// failed at some point.
func (h *HeaderValidator) ValidateString(possibleHeader string) (header *Header, e error) {

	header = &Header{}
	e = InvalidError{}
	var hasError = false

	var matches = h.header.FindStringSubmatch(possibleHeader)
	if matches == nil {
		hasError = true
		return header, InvalidError{}
	}

	// TODO: move these to separate validation functions
	if len(possibleHeader) < h.headerLength[0] || len(possibleHeader) > h.headerLength[1] {
		hasError = true
		e = errors.Join(e, InvalidLengthError{ExpectedMin: h.headerLength[0], ExpectedMax: h.headerLength[1], Received: len(possibleHeader)})
	}

	var scope = matches[h.header.SubexpIndex("scope")]
	if !slices.Contains(h.scopes, scope) {
		hasError = true
		e = errors.Join(e, InvalidScopeError{Expected: h.scopes, Received: scope})
	} else {
		header.Scope = scope
	}

	var desc = matches[h.header.SubexpIndex("description")]
	var verb, content, found = strings.Cut(desc, " ")
	if !found ||
		!slices.Contains(h.verbs, verb) ||
		(len(content) < h.contentLength[0] || len(content) > h.contentLength[1]) {
		hasError = true
		e = errors.Join(e, InvalidDescriptionError{Received: desc, contentMin: h.contentLength[0], contentMax: h.contentLength[1], Verbs: h.verbs})
	} else {
		header.Description = desc
		header.Verb = verb
	}

	if !hasError {
		e = nil
	}
	return header, e
}
