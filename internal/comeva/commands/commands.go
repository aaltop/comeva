// command contains a type that represents a command line command.

package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"comeva/internal/comeva/args"
	"comeva/internal/comeva/config"
	"comeva/internal/comeva/globals"
	"comeva/internal/comeva/io/ansi"
	exitstate "comeva/internal/exitState"
	"comeva/internal/yaml"
	ioUtils "comeva/utils/io"
)

// UsageExample describes a single usage example of a help message.
type UsageExample struct {
	// Invocation is the command line command example.
	Invocation string
	// Description describes the usage example in more detail.
	Description string
}

// HelpMessageUsage describes the usage examples of a help message.
type HelpMessageUsage []UsageExample

// AddExample adds another example. `commandSignature`
// should be the way the associated command is called without
// any common parts: e.g. in a help message "comeva [global flags] help [item]",
// the `commandSignature` is "help [item]", i.e. only the subcommand(s) followed
// by anything they might be followed by (arguments/flags...). `description`
// gives a description of the usage example in more detail.
func (hlpMsgUsg *HelpMessageUsage) AddExample(
	commandSignature, description string,
) {
	*hlpMsgUsg = append(
		*hlpMsgUsg,
		UsageExample{
			Invocation:  fmt.Sprintf("comeva [global flags] %s", commandSignature),
			Description: description,
		})
}

func (usage HelpMessageUsage) String() string {

	// TODO: read Description as lines, indent each line as necessary

	var usageMessage strings.Builder
	for _, us := range usage {
		fmt.Fprintf(
			&usageMessage,
			"%v\n%v\n",
			ioUtils.IndentLines(us.Invocation, 4), ioUtils.IndentLines(us.Description, 8),
		)
	}

	return usageMessage.String()
}

// FlagGroup represents a grouped set of flags, flags of a command that relate
// to specific functionality.
type FlagGroup struct {
	// name for this set of flag options.
	name string
	// The description of the flags, e.g. the output of [flag.PrintDefaults] (flag
	// names, flag defaults, flag descriptions).
	description string
}

func newDefaultFlagGroup() (flgGroup *FlagGroup) {
	return &FlagGroup{}
}

func NewFlagGroup(name, description string) (flgGroup *FlagGroup, e error) {
	flgGroup = newDefaultFlagGroup()
	flgGroup.name = name
	flgGroup.description = description
	return
}

// FlagOptions describes a set of grouped flag options.
type FlagOptions []FlagGroup

// newDefaultFlagOptions returns a [FlagOptions] with global flags already
// included.
func newDefaultFlagOptions() (flgOptions FlagOptions) {
	flgOptions = make(FlagOptions, 0)
	flgOptions.AddGroup("Global", args.GlobalFlagDefaults())
	return
}

// AddGroup adds new [FlagGroup]s to the options.
func (flagOpt *FlagOptions) AddGroup(name, description string) (e error) {
	var group *FlagGroup
	group, e = NewFlagGroup(name, description)
	if e == nil {
		*flagOpt = append(*flagOpt, *group)
	}
	return
}

func (flagOpt FlagOptions) String() string {

	var buffer strings.Builder
	for _, group := range flagOpt {
		fmt.Fprintf(
			&buffer,
			"%v:\n%v\n", ioUtils.IndentLines(group.name, 0), ioUtils.IndentLines(group.description, 0),
		)
	}

	return buffer.String()
}

type HelpMessage struct {
	// synopsis is a short description of the command.
	synopsis string
	// description is a longer description of the command.
	description string
	// usage gives the possible ways to call the command.
	usage HelpMessageUsage
	// options enumerates the possible flags usable with the command.
	options FlagOptions
}

func NewDefaultHelpMessage() (msg *HelpMessage) {
	return &HelpMessage{}
}

