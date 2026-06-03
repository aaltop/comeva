package validators

import (
	"comeva/internal/utils"
	"fmt"
)

// messagePart represents the name of a part of a git commit message.
// Not to be instantiated directly, see [MessageParts].
type messagePart string

// messageParts represents the names for the parts of a git commit message.
var MessageParts = struct {
	Header, Body, Trailer messagePart
}{
	Header: "Header", Body: "Body", Trailer: "Trailer",
}

type IValidatorError interface {
	GetMessagePart() messagePart
	GetLine() uint
}

// ValidatorErrorChild represents children of [ValidatorError].
type ValidatorErrorChild interface {
	error
	IValidatorError
	GetReason() string
}

// ToValidatorErrorChild returns an error that satisfies the interface.
func ToValidatorErrorChild(e error) ValidatorErrorChild {
	return validatorErrorChild{
		Reason: e.Error(),
	}
}

// ValidatorError represents a general validator error, though does not
// implement error itself.
type ValidatorError struct {
	MessagePart messagePart
	Line        uint
}

func (vali ValidatorError) GetMessagePart() messagePart {
	return vali.MessagePart
}

func (vali ValidatorError) GetLine() uint {
	return vali.Line
}

// ErrorString constructs an error message that includes the line and
// part of the message (should they be present).
func ErrorString(e IValidatorError, errorMessage string) (msg string) {
	if e.GetLine() > 0 {
		msg += fmt.Sprintf("Line %d: ", e.GetLine())
	}
	if e.GetMessagePart() != "" {
		msg += fmt.Sprintf("%v: ", e.GetMessagePart())
	}
	msg += errorMessage
	return
}

// ErrorString returns a descriptive error message given a base error
// message.
func (e ValidatorError) ErrorString(errorMessage string) (msg string) {
	return ErrorString(e, errorMessage)
}

type InvalidLineLengthError struct {
	ValidatorError
	Expected utils.Bounds[uint]
	Received uint
}

func (e InvalidLineLengthError) Error() string {
	return e.ErrorString(e.GetReason())
}

func (e InvalidLineLengthError) GetReason() string {
	return fmt.Sprintf("Invalid line length %d, should be %v", e.Received, e.Expected)
}

type InvalidLineError struct {
	ValidatorError
	Reason string
}

func (e InvalidLineError) GetReason() string {
	return e.Reason
}

func (e InvalidLineError) Error() string {
	return e.ErrorString(e.GetReason())
}
