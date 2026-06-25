package flag

import (
	baseFlag "flag"
	"strings"
)

// GetDefaults returns the defaults as printed by [flag.FlagSet.PrintDefaults].
func GetDefaults(flagSet *baseFlag.FlagSet) (usage string) {
	var buffer = &strings.Builder{}
	var oldOutput = flagSet.Output()
	defer flagSet.SetOutput(oldOutput)

	flagSet.SetOutput(buffer)
	flagSet.PrintDefaults()
	return buffer.String()
}

// PassedFlags returns a map that reports whether a given flag was set on the
// command line.
func PassedFlags(flagSet ...*baseFlag.FlagSet) (passed map[string]bool) {

	passed = make(map[string]bool)
	for _, set := range flagSet {
		var setFlags = make(map[string]bool)
		// Find all the flags that were set
		set.Visit(func(f *baseFlag.Flag) { setFlags[f.Name] = true })
		// Determine whether a flag was set
		set.VisitAll(func(f *baseFlag.Flag) {
			var _, isSet = setFlags[f.Name]
			passed[f.Name] = isSet
		})
	}

	return
}

// Combine returns a FlagSet whose flags are a combination of the individual
// Flagset's flags.
func Combine(individual ...*baseFlag.FlagSet) (combined *baseFlag.FlagSet) {
	combined = baseFlag.NewFlagSet("", baseFlag.ContinueOnError)
	for _, set := range individual {
		set.VisitAll(func(f *baseFlag.Flag) {
			combined.Var(f.Value, f.Name, f.Usage)
		})
	}
	return
}
