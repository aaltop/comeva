package validators

import (
	"encoding/json"
	"io"
)

type Validator[T any] interface {
	// Validate validates the passed value, returning an error which
	// is nil if no validation issues were found.
	//
	// Validation, when possible, should be run to completion regardless
	// of whether any errors with validation are encountered; instead,
	// all validation errors should be joined into one.
	Validate(value T) []ValidatorErrorChild
}

// ReaderValidator validates the contents returned by an io.Reader.
type ReaderValidator interface {
	Validator[io.Reader]
}

type StringValidator interface {

	// Validate a string, returning an error which is nil if no validation
	// issues were found.
	//
	// Validation, when possible, should be run to completion regardless
	// of whether any errors with validation are encountered; instead,
	// all validation errors should be joined into one.
	ValidateString(str string) []ValidatorErrorChild
}

// ValidatedContent is the standard format returned by a validator's
// ValidatedContent() method.
type ValidatedContent[T any] struct {
	// Content is the content the validator parsed from the message.
	Content T
	// Errors are the errors encountered specifically by the validator
	// (i.e. not by its children or parents)
	Errors []ValidatorErrorChild
}

type validatedContent[T any] struct {
	Content T
	Errors  []validatorErrorChild
}

func (content ValidatedContent[T]) MarshalJSON() (data []byte, e error) {
	var errs = make([]validatorErrorChild, len(content.Errors))

	var temp validatorErrorChild
	// normalise the errors
	for i, err := range content.Errors {
		temp.Reason = err.GetReason()
		temp.Line = err.GetLine()
		temp.MessagePart = err.GetMessagePart()
		errs[i] = temp
	}

	var jsonContent validatedContent[T]
	jsonContent.Content = content.Content
	jsonContent.Errors = errs

	return json.Marshal(&jsonContent)
}

func (content *ValidatedContent[T]) UnmarshalJSON(data []byte) (e error) {
	var jsonContent validatedContent[T]

	e = json.Unmarshal(data, &jsonContent)

	content.Content = jsonContent.Content
	content.Errors = make([]ValidatorErrorChild, len(jsonContent.Errors))

	for i, err := range jsonContent.Errors {
		content.Errors[i] = err
	}
	return
}

type validatorErrorChild struct {
	Reason      string
	MessagePart messagePart
	Line        uint
}

func (child validatorErrorChild) Error() string {
	return ErrorString(child, child.Reason)
}

func (child validatorErrorChild) GetMessagePart() messagePart {
	return child.MessagePart
}

func (child validatorErrorChild) GetLine() uint {
	return child.Line
}

func (child validatorErrorChild) GetReason() string {
	return child.Reason
}
