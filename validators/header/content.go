package header

import (
	"github.com/aaltop/comeva/validators"
)

func (validator *HeaderValidator) ValidatedContent() validators.ValidatedContent[Header] {
	return validators.ValidatedContent[Header]{
		Content: validator.Header,
		Errors:  validator.Errors,
	}
}
