"""Allow running the package as a module: python -m tools.image_processor"""

from .cli import cli

if __name__ == "__main__":
    cli()

