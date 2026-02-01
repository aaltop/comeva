package trailer

import (
	"bufio"
	"comeva/utils"
	"comeva/validators"
	"errors"
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strings"
)

// Key represents a trailer key.
type Key struct {
	// Value represents the value of the key.
	Value string
	// Info optionally describes the key.
	Info string
}

type Trailer struct {
	Key, Value string
}

func (t Trailer) String() string {
	return fmt.Sprintf("%s: %s", t.Key, t.Value)
}

// KeyMap wraps a map of "string: Key", implementing a more robust way to add values.
type KeyMap map[string]Key

// keyChars matches in regex style the characters allowed in a trailer's key.
var keyChars string = `[A-Za-z]+(?:-[A-Za-z]+)*`

// keyRegex matches a trailer's key.
var keyRegex = regexp.MustCompile(fmt.Sprintf(`\A%s\z`, keyChars))

// Set adds or replaces an item in the map.
func (keyMap *KeyMap) Set(key, info string) (e error) {
	if !keyRegex.MatchString(key) {
		return errors.New("key should contain only ASCII letters with dashes between")
	}
	(*keyMap)[key] = Key{Value: key, Info: info}
	return
}

// Get returns the value at key if found. "ok" reports whether the key was found.
func (keyMap *KeyMap) Get(key string) (value Key, ok bool) {
	value, ok = (*keyMap)[key]
	return
}

// TrailerValidator validates a git commit message's trailer block.
type TrailerValidator struct {
	// continuationRegex matches any continuation lines.
	continuationRegex *regexp.Regexp
	// requiredKeys contains the keys that are required to be found in any commit.
	// If this is empty, there are no required keys.
	requiredKeys KeyMap
	// optionalKeys contains the keys that are optional: if a key is in the trailers,
	// it should exist in either this map or in requiredKeys. If this is empty,
	// any key is valid.
	optionalKeys KeyMap
	// continuationIndent specifies how many spaces a value should be indented by.
	continuationIndent uint
	lineLength         utils.Bounds[uint]

	Trailers []Trailer
}

// Equal reports whether the two TrailerValidators are equal.
func (validator *TrailerValidator) Equal(other *TrailerValidator) bool {
	if validator == nil || other == nil {
		return validator == other
	}

	return validator.continuationRegex.String() == other.continuationRegex.String() &&
		maps.Equal(validator.requiredKeys, other.requiredKeys) &&
		maps.Equal(validator.optionalKeys, other.optionalKeys) &&
		validator.continuationIndent == other.continuationIndent &&
		validator.lineLength == other.lineLength
}

// SetLineLength sets new bounds for line length.
func (validator *TrailerValidator) SetLineLength(min, max uint) (e error) {
	h, e := utils.NewBounds(min, max, false, false)
	if e == nil {
		validator.lineLength = h
	}
	return e
}

// Reset resets any content set during validation.
func (validator *TrailerValidator) Reset() {
	validator.Trailers = []Trailer{}
}

// NewDefaultTrailerValidator creates the base (zero-value) TrailerValidator.
func NewDefaultTrailerValidator() (validator *TrailerValidator) {
	validator, e := NewTrailerValidator(make(KeyMap), make(KeyMap), 2, [2]uint{0, 0})
	if e != nil {
		panic(e)
	}
	return
}

// NewtrailerValidator returns a new TrailerValidator.
//
// For the *Keys arguments, the key should match the Value of the value's Key,
// as in *Keys[<key>] -> Key{ Value: <key>, Info: * }
func NewTrailerValidator(requiredKeys, optionalKeys KeyMap, continuationIndent uint, lineLength [2]uint) (validator *TrailerValidator, e error) {
	validator = &TrailerValidator{}
	var errs []error

	validator.requiredKeys = requiredKeys
	validator.optionalKeys = optionalKeys

	if e = validator.SetLineLength(lineLength[0], lineLength[1]); e != nil {
		errs = append(errs, e)
	}

	// this also sets the continuationRegex
	if e = validator.SetContinuationIndent(continuationIndent); e != nil {
		errs = append(errs, e)
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}

	return validator, e
}

// newContinuationRegex creates a new continuationRegex for TrailerValidator
// based on the new value given as continuationIndent. Should only be
// called by SetContinuationIndent outside testing setups.
func newContinuationRegex(continuationIndent uint) (regex *regexp.Regexp, e error) {
	regex, e = regexp.Compile(fmt.Sprintf(`\A[ ]{%d}\S[^\r\n]*\n?$`, continuationIndent))
	return
}

