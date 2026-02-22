package yaml

import (
	"os"

	goccyYaml "github.com/goccy/go-yaml"
)

type BytesUnmarshaler = goccyYaml.BytesUnmarshaler

// UnMarshalFromFile unmarshals the [BytesUnmarshaler] from the file pointed
// to by `filename`.
func UnMarshalFromFile(filename string, um BytesUnmarshaler) (e error) {
	var data []byte
	data, e = os.ReadFile(filename)
	if e != nil {
		return
	}
	e = um.UnmarshalYAML(data)
	return
}
