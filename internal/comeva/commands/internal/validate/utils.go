package validate

import (
	argus "comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
)

var debugLogger = globals.DebugLogger

type passedLocalArgs struct {
	ValidatorConfigFile, CommitFile, OutputFormat, CommitSeparator bool
}

func getPassedLocalArgs(passedMap map[string]bool) (passed *passedLocalArgs) {
	passed = &passedLocalArgs{}
	passed.ValidatorConfigFile = passedMap[string(flagNames.ValidatorConfigFile)]
	passed.CommitFile = passedMap[string(flagNames.CommitFile)]
	passed.OutputFormat = passedMap[string(flagNames.OutputFormat)]
	passed.CommitSeparator = passedMap[string(flagNames.CommitSeparator)]
	return
}

// passedArgs reports for each argument source whether arguments were passed
// there.
type passedArgs struct {
	Local  passedLocalArgs
	Config config.PassedConfigArgs
	Global argus.PassedGlobalFlags
}

func joinLocalArguments(
	arg *args,
	conf *config.Config,
	passedArgs *passedArgs,
) (joined *args) {

	joined = &args{}
	joined.CommitFile = arg.CommitFile
	joined.OutputFormat = arg.OutputFormat
	joined.ValidatorConfigFile = arg.ValidatorConfigFile
	joined.CommitSeparator = arg.CommitSeparator

	if !passedArgs.Local.CommitFile && passedArgs.Config.CommitFile {
		arg.CommitFile = *conf.CommitFile
	}

	if !passedArgs.Local.ValidatorConfigFile && passedArgs.Config.ValidatorConfigFile {
		arg.ValidatorConfigFile = *conf.ValidatorConfigFile
	}

	return
}
