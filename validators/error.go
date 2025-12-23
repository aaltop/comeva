package validators

import (
	"comeva/utils"
	"fmt"
)

type InvalidLineLengthError struct {
	Expected utils.Bounds[uint]
	Received uint
	Line     uint
}

func (e InvalidLineLengthError) Error() string {
	return fmt.Sprintf("Line %d: Invalid line length %d, should be %v", e.Line, e.Received, e.Expected)
}

type InvalidLineError struct {
	Reason string
	Line   uint
}

func (e InvalidLineError) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.Reason)
}
