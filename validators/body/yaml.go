package body

import (
	goccyYaml "github.com/goccy/go-yaml"

	"github.com/aaltop/comeva/internal/yaml"
)

type bodyValidator struct {
	LineLength [2]uint `yaml:"lineLength"`
}

func CreateCommentMap(prefix string) goccyYaml.CommentMap {
	return yaml.CreateCommentMap(prefix, yaml.SuffixCommentMap{
		"lineLength": {" minimum and maximum for line length, all zeroes means no set limit"},
	})
}

func (validator *BodyValidator) UnmarshalYAML(data []byte) (e error) {
	var temp bodyValidator
	if e = goccyYaml.Unmarshal(data, &temp); e == nil {
		var new *BodyValidator
		new, e = NewBodyValidator(temp.LineLength)
		*validator = *new
	}
	return e
}

func (validator *BodyValidator) MarshalYAML() (data []byte, e error) {
	var temp bodyValidator
	temp.LineLength = [2]uint{validator.lineLength.Lower, validator.lineLength.Upper}

	return yaml.MarshalWithOptions(&temp, goccyYaml.WithComment(CreateCommentMap("")))
}
