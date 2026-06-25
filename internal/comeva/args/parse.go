package args

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	exitstate "comeva/internal/exitState"
	flagUtils "comeva/internal/flag"
	"comeva/internal/io/ansi"
	"comeva/internal/logging"
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

var helpFlag = globalFlagSet.Bool(string(GlobalFlagNames.Help), false, "Print help.")

const verbosityDefault int = globals.VERBOSITY_DEFAULT

var verbosity = globalFlagSet.Int(
	string(GlobalFlagNames.Verbosity), verbosityDefault,
	"Program verbosity, lower means less verbose, higher more verbose.")

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

// GlobalPassedFlags reports which flags of [GlobalFlags] were passed; see
// [GetPassedGlobalFlags].
type PassedGlobalFlags struct {
	Help, Verbosity, ConfigFile, LoggingLevel bool
}

func NewDefaultGlobalFlags() (gFlags *GlobalFlags) {
	return &GlobalFlags{}
}

// ParseGlobalFlags returns the global flags and any remaining command line
// arguments after flag parsing, as well as a map denoting which global flags
// were passed; see [GlobalFlagNames] for help with accessing the latter. `args`
// should be the command line arguments excluding the program name.
func ParseGlobalFlags(args []string) (gFlags *GlobalFlags, passedFlags *PassedGlobalFlags, remainingArgs []string) {
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

	passedFlags = &PassedGlobalFlags{}
	passedFlags = GetPassedGlobalFlags(flagUtils.PassedFlags(globalFlagSet))

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

	var word string
	var flagEncountered bool = false
	for _, word = range args {
		// could use indexing for faster assignments, but it's a little ugly
		// in some cases, and this isn't a very perfomance critical thing
		if flagEncountered || strings.HasPrefix(word, "-") {
			flagEncountered = true
			flgs = append(flgs, word)
		} else {
			cmds = append(cmds, word)
		}
	}

	return
}

// CommandArgs represents command line arguments passed to a program.
type CommandArgs struct {
	// Global flags passed to program. Come before commands.
	GlobalFlags *GlobalFlags
	// PassedGlobalFlags denote which global flags were passed.
	PassedGlobalFlags *PassedGlobalFlags
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

func GetPassedGlobalFlags(global map[string]bool) (passed *PassedGlobalFlags) {
	passed = &PassedGlobalFlags{}
	passed.ConfigFile = global[string(GlobalFlagNames.ConfigFile)]
	passed.Help = global[string(GlobalFlagNames.Help)]
	passed.LoggingLevel = global[string(GlobalFlagNames.LoggingLevel)]
	passed.Verbosity = global[string(GlobalFlagNames.Verbosity)]
	return
}

// joinGlobalArguments joins the global flags and configuration file values.
func JoinGlobalArguments(
	globalArgs *GlobalFlags,
	conf *config.Config,
	passedGlobal *PassedGlobalFlags,
	passedConfig *config.PassedConfigArgs,
) (joined *GlobalFlags) {
	joined = &GlobalFlags{}
	joined.Help = globalArgs.Help
	joined.Verbosity = globalArgs.Verbosity
	joined.ConfigFile = globalArgs.ConfigFile
	joined.LoggingLevel = globalArgs.LoggingLevel

	if !passedGlobal.Verbosity && passedConfig.Verbosity {
		joined.Verbosity = *conf.Verbosity
	}

	if !passedGlobal.LoggingLevel && passedConfig.LoggingLevel {
		joined.LoggingLevel = *conf.LoggingLevel
	}

	return
}
