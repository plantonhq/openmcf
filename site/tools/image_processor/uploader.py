"""
R2 upload utilities.

Handles uploading content to Cloudflare R2 with progress feedback.
"""

import mimetypes
import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Optional

from rich.console import Console
from rich.progress import (
    BarColumn,
    Progress,
    SpinnerColumn,
    TaskProgressColumn,
    TextColumn,
    TimeElapsedColumn,
)

from .config import (
    ASSETS_DIR,
    R2_BASE_PATH,
    R2_BUCKET,
    R2_ENDPOINT,
    R2_PROFILE,
    R2_PUBLIC_URL,
    SYNC_EXCLUDES,
)

console = Console()


@dataclass
class UploadResult:
    """Result of uploading files to R2."""

    success: bool
    files_uploaded: int
    total_bytes: int
    error: str = ""


def get_r2_url(local_path: Path) -> str:
    """
    Get the public R2 URL for a local assets path.

    Args:
        local_path: Path relative to assets directory or absolute path

    Returns:
        Public URL for the file on R2
    """
    # Make path relative to assets directory
    if local_path.is_absolute():
        try:
            rel_path = local_path.relative_to(ASSETS_DIR)
        except ValueError:
            # Path is not under assets directory
            rel_path = local_path
    else:
        rel_path = local_path

    return f"{R2_PUBLIC_URL}/{R2_BASE_PATH}/{rel_path}"


def upload_file(local_path: Path, remote_path: Optional[str] = None) -> bool:
    """
    Upload a single file to R2.

    Args:
        local_path: Path to the local file
        remote_path: Optional remote path (defaults to path relative to assets/)

    Returns:
        True if upload succeeded
    """
    if remote_path is None:
        try:
            rel_path = local_path.relative_to(ASSETS_DIR)
        except ValueError:
            console.print(f"[red]Error: {local_path} is not under assets directory[/red]")
            return False
        remote_path = f"{R2_BASE_PATH}/{rel_path}"

    # Determine content type
    content_type, _ = mimetypes.guess_type(str(local_path))
    if content_type is None:
        content_type = "application/octet-stream"

    s3_uri = f"s3://{R2_BUCKET}/{remote_path}"

    try:
        result = subprocess.run(
            [
                "aws",
                "s3",
                "cp",
                str(local_path),
                s3_uri,
                "--profile",
                R2_PROFILE,
                "--endpoint-url",
                R2_ENDPOINT,
                "--content-type",
                content_type,
            ],
            capture_output=True,
            text=True,
            timeout=120,
        )
        return result.returncode == 0
    except (subprocess.TimeoutExpired, subprocess.SubprocessError):
        return False


def sync_assets() -> UploadResult:
    """
    Sync the assets directory to R2.

    Uses aws s3 sync with excludes for inbox, rules, and markdown files.

    Returns:
        UploadResult with sync statistics
    """
    # Build exclude arguments
    exclude_args = []
    for pattern in SYNC_EXCLUDES:
        exclude_args.extend(["--exclude", pattern])

    s3_uri = f"s3://{R2_BUCKET}/{R2_BASE_PATH}/"

    console.print()
    console.print(f"[bold]Syncing to R2[/bold]")
    console.print(f"  [dim]Source:[/dim] {ASSETS_DIR}")
    console.print(f"  [dim]Destination:[/dim] {R2_PUBLIC_URL}/{R2_BASE_PATH}/")
    console.print()

    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
        BarColumn(),
        TaskProgressColumn(),
        TimeElapsedColumn(),
        console=console,
    ) as progress:
        task = progress.add_task("Uploading...", total=None)

        try:
            result = subprocess.run(
                [
                    "aws",
                    "s3",
                    "sync",
                    str(ASSETS_DIR),
                    s3_uri,
                    "--profile",
                    R2_PROFILE,
                    "--endpoint-url",
                    R2_ENDPOINT,
                    *exclude_args,
                ],
                capture_output=True,
                text=True,
                timeout=600,  # 10 minute timeout for large syncs
            )

            progress.update(task, completed=True)

            if result.returncode != 0:
                return UploadResult(
                    success=False,
                    files_uploaded=0,
                    total_bytes=0,
                    error=result.stderr.strip() or "Sync failed",
                )

            # Parse output to count uploaded files
            # AWS CLI outputs lines like "upload: content/... to s3://..."
            lines = result.stdout.strip().split("\n") if result.stdout.strip() else []
            upload_lines = [l for l in lines if l.startswith("upload:")]

            return UploadResult(
                success=True,
                files_uploaded=len(upload_lines),
                total_bytes=0,  # AWS CLI doesn't easily provide this
            )

        except subprocess.TimeoutExpired:
            return UploadResult(
                success=False,
                files_uploaded=0,
                total_bytes=0,
                error="Sync timed out after 10 minutes",
            )
        except subprocess.SubprocessError as e:
            return UploadResult(
                success=False,
                files_uploaded=0,
                total_bytes=0,
                error=str(e),
            )


def display_upload_result(result: UploadResult) -> None:
    """Display upload result with Rich formatting."""
    console.print()
    if result.success:
        if result.files_uploaded > 0:
            console.print(f"  [green]✓[/green] {result.files_uploaded} files uploaded")
        else:
            console.print("  [green]✓[/green] Already in sync (no changes)")
    else:
        console.print(f"  [red]✗[/red] Upload failed: {result.error}")

