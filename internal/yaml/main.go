package yaml

import (
	"fmt"
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

type SuffixCommentMap = map[string][]string

// CreateCommentMap creates a [goccyYaml.CommentMap]. `prefix` is a prefix
// for a YAML path, while `comments` consists of YAML path suffixes as keys and
// head comment lines as values.
//
// The `prefix` and suffixes of `comments` will be parsed as "$%s.%s".
func CreateCommentMap(prefix string, comments SuffixCommentMap) goccyYaml.CommentMap {
	var getKey = func(suffix string) string {
		return fmt.Sprintf("$%s.%s", prefix, suffix)
	}

	var commentMap = make(goccyYaml.CommentMap, 0)
	for key, value := range comments {
		commentMap[getKey(key)] = []*goccyYaml.Comment{
			goccyYaml.HeadComment(value...),
		}
	}
	return commentMap
}

var defaultMarshalOptions = []goccyYaml.EncodeOption{
	goccyYaml.IndentSequence(true),
}

// Marshal matches [goccyYaml.Marshal] but with default encoding options.
func Marshal(v any) ([]byte, error) {
	return goccyYaml.MarshalWithOptions(&v, defaultMarshalOptions...)
}

// MarshalWithOptions wraps [goccyYaml.MarshalWithOptions], setting some
// default EncodeOptions.
func MarshalWithOptions(v any, opts ...goccyYaml.EncodeOption) ([]byte, error) {
	var options = append(defaultMarshalOptions, opts...)
	return goccyYaml.MarshalWithOptions(&v, options...)
}
