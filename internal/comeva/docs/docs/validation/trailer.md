# Trailer Validation

The trailer section of a commit message is of the following format (see
[git trailers][1] and [Conventional Commits][3]):

```
<key1>: <value1>
[<key2>: <value2>]
...
```

- key
    - Key of the trailer. It starts the line and contains ascii letters separated
    possibly by dashes. A special case is the BREAKING-CHANGE key, discussed
    further below.
    - configuration 
        - required keys
            - If specified, trailers with those keys must be found in the message.
        - optional keys
            - If specified, any key that isn't either in required keys or
            optional keys is not allowed.

- value
    - The value of the trailer. Not particularly limited in what it can contain,
    though it must be non-empty. It can span multiple lines, but any lines
    after the first must be indented by a configurable number of spaces, and no
    line may be empty.
    - configuration
        - continuation indent
            - Specifies how many spaces should be used to indent the continuation
            lines of the value.

The key and value are separated by a colon (:) and a single space. The length of
any line of the trailer can be limited, being by default between 0 and 80
characters. Multiple trailers may be specified, and keys may be repeated.
The trailer section may also be empty depending on the configuration.

## Breaking change trailer

A special hard-coded key is the BREAKING-CHANGE key. This is used to denote,
along with [an exclamation mark in the header](./header.md), that the commit
introduces a breaking change, which equates to a MAJOR version change in 
[Semantic Versioning][2]. 

The BREAKING-CHANGE key is optional, but required
if the header has the exclamation mark that also denotes a breaking change.
Unlike in [Conventional Commits][3], where one can be omitted if the other is
present, the requirement for specifying both at the same time has a two-fold
reason: first, as the exclamation mark is in the header, git commands that
show logs as oneliners (showing the header of each commit log) will still allow
easily seeing where breaking changes were introduced. Second, requiring the
trailer means that the breaking change is more likely to be properly explained.
Further, because trailers allow for more than line, the breaking change can be
explained in further detail than what the header may be able to fit.

[Conventional Commits][3] also allows the BREAKING-CHANGE key to instead be
specified as 'BREAKING CHANGE' (a space instead of the dash). This is not
allowed here, where all keys must instead be specified without spaces. This is
naturally done for consistency's sake.


[1]: https://git-scm.com/docs/git-interpret-trailers

[2]: https://semver.org/ 'Semantic Versioning'

[3]: https://www.conventionalcommits.org/en/v1.0.0/#specification 'Conventional Commits Specification'