# Init command

The init command initialises a configuration directory, '.comeva', in the
directory where it is called. In this directory, configuration files with
defaults are created:

- config.yaml
    - Configuration for specifying certain default arguments.

- validator.yaml
    - Configuration for the validation of commit messages.

If instead using the 'global' argument, the same configuration files are initialised
in the config directory of the user. The config directory is as returned by [UserConfigDir](https://pkg.go.dev/os@go1.26.4#UserConfigDir).