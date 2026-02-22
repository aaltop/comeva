// io contains input-output -related utilities.

package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadFileString reads the contents of the given file, returning them as
// a string.
func ReadFileString(fileName string) (content string, e error) {
	var contentBytes []byte
	contentBytes, e = os.ReadFile(fileName)
	return string(contentBytes), e
}

// PrintStringLines prints each line in str, appending a line number to each line.
func PrintStringLines(str string) {
	var scanner = bufio.NewScanner(strings.NewReader(str))
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%4d: %s\n", i, scanner.Text())
	}
}

// PrintCommitMessage prints the passed message.
func PrintCommitMessage(message string) {
	var delimiterLine = "=================================================="
	fmt.Println(delimiterLine)
	PrintStringLines(message)
	fmt.Println(delimiterLine)
}
