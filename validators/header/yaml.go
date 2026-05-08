package header

import (
	"github.com/goccy/go-yaml"
)

type headerValidator struct {
	Types, Scopes, Verbs []string
	LineLength           [2]uint `yaml:"lineLength"`
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

func (validator *HeaderValidator) MarshalYAML() (data []byte, e error) {
	var temp headerValidator
	temp.Types = validator.types
	temp.Scopes = validator.scopes
	temp.Verbs = validator.verbs
	temp.LineLength = [2]uint{validator.lineLength.Lower, validator.lineLength.Upper}

	return yaml.Marshal(&temp)
}
