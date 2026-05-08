package header

import (
	testingUtils "comeva/internal/comeva/testing"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func errorMsgFmt(expected, received any) string {
	return fmt.Sprintf("Expected: %v Received: %v", expected, received)
}

func TestTypes(t *testing.T) {
	var validator *HeaderValidator = &HeaderValidator{}

	var expected = []string{"type01", "type02"}
	var yamlText = fmt.Sprintf("types: [%s]", strings.Join(expected, ", "))

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Error(e)
	}

	t.Log(validator.types)

	if slices.Compare(validator.types, expected) != 0 {
		t.Error(errorMsgFmt(expected, validator.types))
	}
}
func TestScopes(t *testing.T) {
	var expected = []string{"scope01", "scope02"}
	var yamlText = fmt.Sprintf("scopes: [%s]", strings.Join(expected, ", "))

	var validator HeaderValidator
	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Error(e)
	}

	if slices.Compare(validator.scopes, expected) != 0 {
		t.Error(errorMsgFmt(expected, validator.scopes))
	}
}

func TestVerbs(t *testing.T) {
	var expected = []string{"verb01", "verb02"}
	var yamlText = fmt.Sprintf("verbs: [%s]", strings.Join(expected, ", "))

	var validator HeaderValidator
	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Error(e)
	}

	if slices.Compare(validator.verbs, expected) != 0 {
		t.Error(errorMsgFmt(expected, validator.verbs))
	}
}

func TestLineLength(t *testing.T) {
	var expected = []uint{0, 80}
	var expectedString = []string{fmt.Sprint(expected[0]), fmt.Sprint(expected[1])}
	var yamlText = fmt.Sprintf("lineLength: [%s]", strings.Join(expectedString, ", "))

	var validator HeaderValidator
	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Error(e)
	}

	if validator.lineLength.Lower != expected[0] || validator.lineLength.Upper != expected[1] {
		t.Error(errorMsgFmt(expected, validator.lineLength))
	}
}

// Content can be marshaled and unmarshaled into the same format.
func TestMarshalUnmarshal(t *testing.T) {
	var e error
	var expected = NewHeaderValidatorWithDefaults()

	var data []byte
	data, e = expected.MarshalYAML()

	if e != nil {
		t.Fatalf("Unexpected error: %v", e)
	}

	var received = NewDefaultHeaderValidator()
	received.UnmarshalYAML(data)

	if !expected.Equal(received) {
		t.Log(testingUtils.ValueMismatch("HeaderValidator", *expected, *received))
	}
}
