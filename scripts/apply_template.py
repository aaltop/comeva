from pathlib import Path

root = Path()
docs_dir = root / "docs"


def readme_to_index():
    """
    Copy README.md to docs/index.md.
    """

    origin = root / "README.md"
    destination = docs_dir / "index.md"

    destination.write_text(origin.read_text())


def main():

    readme_to_index()


if __name__ == "__main__":
    main()
