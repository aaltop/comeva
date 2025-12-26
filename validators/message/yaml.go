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
		var new *MessageValidator
		new, e = NewMessageValidator(temp.HeaderValidator, temp.BodyValidator, temp.TrailerValidator)
		*validator = *new
	}
	return
}
