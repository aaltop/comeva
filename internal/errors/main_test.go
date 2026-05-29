package errors

import (
	testingUtils "comeva/internal/testing"
	"errors"
	baseErrors "errors"
	"testing"
)

func handleReturn(in string, err error) (out string, e error) {

	defer func() {
		e = HandleReturn(recover())
	}()

	out = Return2(func() (value string, e error) {
		return in, err
	}())
	return
}

func TestReturn(t *testing.T) {

	var expectedStr = "no error"
	var receivedStr string
	var e error
	receivedStr, e = handleReturn(expectedStr, nil)

	if e != nil {
		t.Errorf("Expected nil error, got %v", e)
	}

	if receivedStr != expectedStr {
		t.Error(testingUtils.ValueMismatch("string", expectedStr, receivedStr))
	}

	var err = baseErrors.New("an error")
	receivedStr, e = handleReturn("an error", err)

	if e != err {
		testingUtils.ValueMismatch("error", err, e)
	}

	if receivedStr != "" {
		t.Error(testingUtils.ValueMismatch("string", "", receivedStr))
	}

}

func handlePanic(in string, err error) (out string, e error) {
	defer func() {
		var panicValue = recover()
		var panicErr, ok = panicValue.(error)
		if ok {
			e = panicErr
		}
	}()

	out = Panic2(func() (value string, e error) {
		return in, err
	}())
	return
}

func TestPanic(t *testing.T) {

	var e error
	var expectedStr, receivedStr string
	expectedStr = "no error"
	receivedStr, e = handlePanic(expectedStr, nil)

	if e != nil {
		t.Errorf("Expected nil error, got %v", e)
	}
	if receivedStr != expectedStr {
		t.Error(testingUtils.ValueMismatch("string", expectedStr, receivedStr))
	}

	var err = baseErrors.New(expectedStr)
	receivedStr, e = handlePanic("an error", err)

	if e != err {
		t.Error(testingUtils.ValueMismatch("error", err, e))
	}

	if receivedStr != "" {
		t.Error(testingUtils.ValueMismatch("string", "", receivedStr))
	}

}

// Unwrapping all errors works.
func TestUnwrapAllErrors(t *testing.T) {
	var expected = []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}
	var wrappedError error = baseErrors.Join(baseErrors.Join(expected[0], expected[1]), expected[2])

	var received []error = UnwrapAll(wrappedError)

	if len(expected) != len(received) {
		t.Fatalf("Expected same length\nExpected: %v\nReceived: %v\n", expected, received)
	}

	for i := range expected {
		if expected[i] != received[i] {
			t.Error(testingUtils.ValueMismatch("error", expected[i], received[i]))
		}
	}
}
