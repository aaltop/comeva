package comeva

import (
	"fmt"
	"os"
)

// exitCode describes the exit code passed to os.Exit. Should be kept in the
// range [0, 125].
type exitCode uint8

const (
	// Program ran successfully.
	SUCCESSFUL exitCode = 0
	// Program ran unsuccessfully.
	PROGRAM_ERROR exitCode = 1
	// Program ran successfully, but a validation error was encountered.
	VALIDATION_ERROR exitCode = 2
	// A panic was encountered for which there is no known course of action.
	UNCAUGHT_ERROR exitCode = 3
)

// exitState describes the exit state of the program.
type exitState struct {
	// Reason describes the reason that caused the exit. This should only
	// be non-nil if the reason was unexpected, i.e. when it would make sense
	// to print the reason to stderr. For example, a validation error is an
	// expected error, so would not be printed to stderr. In contrast, an error reading
	// the config file for the validator is not expected, so should result
	// in the error being printed to stderr, and therefore set here.
	Reason error
	Code   exitCode
}

// Exit exits the program with the given [exitCode], printing first the given
// Reason for exiting, if any.
func (extState *exitState) Exit() {
	if extState.Reason != nil {
		fmt.Fprintln(os.Stderr, extState.Reason)
	}
	os.Exit(int(extState.Code))
}

// handlePanic will set an appropriate [exitState] based on the passed panicValue.
// To be called in a deferred call, where panicValue should be gotten from a
// recover() call made inside the deferred function.
func (extState *exitState) handlePanic(panicValue any) {

	if panicValue == nil {
		return
	}
	var ok bool
	var temp exitState
	temp, ok = panicValue.(exitState)
	if ok {
		*extState = temp
		return
	}
	extState.Reason = fmt.Errorf("%v", panicValue)
	extState.Code = UNCAUGHT_ERROR
}
