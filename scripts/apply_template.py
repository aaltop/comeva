from pathlib import Path

from docs_copy import docs_copy

root = Path()
docs_dir = root / "docs"
deep_docs_dir = root / "internal" / "comeva" / "docs" / "docs"


def readme_to_index():
    """
    Copy README.md to the index of the docs.
    """

    origin = root / "README.md"
    destination = deep_docs_dir / "index.md"

    destination.write_text(origin.read_text())


def main():

    readme_to_index()
    docs_copy(root / "internal" / "comeva" / "docs" / "docs")


if __name__ == "__main__":
    main()
