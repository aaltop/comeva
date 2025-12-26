package trailer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func errorMsgFmt(expected, received any) string {
	return fmt.Sprintf("Expected: %v Received: %v", expected, received)
}

func TestKey(t *testing.T) {
	var expected Key = Key{Value: "Dummy-Key", Info: "a placeholder key"}
	var yamlText = fmt.Sprintf("value: %s\ninfo: %s", expected.Value, expected.Info)

	var key = Key{}
	var e error
	if e = yaml.Unmarshal([]byte(yamlText), &key); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	if expected != key {
		t.Error(errorMsgFmt(expected, key))
	}
}

func TestRequiredKeys(t *testing.T) {
	var validator = TrailerValidator{}

	var expected = []Key{
		{
			Value: "Dummy-Key",
			Info:  "a placeholder key",
		},
	}
	var yamlText = fmt.Sprintf(`
requiredKeys:
  - value: %s
    info: %s`, expected[0].Value, expected[0].Info)

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	var ok bool = true
	for _, expectedKey := range expected {
		var receivedKey Key
		receivedKey, ok = validator.requiredKeys.Get(expectedKey.Value)
		if !ok {
			t.Errorf("Expected key '%s' not found: %v", expectedKey.Value, validator.requiredKeys)
			break
		}
		ok = receivedKey.Info == expectedKey.Info
		if !ok {
			t.Error(errorMsgFmt(expectedKey, receivedKey))
			break
		}
	}

}

func TestOptionaldKeys(t *testing.T) {
	var validator = TrailerValidator{}

	var expected = []Key{
		{
			Value: "Dummy-Key",
			Info:  "a placeholder key",
		},
	}
	var yamlText = fmt.Sprintf(`
optionalKeys:
  - value: %s
    info: %s`, expected[0].Value, expected[0].Info)

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	var ok bool = true
	for _, expectedKey := range expected {
		var receivedKey Key
		receivedKey, ok = validator.optionalKeys.Get(expectedKey.Value)
		if !ok {
			t.Errorf("Expected key '%s' not found: %v", expectedKey.Value, validator.optionalKeys)
			break
		}
		ok = receivedKey.Info == expectedKey.Info
		if !ok {
			t.Error(errorMsgFmt(expectedKey, receivedKey))
			break
		}
	}

}

func TestContinuationIndent(t *testing.T) {
	var validator = TrailerValidator{}

	var expected = 2
	var yamlText = fmt.Sprintf("continuationIndent: %d", expected)

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	if validator.continuationIndent != uint(expected) {
		t.Error(errorMsgFmt(expected, validator.continuationIndent))
	}
}

func TestLineLength(t *testing.T) {
	var validator = TrailerValidator{}

	var expected = [2]uint{0, 80}
	var expectedString = []string{fmt.Sprint(expected[0]), fmt.Sprint(expected[1])}
	var yamlText = fmt.Sprintf("lineLength: [%s]", strings.Join(expectedString, ", "))

	if e := validator.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Unexpected error: %v", e)
	}

	if validator.lineLength.Lower != expected[0] || validator.lineLength.Upper != expected[1] {
		t.Error(errorMsgFmt(expected, validator.lineLength))
	}
}
