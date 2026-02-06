// Utilities for handling command line interaction.

package comeva

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var flagSet = flag.NewFlagSet("", flag.ContinueOnError)

var helpFlag = flagSet.Bool("help", false, "print help")
var configFile = flagSet.String("config-file", "", "file path for configuration file for setting command line values")
var validatorConfigFile = flagSet.String("validator-config-file", "", "file path for configuration of validators")
var commitFile = flagSet.String("commit-file", "", "file path for commit file")
var verbosity = flagSet.Int("verbosity", 0, "program verbosity, lower means less verbose, higher more verbose")

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
