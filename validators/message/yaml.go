package message

import (
	"comeva/validators/body"
	"comeva/validators/header"
	"comeva/validators/trailer"

	"github.com/goccy/go-yaml"
)

type messageValidator struct {
	HeaderValidator  *header.HeaderValidator   `yaml:"header"`
	BodyValidator    *body.BodyValidator       `yaml:"body"`
	TrailerValidator *trailer.TrailerValidator `yaml:"trailer"`
}

func (validator *MessageValidator) UnmarshalYAML(data []byte) (e error) {
	var temp messageValidator
	if e = yaml.Unmarshal(data, &temp); e == nil {
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

func (validator *MessageValidator) MarshalYAML() (data []byte, e error) {
	var temp messageValidator
	temp.HeaderValidator = validator.HeaderValidator
	temp.BodyValidator = validator.BodyValidator
	temp.TrailerValidator = validator.TrailerValidator

	return yaml.Marshal(&temp)
}
