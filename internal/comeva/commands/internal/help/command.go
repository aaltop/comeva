package help

import (
	"fmt"
	"strings"

	argus "github.com/aaltop/comeva/internal/comeva/args"
	"github.com/aaltop/comeva/internal/comeva/config"
	"github.com/aaltop/comeva/internal/comeva/docs"
	exitstate "github.com/aaltop/comeva/internal/exitState"
	flagUtils "github.com/aaltop/comeva/internal/flag"
)

type program struct {
	Args *args
}

func Function(gFlags *argus.GlobalFlags, passedGlobalFlags *argus.PassedGlobalFlags, conf *config.Config, passedConfig *config.PassedConfigArgs) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var e error

	var arg *args
	arg, _ = NewArgs()

	var passedLocalFlags = flagUtils.PassedFlags(FlagSet)

	if passedLocalFlags[string(flagNames.CreateDocs)] {
		e = docs.Create(arg.CreateDocs)
		if e != nil {
			extState.Code = exitstate.PROGRAM_ERROR
			extState.Reason = fmt.Errorf("Error creating docs directory: %w", e)
		}
		return
	}

	e = getDocsForPath(arg.Path...)
	if e != nil {
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = e
	}
	return
}

func getDocsForPath(path ...string) (e error) {
	if len(path) == 0 {
		path = []string{"", ""}
	}

	var content string
	if len(path) == 2 && path[0] == "" && path[1] == "" {
		// get root
		content, e = docs.Get()
	} else {
		content, e = docs.Get(path[1:]...)
	}
	if e != nil {
		return
	}

	var isFile = len(path) > 1 && path[len(path)-1] != ""

	var stringPath = strings.Join(path, "/")
	if isFile {
		fmt.Printf("\nDocumentation for '%v':\n", stringPath)
	} else {
		fmt.Printf("\nContents of '%v':\n", stringPath)
	}
	println(content)
	return
}
