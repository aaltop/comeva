package trailer

import (
	"comeva/validators"
	"fmt"
)

type InvalidKeyError struct {
	validators.ValidatorError
	Expected []string
	Received string
}

func (e InvalidKeyError) Error() string {
	var message string = fmt.Sprintf("Invalid trailer key %s, should be one of %v", e.Received, e.Expected)
	return e.ErrorString(message)
}

type InvalidKeyValueError struct {
	validators.ValidatorError
}

func (e InvalidKeyValueError) Error() string {
	return e.ErrorString("No key-value pair found")
}

type InvalidValueContinuationError struct {
	validators.ValidatorError
	Indent uint
}

func (e InvalidValueContinuationError) Error() string {
	var message string = fmt.Sprintf("Expecting value continuation with indent %d", e.Indent)
	return e.ErrorString(message)
}

type MissingRequiredKeyError struct {
	validators.ValidatorError
	Missing []string
}

func (e MissingRequiredKeyError) Error() string {
	return e.ErrorString(fmt.Sprintf("Required keys %v not found in trailer", e.Missing))
}
