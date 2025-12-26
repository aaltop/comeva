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
