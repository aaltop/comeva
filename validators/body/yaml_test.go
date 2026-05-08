package body

import (
	testingUtils "comeva/internal/comeva/testing"
	"fmt"
	"strings"
	"testing"
)

func errorMsgFmt(expected, received any) string {
	return fmt.Sprintf("Expected: %v Received: %v", expected, received)
}

func TestLineLength(t *testing.T) {
	var validator = &BodyValidator{}

	var expected = [2]uint{0, 80}
	var expectedString = []string{fmt.Sprint(expected[0]), fmt.Sprint(expected[1])}
	var yamlText = fmt.Sprintf("lineLength: [%s]", strings.Join(expectedString, ", "))
	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	if expected[0] != validator.lineLength.Lower || expected[1] != validator.lineLength.Upper {
		t.Error(errorMsgFmt(expected, validator.lineLength))
	}
}

// Content can be marshaled and unmarshaled into the same format.
func TestMarshalUnmarshal(t *testing.T) {
	var e error
	var expected = NewBodyValidatorWithDefaults()

	var data []byte
	data, e = expected.MarshalYAML()

	if e != nil {
		t.Fatalf("Unexpected error: %v", e)
	}

	var received = NewDefaultBodyValidator()
	received.UnmarshalYAML(data)

	if !expected.Equal(received) {
		t.Log(testingUtils.ValueMismatch("BodyValidator", *expected, *received))
	}
}
