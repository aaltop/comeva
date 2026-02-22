package config

import (
	"fmt"
	"testing"
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
		t.Errorf("Expected empty default verbosity value, got %v", config.CommitFile)
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
