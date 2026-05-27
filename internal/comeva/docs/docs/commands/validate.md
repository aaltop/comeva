# Validate command

The validate command validates commit messages. It can use default values for
validation configuration, but it is suggested that a configuration file for
the validator be created with the help of [the init command](./init.md).

The validate command will scan the specified commit file and report any problems
based on the validator configuration. It will return with [an exit code of 2](index.md#general-behaviour)
if it encounters a validation issue. If no validation issue is encountered and
no error occurs during program execution, the exit code will be 0.