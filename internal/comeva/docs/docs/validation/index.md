# Validation

A git commit message is considered to be made up of three parts: the header,
the body, and the trailer. [The header](./header.md) should be at the first
line and span that one line. [The body](./body.md) comes after the header with
a single empty line in between and can be any number of lines long, possibly
being empty. [The trailer](./trailer.md) comes after the body with a single empty
line in between, and is considered to start based on the rules of trailers. The
trailer section may also be empty, depending on the configuration.

Generally, validation of commit messages is based on [Conventional Commits][1].
The intention is that valid messages are valid Conventional Commits, but not
all valid Conventional Commits are necessarily valid messages: validation is more stringent
in some places, such as requiring both a breaking change exclamation mark AND
a breaking change trailer key. In contrast, Conventional Commits allows omitting one
in the case of a breaking change, but does not prohibit both being specified at
the same time.

[1]: https://www.conventionalcommits.org/en/v1.0.0/ 'Conventional Commits'