package commands

import (
	"fmt"

	"comeva/internal/comeva/args"
	exitstate "comeva/internal/exitState"
)

// NewBaseCommand returns a [command] that represents the base command
// of the program.
func NewBaseCommand() (baseCommand *Command, e error) {

	var hlpMsgUsg = &HelpMessageUsage{}

	hlpMsgUsg.AddExample(
		CommandNames.Help,
		"Get help on an item.",
	)
	hlpMsgUsg.AddExample(
		CommandNames.Validate,
		"Validate a git commit message.",
	)

	var subCommands CommandList = make(CommandList)
	subCommands[CommandNames.Help], e = NewHelpCommand()
	subCommands[CommandNames.Validate], e = NewValidateCommand()

	if e != nil {
		return
	}

	var options = newDefaultFlagOptions()
	var hlpMsg *HelpMessage
	hlpMsg, e = NewHelpMessage(
		"CoMeVa (Commit Message Validator) is a tool for validating structured git commit messages.",
		"",
		*hlpMsgUsg,
		options,
	)

	if e != nil {
		return
	}

	baseCommand, _ = NewCommand(
		*hlpMsg,
		func(gFlags *args.GlobalFlags, passedGFlags map[string]bool) (extState *exitstate.ExitState) {
			extState = exitstate.NewDefaultExitState()
			fmt.Println(hlpMsg.String())
			return
		},
		subCommands,
	)
	return
}
