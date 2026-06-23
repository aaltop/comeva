package validate

import (
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	"errors"
	"flag"
	"fmt"
	"slices"
	"strings"

	flagUtils "comeva/internal/utils/flag"
)

var FlagSet *flag.FlagSet
var OptionalFlagSet = flag.NewFlagSet("", flag.ContinueOnError)
var BooleanFlagSet = flag.NewFlagSet("", flag.ContinueOnError)

// var OptionalBoolFlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

// flagNames holds the string name of each command line flag.
var flagNames = struct {
	ValidatorConfigFile,
	CommitFile,
	OutputFormat,
	CommitSeparator FlagString
}{
	ValidatorConfigFile: "validator-config-file",
	CommitFile:          "commit-file",
	OutputFormat:        "output-format",
	CommitSeparator:     "commit-separator",
}
var validatorConfigFile = OptionalFlagSet.String(
	string(flagNames.ValidatorConfigFile), globals.VALIDATOR_CONFIG_PATH,
	"File path for configuration of validators.")
var commitFile = OptionalFlagSet.String(string(flagNames.CommitFile), globals.COMMIT_MESSAGE_PATH, "File path for commit file. The special value '-' indicates stdin, allowing piping.")

var comSep = commitSep("")

func init() {
	BooleanFlagSet.Var(&comSep, string(flagNames.CommitSeparator),
		"Separator when validating multiple commits, `string`-valued. Changes the program to assume that the input will have messages separated by this separator, and process all the messages.")

	FlagSet = flagUtils.Combine(OptionalFlagSet, BooleanFlagSet)
}

// validOutputFormats holds the valid values for the [outputFormat] flag.
var validOutputFormats = struct {
	Human, JSON string
}{
	Human: "human",
	JSON:  "json",
}

// see [outputFormat].
var validOutputFormatsSlice = []string{"human", "json"}
var outputFormat = OptionalFlagSet.String(string(flagNames.OutputFormat), "human", fmt.Sprintf("Output format of validation, %v.", strings.Join(validOutputFormatsSlice, "|")))

// args handles the arguments passed to the program.
type args struct {
	ValidatorConfigFile, CommitFile, OutputFormat string
	CommitSeparator                               *string
}

func NewArgs() (arg *args, e error) {
	arg = &args{}

	arg.ValidatorConfigFile = *validatorConfigFile
	arg.CommitFile = *commitFile
	arg.OutputFormat = *outputFormat
	if comSep != "" {
		arg.CommitSeparator = (*string)(&comSep)
	}
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

type commitSep string

func (sep *commitSep) String() string {
	if sep != nil {
		return fmt.Sprintf("%s", string(*sep))
	}
	return ""
}

func (sep *commitSep) Set(flg string) (e error) {
	if flg == "" || flg == "true" {
		return errors.New("separator should not be empty string")
	}
	*sep = commitSep(flg)
	return
}

func (sep *commitSep) IsBoolFlag() bool {
	return false
}
