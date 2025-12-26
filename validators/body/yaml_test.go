package body

import (
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
