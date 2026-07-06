package validate

import (
	argus "github.com/aaltop/comeva/internal/comeva/args"
	"github.com/aaltop/comeva/internal/comeva/config"
	"github.com/aaltop/comeva/internal/comeva/globals"
	"github.com/aaltop/comeva/internal/flag"
)

var debugLogger = globals.DebugLogger

type passedLocalArgs struct {
	ValidatorConfigFile, CommitFile, OutputFormat, CommitSeparator bool
}

// using this to retain the import of the local 'flag' package, so the comment
// references actually refer to the correct stuff
var _ = flag.Combine

// getPassedLocalArgs reports which of the local arguments were passed. see
// [flag.PassedFlags] for `passedMap`.
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
		joined.CommitFile = *conf.CommitFile
	}

	if !passedArgs.Local.ValidatorConfigFile && passedArgs.Config.ValidatorConfigFile {
		joined.ValidatorConfigFile = *conf.ValidatorConfigFile
	}

	return
}
