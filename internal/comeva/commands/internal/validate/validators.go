// validation contains utilities used in validating the commit message.

package validate

import (
	"errors"
	"fmt"
	"os"

	"comeva/internal/comeva/io/ansi"
	exitState "comeva/internal/exitState"
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
	var errorColor *ansi.ColorScheme = ansi.BasicColorSchemes.Error
	var e error

	messageValidator = &messageValidation.MessageValidator{}
	*messageValidator = *messageValidation.NewDefaultMessageValidator()
	// if validator settings are provided through a file (if they're not, the only
	// other option is the standard setup provided below this block)
	if len(prog.Args.ValidatorConfigFile) > 0 {
		var data []byte
		data, e = os.ReadFile(prog.Args.ValidatorConfigFile)
		if e != nil {
			panic(exitState.ExitState{
				Reason: errors.New(errorColor.ApplyFore(
					"Error reading validator config in '%s': %v\n",
					prog.Args.ValidatorConfigFile, e)),
				Code: exitState.PROGRAM_ERROR,
			})
		}
		e = messageValidator.UnmarshalYAML(data)
		if e != nil {
			panic(exitState.ExitState{
				Reason: errors.New(errorColor.ApplyFore(
					"Error unmarshaling validator config in '%s': %v\n",
					prog.Args.ValidatorConfigFile, e)),
				Code: exitState.PROGRAM_ERROR})
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