// setContinuationRegex sets continuationRegex for TrailerValidator
// based on the new value given as continuationIndent. Should only be
// called by SetContinuationIndent outside testing setups.
func (validator *TrailerValidator) setContinuationRegex(continuationIndent uint) (e error) {
	var regex *regexp.Regexp
	regex, e = newContinuationRegex(continuationIndent)
	if e == nil {
		validator.continuationRegex = regex
	}
	return
}

// SetContinuationIndent sets the continuation indent.
func (validator *TrailerValidator) SetContinuationIndent(indent uint) (e error) {
	if indent < 1 {
		e = errors.New("continuationIndent should be > 0")
		return
	}
	validator.continuationIndent = indent
	validator.setContinuationRegex(indent)
	return e
}

func (validator *TrailerValidator) Validate(reader io.Reader) (e error) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

// ValidateString validates a trailer block.
func (validator *TrailerValidator) ValidateString(possibleTrailer string) (e error) {
	return validator.ValidateStringWithLine(possibleTrailer, 1)
}

// ValidateStringWithLine validates a trailer block.
//
// `startLine` > 0 specifies the line
// on which the string starts in the original message, assuming that the trailer
// is a part of a longer message. This is currently only relevant for accurate
// reporting of the line number on which a validation error occurs.
func (validator *TrailerValidator) ValidateStringWithLine(possibleTrailer string, startLine uint) (e error) {
	return validator.ValidateScannerWithLine(bufio.NewScanner(strings.NewReader(possibleTrailer)), startLine-1)
}

type InvalidTrailerError struct {
	Reason string
	Line   uint
}

func (e InvalidTrailerError) Error() string {
	return fmt.Sprintf("Line %d: %s", e.Line, e.Reason)
}

// GetKeys returns all the keys held by validator. There is no guarantee of
// a specific order.
func (validator *TrailerValidator) GetKeys() (keys []string) {
	for k := range validator.requiredKeys {
		keys = append(keys, k)
	}
	for k := range validator.optionalKeys {
		keys = append(keys, k)
	}
	return
}

type InvalidKeyError struct {
	Expected []string
	Received string
	Line     uint
}

func (e InvalidKeyError) Error() string {
	return fmt.Sprintf("Line %d: Invalid trailer key %s, should be one of %v", e.Line, e.Received, e.Expected)
}

// ValidateKey checks whether the key of a trailer is valid (is included in the
// required or optional keys).
func (validator *TrailerValidator) ValidateKey(possibleKey string, lineNum uint) (e error) {

	var keyError = InvalidKeyError{Expected: validator.GetKeys(), Received: possibleKey, Line: lineNum}

	if !keyRegex.MatchString(possibleKey) {
		return keyError
	}

	var _, ok = validator.requiredKeys.Get(possibleKey)
	if ok {
		return nil
	}

	if len(validator.optionalKeys) == 0 {
		return nil
	}
	_, ok = validator.optionalKeys.Get(possibleKey)
	if ok {
		return nil
	}
	return keyError
}

type InvalidKeyValueError struct {
	Line uint
}

func (e InvalidKeyValueError) Error() string {
	return fmt.Sprintf("Line %d: No key-value pair found", e.Line)
}

func (validator *TrailerValidator) ValidateLineLength(line string, lineNum uint) (e error) {

	var lower, upper uint = validator.lineLength.Lower, validator.lineLength.Upper
	if lower == 0 && upper == 0 {
		return nil
	}

	var lineLength = uint(len(line))
	if !validator.lineLength.Contains(lineLength) {
		return validators.InvalidLineLengthError{
			Expected: validator.lineLength,
			Received: lineLength,
			Line:     lineNum}
	}
	return
}

// keyValueRegex matches the first (and possibly only) line of a git trailer.
var keyValueRegex = regexp.MustCompile(fmt.Sprintf(
	`\A%s: %s`,
	fmt.Sprintf(`(?<key>%s)`, keyChars),
	`(?<value>\S[^\r\n$]*)`))

// ResemblesKeyValue reports whether the string (assumed to be a line of text)
// resembles a git trailer's key-value pair. See [ValidateKeyValue] for proper
// validation.
func (validator *TrailerValidator) ResemblesKeyValue(possibleKeyValue string) bool {
	return keyValueRegex.MatchString(possibleKeyValue)
}

