# Add footer validation

Basically parsing the git trailer -style keys. Need to have a set of
expected keys. Maybe optional set too, where they're not necessary to
have, but if there are keys, they should either be in the expected or
optional set. Empty sets should mean that

1. no expected keys (for expected set)
2. any key is allowed (for optional set)

# Evaluate sensibility of having the io.Reader Validate method.

It's not very
elegant, because it just processes the content in whole, assuming that
the only content is the expected content. Therefore, anyone using this
would have to have a reader that only has the specific content that is
valid (or assumed to be), which is hardly useful. Generally,
the scanner is much more useful anyway, as it processes in a more
predictable way, and could potentially be passed between the three validators, header,
body, and trailer. With the reader, content needs to be processed byte-by-byte,
which just doesn't match how how the commit message would be processed, line-by-line
like the scanner does it. On this point, however, it's also important to make
sure that a _line_ scanner is passed, not some other scanner.

# Combine the validators in a message validator

## Have all the sub-validators accept a CountingScanner for ease of keeping track of line

Passing around the same scanner (ideally) should make it easier to give information
about the line the error is on. The header validator also needs to be changed
to support the counter.

## Require BREAKING-CHANGE trailer and "!" always together

Not quite like this according to the conventional commits
spec, but I'd say the "!" in the header allows for quickly
verifying which commits introduce breaking changes,
while the actual information about the breaking change should
be in trailer.

Why not in the header or body? It might not be
the main thing changing: I mean, I don't know much about professional
development, but you make some change, which on the side happens
to require a breaking change. The main thing in the header should be about
what changed, and if that isn't the breaking change, then what?
Granted, you could just separately make the
breaking change and commit that alone, but this can be a little fiddly
and in my eyes quite unnecessary. Describing the breaking change
in the body also isn't very useful, as it is much more difficult
to then find this information for both a human and a computer.
The very obvious BREAKING-CHANGE trailer is both easy to find
for a person looking with their peepers, and for a computer
parsing the trailers. On this point, it also seems unnecessary
to have the "BREAKING CHANGE" option, as this is not parseable as
a git trailer, so the trailer version, BREAKING-CHANGE, is expected.

# Add config file reading

# Include optional info on type/scope/verb/trailer key

For "feat" type, something like "A feature was added, or changed in such a
way that it affects the end user", and "fix" as "A bugfix was made", etc.
This way, it's possible to print out all the options, and have some info
about each to give an idea of which would be suitable to use for the current
commit.

# ? Add a way to produce/compute stuff based on the commit contents

Stuff like the SemVer, and something about the keys, potentially.

One useful thing is to keep track of the version in the config file
or even somewhere automated, then require it in the trailers. Then,
depending on the type of commit, the increase should be calculated
and compared against the one in the trailer: the trailer is set by
the commit's writer, and should match the current version, while
the program bumps the version in the config file or other and compares
that against the one set by the user in the trailers. Maybe not too
useful with multiple developers, but then again, still useful in some
cases, and the idea of the validator is that it does not force anything,
it just reports potential problems.
