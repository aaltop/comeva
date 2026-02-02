// package config contains utilities to read in configuration for comeva from
// a file.
package config

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

const BASE_PATH string = "./.comeva/"

// Config describes configuration parameters for the program.
type Config struct {
	// The file used to read in the configuration for the validator.
	ValidatorConfigFile string
	// The file that is validated.
	CommitFile string
	// The verbosity of the program, higher is more verbose. Default 0.
	Verbosity int
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
		BASE_PATH+"validator.yaml",
		BASE_PATH+"git_commit.txt",
		0)
	if e != nil {
		panic(e)
	}
	return config
}

func NewConfig(
	validatorConfigFile, commitFile string,
	verbosity int) (config *Config, e error) {
	config = NewDefaultConfig()
	config.ValidatorConfigFile = validatorConfigFile
	config.CommitFile = commitFile
	config.Verbosity = verbosity
	return
}

type configYaml struct {
	ValidatorConfigFile string `yaml:"validatorConfigFile"`
	CommitFile          string `yaml:"commitFile"`
	Verbosity           int    `yaml:"verbosity"`
}

func (config *Config) UnmarshalYAML(data []byte) (e error) {
	var temp configYaml
	fmt.Println(string(data))
	if e = yaml.Unmarshal(data, &temp); e == nil {
		var new *Config
		fmt.Println(temp)
		new, e = NewConfig(temp.ValidatorConfigFile, temp.CommitFile, temp.Verbosity)
		*config = *new
	}
	return
}
