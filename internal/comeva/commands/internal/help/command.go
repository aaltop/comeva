package help

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/docs"
	exitstate "comeva/internal/exitState"
	"fmt"
	"strings"
)

func Function(gFlags *argus.GlobalFlags, passedGlobalFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var e error

	var arg *args
	arg, _ = NewArgs()

	if len(arg.Path) == 0 {
		arg.Path = []string{"", ""}
	}

	var content string
	if len(arg.Path) == 2 && arg.Path[0] == "" && arg.Path[1] == "" {
		// get root
		content, e = docs.Get()
	} else {
		content, e = docs.Get(arg.Path[1:]...)
	}
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = e
		return
	}

	var isFile = len(arg.Path) > 1 && arg.Path[len(arg.Path)-1] != ""

	var path = strings.Join(arg.Path, "/")
	if isFile {
		fmt.Printf("\nDocumentation for '%v':\n", path)
	} else {
		fmt.Printf("\nContents of '%v':\n", path)
	}
	println(content)
	return
}
