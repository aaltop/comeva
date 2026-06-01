package header

import (
	"github.com/goccy/go-yaml"

	yml "comeva/internal/yaml"
)

type headerValidator struct {
	Types, Scopes, Verbs []string
	LineLength           [2]uint `yaml:"lineLength"`
}

// see [yml.CreateCommentMap].
func CreateCommentMap(prefix string) yaml.CommentMap {

	return yml.CreateCommentMap(prefix, yml.SuffixCommentMap{
		"types":      {" type according to conventional commits"},
		"scopes":     {" scope according to conventional commits"},
		"verbs":      {" imperative mood verb that begins the message of a header"},
		"lineLength": {" minimum and maximum for line length, all zeroes means no set limit"},
	})
}

func (validator *HeaderValidator) UnmarshalYAML(data []byte) (e error) {
	var hv headerValidator
	if e = yaml.Unmarshal(data, &hv); e == nil {
		var new *HeaderValidator
		new, e = NewHeaderValidator(
			hv.Types, hv.Scopes, hv.Verbs, hv.LineLength,
		)
		*validator = *new
	}
	return e
}

var comments = CreateCommentMap("")

func (validator *HeaderValidator) MarshalYAML() (data []byte, e error) {
	var temp headerValidator
	temp.Types = validator.types
	temp.Scopes = validator.scopes
	temp.Verbs = validator.verbs
	temp.LineLength = [2]uint{validator.lineLength.Lower, validator.lineLength.Upper}

	return yaml.MarshalWithOptions(&temp, yaml.WithComment(comments))
}
