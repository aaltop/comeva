package help

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/docs"
	exitstate "comeva/internal/exitState"
)

func Function(gFlags *argus.GlobalFlags, passedGlobalFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var e error

	var arg *args
	arg, _ = NewArgs()

	var content string
	content, e = docs.Get(arg.Path...)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = e
		return
	}

	println(content)
	return
}
