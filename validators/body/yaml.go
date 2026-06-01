package body

import (
	"github.com/goccy/go-yaml"

	yml "comeva/internal/yaml"
)

type bodyValidator struct {
	LineLength [2]uint `yaml:"lineLength"`
}

func CreateCommentMap(prefix string) yaml.CommentMap {
	return yml.CreateCommentMap(prefix, yml.SuffixCommentMap{
		"lineLength": {" minimum and maximum for line length, all zeroes means no set limit"},
	})
}

func (validator *BodyValidator) UnmarshalYAML(data []byte) (e error) {
	var temp bodyValidator
	if e = yaml.Unmarshal(data, &temp); e == nil {
		var new *BodyValidator
		new, e = NewBodyValidator(temp.LineLength)
		*validator = *new
	}
	return e
}

func (validator *BodyValidator) MarshalYAML() (data []byte, e error) {
	var temp bodyValidator
	temp.LineLength = [2]uint{validator.lineLength.Lower, validator.lineLength.Upper}

	return yaml.MarshalWithOptions(&temp, yaml.WithComment(CreateCommentMap("")))
}
