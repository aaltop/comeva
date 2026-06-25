package math

import (
	"cmp"
	"fmt"
)

type InvalidBoundsError[T cmp.Ordered] struct {
	Bounds Bounds[T]
	Reason string
}

func (e InvalidBoundsError[T]) Error() string {
	return fmt.Sprintf("Bounds %v are not valid: %s", e.Bounds, e.Reason)
}

// Bounds represents the bounds of an Ordered value.
type Bounds[T cmp.Ordered] struct {
	Lower, Upper T
	// Open signifies whether the corresponding bound is open or not
	Open [2]bool
}

// NewBounds returns a Bounds object, returning a non-nil error if
// the bounds are not valid.
func NewBounds[T cmp.Ordered](lower, upper T, lowerOpen, upperOpen bool) (bounds Bounds[T], e error) {
	bounds.Lower = lower
	bounds.Upper = upper
	bounds.Open = [2]bool{lowerOpen, upperOpen}
	e = bounds.Valid()
	return
}

func (bounds Bounds[T]) String() string {
	var lowerBound = "["
	if bounds.Open[0] {
		lowerBound = "]"
	}
	var upperBound = "]"
	if bounds.Open[1] {
		upperBound = "["
	}
	return fmt.Sprintf("%s%v,%v%s", lowerBound, bounds.Lower, bounds.Upper, upperBound)
}

// Equal reports whether the two Bounds are equal.
func (bounds Bounds[T]) Equal(other Bounds[T]) bool {
	return bounds.Lower == other.Lower &&
		bounds.Upper == other.Upper &&
		bounds.Open == other.Open
}

// Valid tests whether the bounds are valid, returning a non-nil error
// if they are not.
func (bounds Bounds[T]) Valid() (e error) {

	// if the values are equal and either one or both are open, the order
	// is wrong. If for example lower = upper = 1 and lower is open, then the actual
	// infimum > 1, while supremum == 1, which is a contradiction. If also upper is open,
	// then further supremum < 1, which is obviously still wrong.
	if bounds.Lower == bounds.Upper && (bounds.Open[0] || bounds.Open[1]) {
		return InvalidBoundsError[T]{Bounds: bounds, Reason: "lower and upper bound cannot be equal when either is open"}
	}

	if bounds.Lower <= bounds.Upper {
		return nil
	} else {
		return InvalidBoundsError[T]{Bounds: bounds, Reason: "lower bound is larger than upper bound"}
	}
}

func (bounds Bounds[T]) Contains(value T) (contains bool) {
	if bounds.Open[0] {
		contains = bounds.Lower < value
	} else {
		contains = bounds.Lower <= value
	}

	if !contains {
		return
	}

	if bounds.Open[1] {
		contains = value < bounds.Upper
	} else {
		contains = value <= bounds.Upper
	}

	return
}
