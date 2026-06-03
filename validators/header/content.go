package header

import (
	"comeva/validators"
)

func (validator *HeaderValidator) ValidatedContent() validators.ValidatedContent[Header] {
	return validators.ValidatedContent[Header]{
		Content: validator.Header,
		Errors:  validator.Errors,
	}
}
