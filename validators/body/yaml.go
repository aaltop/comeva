package body

import (
	"github.com/goccy/go-yaml"
)

type bodyValidator struct {
	LineLength [2]uint `yaml:"lineLength"`
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

	return yaml.Marshal(&temp)
}
