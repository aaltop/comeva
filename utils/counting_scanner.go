package utils

import "bufio"

// countingScanner wraps bufio.Scanner to add a counter of times
// .Scan() has been called.
type CountingScanner struct {
	Scanner      *bufio.Scanner
	TimesScanned uint
}

func (scanner *CountingScanner) Scan() bool {
	var hadNext = scanner.Scanner.Scan()
	scanner.TimesScanned++
	return hadNext
}

func (scanner *CountingScanner) Text() string {
	return scanner.Scanner.Text()
}
