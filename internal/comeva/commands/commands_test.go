package commands

import (
	"testing"

	testingHelpers "comeva/internal/comeva/testing"
)

// Getting the subcommands of commands works.
func TestGetSubCommand(t *testing.T) {
	// TODO: could test each individually using t.Run?

	var top, sub1, sub2 *Command
	top = NewDefaultCommand()
	sub1, sub2 = NewDefaultCommand(), NewDefaultCommand()
	var commands = []string{"sub1", "sub2"}
	sub1.subCommands[commands[1]] = sub2
	top.subCommands[commands[0]] = sub1
	top.subCommands["notFound"] = NewDefaultCommand()

	var received *Command
	// find the first command in the hierarchy
	received = top.GetSubCommand(commands[0:1])
	if received == nil {
		t.Error("received should not be nil")
	} else if received != sub1 {
		testingHelpers.ValueMismatch("first subcommand", *sub1, *received)
	}

	// find the second command in the hierarchy
	received = top.GetSubCommand(commands)
	if received == nil {
		t.Error("received should not be nil")
	} else if received != sub2 {
		testingHelpers.ValueMismatch("second subcommand", sub2, received)
	}

	// passing empty slice returns nothing
	received = top.GetSubCommand(commands[0:0])
	if received != nil {
		t.Errorf("Expected nil, got %v", *received)
	}

}
