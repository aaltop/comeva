package math

import "testing"

func TestBounds(t *testing.T) {
	var bounds = Bounds[int]{Lower: 1, Upper: 3}

	for _, valAndExpected := range []([2]int){
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 0},
	} {
		var value = valAndExpected[0]
		var received = bounds.Contains(value)
		var expected = valAndExpected[1] == 1
		if received != expected {
			t.Errorf("Mismatch: expected %t, received %t for %d in %v", expected, received, value, bounds)
		}
	}

	bounds.Open = [2]bool{true, true}
	for _, valAndExpected := range []([2]int){
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 0},
		{4, 0},
	} {
		var value = valAndExpected[0]
		var received = bounds.Contains(value)
		var expected = valAndExpected[1] == 1
		if received != expected {
			t.Errorf("Mismatch: expected %t, received %t for %d in %v", expected, received, value, bounds)
		}
	}
}

// NewBounds method correctly checks whether the set bounds are in the correct
// order.
func TestBoundsVali(t *testing.T) {

	var err error
	if _, err = NewBounds(1, 1, false, false); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var bounds Bounds[int]
	if bounds, err = NewBounds(1, 1, true, false); err == nil {
		t.Errorf("Expected error for bounds %v", bounds)
	}
	if bounds, err = NewBounds(1, 1, false, true); err == nil {
		t.Errorf("Expected error for bounds %v", bounds)
	}
	if bounds, err = NewBounds(1, 1, true, true); err == nil {
		t.Errorf("Expected error for bounds %v", bounds)
	}
}
