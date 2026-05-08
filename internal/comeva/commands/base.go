package commands

import (
	"fmt"

	"comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/errors"
	exitstate "comeva/internal/exitState"
)

// NewBaseCommand returns a [Command] that represents the base command
// of the program.
func NewBaseCommand() (baseCommand *Command, e error) {

	defer func() {
		e = errors.HandleReturn(recover())
	}()

	var hlpMsgUsg = &HelpMessageUsage{}

	var subCommands CommandList = make(CommandList)
	subCommands[CommandNames.Help] = errors.Return2(NewHelpCommand())
	subCommands[CommandNames.Validate] = errors.Return2(NewValidateCommand())
	subCommands[CommandNames.Init] = errors.Return2(NewInitCommand())

	for k, v := range subCommands {
		hlpMsgUsg.AddExample(
			k,
			v.HelpMessage.synopsis,
		)
	}

	var options = newDefaultFlagOptions()
	var hlpMsg *HelpMessage
	hlpMsg = errors.Return2(NewHelpMessage(
		"CoMeVa (Commit Message Validator) is a tool for validating structured git commit messages.",
		"",
		*hlpMsgUsg,
		options,
	))

	baseCommand = errors.Return2(NewCommand(
		*hlpMsg,
		func(gFlags *args.GlobalFlags, passedGFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
			extState = exitstate.NewDefaultExitState()
			fmt.Println(hlpMsg.String())
			return
		},
		subCommands,
	))
	return
}
