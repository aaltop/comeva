package init

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	errorHelpers "comeva/internal/errors"
	exitstate "comeva/internal/exitState"
	"comeva/internal/utils/flag"
	messageValidation "comeva/validators/message"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type program struct {
	JoinedLocalArgs   *args
	PassedLocalArgs   *passedArgs
	JoinedGlobalFlags *argus.GlobalFlags
	PassedGlobalFlags *argus.PassedGlobalFlags
	Conf              *config.Config
}

func Function(gFlags *argus.GlobalFlags, passedGlobalFlags *argus.PassedGlobalFlags, conf *config.Config, passedConfig *config.PassedConfigArgs) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var prog = &program{
		JoinedLocalArgs:   NewArgs(),
		PassedLocalArgs:   getPassedArgs(flag.PassedFlags(FlagSet)),
		JoinedGlobalFlags: gFlags,
		PassedGlobalFlags: passedGlobalFlags,
		Conf:              conf}

	if prog.JoinedLocalArgs.Global {
		extState = prog.createGlobalConfig()
	} else {
		extState = prog.createLocalConfig()
	}

	return
}

// createConfig creates the configuration `conf` and other configuration in the specified
// paths.
func (prog *program) createConfig(
	configDir, configFileName, validatorConfigFileName string,
	conf *config.Config,
) (extState *exitstate.ExitState) {
	extState = &exitstate.ExitState{}
	var e error

	e = createConfigDir(configDir, 0750)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = e
		return
	}

	e = prog.createConfigFile(configFileName, conf)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Error creating config file '%v': %v", configFileName, e)
		return
	}

	e = prog.createValidatorFile(validatorConfigFileName)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Error creating validator file '%v': %v", validatorConfigFileName, e)
		return
	}

	return
}

func (prog *program) createLocalConfig() (extState *exitstate.ExitState) {
	extState = &exitstate.ExitState{}

	var conf = config.NewConfigWithDefaults()
	return prog.createConfig(
		globals.CONFIG_BASE_PATH,
		globals.CONFIG_PATH,
		globals.VALIDATOR_CONFIG_PATH,
		conf,
	)
}

func createConfigDir(path string, mode os.FileMode) (e error) {
	e = os.Mkdir(path, mode)
	if e != nil && !errors.Is(e, fs.ErrExist) {
		e = fmt.Errorf("Error creating config directory: %v", e)
	} else {
		e = nil
	}
	return
}

func (prog *program) createGlobalConfig() (extState *exitstate.ExitState) {
	extState = &exitstate.ExitState{}
	var e error

	var configPaths *globals.ConfigVariables
	configPaths, e = globals.NewGlobalConfigVariables()
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = e
		return
	}
	var conf = config.NewConfigWithDefaults()
	*conf.ValidatorConfigFile = configPaths.ValidatorConfigPath
	prog.createConfig(
		configPaths.BasePath,
		configPaths.ConfigPath,
		configPaths.ValidatorConfigPath,
		conf,
	)
	return
}

func (prog *program) createConfigFile(filename string, conf *config.Config) (e error) {
	_, e = os.Stat(filename)

	if e == nil {
		fmt.Printf("Config file '%v' already exists, skipping creation\n", filename)
	} else if !errors.Is(e, fs.ErrNotExist) {
		globals.ErrorLogger.Warning().Printf("Warning: Error with config file: %v\n", e)
	} else {
		var configFile *os.File
		configFile = errorHelpers.Panic2(os.Create(filename))
		var data []byte
		data = errorHelpers.Panic2(conf.MarshalYAML())
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
		globals.ErrorLogger.Warning().Printf("Warning: Error with validator config file: %v\n", e)
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
