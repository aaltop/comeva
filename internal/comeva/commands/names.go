package commands

// CommandNames contains the names of all the commands that are
// available to execute.
var CommandNames = struct {
	Help, Validate string
}{
	Help:     "help",
	Validate: "validate",
}
