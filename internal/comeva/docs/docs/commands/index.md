# Commands

Commands generally have the following structure:

```
comeva [global flags] [<subcommand1> [<subcommand2>...] [command-specifc flags]]
```

- comeva is the "base" command. Global flags come after the base command and are
followed by "sub-commands" which are followed by command-specific flags.

- Flags are considered the only arguments to the commands. This means anything
passed without a flag to a command is considered to be a "sub-command".

- Global flags are flags that *could* be globally used, though they might not
always have an effect on the executed command. Command-specific flags are the
arguments to the last specified sub-command.

## Flag types

A distinction is made between different types of flags. Two types of flags (outside
the global/command-specific) exist: *boolean*[^1] and *non-boolean*. This distinction is
based on whether the flag is specified as a boolean flag, a flag that does not
require a value to be specified when the flag is passed. Note specifically that
a boolean flag does not necessarily have a boolean value; rather, the flag
**can be passed without a value** in which case a default
value is used. A further distinction[^2] is that it is an
optional flag whose related functionality takes effect **only when the flag is
passed**. In more simple terms, a *boolean* flag is a toggle-with-a-default
flag. A *non-boolean* flag is one that **must have a value specified** when it is
passed, but whose presence or lack thereof **may or may not change the behaviour**
of the command it is passed to.

Whether a flag is boolean or not is marked in help messages. When using a boolean
flag and wanting to specify a value for it, it is critical to specify the value
as
```
--boolean-flag=value
```
rather than
```
--boolean-flag value
```
The former format works for all flags regardless of whether they're boolean,
so can be a safe format to use if not wanting to remember or check which flags
are boolean.

[^1]: Use of the term 'boolean' comes from how boolean flags are usually parsed:
they can be passed without a value, and their existence is considered to
indicate a 'true' for the flag.

[^2]: This distinction arises because otherwise, such a flag would be identical
in functionality to *a flag with a default*, a flag whose value is
**always used** but which has a default in case it is not passed. The only
difference would be that the boolean flag would have to always be passed
regardless of whether its value is specified (or alternatively, the boolean
functionality of the flag would be pointless).

## General behaviour

### Exit states

The program utilises four exit states:

- successful
    - exit code 0. The program ran successfully.

- program error
    - exit code 1. The program encountered an expected error (e.g. a passed flag
    was found to be invalid in syntax and could not be parsed) that lead to
    program termination.

- validation error
    - exit code 2. The program validated a commit message and ran succesfully,
    but validation errors were encountered.

- uncaught error
    - exit code 3. The program encountered an unexpected error that could not
    be handled.

### Argument priority

Arguments can be passed from a number of places. The arguments are parsed in the
following order (in descending order of priority):

1. Command line arguments
2. Local configuration (e.g. `./.comeva/*`)
3. Global configuration (e.g. `/home/<user>/.config/comeva/*` on Linux)

See [the init command](../commands/init.md) for how to create the configuration setup
in the local and global locations. Not all arguments may be passable as command line
arguments, nor all as configuration (i.e. the set of command line arguments is not
necessarily a subset of the configuration file arguments nor vice versa).

## Parsing behaviour motivation

Parsing the command line in the way that it is has some benefits, some from a development
perspective and some from the user perspective:

- Having global flags at the start means that they're easy to parse as they're
not mixed in with command-specific flags.
    - Having global flags in general is useful as they don't have to be specified
    for each command separately. 

- Having non-flag arguments only be commands rather than any positional
arguments to commands makes it easy to determine whether the passed-in
hierarchy of values does actually point to a command: if the hierarchy doesn't
match a sub-command, there's no need to separately check whether we stopped
specifying commands at some point and started specifying positional arguments
to those commands, and at what point that could have happened. However, this assumes
that a non-leaf sub-command could potentially be invoked. If this wasn't
the case, it would only be necessary to check whether a leaf sub-command has
been specified.
    - Having commands come after the global flags and before command-specific flags
    again makes parsing fairly obvious.

The above considerations naturally lead to having to specify any command-specific
arguments as flags. While this may seem cumbersome, it leads to more obvious
command invocations where it is easier to tell what the command is doing. It
also makes it possible to specify arguments more optionally and without a
particular order. It may also help if a breaking change to a command's
invocation is introduced: if a new required argument is added, it is easy to tell
when it is not passed, and it is easier to create a descriptive error message.
In essence, the difference is the same as using positional arguments versus
using keyword arguments.

Another thing considering command-specific arguments is that there are no
shorthand flags. For one, for the developer, this prevents having to specify
shorthands in addition to the existing full argument and dealing with both.
Second, again having more explicit keyword arguments makes invocations more clear.