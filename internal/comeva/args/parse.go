package args

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"comeva/internal/comeva/globals"
	exitstate "comeva/internal/exitState"
	"comeva/internal/io/ansi"
	"comeva/internal/logging"
	flagUtils "comeva/internal/utils/flag"
)

// flagString represents the string name of a command line flag.
// Should not be instantiated directly, instead see [GlobalFlagNames].
type flagString string

// // GlobalFlagNames holds the string name of each global command line flag.
var GlobalFlagNames = struct {
	Help, Verbosity, ConfigFile, LoggingLevel flagString
}{
	Help:         "help",
	Verbosity:    "verbosity",
	ConfigFile:   "config-file",
	LoggingLevel: "logging-level",
}

var globalFlagSet = flag.NewFlagSet("", flag.ContinueOnError)

var helpFlag = globalFlagSet.Bool(string(GlobalFlagNames.Help), false, "print help")

const verbosityDefault int = globals.VERBOSITY_DEFAULT

var verbosity = globalFlagSet.Int(
	string(GlobalFlagNames.Verbosity), verbosityDefault,
	"program verbosity, lower means less verbose, higher more verbose")

var configFile = globalFlagSet.String(
	string(GlobalFlagNames.ConfigFile), globals.CONFIG_PATH,
	"File path for configuration file for setting command line values. Any values given on the command line take precedence.")

var loggingLevel = globalFlagSet.Int(
	string(GlobalFlagNames.LoggingLevel), logging.ERROR,
	fmt.Sprintf(
		"Level set for logging. The levels are: debug: %d; info: %d; warning: %d; error: %d; critical: %d",
		logging.DEBUG, logging.INFO, logging.WARNING, logging.ERROR, logging.CRITICAL,
	),
)

type GlobalFlags struct {
	// Help reports whether a help flag was passed.
	Help bool
	// Verbosity indicates the level of verbosity of the program.
	Verbosity int
	// ConfigFile is the file path for a configuration file for setting command line values.
	ConfigFile string
	// LoggingLevel is the logging level.
	LoggingLevel int
}

func NewDefaultGlobalFlags() (gFlags *GlobalFlags) {
	return &GlobalFlags{}
}

// ParseGlobalFlags returns the global flags and any remaining command line
// arguments after flag parsing, as well as a map denoting which global flags
// were passed; see [GlobalFlagNames] for help with accessing the latter. `args`
// should be the command line arguments excluding the program name.
func ParseGlobalFlags(args []string) (gFlags *GlobalFlags, passedFlags map[string]bool, remainingArgs []string) {
	globalFlagSet.Usage = func() {}
	var e error

	var colorSchemes = ansi.BasicColorSchemes

	e = globalFlagSet.Parse(args)

	if e != nil {
		exitstate.NewExitState(
			errors.New(colorSchemes.Error.ApplyForef("Error parsing global flags: %v", e)),
			exitstate.PROGRAM_ERROR,
		).Panic()
	}

	passedFlags = flagUtils.PassedFlags(globalFlagSet)

	gFlags = NewDefaultGlobalFlags()
	gFlags.Help = *helpFlag
	gFlags.Verbosity = *verbosity
	gFlags.ConfigFile = *configFile
	gFlags.LoggingLevel = *loggingLevel

	return gFlags, passedFlags, globalFlagSet.Args()
}

// SplitCommandsAndFlags returns command arguments and flag arguments
// based on `args`. `args` is expected to not have any preceding flags.
func SplitCommandsAndFlags(args []string) (cmds []string, flgs []string) {

	var idx int
	var word string
	for idx, word = range args {
		if strings.HasPrefix(word, "-") {
			break
		}
	}

	if idx+1 < len(args) {
		cmds = args[:idx]
		flgs = args[idx:]
	} else {
		cmds = args
	}

	return
}

// CommandArgs represents command line arguments passed to a program.
type CommandArgs struct {
	// Global flags passed to program. Come before commands.
	GlobalFlags *GlobalFlags
	// PassedGlobalFlags denote which global flags were passed; see [GlobalFlagNames] for help with accessing.
	PassedGlobalFlags map[string]bool
	// Commands passed to program.
	Commands []string
	// The rest of the flags passed to the program, expected to be
	// specific to the commands (non-global). Come after commands.
	Flags []string
}

func newDefaultCommandArgs() (cmdArgs *CommandArgs) {
	cmdArgs = &CommandArgs{}
	cmdArgs.GlobalFlags = NewDefaultGlobalFlags()
	return
}

// ParseCommandArgs returns the content of the command line. `args` is assumed
// to be the content passed on the command line excluding the program
// name.
func ParseCommandArgs(args []string) (cmdArgs *CommandArgs) {
	cmdArgs = newDefaultCommandArgs()
	var remainingArgs []string
	cmdArgs.GlobalFlags, cmdArgs.PassedGlobalFlags, remainingArgs = ParseGlobalFlags(args)
	cmdArgs.Commands, cmdArgs.Flags = SplitCommandsAndFlags(remainingArgs)
	return
}

// GlobalFlagDefaults returns the value that would be printed by [flag.FlagSet.PrintDefaults]
// of the global flags.
func GlobalFlagDefaults() (defaults string) {
	defaults = flagUtils.GetDefaults(globalFlagSet)
	return
}
