package validators

import "io"

type Validator interface {
	Validate(reader io.Reader) error
}
