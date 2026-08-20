"""
Image compression utilities.

Uses jpegoptim for JPEG and pngquant for PNG compression.
"""

import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path

from rich.console import Console

from .config import (
    COMPRESSIBLE_FORMATS,
    JPEG_QUALITY_MAX,
    PNG_QUALITY_RANGE,
    PNG_SPEED,
)

console = Console()


@dataclass
class CompressionResult:
    """Result of compressing an image."""

    success: bool
    original_size: int
    compressed_size: int
    output_path: Path
    error: str = ""

    @property
    def saved_bytes(self) -> int:
        """Bytes saved by compression."""
        return self.original_size - self.compressed_size

    @property
    def saved_percent(self) -> float:
        """Percentage of bytes saved."""
        if self.original_size == 0:
            return 0.0
        return (self.saved_bytes / self.original_size) * 100


def format_size(size_bytes: int) -> str:
    """Format byte size as human-readable string."""
    if size_bytes < 1024:
        return f"{size_bytes} B"
    elif size_bytes < 1024 * 1024:
        return f"{size_bytes / 1024:.1f} KB"
    else:
        return f"{size_bytes / (1024 * 1024):.1f} MB"


def compress_jpeg(input_path: Path, output_path: Path) -> CompressionResult:
    """
    Compress a JPEG image using jpegoptim.

    Args:
        input_path: Path to the input JPEG file.
        output_path: Path where the compressed file should be saved.

    Returns:
        CompressionResult with compression statistics.
    """
    original_size = input_path.stat().st_size

    try:
        # jpegoptim modifies in-place, so we need to copy first
        import shutil

        shutil.copy2(input_path, output_path)

        result = subprocess.run(
            [
                "jpegoptim",
                "--strip-all",  # Remove all metadata
                f"--max={JPEG_QUALITY_MAX}",  # Max quality
                "--quiet",
                str(output_path),
            ],
            capture_output=True,
            text=True,
            timeout=60,
        )

        if result.returncode != 0:
            return CompressionResult(
                success=False,
                original_size=original_size,
                compressed_size=original_size,
                output_path=output_path,
                error=result.stderr.strip() or "jpegoptim failed",
            )

        compressed_size = output_path.stat().st_size
        return CompressionResult(
            success=True,
            original_size=original_size,
            compressed_size=compressed_size,
            output_path=output_path,
        )

    except subprocess.TimeoutExpired:
        return CompressionResult(
            success=False,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
            error="Compression timed out",
        )
    except Exception as e:
        return CompressionResult(
            success=False,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
            error=str(e),
        )


def compress_png(input_path: Path, output_path: Path) -> CompressionResult:
    """
    Compress a PNG image using pngquant.

    Args:
        input_path: Path to the input PNG file.
        output_path: Path where the compressed file should be saved.

    Returns:
        CompressionResult with compression statistics.
    """
    original_size = input_path.stat().st_size

    try:
        result = subprocess.run(
            [
                "pngquant",
                f"--quality={PNG_QUALITY_RANGE}",
                f"--speed={PNG_SPEED}",
                "--force",  # Overwrite output if exists
                "--output",
                str(output_path),
                str(input_path),
            ],
            capture_output=True,
            text=True,
            timeout=60,
        )

        # pngquant returns 99 if quality can't be achieved (not an error)
        if result.returncode not in (0, 99):
            # If pngquant fails, just copy the original
            import shutil

            shutil.copy2(input_path, output_path)
            return CompressionResult(
                success=True,
                original_size=original_size,
                compressed_size=original_size,
                output_path=output_path,
                error="pngquant skipped (quality threshold not met, kept original)",
            )

        compressed_size = output_path.stat().st_size
        return CompressionResult(
            success=True,
            original_size=original_size,
            compressed_size=compressed_size,
            output_path=output_path,
        )

    except subprocess.TimeoutExpired:
        return CompressionResult(
            success=False,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
            error="Compression timed out",
        )
    except Exception as e:
        return CompressionResult(
            success=False,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
            error=str(e),
        )


def compress_image(input_path: Path, output_path: Path) -> CompressionResult:
    """
    Compress an image based on its format.

    Args:
        input_path: Path to the input image file.
        output_path: Path where the compressed file should be saved.

    Returns:
        CompressionResult with compression statistics.
    """
    suffix = input_path.suffix.lower()
    original_size = input_path.stat().st_size

    if suffix not in COMPRESSIBLE_FORMATS:
        # For non-compressible formats (GIF, WebP), just copy
        import shutil

        shutil.copy2(input_path, output_path)
        return CompressionResult(
            success=True,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
            error=f"Format {suffix} not compressed (copied as-is)",
        )

    if suffix in (".jpg", ".jpeg"):
        return compress_jpeg(input_path, output_path)
    elif suffix == ".png":
        return compress_png(input_path, output_path)
    else:
        # Shouldn't reach here, but handle gracefully
        import shutil

        shutil.copy2(input_path, output_path)
        return CompressionResult(
            success=True,
            original_size=original_size,
            compressed_size=original_size,
            output_path=output_path,
        )

