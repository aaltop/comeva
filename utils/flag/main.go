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
func PassedFlags(flagSet *baseFlag.FlagSet) (passed map[string]bool) {
	passed = make(map[string]bool)
	var setFlags = make(map[string]bool)
	// Find all the flags that were set
	flagSet.Visit(func(f *baseFlag.Flag) { setFlags[f.Name] = true })
	// Determine whether a flag was set
	flagSet.VisitAll(func(f *baseFlag.Flag) {
		var _, isSet = setFlags[f.Name]
		passed[f.Name] = isSet
	})
	return
}
