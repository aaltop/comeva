package validators

import (
	"io"
)

type Validator[T any] interface {
	// Validate validates the passed value, returning an error which
	// is nil if no validation issues were found.
	//
	// Validation, when possible, should be run to completion regardless
	// of whether any errors with validation are encountered; instead,
	// all validation errors should be joined into one.
	Validate(value T) error
}

// ReaderValidator validates the contents returned by an io.Reader.
type ReaderValidator interface {
	Validator[io.Reader]
}

type StringValidator interface {
	ReaderValidator

	// Validate a string, returning an error which is nil if no validation
	// issues were found.
	//
	// Validation, when possible, should be run to completion regardless
	// of whether any errors with validation are encountered; instead,
	// all validation errors should be joined into one.
	ValidateString(str string) error
}
