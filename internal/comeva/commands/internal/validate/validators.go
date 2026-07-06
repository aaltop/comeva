// validation contains utilities used in validating the commit message.

package validate

import (
	baseErrors "errors"
	"fmt"
	"os"

	"github.com/aaltop/comeva/internal/errors"
	exitState "github.com/aaltop/comeva/internal/exitState"
	"github.com/aaltop/comeva/internal/io/ansi"
	"github.com/aaltop/comeva/validators"
	bodyValidation "github.com/aaltop/comeva/validators/body"
	headerValidation "github.com/aaltop/comeva/validators/header"
	messageValidation "github.com/aaltop/comeva/validators/message"
	trailerValidation "github.com/aaltop/comeva/validators/trailer"
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
	if prog.PassedArgs.Local.ValidatorConfigFile || prog.PassedArgs.Config.ValidatorConfigFile {

		var validatorConfigFile string = prog.JoinedLocalArgs.ValidatorConfigFile
		if !prog.PassedArgs.Local.ValidatorConfigFile {
			debugLogger.Debug().Printf(
				"Using validator config file location '%v' as specified in the config file",
				validatorConfigFile,
			)
		} else {
			debugLogger.Debug().Printf("Using validator config file location '%v' as specified on the command line", validatorConfigFile)
		}

		var data []byte
		data, e = os.ReadFile(validatorConfigFile)
		if e != nil {
			panic(exitState.ExitState{
				Reason: baseErrors.New(errorColor.ApplyForef(
					"Error reading validator config in '%s': %v",
					validatorConfigFile, e)),
				Code: exitState.PROGRAM_ERROR,
			})
		}
		e = messageValidator.UnmarshalYAML(data)
		if e != nil {
			panic(exitState.ExitState{
				Reason: baseErrors.New(errorColor.ApplyForef(
					"Error unmarshaling validator config in '%s': %v",
					validatorConfigFile, e)),
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

	var validator = prog.getMessageValidator()
	e = errors.Join(validator.ValidateString(message)...)
	return e
}

func (prog *program) getValidatedContent(message string) (
	validatedContent validators.ValidatedContent[messageValidation.ValidatedMessage],
	e error,
) {
	var validator = prog.getMessageValidator()
	e = errors.Join(validator.ValidateString(message)...)
	validatedContent = validator.ValidatedContent()
	return
}
