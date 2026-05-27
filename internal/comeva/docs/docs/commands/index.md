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

Parsing the command line like this has some benefits, some from a development
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
to those commands, and at what point that could have happened. Obviously the idea
here is that a non-leaf sub-command could potentially be invoked. If this wasn't
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

## General behaviour

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