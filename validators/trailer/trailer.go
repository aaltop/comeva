package trailer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"maps"
	baseRegexp "regexp"
	"slices"
	"strings"

	"comeva/internal/regexp"
	"comeva/internal/utils"
	"comeva/validators"
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
var keyRegex = baseRegexp.MustCompile(fmt.Sprintf(`\A%s\z`, keyChars))

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
	continuationRegex *baseRegexp.Regexp
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

	Errors []validators.ValidatorErrorChild

	Trailers []Trailer
}

// NewDefaultTrailerValidator creates the base TrailerValidator.
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

// NewTrailerValidatorWithDefaults returns a [TrailerValidator] with default values.
func NewTrailerValidatorWithDefaults() (validator *TrailerValidator) {
	var e error
	var requiredKeys KeyMap
	var optionalKeys KeyMap
	validator, e = NewTrailerValidator(
		requiredKeys,
		optionalKeys,
		2,
		[2]uint{0, 80},
	)
	if e != nil {
		panic(e)
	}
	return
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
	validator.Errors = make([]validators.ValidatorErrorChild, 0)
}

// newContinuationRegex creates a new continuationRegex for TrailerValidator
// based on the new value given as continuationIndent. Should only be
// called by SetContinuationIndent outside testing setups.
func newContinuationRegex(continuationIndent uint) (regex *baseRegexp.Regexp, e error) {
	regex, e = baseRegexp.Compile(fmt.Sprintf(`\A[ ]{%d}\S[^\r\n]*\n?$`, continuationIndent))
	return
}

// setContinuationRegex sets continuationRegex for TrailerValidator
// based on the new value given as continuationIndent. Should only be
// called by SetContinuationIndent outside testing setups.
func (validator *TrailerValidator) setContinuationRegex(continuationIndent uint) (e error) {
	var regex *baseRegexp.Regexp
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

func (validator *TrailerValidator) ValidatedContent() validators.ValidatedContent[[]Trailer] {
	return validators.ValidatedContent[[]Trailer]{
		Content: validator.Trailers,
		Errors:  validator.Errors,
	}
}

func (validator *TrailerValidator) Validate(reader io.Reader) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateScanner(bufio.NewScanner(reader))
}

// ValidateString validates a trailer block.
func (validator *TrailerValidator) ValidateString(possibleTrailer string) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateStringWithLine(possibleTrailer, 1)
}

// ValidateStringWithLine validates a trailer block.
//
// `startLine` > 0 specifies the line
// on which the string starts in the original message, assuming that the trailer
// is a part of a longer message. This is currently only relevant for accurate
// reporting of the line number on which a validation error occurs.
func (validator *TrailerValidator) ValidateStringWithLine(possibleTrailer string, startLine uint) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateScannerWithLine(bufio.NewScanner(strings.NewReader(possibleTrailer)), startLine-1)
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

// validateKey checks whether the key of a trailer is valid (is included in the
// required or optional keys or is BREAKING-CHANGE). Returns [InvalidKeyError] if `e` is non-nil.
func (validator *TrailerValidator) validateKey(possibleKey regexp.SubMatch, lineNum uint) (errs []validators.ValidatorErrorChild) {

	if possibleKey.Match == "BREAKING-CHANGE" {
		return
	}

	var tempErrs = append(errs, InvalidKeyError{
		Expected: validator.GetKeys(),
		Received: possibleKey.Match,
		ValidatorError: validators.ValidatorError{
			MessagePart: validators.MessageParts.Trailer,
			Line:        lineNum,
			Cols:        possibleKey.Cols,
		},
	})

	if !keyRegex.MatchString(possibleKey.Match) {
		return tempErrs
	}

	var _, ok = validator.requiredKeys.Get(possibleKey.Match)
	if ok {
		return
	}

	if len(validator.optionalKeys) == 0 {
		return
	}
	_, ok = validator.optionalKeys.Get(possibleKey.Match)
	if ok {
		return
	}
	return tempErrs
}

func (validator *TrailerValidator) validateLineLength(line string, lineNum uint) (errs []validators.ValidatorErrorChild) {

	var lower, upper uint = validator.lineLength.Lower, validator.lineLength.Upper
	if lower == 0 && upper == 0 {
		return
	}

	var lineLength = uint(len(line))
	if !validator.lineLength.Contains(lineLength) {
		errs = append(errs, validators.InvalidLineLengthError{
			Expected: validator.lineLength,
			Received: lineLength,
			ValidatorError: validators.ValidatorError{
				Line:        lineNum,
				MessagePart: validators.MessageParts.Trailer,
			}})
	}
	return
}

// keyValueRegex matches the first (and possibly only) line of a git trailer.
var keyValueRegex = baseRegexp.MustCompile(fmt.Sprintf(
	`\A%s: %s`,
	fmt.Sprintf(`(?<key>%s)`, keyChars),
	`(?<value>\S[^\r\n$]*)`))

