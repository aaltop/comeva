package body

import (
	"comeva/validators"
)

func (validator *BodyValidator) ValidatedContent() validators.ValidatedContent[Body] {
	return validators.ValidatedContent[Body]{
		Content: validator.Body,
		Errors:  validator.Errors,
	}
}
