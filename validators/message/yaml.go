package message

import (
	"maps"

	"github.com/aaltop/comeva/validators/body"
	"github.com/aaltop/comeva/validators/header"
	"github.com/aaltop/comeva/validators/trailer"

	goccyYaml "github.com/goccy/go-yaml"

	"github.com/aaltop/comeva/internal/yaml"
)

type messageValidator struct {
	HeaderValidator  *header.HeaderValidator   `yaml:"header"`
	BodyValidator    *body.BodyValidator       `yaml:"body"`
	TrailerValidator *trailer.TrailerValidator `yaml:"trailer"`
}

func (validator *MessageValidator) UnmarshalYAML(data []byte) (e error) {
	var temp messageValidator
	if e = goccyYaml.Unmarshal(data, &temp); e == nil {
		var new *MessageValidator = NewDefaultMessageValidator()

		// Maybe would be easier to have the sub-validators be by value?
		// don't think would need to do this then, and it's not like it's
		// that big an issue, shouldn't need to have a bunch of validators.
		if temp.HeaderValidator != nil {
			*new.HeaderValidator = *temp.HeaderValidator
		}
		if temp.BodyValidator != nil {
			*new.BodyValidator = *temp.BodyValidator
		}
		if temp.TrailerValidator != nil {
			*new.TrailerValidator = *temp.TrailerValidator
		}

		*validator = *new
	}
	return
}

var headerComments = header.CreateCommentMap(".header")
var bodyComments = body.CreateCommentMap(".body")
var trailerComments = trailer.CreateCommentMap(".trailer")

func (validator *MessageValidator) MarshalYAML() (data []byte, e error) {
	var temp messageValidator
	temp.HeaderValidator = validator.HeaderValidator
	temp.BodyValidator = validator.BodyValidator
	temp.TrailerValidator = validator.TrailerValidator

	var comments = yaml.CreateCommentMap("", yaml.SuffixCommentMap{
		"header":  {" the header spans the first line of the message"},
		"body":    {" the body is between the header and the trailer"},
		"trailer": {" the trailer block is at the end of the message"},
	})

	maps.Copy(comments, headerComments)
	maps.Copy(comments, bodyComments)
	maps.Copy(comments, trailerComments)
	return yaml.MarshalWithOptions(
		&temp,
		goccyYaml.WithComment(comments),
	)
}
