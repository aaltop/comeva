// validation contains utilities used in validating the commit message.

package comeva

import (
	"fmt"
	"os"

	bodyValidation "comeva/validators/body"
	headerValidation "comeva/validators/header"
	messageValidation "comeva/validators/message"
	trailerValidation "comeva/validators/trailer"
)

func getHeaderValidator() (headerValidator *headerValidation.HeaderValidator) {

	var e error
	headerValidator, e = headerValidation.NewHeaderValidator(
		[]string{"feat", "fix"},
		[]string{"main", "validators"},
		[]string{"Add", "Change", "Remove", "Update", "Fix"},
		[2]uint{0, 80})
	if e != nil {
		panic(fmt.Sprintf("HeaderValidator should be valid, got error: %v", e))
	}
	return headerValidator
}

func getBodyValidator() (bodyValidator *bodyValidation.BodyValidator) {
	var e error
	bodyValidator, e = bodyValidation.NewBodyValidator([2]uint{0, 80})
	if e != nil {
		panic(fmt.Sprintf("BodyValidator should be valid, got error: %v", e))
	}
	return bodyValidator
}

func getTrailerValidator() (trailerValidator *trailerValidation.TrailerValidator) {
	var e error
	var requiredKeys trailerValidation.KeyMap
	var optionalKeys trailerValidation.KeyMap
	trailerValidator, e = trailerValidation.NewTrailerValidator(
		requiredKeys,
		optionalKeys,
		2,
		[2]uint{0, 80},
	)
	if e != nil {
		panic(fmt.Sprintf("TrailerValidator should be valid, got error: %v", e))
	}
	return trailerValidator
}

func (prog *program) getMessageValidator() (messageValidator *messageValidation.MessageValidator) {
	var e error

	messageValidator = &messageValidation.MessageValidator{}
	*messageValidator = *messageValidation.NewDefaultMessageValidator()
	// if validator settings are provided through a file
	if len(prog.args.ValidatorConfigFile) > 0 {
		var data []byte
		data, e = os.ReadFile(prog.args.ValidatorConfigFile)
		if e != nil {
			panic(exitState{Reason: fmt.Errorf("Error reading validator config in '%s': %v\n", prog.args.ValidatorConfigFile, e), Code: PROGRAM_ERROR})
		}
		e = messageValidator.UnmarshalYAML(data)
		if e != nil {
			panic(exitState{
				Reason: fmt.Errorf(
					"Error unmarshaling validator config in '%s': %v\n",
					prog.args.ValidatorConfigFile, e),
				Code: PROGRAM_ERROR})
		}
		return messageValidator
	}

	messageValidator, e = messageValidation.NewMessageValidator(
		getHeaderValidator(),
		getBodyValidator(),
		getTrailerValidator(),
	)
	if e != nil {
		// not supposed to happen, so don't use exitState directly, at least no
		// real reason to
		panic(fmt.Errorf("Unexpected error while creating message validator: %v\n", e))
	}
	return messageValidator
}

func (prog *program) validateCommitMessage(message string) (e error) {

	var validator *messageValidation.MessageValidator
	validator = prog.getMessageValidator()
	e = validator.ValidateString(message)
	return e
}
