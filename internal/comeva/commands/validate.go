package commands

import (
	"fmt"

	validate "comeva/internal/comeva/commands/internal/validate"
)

// NewValidateCommand returns a [Command] that represents the validator command.
func NewValidateCommand() (com *Command, e error) {

	var usage = HelpMessageUsage{}
	usage.AddExample(
		fmt.Sprintf("%s [flags]", CommandNames.Validate),
		`Validate a commit message that is in a file. If a given flag is not provided,
any possible defaults are used instead.`,
	)

	var options = newDefaultFlagOptions()
	options.AddGroupFlagSet("Optional", validate.OptionalFlagSet)
	options.AddGroupFlagSet("Boolean", validate.BooleanFlagSet)

	var helpMessage, _ = NewHelpMessage(
		"Validate a commit message.",
		"",
		usage,
		options,
	)

	com, e = NewCommand(
		*helpMessage,
		validate.Function,
		make(CommandList),
	)
	com.CommandFlags = validate.FlagSet
	return
}
