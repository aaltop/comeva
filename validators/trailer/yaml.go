package trailer

import (
	goccyYaml "github.com/goccy/go-yaml"

	"github.com/aaltop/comeva/internal/yaml"
)

type trailerValidator struct {
	RequiredKeys       []Key   `yaml:"requiredKeys"`
	OptionalKeys       []Key   `yaml:"optionalKeys"`
	ContinuationIndent uint    `yaml:"continuationIndent"`
	LineLength         [2]uint `yaml:"lineLength"`
}

// see [yaml.CreateCommentMap].
func CreateCommentMap(prefix string) goccyYaml.CommentMap {

	return yaml.CreateCommentMap(prefix, yaml.SuffixCommentMap{
		"requiredKeys": {" keys that have to exist in the trailer"},
		"optionalKeys": {
			" Keys that are optional.",
			" If specified, any key must match either these or requiredKeys.",
		},
		"continuationIndent": {" how many spaces to require a trailer's value to be indented by"},
		"lineLength":         {" minimum and maximum for line length, all zeroes means no set limit"},
	})

}

func (validator *TrailerValidator) UnmarshalYAML(data []byte) (e error) {
	var temp = trailerValidator{}
	if e = goccyYaml.Unmarshal(data, &temp); e == nil {
		var new = &TrailerValidator{}

		// indent is expected to be at least one
		var continuationIndent uint = 1
		if temp.ContinuationIndent > 0 {
			continuationIndent = temp.ContinuationIndent
		}

		new, e = NewTrailerValidator(make(KeyMap), make(KeyMap), continuationIndent, temp.LineLength)

		// set the keys
		for _, key := range temp.RequiredKeys {
			new.requiredKeys.Set(key.Value, key.Info)
		}
		for _, key := range temp.OptionalKeys {
			new.optionalKeys.Set(key.Value, key.Info)
		}

		*validator = *new
	}
	return e
}

var comments = CreateCommentMap("")

func (validator *TrailerValidator) MarshalYAML() (data []byte, e error) {
	var temp trailerValidator

	for _, v := range validator.requiredKeys {
		temp.RequiredKeys = append(temp.RequiredKeys, v)
	}

	for _, v := range validator.optionalKeys {
		temp.OptionalKeys = append(temp.OptionalKeys, v)
	}

	temp.ContinuationIndent = validator.continuationIndent
	temp.LineLength = [2]uint{validator.lineLength.Lower, validator.lineLength.Upper}

	return yaml.MarshalWithOptions(&temp, goccyYaml.WithComment(comments))
}
