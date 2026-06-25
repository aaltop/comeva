// Contains the command function.

package validate

import (
	"bufio"
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	comevaIo "comeva/internal/comeva/io"
	"comeva/internal/errors"
	exitstate "comeva/internal/exitState"
	"comeva/internal/io/ansi"

	"comeva/internal/flag"
	"comeva/internal/json"
	"fmt"
	baseIo "io"
	"os"
	"strings"
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

	var passedLocal passedLocalArgs = *getPassedLocalArgs(flag.PassedFlags(FlagSet))

	var allPassed = &passedArgs{
		Local:  passedLocal,
		Global: *passedGlobalFlags,
		Config: *passedConfig,
	}

	var _arg *args
	_arg, e = NewArgs()
	if e != nil {
		extState.Reason = fmt.Errorf("Error with validate arguments: %w", e)
		extState.Code = exitstate.PROGRAM_ERROR
		return
	}
	var joinedLocal = joinLocalArguments(_arg, conf, allPassed)

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

	// get a reader for the message content
	var messageReader baseIo.Reader
	if joinedLocal.CommitFile == "-" {
		messageReader = os.Stdin
	} else {
		messageReader, e = os.Open(joinedLocal.CommitFile)
		if e != nil {
			extState.Reason = fmt.Errorf("Error opening commit file '%v' for reading: %w", joinedLocal.CommitFile, e)
			extState.Code = exitstate.PROGRAM_ERROR
			return
		}
	}

	// messageHandler operates on the commit message depending on passed flags.
	var messageHandler func(msg string) *exitstate.ExitState
	// get the correct handler for the message
	switch joinedLocal.OutputFormat {
	case validOutputFormats.Human:
		messageHandler = prog.humanOutput
	case validOutputFormats.JSON:
		messageHandler = prog.jsonOutput
	default:
		extState.Code = exitstate.PROGRAM_ERROR
		extState.Reason = fmt.Errorf("Unexpected output format '%v', should be one of %v", joinedLocal.OutputFormat, validOutputFormatsSlice)
		return
	}

	// assume one message (no separator)
	if joinedLocal.CommitSeparator == nil {
		var data []byte
		data, e = baseIo.ReadAll(messageReader)
		if e != nil {
			extState.Code = exitstate.PROGRAM_ERROR
			if joinedLocal.CommitFile == "-" {
				extState.Reason = fmt.Errorf("Error reading commit file from stdin: %v", e)
			} else {
				extState.Reason = fmt.Errorf("Error reading commit file from file '%v': %v", joinedLocal.CommitFile, e)
			}
			return
		}
		commitMessage = string(data)
		*extState = *messageHandler(commitMessage)
	} else {
		var commitMessageGenerator = make(chan string)
		go getCommitMessages(messageReader, *joinedLocal.CommitSeparator, commitMessageGenerator)
		var msg string
		var tempExtState *exitstate.ExitState
		for msg = range commitMessageGenerator {
			tempExtState = messageHandler(msg)
			switch tempExtState.Code {
			case exitstate.SUCCESSFUL:
				// these are fine, continue
			case exitstate.VALIDATION_ERROR:
				extState.Code = exitstate.VALIDATION_ERROR
				if extState.Reason == nil {
					extState.Reason = tempExtState.Reason
				} else {
					extState.Reason = fmt.Errorf("%w\n%w", extState.Reason, tempExtState.Reason)
				}
			default:
				*extState = *tempExtState
				return
			}
		}

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

func getCommitMessages(reader baseIo.Reader, separator string, c chan string) {
	defer func() {
		close(c)
	}()

	var scanner = bufio.NewScanner(reader)
	var builder = strings.Builder{}
	var after string
	var found bool
	for scanner.Scan() {
		after, found = strings.CutPrefix(scanner.Text(), separator)
		if found && after == "" {
			c <- builder.String()
			builder.Reset()
		} else {
			fmt.Fprintln(&builder, scanner.Text())
		}
	}

	// if messages didn't end in separator, output the final message
	if builder.String() != "" {
		c <- builder.String()
	}
}
