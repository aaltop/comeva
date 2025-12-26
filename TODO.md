# Evaluate sensibility of having the io.Reader Validate method.

It's not very
elegant, because it just processes the content in whole, assuming that
the only content is the expected content. Therefore, anyone using this
would have to have a reader that only has the specific content that is
valid (or assumed to be), which is hardly useful. Generally,
the scanner is much more useful anyway, as it processes in a more
predictable way, and could potentially be passed between the three validators, header,
body, and trailer. With the reader, content needs to be processed byte-by-byte,
which just doesn't match how the commit message would be processed, line-by-line
like the scanner does it. On this point, however, it's also important to make
sure that a _line_ scanner is passed, not some other scanner.

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

# Add help messages for validators, include in program

Like in the header, something about what is expected of each part (line length etc.)

# Add info to errors about part (e.g. Body: Line X: (issue with body))

# ? Save errors as per-line and print them out next to the line of text
