package header

import (
	"comeva/validators"
	"fmt"
)

// InvalidError is returned when the header cannot be validated at all.
type InvalidError struct {
	validators.ValidatorError
}

func (e InvalidError) Error() string {
	return e.ErrorString("Invalid header, did not match regular expression")
}

type InvalidScopeError struct {
	validators.ValidatorError
	Expected []string
	Received string
}

func (e InvalidScopeError) Error() string {
	return e.ErrorString(fmt.Sprintf("Invalid scope '%s', should be one of %v", e.Received, e.Expected))
}

type InvalidTypeError struct {
	validators.ValidatorError
	Expected []string
	Received string
}

func (e InvalidTypeError) Error() string {
	return e.ErrorString(fmt.Sprintf("Invalid type '%s', should be one of %v", e.Received, e.Expected))
}

type InvalidDescriptionError struct {
	validators.ValidatorError
	Verbs []string
	// Received is the received description
	Received string
}

func (e InvalidDescriptionError) Error() string {
	var message string = fmt.Sprintf("Invalid description '%s', ", e.Received) +
		fmt.Sprintf("should be '<verb> <content>', where <verb> is one of %v ", e.Verbs)

	return e.ErrorString(message)
}

type InvalidVerbError struct {
	validators.ValidatorError
	Expected []string
	Received string
}

func (e InvalidVerbError) Error() string {
	return e.ErrorString(fmt.Sprintf("Invalid verb '%s', should be one of %v", e.Received, e.Expected))
}
