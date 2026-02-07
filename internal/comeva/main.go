package comeva

import (
	"fmt"

	"comeva/misc/config"
)

type program struct {
	args   *args
	config *config.Config
}

// printAndValidateMessage prints the commit message if the verbosity is high
// enough, then validates the commit message.
func (prog *program) printAndValidateMessage(message string) (e error) {
	prog.printCommitMessage(message)
	return prog.validateCommitMessage(message)
}

func (prog *program) main() (extState *exitState) {
	extState = &exitState{}
	var e error

	// allows nested functions to call panic to stop the execution of the program
	// while passing an exit code, but without using `os.Exit` directly.
	defer func() {
		var panicValue any = recover()
		extState.handlePanic(panicValue)
	}()

	prog.args, e = newArgs()
	if e != nil {
		extState.Reason = e
		extState.Code = PROGRAM_ERROR
		return
	}

	prog.args.readFromConfig()

	if prog.args.HelpFlag {
		helpMessage()
		extState.Code = SUCCESSFUL
		return
	}

	if e = prog.args.Validate(); e != nil {
		helpMessage()
		extState.Reason = e
		extState.Code = PROGRAM_ERROR
		return
	}

	var commitMessage string
	if prog.args.NumNonFlag == 1 {
		commitMessage = prog.args.Raw[0]
		e = prog.printAndValidateMessage(commitMessage)
		if e != nil {
			fmt.Printf("Found issues validating message:\n%v\n", e)
		} else {
			fmt.Println("No issues found validating message.")
		}

	} else {

		commitMessage, e = readFileString(prog.args.CommitFile)
		if e != nil {
			extState.Reason = fmt.Errorf("Error reading commit file: %v\n", e)
			extState.Code = PROGRAM_ERROR
			return extState
		}
		e = prog.printAndValidateMessage(commitMessage)
		if e != nil {
			fmt.Printf("Found issues validating file '%s':\n%v\n", prog.args.CommitFile, e)
		} else {
			fmt.Printf("No issues found with file '%s'\n", prog.args.CommitFile)
		}
	}

	if e != nil {
		extState.Code = VALIDATION_ERROR
	}
	return
}

// Main runs the program.
func Main() {
	var p = program{}
	var exitState = p.main()
	exitState.Exit()
}
