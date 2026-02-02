# ? Evaluate sensibility of having the io.Reader Validate method.

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

I guess it's fine? doesn't really matter too much.

# Add config file

## comeva config file

Config file to configure the default values (commit file, validator config file)
to be passed to comeva.

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

# add standard settings file

Specify standard commit file, verbosity etc.

# ? Add use of colors in console

# ? Save errors as per-line and print them out next to the line of text

# ? add parsing of trailer's values with expected syntax

Something like having an "Effect" key that denotes the effect a breaking change
has, then being able to demand a comma-separated list of specific allowed values
for that. A good bit of effort to implement, I imagine.

Should be part of a "second step" in validation? First step is largely
formatting related, like with key-value pairs having a colon and space in between,
certain line length, and a certain indentation. After this has been correctly
parsed is the value actually checked further for correct formatting. First step
does still have similar checks, though, so there's not currently any particular
separation into steps even now. Could potentially make a change, so that there
is a FormattingValidator that ensures that the content looks right, then another
that checks for expected values. However, much of the current validation is
already the second step anyway. Regardless, should still at least make a separate
TrailerValueValidator if nothing else, so that the main TrailerValidator doesn't
become too bloated.

# ? Add option for asking whether the values are correct

If "feat" is passed, the program would ask "are you adding a new feature?"
or something along those lines, and then the user would confirm. This
is to validate that the values actually also match the intention, in addition
to being otherwise in the accepted set/range of values. This should ultimately
be naturally part of the commit message writing process, so whether the
developer should really need to be reminded to check the accuracy is
questionable. Having this feature might just mean that it is mindlessly
mashed through (even though it would be optional and off by default),
so deving it might not be that worthwhile.
