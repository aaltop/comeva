// Utilities for handling command line interaction.

package comeva

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"comeva/misc/config"
)

var flagSet = flag.NewFlagSet("", flag.ContinueOnError)

// flagString represents a the string name of a command line flag.
// Should not be instantiated directly, instead see [flagNames].
type flagString string

// flagNames holds the string name of each command line flag.
var flagNames = struct {
	Help, ConfigFile, ValidatorConfigFile,
	CommitFile, Verbosity flagString
}{
	Help:                "help",
	ConfigFile:          "config-file",
	ValidatorConfigFile: "validator-config-file",
	CommitFile:          "commit-file",
	Verbosity:           "verbosity",
}

var helpFlag = flagSet.Bool(string(flagNames.Help), false, "print help")
var configFile = flagSet.String(
	string(flagNames.ConfigFile), "",
	"file path for configuration file for setting command line values")
var validatorConfigFile = flagSet.String(
	string(flagNames.ValidatorConfigFile), "",
	"file path for configuration of validators")
var commitFile = flagSet.String(string(flagNames.CommitFile), "", "file path for commit file")

const verbosityDefault int = 0

var verbosity = flagSet.Int(
	string(flagNames.Verbosity), verbosityDefault,
	"program verbosity, lower means less verbose, higher more verbose")

// flagPassed reports whether the given flagString was passed on the
// command line.
func flagPassed(flagStr flagString) (passed bool) {
	// could also just run once for all, set in map, read from that?

	passed = false
	var checkFlag = func(fl *flag.Flag) {
		// found the flag, ignore rest
		if passed {
			return
		}
		if fl.Name == string(flagStr) {
			passed = true
		}
	}

	flagSet.Visit(checkFlag)
	return
}

// args handles the arguments passed to the program.
type args struct {
	HelpFlag                                    bool
	ConfigFile, ValidatorConfigFile, CommitFile string
	Verbosity                                   int

	// Number of non-flag arguments passed to the program.
	NumNonFlag int

	// Contains the raw non-flag args as they were passed at the command line.
	Raw []string
}

func newArgs() (arg *args, e error) {
	arg = &args{}

	flagSet.Usage = usageMessage
	// TODO: is it possible to suppress the "flag provided but not defined"?
	e = flagSet.Parse(os.Args[1:])
	if e != nil {
		return arg, fmt.Errorf("Error parsing arguments: %v\n", e)
	}

	arg.HelpFlag = *helpFlag
	arg.ConfigFile = *configFile
	arg.ValidatorConfigFile = *validatorConfigFile
	arg.CommitFile = *commitFile
	arg.Raw = flag.Args()
	arg.NumNonFlag = len(arg.Raw)
	arg.Verbosity = *verbosity

	return
}

// setFromConfig sets values in arg based on the config values in conf
// if arg does not have the value set yet (values passed directly
// overwrite those passed in the config file).
func (arg *args) setFromConfig(conf config.Config) (e error) {
	if !flagPassed(flagNames.CommitFile) {
		arg.CommitFile = conf.CommitFile
	}

	if !flagPassed(flagNames.ValidatorConfigFile) {
		arg.ValidatorConfigFile = conf.ValidatorConfigFile
	}

	if !flagPassed(flagNames.Verbosity) {
		arg.Verbosity = conf.Verbosity
	}

	return
}

// readFromConfig sets arg values from a config file if a config file
// was specified.
func (arg *args) readFromConfig() (e error) {
	if arg.ConfigFile == "" {
		return nil
	}

	var data []byte
	data, e = os.ReadFile(arg.ConfigFile)
	if e != nil {
		return fmt.Errorf("Error reading config file %s: %v", arg.ConfigFile, e)
	}

	var conf = config.NewDefaultConfig()
	e = conf.UnmarshalYAML(data)
	if e != nil {
		return fmt.Errorf("Error unmarshaling config file %s: %v", arg.ConfigFile, e)
	}

	arg.setFromConfig(*conf)
	return
}

// CommitFileSpecified reports whether the commitFile was specified as an argument.
func (arg *args) CommitFileSpecfied() (b bool) {
	return len(arg.CommitFile) > 0
}

// Validate validates the passed arguments.
func (arg *args) Validate() (e error) {

	var numArgs = arg.NumNonFlag
	var commitFileSpecified = arg.CommitFileSpecfied()

	switch {
	case (numArgs < 1 && !commitFileSpecified) || (numArgs == 1 && commitFileSpecified):
		e = errors.New("Error: specify either a commit message string or a commit message file.\n\n")
	case numArgs != 1 && !commitFileSpecified:
		e = fmt.Errorf("Error: expected one argument, got %d\n\n", numArgs)
	}
	return
}

func usageMessage() {
	var argsOutput = flagSet.Output()
	fmt.Fprintln(argsOutput, "Usage:")
	fmt.Fprintln(argsOutput, "  comeva [flags] --commit-file <commit_file>")
	fmt.Fprintln(argsOutput, "  	Validate the commit message in the file <commit_file>.")
	fmt.Fprintln(argsOutput, "  comeva [flags] <commit_message>")
	fmt.Fprintln(argsOutput, "  	Validate the commit message <commit_message>, passed as a string.")
	fmt.Fprintln(argsOutput, "  comeva --help")
	fmt.Fprintln(argsOutput, "  	Print help.")
}

// helpMessage prints the help message.
func helpMessage() {
	// TODO: should this maybe just create the message, then it could be output
	// to different streams depending on whether it's error output or --help output?
	var argsOutput = flagSet.Output()

	fmt.Fprint(argsOutput, "\nCoMeVa (Commit Message Validator) is a tool for validating git commit messages.\n\n")

	usageMessage()

	fmt.Fprintln(argsOutput)

	fmt.Fprint(argsOutput, "Options:\n")
	flagSet.PrintDefaults()
}
