# Header Validation

The header of a commit message is of the following format (see
[Conventional Commits][1]):

```
<type>[(<scope>)]: <verb> <content>
```

- type
    - The type of change the commit is making. [Conventional Commits][1]
    suggests 'feat' and 'fix' as the defaults, also used by default by comeva.
    'feat' specifies a MINOR change in [Semantic Versioning][2], while 'fix'
    specifies a PATCH change.

    - configuration: empty means any value is valid. Otherwise, one of the
    supplied values must be present.

- scope
    - The scope of the commit. Can be any ['additional contextual information'][3].
    In practice specified by the user, with anything allowed by default.
    
    - configuration: empty means any value is valid. Otherwise, when a scope
    is encountered, it must be one of the supplied values.

- description
    - The 'message' part of the header, tells what happened (or what will happen
    if this commit is to be applied). Comes after the type and optional scope,
    from which it is delimited by a colon (:) and a single space.
    Expected format `<verb> <content>`.
        - verb
            - [The imperative][4] that begins the description. A verb in the
            imperative mood describing what the commit will do if applied.
            For example, 'Add' when something is added, or 'Change' when
            something already existing is changed. Refers to the main contribution
            of the commit, as naturally multiple things may change or be added
            as a result of the commit.

            - configuration: empty means any value is valid. Otherwise, one of
            the supplied values must begin the description.

        - content
            - The rest of the description after the verb. Not particularly limited
            in what it can contain.

Besides the content of the header, the length of the header can also be limited.
By default, it is set to be between 0 and 80 characters, and is configurable.




[1]: https://www.conventionalcommits.org/en/v1.0.0/ 'Conventional Commits'

[2]: https://semver.org/ 'Semantic Versioning'

[3]: https://www.conventionalcommits.org/en/v1.0.0/#:~:text=additional%20contextual%20information

[4]: https://cbea.ms/git-commit/#:~:text=Imperative%20mood%20just%20means%20%E2%80%9Cspoken%20or%20written%20as%20if%20giving%20a%20command%20or%20instruction%E2%80%9D%2E 'Imperative mood'