package io

import (
	"bufio"
	"fmt"
	"strings"
)

// stringScanner creates a default [bufio.NewScanner] based on the passed
// string.
func stringScanner(str string) (scanner *bufio.Scanner) {
	return bufio.NewScanner(strings.NewReader(str))
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
