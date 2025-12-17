package validators

import "io"

type Validator interface {
	// Validate validates the reader's contents, returning an error which
	// is nil if no validation issues were found.
	// 
	// Validation, when possible, should be run to completion regardless
	// of whether any errors with validation are encountered; instead,
	// all validation errors should be joined into one.
	Validate(reader io.Reader) error
}
