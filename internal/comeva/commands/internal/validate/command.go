// Contains the command function.

package validate

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	comevaIo "comeva/internal/comeva/io"
	exitstate "comeva/internal/exitState"
	io "comeva/internal/io"
	"comeva/internal/io/ansi"
	"comeva/internal/utils/flag"
	"errors"
	"fmt"
)

// program acts as the state of the command.
type program struct {
	Args              *args
	passedLocalFlags  map[string]bool
	passedGlobalFlags map[string]bool
	conf              *config.Config
}

func Function(gFlags *argus.GlobalFlags, passedGlobalFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState) {
	var colorSchemes = ansi.BasicColorSchemes

	extState = exitstate.NewDefaultExitState()
	var e error
	var arg *args
	arg, e = NewArgs()
	var passedLocalFlags map[string]bool = flag.PassedFlags(FlagSet)

	// SORT OUT FLAGS
	// ---------------------------------------------------

	// verbosity contains either the verbosity as passed as a flag, or the value
	// from the config file (or the default value if neither is passed).
	var verbosity int = gFlags.Verbosity
	if conf.Verbosity != nil && !passedGlobalFlags[string(argus.GlobalFlagNames.Verbosity)] {
		verbosity = *conf.Verbosity
	}

	if !passedLocalFlags[string(flagNames.CommitFile)] && conf.CommitFile != nil {
		arg.CommitFile = *conf.CommitFile
	}

	if !passedLocalFlags[string(flagNames.ValidatorConfigFile)] && conf.ValidatorConfigFile != nil {
		arg.ValidatorConfigFile = *conf.ValidatorConfigFile
	}

	// SORT OUT FLAGS
	// ===================================================

	// VALIDATE MESSAGE
	// ---------------------------------------------------

	var prog *program = &program{}
	prog.Args = arg
	prog.passedGlobalFlags = passedGlobalFlags
	prog.passedLocalFlags = passedLocalFlags
	prog.conf = conf

	var commitMessage string
	commitMessage, e = io.ReadFileString(arg.CommitFile)

	if e != nil {
		extState.Reason = errors.New(colorSchemes.Error.ApplyForef("Error reading commit message file: %v", e))
		extState.Code = exitstate.PROGRAM_ERROR
		return
	}

	if verbosity > globals.VERBOSITY_DEFAULT {
		comevaIo.PrintCommitMessage(commitMessage)
	}
	e = prog.validateCommitMessage(commitMessage)
	if e != nil {
		var unWrappable, ok = e.(interface{ Unwrap() []error })
		// unwrap and show the individual errors
		if ok {
			var unWrapped []error = unWrappable.Unwrap()

			var prob string = "problem"
			if len(unWrapped) != 1 {
				prob = "problems"
			}

			fmt.Println(colorSchemes.Error.ApplyForef("%d %v found:", len(unWrapped), prob))
			for i, err := range unWrapped {
				print(colorSchemes.Error.ApplyForef("%d: ", i+1))
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Printf("%v\n", e)
		}
		extState.Code = exitstate.VALIDATION_ERROR
		return
	} else {
		fmt.Println("No problems found.")
	}

	// VALIDATE MESSAGE
	// ===================================================

	return
}
