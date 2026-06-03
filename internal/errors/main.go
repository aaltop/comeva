// package errors contains utitilies for handling errors.
package errors

import "errors"

type handledError struct {
	Error error
}

// HandleReturn is used in a defer function to handle panics caused by
// calls of ReturnX (e.g. [Return] or [Return2]). This NEEDS to be called
// if one of the ReturnX calls is made. Example defer:
//
//	defer func() {
//	  // e is the error returned by the function in which this defer is
//	  // called
//	  e = HandleReturn(recover())
//	}()
func HandleReturn(panicValue any) (e error) {
	var handled, ok = panicValue.(handledError)
	if ok {
		return handled.Error
	} else if panicValue != nil {
		panic(panicValue)
	}
	return
}

// Return sets (with the help of [HandleReturn]) the return error of
// the function where this is called to `error` if error is not nil. See
// Also other ReturnX functions: [Return2].
func Return(e error) {
	if e != nil {
		panic(handledError{Error: e})
	}
}

// Return2 sets (with the help of [HandleReturn]) the return error of
// the function where this is called to `error` if error is not nil, otherwise
// returning the arguments -- apart from the error -- as-is.
func Return2[T any](in T, e error) T {
	Return(e)
	return in
}

// Panic panics with the error if it is not nil. See also other PanicX functions:
// [Panic2].
func Panic(e error) {
	if e != nil {
		panic(e)
	}
}

// Panic2 panics with the error if it is not nil, otherwise returning the arguments
// -- apart from the error -- as-is.
func Panic2[T any](in T, e error) T {
	Panic(e)
	return in
}

// UnwrapAll returns the deeply unwrapped error.
func UnwrapAll(wrappedError error) (errors []error) {
	if unWrappable, ok := wrappedError.(interface{ Unwrap() []error }); ok {
		for _, e := range unWrappable.Unwrap() {
			errors = append(errors, UnwrapAll(e)...)
		}
	} else {
		errors = append(errors, wrappedError)
	}
	return
}

// Cast casts slices of types that implement the error interface into slices
// of pure errors.
func Cast[T error](errs ...T) []error {
	var cast = make([]error, len(errs))
	for i, e := range errs {
		cast[i] = e
	}
	return cast
}

// Join wraps [errors.Join], allowing values that implement the error interface
// to be passed.
func Join[T error](errs ...T) error {
	var cast = Cast(errs...)
	return errors.Join(cast...)
}
