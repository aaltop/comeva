package message

import (
	"comeva/validators"
)

type UnexpectedEOFError struct {
	validators.ValidatorError
}

func (e UnexpectedEOFError) Error() string {
	return e.ErrorString("Unexpected EOF")
}

// BreakingChangeError is returned when breaking changes are not properly
// marked in a commit message.
type BreakingChangeError struct {
	validators.ValidatorError
}

func (e BreakingChangeError) Error() string {
	var message string = "An exclamation mark denoting a breaking change should be accompanied by a BREAKING-CHANGE key in the trailer block and vice versa"
	return e.ErrorString(message)
}
