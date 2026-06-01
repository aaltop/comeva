// package config contains utilities to read in configuration for comeva from
// a file.
package config

import (
	"comeva/internal/comeva/globals"
	"comeva/internal/equal"

	"comeva/internal/yaml"

	goccyYaml "github.com/goccy/go-yaml"
)

// Config describes configuration parameters for the program.
type Config struct {
	// The file used to read in the configuration for the validator.
	ValidatorConfigFile *string
	// The file that is validated.
	CommitFile *string
	// The verbosity of the program, higher is more verbose. Default 0.
	Verbosity *int
	// The logging level.
	LoggingLevel *int
}

// NewDefaultConfig creates the base Config.
func NewDefaultConfig() (config *Config) {
	config = &Config{}
	return
}

// NewConfigWithDefaults creates a Config with default values.
func NewConfigWithDefaults() (config *Config) {
	var e error
	config, e = NewConfig(
		globals.VALIDATOR_CONFIG_PATH,
		globals.COMMIT_MESSAGE_PATH,
		globals.VERBOSITY_DEFAULT,
		globals.LOGGING_LEVEL_DEFAULT,
	)
	if e != nil {
		panic(e)
	}
	return config
}

func NewConfig(
	validatorConfigFile, commitFile string,
	verbosity int, loggingLevel int,
) (config *Config, e error) {
	config = NewDefaultConfig()
	config.ValidatorConfigFile = &validatorConfigFile
	config.CommitFile = &commitFile
	config.Verbosity = &verbosity
	config.LoggingLevel = &loggingLevel
	return
}

// Equal reports whether the two Configs are equal.
func (config *Config) Equal(other *Config) bool {
	if config == nil || other == nil {
		return config == other
	}

	return equal.NilEqual(config.ValidatorConfigFile, other.ValidatorConfigFile) &&
		equal.NilEqual(config.CommitFile, other.CommitFile) &&
		equal.NilEqual(config.Verbosity, other.Verbosity) &&
		equal.NilEqual(config.LoggingLevel, other.LoggingLevel)
}

type configYaml struct {
	ValidatorConfigFile *string `yaml:"validatorConfigFile"`
	CommitFile          *string `yaml:"commitFile"`
	Verbosity           *int    `yaml:"verbosity"`
	LoggingLevel        *int    `yaml:"loggingLevel"`
}

func CreateCommentMap(prefix string) goccyYaml.CommentMap {
	return yaml.CreateCommentMap(prefix, yaml.SuffixCommentMap{
		"validatorConfigFile": {" path to the config of the validator"},
		"commitFile":          {" path to the commit message file"},
		"verbosity":           {" default verbosity to use"},
		"loggingLevel":        {" default logging level to use"},
	})
}

func (config *Config) UnmarshalYAML(data []byte) (e error) {
	var temp configYaml
	if e = goccyYaml.Unmarshal(data, &temp); e == nil {
		config.CommitFile = temp.CommitFile
		config.ValidatorConfigFile = temp.ValidatorConfigFile
		config.Verbosity = temp.Verbosity
		config.LoggingLevel = temp.LoggingLevel
	}
	return
}

var comments = CreateCommentMap("")

func (config *Config) MarshalYAML() (data []byte, e error) {
	var temp configYaml
	temp.CommitFile = config.CommitFile
	temp.ValidatorConfigFile = config.ValidatorConfigFile
	temp.Verbosity = config.Verbosity
	temp.LoggingLevel = config.LoggingLevel

	return goccyYaml.MarshalWithOptions(&temp, goccyYaml.WithComment(comments))
}

// UnmarshalYAMLFile unmarshals the [Config] from the file.
func (config *Config) UnmarshalYAMLFile(filename string) (e error) {
	e = yaml.UnMarshalFromFile(filename, config)
	return
}
