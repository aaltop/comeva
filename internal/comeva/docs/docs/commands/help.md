# Help command

The help command is used to access the documentation of comeva from the command
line. The documentation is mapped as its directory structure to the help command.
Documentation consists of Markdown files.

Without arguments, the documentation structure is shown. Using the path argument,
a documentation path can be specified: directory contents will be shown when the passed
argument ends with a slash, and file contents when the argument does *not* end with
a slash, assuming that the path matches the respective type of file system object.
Some examples of paths:

- `/`
    - The root of the docs, shows the whole documentation directory structure.

- `/commands/`
    - Shows the /commands/ directory.

- `/commands/help`
    - shows this documentation.

The documentation directory can also be recreated locally using the 'create-docs'
argument.