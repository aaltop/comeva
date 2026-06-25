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

// IndentLines adds indentation to each line in `linedString`.
// Indentation is one space per `indentation`.
func IndentLines(linedString string, indentationAmount uint) string {
	return IndentLinesWithString(linedString, indentationAmount, " ")
}

// IndentLinesWithString works like [IndentLines], but instead of spaces, it
// uses `indentWith` for indentation.
func IndentLinesWithString(linedString string, indentationAmount uint, indentWith string) string {
	var indentation = strings.Repeat(indentWith, int(indentationAmount))
	var builder = strings.Builder{}
	for line := range strings.Lines(linedString) {
		fmt.Fprintf(&builder, "%v%v", indentation, line)
	}
	return builder.String()
}
