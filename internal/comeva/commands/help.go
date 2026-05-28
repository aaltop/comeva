package commands

import (
	help "comeva/internal/comeva/commands/internal/help"
	flagUtils "comeva/internal/utils/flag"
	"fmt"
)

// NewHelpCommand creates a help command.
func NewHelpCommand() (helpCommand *Command, e error) {
	var helpMsg *HelpMessage

	var usage = HelpMessageUsage{}
	usage.AddExample(
		fmt.Sprintf("%s", CommandNames.Help),
		"Show the documentation structure. This represents a directory tree which the path flag should be based on.",
	)
	usage.AddExample(
		fmt.Sprintf("%s --path <path>", CommandNames.Help),
		`Get help on the given item specified by the path.`,
	)

	usage.AddExample(
		fmt.Sprintf("%s --create-docs=<path>", CommandNames.Help),
		"Recreate the entire documentation directory.",
	)

	var options = newDefaultFlagOptions()
	options.AddGroup("Local", flagUtils.GetDefaults(help.OptionalFlagSet))
	options.AddGroup("Local/Boolean", flagUtils.GetDefaults(help.OptionalBoolFlagSet))
	helpMsg, e = NewHelpMessage(
		"Access the documentation.",
		"",
		usage,
		options,
	)
	if e != nil {
		return
	}
	helpCommand, e = NewCommand(
		*helpMsg,
		help.Function,
		make(CommandList),
	)
	if e != nil {
		return
	}
	helpCommand.CommandFlags = help.FlagSet
	return
}
