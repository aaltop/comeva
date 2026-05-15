// globals contains global variables and other globally usable values.
package globals

import (
	"comeva/internal/logging"
	"log"
	"os"
)

// CONFIG_BASE_PATH is the base path of the config directory for the program.
const CONFIG_BASE_PATH string = "./.comeva/"

// VERBOSITY_DEFAULT is the default value for the verbosity of the program.
const VERBOSITY_DEFAULT int = 0

// VERBOSITY_WARNING is the verbosity level at and beyond which warnings
// should be printed.
const VERBOSITY_WARNING int = 10

// DebugLogger is a logger with an output particularly suitable for
// debugging.
var DebugLogger = logging.CreateColoredLogger()

// ErrorLogger suitable for logging error-related messages. Outputs
// to [os.stderr].
var ErrorLogger = logging.CreateColoredLogger()

func init() {
	DebugLogger.Logger = logging.NewDebugLogger(log.Writer())
	var debugLevel int = 0
	DebugLogger.SetLevel(&debugLevel)
	DebugLogger.SetOutput(log.Writer())
	ErrorLogger.SetOutput(os.Stderr)
}
