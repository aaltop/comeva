package comeva

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"comeva/misc/config"
	bodyValidation "comeva/validators/body"
	headerValidation "comeva/validators/header"
	messageValidation "comeva/validators/message"
	trailerValidation "comeva/validators/trailer"
)

// exitCode describes the exit code passed to os.Exit. Should be kept in the
// range [0, 125].
type exitCode uint8

const (
	// Program ran successfully.
	SUCCESSFUL exitCode = 0
	// Program ran unsuccessfully.
	PROGRAM_ERROR exitCode = 1
	// Program ran successfully, but a validation error was encountered.
	VALIDATION_ERROR exitCode = 2
	// A panic was encountered for which there is no known course of action.
	UNCAUGHT_ERROR exitCode = 3
)

// exitState describes the exit state of the program.
type exitState struct {
	// Reason describes the reason that caused the exit. This should only
	// be non-nil if the reason was unexpected, i.e. when it would make sense
	// to print the reason to stderr. For example, a validation error is an
	// expected error, so would not be printed to stderr. In contrast, an error reading
	// the config file for the validator is not expected, so should result
	// in the error being printed to stderr, and therefore set here.
	Reason error
	Code   exitCode
}

// Exit exits the program with the given ExitCode, printing first the given
// Reason for exiting, if any.
func (extState *exitState) Exit() {
	if extState.Reason != nil {
		fmt.Fprintln(os.Stderr, extState.Reason)
	}
	os.Exit(int(extState.Code))
}

// handlePanic will return an appropriate exitState for the program. To be called in a
// deferred call.
func (extState *exitState) handlePanic(panicValue any) {

	if panicValue == nil {
		return
	}
	var ok bool
	var temp exitState
	temp, ok = panicValue.(exitState)
	if ok {
		*extState = temp
		return
	}
	extState.Reason = fmt.Errorf("%v", panicValue)
	extState.Code = UNCAUGHT_ERROR
}

func getHeaderValidator() (headerValidator *headerValidation.HeaderValidator) {

	var e error
	headerValidator, e = headerValidation.NewHeaderValidator(
		[]string{"feat", "fix"},
		[]string{"main", "validators"},
		[]string{"Add", "Change", "Remove", "Update", "Fix"},
		[2]uint{0, 80})
	if e != nil {
		panic(fmt.Sprintf("HeaderValidator should be valid, got error: %v", e))
	}
	return headerValidator
}

func getBodyValidator() (bodyValidator *bodyValidation.BodyValidator) {
	var e error
	bodyValidator, e = bodyValidation.NewBodyValidator([2]uint{0, 80})
	if e != nil {
		panic(fmt.Sprintf("BodyValidator should be valid, got error: %v", e))
	}
	return bodyValidator
}

func getTrailerValidator() (trailerValidator *trailerValidation.TrailerValidator) {
	var e error
	var requiredKeys trailerValidation.KeyMap
	var optionalKeys trailerValidation.KeyMap
	trailerValidator, e = trailerValidation.NewTrailerValidator(
		requiredKeys,
		optionalKeys,
		2,
		[2]uint{0, 80},
	)
	if e != nil {
		panic(fmt.Sprintf("TrailerValidator should be valid, got error: %v", e))
	}
	return trailerValidator
}

func (prog *program) getMessageValidator() (messageValidator *messageValidation.MessageValidator) {
	var e error

	messageValidator = &messageValidation.MessageValidator{}
	*messageValidator = *messageValidation.NewDefaultMessageValidator()
	// if validator settings are provided through a file
	if len(prog.args.ValidatorConfigFile) > 0 {
		var data []byte
		data, e = os.ReadFile(prog.args.ValidatorConfigFile)
		if e != nil {
			panic(exitState{Reason: fmt.Errorf("Error reading validator config in '%s': %v\n", prog.args.ValidatorConfigFile, e), Code: PROGRAM_ERROR})
		}
		e = messageValidator.UnmarshalYAML(data)
		if e != nil {
			panic(exitState{
				Reason: fmt.Errorf(
					"Error unmarshaling validator config in '%s': %v\n",
					prog.args.ValidatorConfigFile, e),
				Code: PROGRAM_ERROR})
		}
		return messageValidator
	}

	messageValidator, e = messageValidation.NewMessageValidator(
		getHeaderValidator(),
		getBodyValidator(),
		getTrailerValidator(),
	)
	if e != nil {
		// not supposed to happen, so don't use exitState directly, at least no
		// real reason to
		panic(fmt.Errorf("Unexpected error while creating message validator: %v\n", e))
	}
	return messageValidator
}

// readFileString reads the contents of the given file, returning them as
// a string.
func readFileString(fileName string) (content string, e error) {
	var contentBytes []byte
	contentBytes, e = os.ReadFile(fileName)
	return string(contentBytes), e
}

func (prog *program) validateCommitMessage(message string) (e error) {

	var validator *messageValidation.MessageValidator
	validator = prog.getMessageValidator()
	e = validator.ValidateString(message)
	return e
}

func printStringLines(str string) {
	var scanner = bufio.NewScanner(strings.NewReader(str))
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%4d: %s\n", i, scanner.Text())
	}
}

func (prog *program) printCommitMessage(message string) {
	if prog.args.Verbosity < 1 {
		return
	}
	var delimiterLine = "=================================================="
	fmt.Println(delimiterLine)
	printStringLines(message)
	fmt.Println(delimiterLine)
}

// printAndValidateMessage prints the commit message if the verbosity is high
// enough, then validates the commit message.
func (prog *program) printAndValidateMessage(message string) (e error) {
	prog.printCommitMessage(message)
	return prog.validateCommitMessage(message)
}

type program struct {
	args   *args
	config *config.Config
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
