package args

import (
	"fmt"
	"slices"
	"testing"

	stringUtils "comeva/internal/strings"
	testingHelpers "comeva/internal/testing"
)

var globalFlag = "--help --verbosity 1 --config-file config.yaml --logging-level 14"
var commandsString = "command1 command2"

// Get
func FixtureSubFlag() []string {
	return []string{
		`--subFlag1 123 --subFlag2 "some flag value" --bool-flag`,
		`--bool-flag`,
	}
}

// SplitCommandLine works to split the contents of the command line
// into global flags, commands, and command flags.
func TestParse(t *testing.T) {

	for i, subFlag := range FixtureSubFlag() {
		t.Run(fmt.Sprintf("test %d", i+1), func(t *testing.T) {
			var args []string = stringUtils.SplitWords(" ", globalFlag, commandsString, subFlag)
			var cmdArgs *CommandArgs = ParseCommandArgs(args)

			var expectedGlobalFlags = &GlobalFlags{
				Help:         true,
				Verbosity:    1,
				ConfigFile:   "config.yaml",
				LoggingLevel: 14,
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
		})
	}

}

func TestSplitCommandsAndFlags(t *testing.T) {

	for i, subFlag := range FixtureSubFlag() {
		t.Run(fmt.Sprintf("test %d", i+1), func(t *testing.T) {
			var expectedCommands []string = stringUtils.SplitWords(" ", commandsString)
			var expectedSubFlags []string = stringUtils.SplitWords(" ", subFlag)
			var args []string = stringUtils.SplitWords(" ", commandsString, subFlag)

			var receivedCommands, receivedSubFlags []string = SplitCommandsAndFlags(args)

			if slices.Compare(expectedCommands, receivedCommands) != 0 {
				t.Error(testingHelpers.ValueMismatch("Commands", expectedCommands, receivedCommands))
			}

			if slices.Compare(expectedSubFlags, receivedSubFlags) != 0 {
				t.Error(testingHelpers.ValueMismatch("sub-flags", expectedSubFlags, receivedSubFlags))
			}
		})
	}

}
