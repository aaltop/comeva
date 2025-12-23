package main

import (
	"bufio"
	bodyValidation "comeva/validators/body"
	headerValidation "comeva/validators/header"
	messageValidation "comeva/validators/message"
	trailerValidation "comeva/validators/trailer"
	"flag"
	"fmt"
	"os"
	"strings"
)

var flagSet = flag.NewFlagSet("", flag.ExitOnError)

var helpFlag = flagSet.Bool("help", false, "print help")
var configFile = flagSet.String("config-file", "", "file path for configuration file")
var commitFile = flagSet.String("commit-file", "", "file path for commit file")

func usageMessage() {
	var argsOutput = flagSet.Output()
	fmt.Fprintln(argsOutput, "Usage:")
	fmt.Fprintln(argsOutput, "  comeva [flags] --commit-file <commit_file>")
	fmt.Fprintln(argsOutput, "  	Validate the commit message in the file <commit_file>.")
	fmt.Fprintln(argsOutput, "  comeva [flags] <commit_message>")
	fmt.Fprintln(argsOutput, "  	Validate the commit message <commit_message>, passed as a string.")
	fmt.Fprintln(argsOutput, "  comeva --help")
	fmt.Fprintln(argsOutput, "  	Print help.")
}

func helpMessage() {
	var argsOutput = flagSet.Output()

	fmt.Fprint(argsOutput, "\nCoMeVa (Commit Message Validator) is a tool for validating git commit messages.\n\n")

	usageMessage()

	fmt.Fprintln(argsOutput)

	fmt.Fprint(argsOutput, "Options:\n")
	flagSet.PrintDefaults()
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

func getMessageValidator() (messageValidator *messageValidation.MessageValidator) {
	var e error
	messageValidator, e = messageValidation.NewMessageValidator(
		*getHeaderValidator(),
		*getBodyValidator(),
		*getTrailerValidator(),
	)
	if e != nil {
		panic(fmt.Sprintf("MessageValidator should be valid, got error: %v\n", e))
	}
	return messageValidator
}

func readFileString(fileName string) (content string, e error) {
	var contentBytes []byte
	contentBytes, e = os.ReadFile(fileName)
	return string(contentBytes), e
}

func validateCommitMessage(message string) (e error) {
	var validator *messageValidation.MessageValidator = getMessageValidator()
	e = validator.ValidateString(message)
	return e
}

func printStringLines(str string) {
	var scanner = bufio.NewScanner(strings.NewReader(str))
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%4d: %s\n", i, scanner.Text())
	}
}

func printCommitMessage(message string) {
	var delimiterLine = "=================================================="
	fmt.Println(delimiterLine)
	printStringLines(message)
	fmt.Println(delimiterLine)
}

func printAndValidateMessage(message string) (e error) {
	printCommitMessage(message)
	return validateCommitMessage(message)
}

func main() {
	var argsOutput = flagSet.Output()

	flagSet.Usage = usageMessage
	flagSet.Parse(os.Args[1:])

	var args []string = flagSet.Args()
	var commitFileSpecified = len(*commitFile) > 0
	var numArgs = len(args)

	// check validity of input
	switch {
	case *helpFlag:
		helpMessage()
		return
	case (numArgs < 1 && !commitFileSpecified) || (numArgs == 1 && commitFileSpecified):
		fmt.Fprint(argsOutput, "Error: specify either a commit message string or a commit message file.\n\n")
		helpMessage()
		return
	case numArgs != 1 && !commitFileSpecified:
		fmt.Fprintf(argsOutput, "Error: expected one argument, got %d\n\n", len(args))
		helpMessage()
		return
	}

	var e error
	var commitMessage string
	if numArgs == 1 {
		commitMessage = args[0]
		e = printAndValidateMessage(commitMessage)
		if e != nil {
			fmt.Printf("Found issues validating message:\n%v\n", e)
		} else {
			fmt.Println("No issues found validating message.")
		}

	} else {

		commitMessage, e = readFileString(*commitFile)
		if e != nil {
			fmt.Println(e)
		} else {
			e = printAndValidateMessage(commitMessage)
			if e != nil {
				fmt.Printf("Found issues validating file '%s':\n%v\n", *commitFile, e)
			} else {
				fmt.Printf("No issues found with file '%s'\n", *commitFile)
			}
		}
	}

	if e != nil {
		os.Exit(1)
	}
}
