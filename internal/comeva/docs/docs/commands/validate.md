# Validate command

The validate command validates commit messages. It can use default values for
validation configuration, but it is suggested that a configuration file for
the validator be created with the help of [the init command](./init.md).

The validate command will scan the specified commit file and report any problems
based on the validator configuration. It will return with [an exit code of 2](index.md#general-behaviour)
if it encounters a validation issue. If no validation issue is encountered and
no error occurs during program execution, the exit code will be 0.

Multiple commit messages may be validated at once using the 'commit-separator' argument.
Passing this value changes the program to attempt processing multiple commits from the
passed in file, assuming that these commits are separated in the file by the passed
separator value. Particularly useful for parsing the entire git log history for further
processing, for example:
```sh
git log --format=format:"%B-|-" | comeva validate --commit-file "-" --output-format json --commit-separator "-|-" > comeva_output.txt
```

## Output format

The output format can be controlled using the 'output-format' argument, which
accepts one of two values:

- human
    - A human-friendly format. Prints the validation problems on the command
    line in a neat order. With some verbosity settings, shows the validated
    message as well.

- json
    - JSON output.

The JSON structure is the following:

- <object root\>
    - Content
        - Header
            - Content
                - Type (`string`)
                - Scope (`string`)
                - Breaking (`boolean`)
                - Description
                    - Verb (`string`)
                    - Content (`string`)
            
            - Errors

        - Body
            - Content (`[]string`)
            - Errors

        - Trailer
            - Content (`[]`)
                - Key (`string`)
                - Value (`string`)

            - Errors

    - Errors

See the respective files for each [validator](../validation/). The **Errors**
are of the following format:

- Errors (`[]`)
    - MessagePart (`string`, "Header"|"Body"|"Trailer"|"")
    - Line (`uint`)
    - Cols (`[2]int`)
    - Reason (`string`)

**MessagePart** tells in which part of a message the problem was considered to be,
**Line** tells on which line in the message the problem was found (0 means no line
specified), **Cols** tells on which columns the error was (when `Cols[0]` and `Cols[1]` are
equal, the columns are unspecified), and **Reason** describes the problem in a free-form way.