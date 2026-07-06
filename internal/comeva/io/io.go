// io contains input-output -related utilities.

package io

import (
	"fmt"

	"github.com/aaltop/comeva/internal/io"
	"github.com/aaltop/comeva/internal/io/ansi"
)

// PrintCommitMessage prints the passed message.
func PrintCommitMessage(message string) {

	var color = ansi.NewColorScheme(ansi.BasicColors.SkyBlue, ansi.BasicColors.Black)

	var delimiterLine = "=================================================="
	fmt.Println(color.ApplyFore(delimiterLine))
	io.PrintStringLines(message)
	fmt.Println(color.ApplyFore(delimiterLine))
}
