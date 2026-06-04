package validate

import (
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	"flag"
	"fmt"
	"slices"
	"strings"
)

var FlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

// flagNames holds the string name of each command line flag.
var flagNames = struct {
	ValidatorConfigFile,
	CommitFile,
	OutputFormat FlagString
}{
	ValidatorConfigFile: "validator-config-file",
	CommitFile:          "commit-file",
	OutputFormat:        "output-format",
}
var validatorConfigFile = FlagSet.String(
	string(flagNames.ValidatorConfigFile), globals.VALIDATOR_CONFIG_PATH,
	"file path for configuration of validators")
var commitFile = FlagSet.String(string(flagNames.CommitFile), globals.COMMIT_MESSAGE_PATH, "File path for commit file. The special value '-' indicates stdin, allowing piping.")

// validOutputFormats holds the valid values for the [outputFormat] flag.
var validOutputFormats = struct {
	Human, JSON string
}{
	Human: "human",
	JSON:  "json",
}

// see [outputFormat].
var validOutputFormatsSlice = []string{"human", "json"}
var outputFormat = FlagSet.String(string(flagNames.OutputFormat), "human", fmt.Sprintf("output format of validation, %v", strings.Join(validOutputFormatsSlice, "|")))

// args handles the arguments passed to the program.
type args struct {
	ValidatorConfigFile, CommitFile, OutputFormat string
}

func NewArgs() (arg *args, e error) {
	arg = &args{}

	arg.ValidatorConfigFile = *validatorConfigFile
	arg.CommitFile = *commitFile
	arg.OutputFormat = *outputFormat
	if !slices.Contains(validOutputFormatsSlice, arg.OutputFormat) {
		e = fmt.Errorf("Given output format '%v' not in valid output formats %v", arg.OutputFormat, validOutputFormats)
		return
	}

	return
}

// argsFromConfig returns [args] based on the contents of the passed [config.Config].
func argsFromConfig(conf *config.Config) (args *args, e error) {
	args, e = NewArgs()
	if e != nil {
		return
	}
	if conf.CommitFile != nil {
		args.CommitFile = *conf.CommitFile
	}

	if conf.ValidatorConfigFile != nil {
		args.ValidatorConfigFile = *conf.ValidatorConfigFile
	}
	return
}
