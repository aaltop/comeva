package message

import (
	"comeva/validators/body"
	"comeva/validators/header"
	"comeva/validators/trailer"
	"fmt"
	"strings"
	"testing"

	testingUtils "comeva/internal/testing"
)

func errorMsgFmt(expected, received any) string {
	return fmt.Sprintf("\nExpected: %#v\nReceived: %#v", expected, received)
}

func TestHeader(t *testing.T) {
	var validator = MessageValidator{}

	var expectedLineLength = [2]uint{0, 80}
	var expected, _ = header.NewHeaderValidator([]string{}, []string{}, []string{}, expectedLineLength)
	var expectedLineLengthString = []string{fmt.Sprint(expectedLineLength[0]), fmt.Sprint(expectedLineLength[1])}
	var yamlText = `
header:
  lineLength: [%s]`
	yamlText = fmt.Sprintf(yamlText, strings.Join(expectedLineLengthString, ", "))

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}
	// ensure that validation works when only one sub-validator is defined in config
	validator.ValidateString("Any old string")

	var received = validator.HeaderValidator

	if !expected.Equal(received) {
		t.Error(errorMsgFmt(*expected, received))
	}
}

func TestBody(t *testing.T) {
	var validator = MessageValidator{}

	var expectedLineLength = [2]uint{0, 80}
	var expected, _ = body.NewBodyValidator(expectedLineLength)
	var expectedLineLengthString = []string{fmt.Sprint(expectedLineLength[0]), fmt.Sprint(expectedLineLength[1])}
	var yamlText = `
body:
  lineLength: [%s]`
	yamlText = fmt.Sprintf(yamlText, strings.Join(expectedLineLengthString, ", "))

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}
	// ensure that validation works when only one sub-validator is defined in config
	validator.ValidateString("Any old string")

	var received = validator.BodyValidator

	if *expected != *received {
		t.Error(errorMsgFmt(*expected, received))
	}
}

func TestTrailer(t *testing.T) {
	var validator = MessageValidator{}

	var expectedLineLength = [2]uint{0, 80}
	var continuationIndent uint = 2
	var expected, _ = trailer.NewTrailerValidator(make(trailer.KeyMap), make(trailer.KeyMap), continuationIndent, expectedLineLength)
	var expectedLineLengthString = []string{fmt.Sprint(expectedLineLength[0]), fmt.Sprint(expectedLineLength[1])}
	var yamlText = `
trailer:
  lineLength: [%s]`
	yamlText = fmt.Sprintf(yamlText, strings.Join(expectedLineLengthString, ", "))

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}
	// ensure that validation works when only one sub-validator is defined in config
	validator.ValidateString("Any old string")

	var received = validator.TrailerValidator
	received.SetContinuationIndent(continuationIndent)

	if !expected.Equal(received) {
		t.Error(errorMsgFmt(*expected, received))
	}
}

// Content can be marshaled and unmarshaled into the same format.
func TestMarshalUnmarshal(t *testing.T) {
	var e error

	var expected = NewMessageValidatorWithDefaults()

	var data []byte
	data, e = expected.MarshalYAML()
	if e != nil {
		t.Fatalf("Unexpected error: %v", e)
	}

	var received = NewDefaultMessageValidator()
	received.UnmarshalYAML(data)

	if !expected.Equal(received) {
		t.Log(string(data))
		t.Logf("%v, %v", expected.HeaderValidator, received.HeaderValidator)
		t.Error(testingUtils.ValueMismatch("MessageValidator", *expected, *received))
	}
}
