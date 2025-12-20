package main

import (
	"bufio"
	headerValidation "comeva/validators/header"
	"flag"
	"fmt"
	"os"
)

// processCommitMessage prints information about invalid formatting
// in a commit message, and returns a (possibly empty) slice of lines
// that may contain yaml content: if the commit message contains three dashes,
// everything after these are assumed (but not tested) to be yaml content.
func processCommitMessage(scanner *bufio.Scanner) (yamlLines []string) {
	var line string
	var possibleYaml bool
	yamlLines = make([]string, 0)
	for i := 0; scanner.Scan(); i++ {
		line = scanner.Text()

		if line == "---" {
			if possibleYaml {
				fmt.Println("yaml content should not contain another document start, found one on line", i)
				continue
			}
			fmt.Println("Assumed yaml start on line", i+1)
			possibleYaml = true
			continue
		}

		// yaml content should be at the end of file, so just add the rest
		// of the lines once the three dashes have been encountered
		if possibleYaml {
			yamlLines = append(yamlLines, line)
			continue
		}

		switch i {
		case 0:
			var err error
			headerValidator, err := headerValidation.NewHeaderValidator(
				[]string{"feat", "fix"},
				[]string{"main", "validators"},
				[]string{"Add", "Change", "Remove", "Fix"},
				[2]uint{})
			if err != nil {
				panic("HeaderValidator should be valid")
			}
			err = headerValidator.ValidateString(line)
			fmt.Printf("Header: %v\nerror: %v\n", headerValidator.Header, err)
			fmt.Printf("Header values: %#v\n", headerValidator.Header)
			// // TODO: do this with a regex
			// var location, heading, hasLocation = strings.Cut(line, ":")
			// location = strings.ReplaceAll(location, " ", "")
			// if !hasLocation {
			// 	fmt.Println("Warning: no location set for heading")
			// 	fmt.Printf("Heading: %v\n", location)
			// } else {
			// 	heading = strings.Replace(heading, " ", "", 1)
			// 	fmt.Printf("Location: %v Heading: %v\n", location, heading)
			// }
		case 1:
			if len(line) != 0 {
				fmt.Println("Line after heading should be empty, found:", line)
			}
		}

	}
	return
}

// processCommitMessageFile prints information about invalid formatting
// in a commit message file, and returns a (possibly empty) slice of lines
// that may contain yaml content: if the commit message contains three dashes,
// everything after these are assumed (but not tested) to be yaml content.
func processCommitMessageFile(fileName string) (yamlLines []string) {
	var file, err = os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: Unable to open file: %v\n", err)
		return
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	yamlLines = processCommitMessage(scanner)
	return
}

var flagSet = flag.NewFlagSet("", flag.ExitOnError)

var helpFlag = flagSet.Bool("help", false, "print help")
var configFile = flagSet.String("config-file", "", "file path for configuration file")
var commitFile = flagSet.String("commit-file", "", "file path for commit file")

func main() {
	var argsOutput = flagSet.Output()

	var usageMessage = func() {

		fmt.Fprintln(argsOutput, "Usage:")
		fmt.Fprintln(argsOutput, "  comeva [flags] --commit-file <commit_file>")
		fmt.Fprintln(argsOutput, "  	Validate the commit message in the file <commit_file>.")
		fmt.Fprintln(argsOutput, "  comeva [flags] <commit_message>")
		fmt.Fprintln(argsOutput, "  	Validate the commit message <commit_message>, passed as a string.")
		fmt.Fprintln(argsOutput, "  comeva --help")
		fmt.Fprintln(argsOutput, "  	Print help.")
	}

	var helpMessage = func() {
		fmt.Fprint(argsOutput, "\nCoMeVa (Commit Message Validator) is a tool for validating git commit messages.\n\n")

		usageMessage()

		fmt.Fprintln(argsOutput)

		fmt.Fprint(argsOutput, "Options:\n")
		flagSet.PrintDefaults()
	}

	flagSet.Usage = usageMessage
	flagSet.Parse(os.Args[1:])

	var args []string = flagSet.Args()
	var commitFileSpecified = len(*commitFile) > 0
	var numArgs = len(args)
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

	if numArgs == 1 {
		var commitMessage string = args[0]
		fmt.Printf("Message:\n%s\n", commitMessage)
	}

	if len(*commitFile) < 1 {
		return
	}
	var yamlLines []string = processCommitMessageFile(*commitFile)
	fmt.Println("yaml:")
	for _, yamlLine := range yamlLines {
		fmt.Println(yamlLine)
	}
}
