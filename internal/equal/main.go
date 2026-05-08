// package equal contains utilities for testing for equality.
package equal

// NilEqual reports whether the two passed pointer values are equal, also considering
// nil values.
func NilEqual[T comparable](first, second *T) bool {
	if first == nil || second == nil {
		return first == second
	}
	return *first == *second
}