// ValidateKeyValue validates a trailer key-value pair. A key-value pair does
// not, by the validation definition used here, continue to another line.
func (validator *TrailerValidator) ValidateKeyValue(possibleKeyValue string, lineNum uint) (t Trailer, e error) {
	var errs []error

	var matches = keyValueRegex.FindStringSubmatch(possibleKeyValue)
	if matches == nil {
		errs = append(errs, InvalidKeyValueError{Line: lineNum})
		// can't really do anything else, so return early.
		return t, errors.Join(errs...)
	}

	// there isn't really any hard rules as to what the value in a trailer should
	// be, so not particularly testing that here.
	t.Key, t.Value = matches[keyValueRegex.SubexpIndex("key")], matches[keyValueRegex.SubexpIndex("value")]
	if e = validator.ValidateKey(t.Key, lineNum); e != nil {
		errs = append(errs, e)
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}
	return t, e
}

type InvalidValueContinuationError struct {
	Indent uint
	Line   uint
}

func (e InvalidValueContinuationError) Error() string {
	return fmt.Sprintf("Line %d: expecting value continuation with indent %d", e.Line, e.Indent)
}

// ValidateValueContinuation checks whether the line can be a valid continuation
// of the value of a key-value pair.
func (validator *TrailerValidator) ValidateValueContinuation(possibleContinuation string, lineNum uint) (e error) {

	var errs []error

	if e := validator.ValidateLineLength(possibleContinuation, lineNum); e != nil {
		errs = append(errs, e)
	}

	if !validator.continuationRegex.MatchString(possibleContinuation) {
		errs = append(
			errs,
			InvalidValueContinuationError{
				Indent: validator.continuationIndent,
				Line:   lineNum})
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}

	return e
}

type MissingRequiredKeyError struct {
	Missing []string
}

func (e MissingRequiredKeyError) Error() string {
	return fmt.Sprintf("Required keys %v not found in trailer", e.Missing)
}

// ValidateScanner validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner. `timesScanned` specifies
// the number of times .Scan() has been called on `scanner`, representing the
// line at which the the scanner is.
func (validator *TrailerValidator) ValidateScanner(scanner *bufio.Scanner) (e error) {
	return validator.ValidateScannerWithLine(scanner, 0)
}

// EnsureRequiredKeys checks that all required keys are found in the trailers
// after validation has been performed.
func (validator *TrailerValidator) EnsureRequiredKeys() (e error) {

	var missingRequired []string
	var presentKeys []string
	for _, trailer := range validator.Trailers {
		presentKeys = append(presentKeys, trailer.Key)
	}
	for key := range validator.requiredKeys {
		if !slices.Contains(presentKeys, key) {
			missingRequired = append(missingRequired, key)
		}
	}

	if len(missingRequired) > 0 {
		e = MissingRequiredKeyError{Missing: missingRequired}
	}
	return
}

// ValidateScanner validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner. `timesScanned` specifies
// the number of times .Scan() has been called on `scanner`, representing the
// line at which the the scanner is.
func (validator *TrailerValidator) ValidateScannerWithLine(scanner *bufio.Scanner, timesScanned uint) (e error) {

	// empty the Trailers in case this function has been called previously
	validator.Reset()

	var scner = utils.CountingScanner{Scanner: scanner}
	scner.TimesScanned = timesScanned
	var errs []error

	var trailer Trailer = Trailer{}
	var keyValueError, continuationError error
	for scner.Scan() {
		var line string = scner.Text()
		var tempTrailer Trailer

		if e = validator.ValidateLineLength(line, scner.TimesScanned); e != nil {
			errs = append(errs, e)
		}

		tempTrailer, keyValueError = validator.ValidateKeyValue(line, scner.TimesScanned)
		if keyValueError == nil {
			// valid key-value pair
			trailer = tempTrailer
			validator.Trailers = append(validator.Trailers, tempTrailer)
		} else {
			// invalid key-value pair

			// might be continuation instead
			if trailer.Key != "" {
				continuationError = validator.ValidateValueContinuation(line, scner.TimesScanned)
				if continuationError == nil {
					// is continuation, just add it to the value as-is
					validator.Trailers[len(validator.Trailers)-1].Value += "\n" + line
					continue
				} else {
					// not continuation
					errs = append(errs, continuationError)
				}

			} else {
				// First scanned line was not a valid key-value pair
				errs = append(errs, keyValueError)
			}

			// reset trailer so that next loop will not test for line continuation
			// if a valid key-value pair is not found
			//
			// Consequently, looping will next continue until a valid key-value pair
			// is found -- at which point the above testing is performed again --
			// or end-of-content is reached
			trailer = Trailer{}
		}

	}

	if e = validator.EnsureRequiredKeys(); e != nil {
		errs = append(errs, e)
	}

	if len(errs) > 0 {
		e = errors.Join(errs...)
	}
	return e
}
