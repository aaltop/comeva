package json

import (
	"bytes"
	"encoding/json"
)

type encoder struct {
	Encoder *json.Encoder
	Buffer  *bytes.Buffer
}

func newDefaultEncoder() (enc *encoder) {
	var buf = &bytes.Buffer{}
	var baseEncoder = json.NewEncoder(buf)
	baseEncoder.SetEscapeHTML(false)
	enc = &encoder{
		Encoder: baseEncoder,
		Buffer:  buf,
	}
	return
}

// Marshal matches [json.Marshal], but with default encoding
// options.
func Marshal(v any) ([]byte, error) {
	var enc = newDefaultEncoder()
	enc.Encoder.Encode(v)
	// for some reason, the Encoder's .Encode adds a newline, which isn't
	// done by just the base json.Marshal, I don't think. So, just trim any
	// space from the start and end (start shouldn't have any to begin with,
	// but then it's fine to trim it anyway)
	return bytes.TrimSpace(enc.Buffer.Bytes()), nil
}

// Unmarshal matches [json.Unmarshal], but with default encoding options.
func Unmarshal(data []byte, v any) error {
	// Don't want to actually do anything else for now, so just call the
	// original
	return json.Unmarshal(data, v)
}
