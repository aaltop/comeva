// validation contains utilities used in validating the commit message.

package validate

import (
	"errors"
	"fmt"
	"os"

	exitState "comeva/internal/exitState"
	"comeva/internal/io/ansi"
	bodyValidation "comeva/validators/body"
	headerValidation "comeva/validators/header"
	messageValidation "comeva/validators/message"
	trailerValidation "comeva/validators/trailer"
)

func getHeaderValidator() (headerValidator *headerValidation.HeaderValidator) {
	return headerValidation.NewHeaderValidatorWithDefaults()
}

func getBodyValidator() (bodyValidator *bodyValidation.BodyValidator) {
	return bodyValidation.NewBodyValidatorWithDefaults()
}

func getTrailerValidator() (trailerValidator *trailerValidation.TrailerValidator) {
	return trailerValidation.NewTrailerValidatorWithDefaults()
}

func (prog *program) getMessageValidator() (messageValidator *messageValidation.MessageValidator) {
	var errorColor *ansi.ColorScheme = ansi.BasicColorSchemes.Error
	var e error

	messageValidator = &messageValidation.MessageValidator{}
	*messageValidator = *messageValidation.NewDefaultMessageValidator()
	// if validator settings are provided through a file (if they're not, the only
	// other option is the standard setup provided below this block)
	if prog.passedLocalFlags[string(flagNames.ValidatorConfigFile)] || prog.conf.ValidatorConfigFile != nil {

		var validatorConfigFile string = prog.Args.ValidatorConfigFile
		if !prog.passedLocalFlags[string(flagNames.ValidatorConfigFile)] {
			debugLogger.Debug().Printf(
				"Using validator config file location '%v' as specified in the config file",
				*prog.conf.ValidatorConfigFile,
			)
			validatorConfigFile = *prog.conf.ValidatorConfigFile
		} else {
			debugLogger.Debug().Printf("Using validator config file location '%v' as specified on the command line", prog.Args.ValidatorConfigFile)
		}

		var data []byte
		data, e = os.ReadFile(validatorConfigFile)
		if e != nil {
			panic(exitState.ExitState{
				Reason: errors.New(errorColor.ApplyForef(
					"Error reading validator config in '%s': %v",
					prog.Args.ValidatorConfigFile, e)),
				Code: exitState.PROGRAM_ERROR,
			})
		}
		e = messageValidator.UnmarshalYAML(data)
		if e != nil {
			panic(exitState.ExitState{
				Reason: errors.New(errorColor.ApplyForef(
					"Error unmarshaling validator config in '%s': %v",
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
