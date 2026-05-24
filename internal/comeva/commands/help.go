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
		"Show the documentation structure. This represents a directory tree which the path flag should be based on, not including the <root>.",
	)
	usage.AddExample(
		fmt.Sprintf("%s --path <path>", CommandNames.Help),
		`Get help on the given item specified by the path.`,
	)

	var options = newDefaultFlagOptions()
	options.AddGroup("Other", flagUtils.GetDefaults(help.FlagSet))
	helpMsg, e = NewHelpMessage(
		"Get help on a given item.",
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
