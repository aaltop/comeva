package logging

import (
	"io"
	"log"
)

// NewDebugLogger creates a logger useful for debugging purposes.
func NewDebugLogger(out io.Writer) (logger *log.Logger) {
	return log.New(out, "", log.Llongfile)
}
