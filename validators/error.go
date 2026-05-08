package validators

import (
	"comeva/internal/utils"
	"fmt"
)

// messagePart represents the name of a part of a git commit message.
// Not to be instantiated directly, see [messageParts].
type messagePart string

// messageParts represents the names for the parts of a git commit message.
var MessageParts = struct {
	Header, Body, Trailer messagePart
}{
	Header: "Header", Body: "Body", Trailer: "Trailer",
}

// ValidatorError represents a general error.
type ValidatorError struct {
	MessagePart messagePart
	Line        uint
}

// ErrorString returns a descriptive error message given a base error
// message.
func (e ValidatorError) ErrorString(errorMessage string) (msg string) {
	if e.Line > 0 {
		msg += fmt.Sprintf("Line %d: ", e.Line)
	}
	if e.MessagePart != "" {
		msg += fmt.Sprintf("%v: ", e.MessagePart)
	}
	msg += errorMessage
	return
}

type InvalidLineLengthError struct {
	ValidatorError
	Expected utils.Bounds[uint]
	Received uint
}

func (e InvalidLineLengthError) Error() string {
	return e.ErrorString(fmt.Sprintf("Invalid line length %d, should be %v", e.Received, e.Expected))
}

type InvalidLineError struct {
	ValidatorError
	Reason string
	Line   uint
}

func (e InvalidLineError) Error() string {
	return e.ErrorString(e.Reason)
}