// ResemblesKeyValue reports whether the string (assumed to be a line of text)
// resembles a git trailer's key-value pair. See [validateKeyValue] for proper
// validation.
func (validator *TrailerValidator) ResemblesKeyValue(possibleKeyValue string) bool {
	return keyValueRegex.MatchString(possibleKeyValue)
}

// validateKeyValue validates a trailer key-value pair. A key-value pair does
// not, by the validation definition used here, continue to another line.
// The [Trailer] `t` will have empty strings for its values if a regexp match
// is not found, but may otherwise have non-empty strings for its values even
// if errors are returned.
func (validator *TrailerValidator) validateKeyValue(possibleKeyValue string, lineNum uint) (t Trailer, errs []validators.ValidatorErrorChild) {
	var ve validators.ValidatorErrorChild

	// keyValueRegex might encounter something with a correct key, but no value
	// specified. In this case, it's going to report that it did not find
	// a key-value pair, which is technically correct, but it might be nicer
	// to hint that something that could be a key-value pair (the key was
	// correct) was encountered.

	// var matches = keyValueRegex.FindStringSubmatch(possibleKeyValue)
	var matches = regexp.GetSubMatches(keyValueRegex, possibleKeyValue)
	if matches["key"].Cols == [2]int{-1, 1} {
		ve = InvalidKeyValueError{
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Trailer,
				Line:        lineNum,
			},
		}
		errs = append(errs, ve)
		// can't really do anything else, so return early.
		return
	}

	// there isn't really any hard rules as to what the value in a trailer should
	// be, so not particularly testing that here.
	t.Key, t.Value = matches["key"].Match, matches["value"].Match
	errs = append(errs, validator.validateKey(matches["key"], lineNum)...)

	return
}

// validateValueContinuation checks whether the line can be a valid continuation
// of the value of a key-value pair.
func (validator *TrailerValidator) validateValueContinuation(possibleContinuation string, lineNum uint) (errs []validators.ValidatorErrorChild) {
	var ve validators.ValidatorErrorChild

	errs = append(errs, validator.validateLineLength(possibleContinuation, lineNum)...)

	if !validator.continuationRegex.MatchString(possibleContinuation) {

		ve = InvalidValueContinuationError{
			Indent: validator.continuationIndent,
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Trailer,
				Line:        lineNum,
			},
		}
		errs = append(
			errs,
			ve,
		)
	}

	return
}

// ValidateScanner validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner. `timesScanned` specifies
// the number of times .Scan() has been called on `scanner`, representing the
// line at which the the scanner is.
func (validator *TrailerValidator) ValidateScanner(scanner *bufio.Scanner) (errs []validators.ValidatorErrorChild) {
	return validator.ValidateScannerWithLine(scanner, 0)
}

// EnsureRequiredKeys checks that all required keys are found in the trailers
// after validation has been performed.
func (validator *TrailerValidator) EnsureRequiredKeys() (errs []validators.ValidatorErrorChild) {

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
		errs = append(errs, MissingRequiredKeyError{
			Missing: missingRequired,
			ValidatorError: validators.ValidatorError{
				MessagePart: validators.MessageParts.Trailer,
			},
		})
	}
	return
}

// ValidateScannerWithLine validates the content returned by the scanner. `scanner` is expected
// to be a line-by-line scanner. `timesScanned` specifies
// the number of times .Scan() has been called on `scanner`, representing the
// line at which the the scanner is.
func (validator *TrailerValidator) ValidateScannerWithLine(scanner *bufio.Scanner, timesScanned uint) (errs []validators.ValidatorErrorChild) {

	// empty the Trailers in case this function has been called previously
	validator.Reset()

	// CountingScanner, namely TimesScanned, is used to keep track of the line
	// so the errors can report the correct line.
	var scner = utils.CountingScanner{Scanner: scanner}
	scner.TimesScanned = timesScanned

	defer func() {
		validator.Errors = errs
	}()

	var keyValueErrors, continuationErrors []validators.ValidatorErrorChild
	for scner.Scan() {
		var line string = scner.Text()
		var tempTrailer Trailer

		keyValueErrors, continuationErrors = make([]validators.ValidatorErrorChild, 0), make([]validators.ValidatorErrorChild, 0)

		errs = append(errs, validator.validateLineLength(line, scner.TimesScanned)...)

		tempTrailer, keyValueErrors = validator.validateKeyValue(line, scner.TimesScanned)

		if tempTrailer.Key != "" {
			// resembles a key-value pair, but with potentially invalid values
			// for the key (or value, if that is implemented)

			// add regardless, the errors will tell whether it's incorrect
			validator.Trailers = append(validator.Trailers, tempTrailer)
			errs = append(errs, keyValueErrors...)

		} else {
			// not possible to be key-value pair

			var prev = &(validator.Trailers[len(validator.Trailers)-1])

			// because it wasn't validated as a key-value pair, assume that
			// it is a continuation that has a problem
			prev.Value += "\n" + line
			continuationErrors = validator.validateValueContinuation(line, scner.TimesScanned)
			errs = append(errs, continuationErrors...)
		}

	}

	errs = append(errs, validator.EnsureRequiredKeys()...)

	return
}
