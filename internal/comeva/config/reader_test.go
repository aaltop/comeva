package config

import (
	"fmt"
	"testing"

	testingHelpers "github.com/aaltop/comeva/internal/testing"
)

func FixtureConfig() (config *Config) {
	return NewDefaultConfig()
}

// Validator config file can be unmarshaled from YAML content.
func TestUnmarshalValidatorConfigFile(t *testing.T) {
	var config = NewDefaultConfig()

	if config.ValidatorConfigFile != nil {
		t.Errorf("Expected empty default config file value, got %v", config.ValidatorConfigFile)
	}

	var validatorConfigFile = "./.comeva/validator_config.yaml"
	var yamlText = `
validatorConfigFile: %s`

	yamlText = fmt.Sprintf(yamlText, validatorConfigFile)

	if e := config.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Valid config file was found to be invalid: %v\n", e)
	}

	if *config.ValidatorConfigFile != validatorConfigFile {
		t.Errorf("validator config file did not match,\nExpected:\n%s\nReceived:\n%s\n", validatorConfigFile, *config.ValidatorConfigFile)
	}
}

// Commit file can be unmarshaled from YAML content.
func TestUnmarshalCommitFile(t *testing.T) {
	var config = NewDefaultConfig()

	if config.CommitFile != nil {
		t.Errorf("Expected empty default commit file value, got %v", config.CommitFile)
	}

	var commitFile = "./git_commit.txt"
	var yamlText = `
commitFile: %s`

	yamlText = fmt.Sprintf(yamlText, commitFile)

	if e := config.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Valid config file was found to be invalid: %v\n", e)
	}

	if *config.CommitFile != commitFile {
		t.Errorf("validator commit file did not match,\nExpected:\n%s\nReceived:\n%s\n", commitFile, *config.CommitFile)
	}
}

// Verbosity can be unmarshaled from YAML content.
func TestUnmarshalVerbosity(t *testing.T) {
	var config = NewDefaultConfig()

	if config.Verbosity != nil {
		t.Errorf("Expected empty default verbosity value, got %v", config.Verbosity)
	}

	var verbosity = 11
	var yamlText = `
verbosity: %d`

	yamlText = fmt.Sprintf(yamlText, verbosity)

	if e := config.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Valid config file was found to be invalid: %v\n", e)
	}

	if *config.Verbosity != verbosity {
		t.Errorf("validator verbosity did not match,\nExpected:\n%d\nReceived:\n%d\n", verbosity, config.Verbosity)
	}
}

func TestUnmarshalLoggingLevel(t *testing.T) {
	var config = NewDefaultConfig()

	if config.LoggingLevel != nil {
		t.Errorf("Expected empty default logging level value, got %v", config.LoggingLevel)
	}

	var loggingLevel = 11
	var yamlText = `
loggingLevel: %d`

	yamlText = fmt.Sprintf(yamlText, loggingLevel)

	if e := config.UnmarshalYAML([]byte(yamlText)); e != nil {
		t.Errorf("Valid config file was found to be invalid: %v\n", e)
	}

	if *config.LoggingLevel != loggingLevel {
		t.Errorf("validator logging level did not match,\nExpected:\n%d\nReceived:\n%d\n", loggingLevel, *config.LoggingLevel)
	}
}

func TestEqual(t *testing.T) {
	var conf, equal = FixtureConfig(), FixtureConfig()
	var notEqual, e = NewConfig("notliekly", "meniether", -999, -543)
	if e != nil {
		t.Fatalf("Unexpected error: %v\n", e)
	}

	if !conf.Equal(equal) {
		t.Errorf("Equal Configs found to be inequal: %v\n%v", conf, equal)
	}

	if conf.Equal(notEqual) {
		t.Errorf("Inequal Configs found to be equal: %v\n%v", conf, notEqual)
	}
}

// Content can be marshaled and unmarshaled into the same state.
func TestMarshalUnmarshal(t *testing.T) {
	var config, e = NewConfig("valida", "commi", 123, 321)
	if e != nil {
		t.Fatalf("Unexpected error: %v\n", e)
	}

	var data []byte
	data, e = config.MarshalYAML()
	if e != nil {
		t.Fatalf("Unexpected error: %v\n", e)
	}
	var newConfig = NewDefaultConfig()
	e = newConfig.UnmarshalYAML(data)
	if e != nil {
		t.Fatalf("Unexpected error: %v\n", e)
	}

	if !config.Equal(newConfig) {
		t.Error(testingHelpers.ValueMismatch("Config", *config, *newConfig))
	}
}