// NewHelpMessage creates a new [HelpMessage].
//   - synopsis is a short description of the command.
//   - description is a longer description of the command.
//   - usage gives the possible ways to call the command.
//   - options enumerates the possible flags usable with the command.
func NewHelpMessage(
	synopsis string, description string, usage HelpMessageUsage, options FlagOptions,
) (msg *HelpMessage, e error) {
	msg = NewDefaultHelpMessage()

	msg.synopsis = synopsis
	msg.description = description
	msg.usage = usage
	msg.options = options

	return
}

func (msg HelpMessage) String() (msgString string) {
	return fmt.Sprintf(`
%v

%v

Usage:
%v
Options:

%v
`, msg.synopsis, msg.description, msg.usage, msg.options)
}

type CommandFunc func(gFlags *args.GlobalFlags, passedGFlags map[string]bool, conf *config.Config) (extState *exitstate.ExitState)
type CommandList map[string]*Command

// Command represents a command line command.
type Command struct {
	// HelpMessage represents the message to be printed e.g. when
	// the help flag is used for this command.
	HelpMessage HelpMessage

	// CommandFlags defines the flags specific to this command.
	CommandFlags *flag.FlagSet
	// [Command.function] executes the main behaviour of the command. It will be called by
	// [Command.Execute], which should be used to actually execute the command.
	// Inside [Command.function], [Command.CommandFlags] will have been parsed,
	// so the flagset's values are usable.
	function CommandFunc
	// subCommands contains the subCommands of this command, if any.
	subCommands CommandList
}

func NewDefaultCommand() (com *Command) {
	com = &Command{}
	com.function = func(gFlags *args.GlobalFlags, passedGFlags map[string]bool, conf *config.Config) (e *exitstate.ExitState) {
		return
	}
	com.subCommands = make(CommandList)

	return
}

// NewCommand returns a new CLI command.
func NewCommand(
	helpMessage HelpMessage,
	function CommandFunc,
	subcommands CommandList,
) (com *Command, e error) {
	com = NewDefaultCommand()

	// No flags because it's optional, there might not be any.
	com.HelpMessage = helpMessage
	com.function = function
	com.subCommands = subcommands

	return
}

func (com *Command) Execute(cmdArgs *args.CommandArgs) (extState *exitstate.ExitState) {
	extState = exitstate.NewDefaultExitState()
	var e error

	var colorSchemes = ansi.BasicColorSchemes

	if com.CommandFlags != nil {
		com.CommandFlags.Parse(cmdArgs.Flags)
	}

	if cmdArgs.GlobalFlags.Help {
		println(com.HelpMessage.String())
		return
	}

	var conf = config.NewDefaultConfig()
	if cmdArgs.GlobalFlags.ConfigFile != "" {
		e = yaml.UnMarshalFromFile(cmdArgs.GlobalFlags.ConfigFile, conf)
		if e != nil && cmdArgs.GlobalFlags.Verbosity >= globals.VERBOSITY_WARNING {
			fmt.Fprint(os.Stderr, colorSchemes.Warning.ApplyFore("Warning: error reading config file: %v", e))
		}
	}

	return com.function(cmdArgs.GlobalFlags, cmdArgs.PassedGlobalFlags, conf)
}

// GetSubCommand returns a subcommand of this command and further
// sub-commands therein based on the [Command.subCommands] field of `com`
// and the command string equivalents passed in `commands`.
//
// For example, if the current [Command] has [Command.subCommands] contain
// a value behind a key 'command1', and that value (another [Command])
// contains a value behind a key 'command2', GetSubCommand finds the
// deepest subcommand given `commands` as [command1, command2].
//
// If no command is found in the chain based on `commands`,
// the returned value will be nil.
func (com *Command) GetSubCommand(commands []string) (subCom *Command) {
	if len(commands) < 1 {
		return nil
	}
	var ok bool
	subCom, ok = com.subCommands[commands[0]]
	if !ok {
		return nil
	}

	// found the last command
	if len(commands) == 1 {
		return
	}

	subCom = subCom.GetSubCommand(commands[1:])
	return
}
