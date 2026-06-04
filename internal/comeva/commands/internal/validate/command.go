// Contains the command function.

package validate

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	comevaIo "comeva/internal/comeva/io"
	"comeva/internal/errors"
	exitstate "comeva/internal/exitState"
	io "comeva/internal/io"
	"comeva/internal/io/ansi"
	"comeva/internal/utils/flag"
	"encoding/json"
	"fmt"
)

// program acts as the state of the command.
type program struct {
	JoinedLocalArgs  *args
	JoinedGlobalArgs *argus.GlobalFlags
	Config           *config.Config
	PassedArgs       *passedArgs
}

var colorSchemes = ansi.BasicColorSchemes

func Function(
	gFlags *argus.GlobalFlags,
	passedGlobalFlags *argus.PassedGlobalFlags,
	conf *config.Config,
	passedConfig *config.PassedConfigArgs,
) (extState *exitstate.ExitState) {

	extState = exitstate.NewDefaultExitState()
	var e error
	var arg *args
	arg, e = NewArgs()
	if e != nil {
		extState.Reason = fmt.Errorf("Error with validate arguments: %w", e)
		extState.Code = exitstate.PROGRAM_ERROR
		return
	}

	var passedLocal passedLocalArgs = *getPassedLocalArgs(flag.PassedFlags(FlagSet))

	var allPassed = &passedArgs{
		Local:  passedLocal,
		Global: *passedGlobalFlags,
		Config: *passedConfig,
	}

	var joinedLocal = joinLocalArguments(arg, conf, allPassed)

	var prog = &program{
		JoinedLocalArgs:  joinedLocal,
		JoinedGlobalArgs: gFlags,
		PassedArgs:       allPassed,
	}

	// SORT OUT FLAGS
	// ===================================================

	// VALIDATE MESSAGE
	// ---------------------------------------------------

	var commitMessage string
	commitMessage, e = io.ReadFileString(arg.CommitFile)

	switch joinedLocal.OutputFormat {
	case validOutputFormats.Human:
		prog.humanOutput(commitMessage)
	case validOutputFormats.JSON:
		prog.jsonOutput(commitMessage)
	default:
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Unexpected output format '%v', should be one of %v", joinedLocal.OutputFormat, validOutputFormatsSlice)
	}
	return
}

func (prog *program) humanOutput(commitMessage string) (extState *exitstate.ExitState) {
	extState = &exitstate.ExitState{}
	var e error
	if prog.JoinedGlobalArgs.Verbosity > globals.VERBOSITY_DEFAULT {
		comevaIo.PrintCommitMessage(commitMessage)
	}
	e = prog.validateCommitMessage(commitMessage)
	if e != nil {

		var unWrapped []error = errors.UnwrapAll(e)

		var prob string = "problem"
		if len(unWrapped) != 1 {
			prob = "problems"
		}

		fmt.Println(colorSchemes.Error.ApplyForef("%d %v found:", len(unWrapped), prob))
		for i, err := range unWrapped {
			print(colorSchemes.Error.ApplyForef("%d: ", i+1))
			fmt.Printf("%v\n", err)
		}
		extState.Code = exitstate.VALIDATION_ERROR
		return
	} else {
		fmt.Println("No problems found.")
	}

	return
}

func (prog *program) jsonOutput(commitMessage string) (extState *exitstate.ExitState) {
	extState = &exitstate.ExitState{}

	var validatedContent, e = prog.getValidatedContent(commitMessage)
	if e != nil {
		extState.Code = exitstate.VALIDATION_ERROR
	}
	var data []byte
	data = errors.Panic2(json.Marshal(validatedContent))
	fmt.Println(string(data))
	return
}
