package validate

import (
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	"flag"
)

var FlagSet = flag.NewFlagSet("", flag.ContinueOnError)

type FlagString string

// flagNames holds the string name of each command line flag.
var flagNames = struct {
	ValidatorConfigFile,
	CommitFile FlagString
}{
	ValidatorConfigFile: "validator-config-file",
	CommitFile:          "commit-file",
}
var validatorConfigFile = FlagSet.String(
	string(flagNames.ValidatorConfigFile), globals.CONFIG_BASE_PATH+"validator.yaml",
	"file path for configuration of validators")
var commitFile = FlagSet.String(string(flagNames.CommitFile), "./git_commit.txt", "file path for commit file")

// args handles the arguments passed to the program.
type args struct {
	ConfigFile, ValidatorConfigFile, CommitFile string
}

func NewArgs() (arg *args, e error) {
	arg = &args{}

	arg.ValidatorConfigFile = *validatorConfigFile
	arg.CommitFile = *commitFile

	return
}

// argsFromConfig returns [args] based on the contents of the passed [config.Config].
func argsFromConfig(conf *config.Config) (args *args, e error) {
	args, e = NewArgs()
	if conf.CommitFile != nil {
		args.CommitFile = *conf.CommitFile
	}

	if conf.ValidatorConfigFile != nil {
		args.ValidatorConfigFile = *conf.ValidatorConfigFile
	}
	return
}
