from collections.abc import Generator

from pathlib import Path
import sys


def flatten_contents(path: Path) -> Generator[Path, None, None]:
    """
    Get the contents of the path, expected to be a directory, recursively.
    """

    for p in path.iterdir():
        if p.is_file():
            yield p
        elif p.is_dir():
            yield p
            for sub_p in flatten_contents(p):
                yield sub_p
        else:
            raise ValueError("Expected file or directory")


_root = Path()


# the idea of having the docs in src/ first is that that way, the up-to-date
# docs are guaranteed to be embedded in the tool, which is more important than
# having the docs/ be up-to-date.
def docs_copy(docs: Path):
    """
    Copy docs from `docs` to docs/.
    """

    deep_docs = docs
    root_docs = _root / "docs"

    files: list[tuple[Path, Path]] = []
    for p in flatten_contents(deep_docs):
        docs_relative = p.relative_to(deep_docs.parent)
        if p.is_dir():
            for dir in list(reversed(list(docs_relative.parents))) + [docs_relative]:
                # make only those directories that that are relative to
                # the docs/ directory at the root (making parents (parents=True) always
                # feels a little sketchy)
                if dir.is_relative_to(root_docs):
                    dir.mkdir(exist_ok=True)
        else:
            files.append((p, docs_relative))

    for origin, destination in files:
        destination.write_text(origin.read_text())


def main():

    deep_docs = _root.joinpath(sys.argv[1])
    if not deep_docs.exists():
        raise ValueError("Passed docs location should exist")

    if not deep_docs.is_relative_to(_root):
        raise ValueError("Passed docs location should be relative to the root")

    docs_copy(deep_docs)


if __name__ == "__main__":
    main()
