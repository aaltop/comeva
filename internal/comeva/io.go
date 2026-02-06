// io contains input-output -related utilities.

package comeva

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readFileString reads the contents of the given file, returning them as
// a string.
func readFileString(fileName string) (content string, e error) {
	var contentBytes []byte
	contentBytes, e = os.ReadFile(fileName)
	return string(contentBytes), e
}

// printStringLines prints each line in str, appending a line number to each line.
func printStringLines(str string) {
	var scanner = bufio.NewScanner(strings.NewReader(str))
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%4d: %s\n", i, scanner.Text())
	}
}

// printCommitMessage prints the passed message if [program.args.Verbosity] is high enough.
func (prog *program) printCommitMessage(message string) {
	if prog.args.Verbosity < 1 {
		return
	}
	var delimiterLine = "=================================================="
	fmt.Println(delimiterLine)
	printStringLines(message)
	fmt.Println(delimiterLine)
}
