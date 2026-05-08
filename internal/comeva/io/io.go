// io contains input-output -related utilities.

package io

import (
	"comeva/internal/io"
	"fmt"
)

// PrintCommitMessage prints the passed message.
func PrintCommitMessage(message string) {
	var delimiterLine = "=================================================="
	fmt.Println(delimiterLine)
	io.PrintStringLines(message)
	fmt.Println(delimiterLine)
}
