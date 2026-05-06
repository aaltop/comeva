package commands

import (
	"comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	exitstate "comeva/internal/exitState"
	"fmt"
)

// NewHelpCommand creates a help command.
func NewHelpCommand() (helpCommand *Command, e error) {
	var helpMsg *HelpMessage

	var usage = HelpMessageUsage{}
	usage.AddExample(
		fmt.Sprintf("%s <item>", CommandNames.Help),
		`Get help on the given item.`,
	)

	var options = newDefaultFlagOptions()
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
		func(gFlags *args.GlobalFlags, passedGFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
			extState = exitstate.NewDefaultExitState()
			fmt.Println(helpMsg.String())
			return
		},
		make(CommandList),
	)
	if e != nil {
		return
	}
	return
}
