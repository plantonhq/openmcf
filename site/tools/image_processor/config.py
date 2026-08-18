"""
Configuration constants for the image processor.
"""

from pathlib import Path

# Project paths
PROJECT_ROOT = Path(__file__).parent.parent.parent
CONTENT_DIR = PROJECT_ROOT / "content"
ASSETS_DIR = CONTENT_DIR / "assets"
INBOX_DIR = ASSETS_DIR / "_inbox"
IMAGES_DIR = ASSETS_DIR / "images"

# R2 configuration
R2_BUCKET = "planton-assets"
R2_ENDPOINT = "https://074755a78d8e8f77c119a90a125e8a06.r2.cloudflarestorage.com"
R2_PROFILE = "r2"
R2_BASE_PATH = "site"
R2_PUBLIC_URL = "https://assets.planton.ai"

# Compression settings
JPEG_QUALITY_MAX = 90  # jpegoptim --max value
PNG_QUALITY_RANGE = "80-95"  # pngquant quality range
PNG_SPEED = 1  # pngquant speed (1=slowest/best, 11=fastest)

# Supported image formats
SUPPORTED_FORMATS = {
    ".jpg": "jpeg",
    ".jpeg": "jpeg",
    ".png": "png",
    ".gif": "gif",
    ".webp": "webp",
}

# Formats that should be compressed
COMPRESSIBLE_FORMATS = {".jpg", ".jpeg", ".png"}

# Files to exclude from sync
SYNC_EXCLUDES = [
    "_inbox/*",
    "_rules/*",
    "*.md",
    ".gitignore",
    ".DS_Store",
]

