package internal

import (
	"github.com/aaltop/comeva/internal/comeva/args"
	"github.com/aaltop/comeva/internal/comeva/config"
)

// JoinGlobalArguments joins the global flags and configuration file values.
func JoinGlobalArguments(
	globalArgs *args.GlobalFlags,
	conf *config.Config,
	passedGlobal *args.PassedGlobalFlags,
	passedConfig *config.PassedConfigArgs,
) (joined *args.GlobalFlags) {
	joined = &args.GlobalFlags{}
	joined.Help = globalArgs.Help
	joined.Verbosity = globalArgs.Verbosity
	joined.ConfigFile = globalArgs.ConfigFile
	joined.LoggingLevel = globalArgs.LoggingLevel

	if !passedGlobal.Verbosity && passedConfig.Verbosity {
		joined.Verbosity = *conf.Verbosity
	}

	if !passedGlobal.LoggingLevel && passedConfig.LoggingLevel {
		joined.LoggingLevel = *conf.LoggingLevel
	}

	return
}
