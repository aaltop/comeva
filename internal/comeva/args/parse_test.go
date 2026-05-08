package args

import (
	"slices"
	"testing"

	testingHelpers "comeva/internal/testing"
	stringUtils "comeva/internal/utils/strings"
)

var globalFlag = "--help --verbosity 1 --config-file config.yaml"
var commandsString = "command1 command2"
var subFlag = `--subFlag1 123 --subFlag2 "some flag value"`

// SplitCommandLine works to split the contents of the command line
// into global flags, commands, and command flags.
func TestParse(t *testing.T) {

	var args []string = stringUtils.SplitWords(" ", globalFlag, commandsString, subFlag)
	var cmdArgs *CommandArgs = ParseCommandArgs(args)

	var expectedGlobalFlags = &GlobalFlags{
		Help:       true,
		Verbosity:  1,
		ConfigFile: "config.yaml",
	}

	var receivedGlobalFlags = cmdArgs.GlobalFlags

	if *expectedGlobalFlags != *receivedGlobalFlags {
		t.Error(testingHelpers.ValueMismatch("GlobalFlags", *expectedGlobalFlags, *receivedGlobalFlags))
	}

	var expectedCommands = []string{"command1", "command2"}
	var receivedCommands = cmdArgs.Commands

	if slices.Compare(expectedCommands, receivedCommands) != 0 {
		t.Error(testingHelpers.ValueMismatch("Commands", expectedCommands, receivedCommands))
	}

}

func TestSplitCommandsAndFlags(t *testing.T) {

	var expectedCommands []string = stringUtils.SplitWords(" ", commandsString)
	var expectedSubFlags []string = stringUtils.SplitWords(" ", subFlag)
	var args []string = stringUtils.SplitWords(" ", commandsString, subFlag)

	var receivedCommands, receivedSubFlags []string = SplitCommandsAndFlags(args)

	if slices.Compare(expectedCommands, receivedCommands) != 0 {
		t.Error(testingHelpers.ValueMismatch("commands", expectedCommands, receivedCommands))
	}

	if slices.Compare(expectedSubFlags, receivedSubFlags) != 0 {
		t.Error(testingHelpers.ValueMismatch("sub-flags", expectedSubFlags, receivedSubFlags))
	}

}
