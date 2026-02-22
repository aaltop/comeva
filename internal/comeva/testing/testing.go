package testing

import "fmt"

// ValueMismatch returns a string that compares an expected
// and received value, which are expected to have a mismatch. `typ`
// should be a word used to describe the types `expected` and `received`.
func ValueMismatch[T any](typ string, expected, received T) string {
	return fmt.Sprintf("%s mismatch:\nExpected:\n%v\nReceived:\n%v\n", typ, expected, received)
}
