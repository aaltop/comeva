package message

import (
	"comeva/validators/body"
	"comeva/validators/header"
	"comeva/validators/trailer"
)

// Message represents the commit message.
type Message struct {
	Header  header.Header
	Body    body.Body
	Trailer []trailer.Trailer
}

// ValidatedContent provides easy access to the content parsed and validated by
// the validator.
func (validator *MessageValidator) ValidatedContent() (cont *Message) {
	cont = &Message{}
	cont.Header = *validator.HeaderValidator.ValidatedContent()
	cont.Body = validator.BodyValidator.ValidatedContent()
	cont.Trailer = validator.TrailerValidator.ValidatedContent()

	return
}
