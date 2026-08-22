"""
Image organization utilities.

Handles naming conventions and file organization.
"""

import re
import shutil
from datetime import datetime
from pathlib import Path
from typing import Optional

from rich.console import Console

from .config import IMAGES_DIR, INBOX_DIR, SUPPORTED_FORMATS

console = Console()


def get_inbox_images() -> list[Path]:
    """Get all images from the inbox directory."""
    if not INBOX_DIR.exists():
        return []

    images = []
    for path in INBOX_DIR.iterdir():
        if path.is_file() and path.suffix.lower() in SUPPORTED_FORMATS:
            images.append(path)

    # Sort by modification time (newest first)
    images.sort(key=lambda p: p.stat().st_mtime, reverse=True)
    return images


def sanitize_name(name: str) -> str:
    """
    Sanitize a string for use in a filename.

    Converts to lowercase, replaces spaces and special chars with hyphens,
    removes consecutive hyphens, and strips leading/trailing hyphens.
    """
    # Convert to lowercase
    name = name.lower()
    # Replace common separators with hyphens
    name = re.sub(r"[\s_\.]+", "-", name)
    # Remove any character that's not alphanumeric or hyphen
    name = re.sub(r"[^a-z0-9\-]", "", name)
    # Replace multiple consecutive hyphens with single hyphen
    name = re.sub(r"-+", "-", name)
    # Strip leading/trailing hyphens
    name = name.strip("-")
    return name


def generate_filename(
    context: str,
    description: str,
    original_suffix: str,
    timestamp: Optional[datetime] = None,
) -> str:
    """
    Generate a filename following the naming convention.

    Format: YYYY-MM-DD-HHMMSS-context-description.ext

    Args:
        context: The context/feature area (e.g., "service-hub", "kubernetes-dashboard")
        description: Description of the image content
        original_suffix: The original file extension (e.g., ".png")
        timestamp: Optional timestamp (defaults to now)

    Returns:
        Generated filename string
    """
    if timestamp is None:
        timestamp = datetime.now()

    # Format timestamp
    ts_str = timestamp.strftime("%Y-%m-%d-%H%M%S")

    # Sanitize context and description
    ctx = sanitize_name(context)
    desc = sanitize_name(description)

    # Ensure suffix starts with a dot
    suffix = original_suffix.lower()
    if not suffix.startswith("."):
        suffix = f".{suffix}"

    # Build filename
    parts = [ts_str]
    if ctx:
        parts.append(ctx)
    if desc:
        parts.append(desc)

    return "-".join(parts) + suffix


def determine_destination(context: str) -> Path:
    """
    Determine the destination directory based on context.

    Args:
        context: The context/feature area string

    Returns:
        Path to the destination directory
    """
    # Sanitize context for use as directory path
    ctx = sanitize_name(context)

    # Handle nested contexts (e.g., "service-hub/deployment-environments")
    if "/" in context:
        parts = [sanitize_name(p) for p in context.split("/")]
        dest = IMAGES_DIR.joinpath(*parts)
    else:
        dest = IMAGES_DIR / ctx if ctx else IMAGES_DIR

    return dest


def move_image(
    source: Path,
    destination_dir: Path,
    new_filename: str,
) -> Path:
    """
    Move an image to its destination with a new filename.

    Args:
        source: Path to the source image
        destination_dir: Directory to move the image to
        new_filename: New filename for the image

    Returns:
        Path to the moved image
    """
    # Create destination directory if it doesn't exist
    destination_dir.mkdir(parents=True, exist_ok=True)

    destination = destination_dir / new_filename

    # Handle filename collision
    if destination.exists():
        # Append a counter to make unique
        base = destination.stem
        suffix = destination.suffix
        counter = 1
        while destination.exists():
            destination = destination_dir / f"{base}-{counter}{suffix}"
            counter += 1

    shutil.move(str(source), str(destination))
    return destination


def clear_inbox() -> int:
    """
    Clear all images from the inbox.

    Returns:
        Number of files removed
    """
    count = 0
    if INBOX_DIR.exists():
        for path in INBOX_DIR.iterdir():
            if path.is_file() and path.suffix.lower() in SUPPORTED_FORMATS:
                path.unlink()
                count += 1
    return count

