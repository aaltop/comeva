package main

import (
	"fmt"
	"os"

	"github.com/aaltop/comeva/internal/comeva/args"
	"github.com/aaltop/comeva/internal/comeva/commands"
	"github.com/aaltop/comeva/internal/comeva/globals"
	"github.com/aaltop/comeva/internal/errors"
	exitstate "github.com/aaltop/comeva/internal/exitState"
)

type program struct {
}

func (prog *program) main() (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()

	defer func() {
		var panicValue any = recover()
		extState.HandlePanic(panicValue)
	}()

	var cmdArgs *args.CommandArgs = args.ParseCommandArgs(os.Args[1:])

	var base = errors.Panic2(commands.NewBaseCommand())
	var sub *commands.Command = base.GetSubCommand(cmdArgs.Commands)
	if sub != nil {
		extState = sub.Execute(cmdArgs)
	} else {
		if cmdArgs.PassedGlobalFlags.Help {
			fmt.Println(base.Help())
			return
		}
		// not unsuccessful per se, so don't set exitState
		globals.ErrorLogger.Error().Printf("Sub-command '%v' not found.", cmdArgs.Commands)
		fmt.Println(base.Help())
	}
	return
}

func main() {
	var p = program{}
	var state = p.main()
	state.Exit()
}
