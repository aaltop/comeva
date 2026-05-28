package main

import (
	"comeva/internal/comeva/args"
	"comeva/internal/comeva/commands"
	"comeva/internal/comeva/globals"
	"comeva/internal/errors"
	exitstate "comeva/internal/exitState"
	"fmt"
	"os"
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
