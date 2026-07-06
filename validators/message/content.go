package message

import (
	"github.com/aaltop/comeva/validators"
	"github.com/aaltop/comeva/validators/body"
	"github.com/aaltop/comeva/validators/header"
	"github.com/aaltop/comeva/validators/trailer"
)

// ValidatedMessage represents the validated commit message.
type ValidatedMessage struct {
	Header  validators.ValidatedContent[header.Header]
	Body    validators.ValidatedContent[body.Body]
	Trailer validators.ValidatedContent[[]trailer.Trailer]
}

// ValidatedContent provides easy access to the content parsed and validated by
// the validator.
func (validator *MessageValidator) ValidatedContent() (cont validators.ValidatedContent[ValidatedMessage]) {
	var msg = ValidatedMessage{}
	msg.Header = validator.HeaderValidator.ValidatedContent()
	msg.Body = validator.BodyValidator.ValidatedContent()
	msg.Trailer = validator.TrailerValidator.ValidatedContent()

	cont = validators.ValidatedContent[ValidatedMessage]{
		Content: msg,
		Errors:  validator.Errors,
	}

	return
}
