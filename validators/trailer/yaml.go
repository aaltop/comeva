package trailer

import "github.com/goccy/go-yaml"

type trailerValidator struct {
	RequiredKeys       []Key   `yaml:"requiredKeys"`
	OptionalKeys       []Key   `yaml:"optionalKeys"`
	ContinuationIndent uint    `yaml:"continuationIndent"`
	LineLength         [2]uint `yaml:"lineLength"`
}

func (validator *TrailerValidator) UnmarshalYAML(data []byte) (e error) {
	var temp = trailerValidator{}
	if e = yaml.Unmarshal(data, &temp); e == nil {
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
