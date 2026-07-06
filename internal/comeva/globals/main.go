// globals contains global variables and other globally usable values.
package globals

import (
	baseErrors "errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aaltop/comeva/internal/logging"
)

// CONFIG_BASE_PATH is the base path of the config directory for the program.
const CONFIG_BASE_PATH string = "./.comeva/"

// VALIDATOR_CONFIG_NAME is the default name for the validator configuration file.
const VALIDATOR_CONFIG_NAME = "validator.yaml"

// VALIDATOR_CONFIG_PATH is the default path for the validator configuration file.
const VALIDATOR_CONFIG_PATH = CONFIG_BASE_PATH + VALIDATOR_CONFIG_NAME

// CONFIG_NAME is the default name for the configuration file.
const CONFIG_NAME = "config.yaml"

// CONFIG_PATH is the default path for the configuration file.
const CONFIG_PATH = CONFIG_BASE_PATH + CONFIG_NAME

// COMMIT_MESSAGE_PATH is the default path for the git commit message file.
const COMMIT_MESSAGE_PATH = "./git_commit.txt"

// LOGGING_LEVEL_DEFAULT is the default level for logging.
const LOGGING_LEVEL_DEFAULT = logging.ERROR

// VERBOSITY_DEFAULT is the default value for the verbosity of the program.
const VERBOSITY_DEFAULT int = 0

// DebugLogger is a logger with an output particularly suitable for
// debugging.
var DebugLogger = logging.CreateColoredLogger()

// ErrorLogger suitable for logging error-related messages. Outputs
// to [os.stderr].
var ErrorLogger = logging.CreateColoredLogger()

func init() {
	DebugLogger.Logger = logging.NewDebugLogger(log.Writer())
	DebugLogger.SetOutput(log.Writer())
	ErrorLogger.SetOutput(os.Stderr)
}

var initialized = false

// InitArgs is the arguments passed to [Init]. Arguments that are left
// nil are ignored.
type InitArgs struct {
	LoggingLevel *int
	Debug        *bool
}

// Init allows setting of certain initial state. Will panic if called more
// than once.
func Init(args InitArgs) {

	if initialized {
		panic(baseErrors.New("Already initialised"))
	}
	initialized = true

	var debugLevel = 1000
	if args.Debug != nil && *args.Debug {
		debugLevel = 0
	}
	DebugLogger.SetLevel(&debugLevel)

	if args.LoggingLevel != nil {
		var level int = *args.LoggingLevel
		ErrorLogger.SetLevel(&level)
	}

}

type ConfigVariables struct {
	BasePath, ConfigPath, ValidatorConfigPath string
}

// NewGlobalConfigVariables returns the variables related to the default global
// configuration.
func NewGlobalConfigVariables() (paths *ConfigVariables, e error) {
	var base string
	base, e = os.UserConfigDir()
	if e != nil {
		e = fmt.Errorf("Error getting user config directory: %v", e)
		return
	}
	base = filepath.Join(base, "comeva")
	paths = &ConfigVariables{
		BasePath:            base,
		ConfigPath:          filepath.Join(base, CONFIG_NAME),
		ValidatorConfigPath: filepath.Join(base, VALIDATOR_CONFIG_NAME),
	}
	return
}
