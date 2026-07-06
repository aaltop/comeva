package commands

import (
	initt "github.com/aaltop/comeva/internal/comeva/commands/internal/init"
)

// NewInitCommand returns a [Command] that represents an initialisation command
// used to initialise directory state related to the program.
func NewInitCommand() (initCommand *Command, e error) {

	var hlpMsgUsage = &HelpMessageUsage{}
	hlpMsgUsage.AddExample("init", "Use default values to initialise.")
	var hlpMsg *HelpMessage
	var flagOptions = newDefaultFlagOptions()
	flagOptions.AddGroupFlagSet("Optional/Boolean", initt.FlagSet)
	hlpMsg, e = NewHelpMessage(
		"Initialise directory state related to the program.",
		`Creates any missing configuration files in the current directory and fills them with
default values.`,
		*hlpMsgUsage,
		flagOptions,
	)

	initCommand, e = NewCommand(
		*hlpMsg,
		initt.Function,
		make(CommandList),
	)

	initCommand.CommandFlags = initt.FlagSet
	return
}
