package init

import (
	"comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	errorHelpers "comeva/internal/errors"
	exitstate "comeva/internal/exitState"
	messageValidation "comeva/validators/message"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type program struct {
	gFlags            *args.GlobalFlags
	passedGlobalFlags map[string]bool
	conf              *config.Config
}

func Function(gFlags *args.GlobalFlags, passedGlobalFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var e error

	e = os.Mkdir(globals.CONFIG_BASE_PATH, 0775)
	if e != nil && !errors.Is(e, fs.ErrExist) {
		extState.Reason = errors.New(colorSchemes.Error.ApplyForef("Error creating config directory: %v", e))
		extState.Code = exitstate.PROGRAM_ERROR
	}

	var prog = &program{gFlags: gFlags, passedGlobalFlags: passedGlobalFlags, conf: conf}

	var configFileName = globals.CONFIG_PATH
	e = prog.createConfigFile(configFileName)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Error creating config file '%v': %v", configFileName, e)
		return
	}

	var validatorConfigFileName = globals.VALIDATOR_CONFIG_PATH
	e = prog.createValidatorFile(validatorConfigFileName)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Error creating validator file '%v': %v", validatorConfigFileName, e)
		return
	}

	return
}

func (prog *program) createConfigFile(filename string) (e error) {
	_, e = os.Stat(filename)

	if e == nil {
		fmt.Printf("Config file '%v' already exists, skipping creation\n", filename)
	} else if !errors.Is(e, fs.ErrNotExist) {
		if prog.gFlags.Verbosity >= globals.VERBOSITY_WARNING {
			globals.ErrorLogger.Warning().Printf("Warning: Error with config file: %v\n", e)
		}
	} else {
		var configFile *os.File
		configFile = errorHelpers.Panic2(os.Create(filename))
		var data []byte
		data = errorHelpers.Panic2(config.NewConfigWithDefaults().MarshalYAML())
		errorHelpers.Panic2(configFile.Write(data))
	}
	e = nil
	return
}

func (prog *program) createValidatorFile(filename string) (e error) {
	_, e = os.Stat(filename)

	if e == nil {
		fmt.Printf("Config file '%v' already exists, skipping creation\n", filename)
	} else if !errors.Is(e, fs.ErrNotExist) {
		if prog.gFlags.Verbosity >= globals.VERBOSITY_WARNING {
			globals.ErrorLogger.Warning().Printf("Warning: Error with validator config file: %v\n", e)
		}
	} else {
		var validatorFile *os.File
		validatorFile = errorHelpers.Panic2(os.Create(filename))
		var data []byte
		data = errorHelpers.Panic2(messageValidation.NewMessageValidatorWithDefaults().MarshalYAML())
		errorHelpers.Panic2(validatorFile.Write(data))
	}
	e = nil
	return
}
