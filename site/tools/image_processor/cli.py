#!/usr/bin/env python3
"""
Planton Image Processor CLI

A tool for processing, compressing, and uploading images to Cloudflare R2.
"""

import sys
import tempfile
from pathlib import Path
from typing import Optional

import click
from rich.console import Console
from rich.panel import Panel
from rich.table import Table

from . import __version__
from .compressor import CompressionResult, compress_image, format_size
from .config import IMAGES_DIR, INBOX_DIR, R2_PUBLIC_URL, R2_BASE_PATH, CONTENT_DIR
from .organizer import (
    determine_destination,
    generate_filename,
    get_inbox_images,
    move_image,
)
from .preflight import run_preflight
from .uploader import display_upload_result, get_r2_url, sync_assets

console = Console()


def print_header():
    """Print the application header."""
    console.print()
    console.print(
        Panel(
            "[bold cyan]Planton Image Processor[/bold cyan]",
            subtitle=f"v{__version__}",
            border_style="cyan",
        )
    )


@click.group()
@click.version_option(version=__version__)
def cli():
    """Planton Image Processor - Compress, organize, and upload images to R2."""
    pass


@cli.command()
def check():
    """Run preflight checks to verify tools and R2 access."""
    print_header()
    all_passed = run_preflight()

    if all_passed:
        console.print("[green]All checks passed![/green]")
        console.print()
        sys.exit(0)
    else:
        console.print("[red]Some checks failed. Please fix the issues above.[/red]")
        console.print()
        sys.exit(1)


@cli.command()
def sync():
    """Sync assets folder to R2 (without processing)."""
    print_header()

    # Run preflight first
    if not run_preflight():
        console.print("[red]Preflight checks failed. Cannot sync.[/red]")
        console.print()
        sys.exit(1)

    result = sync_assets()
    display_upload_result(result)
    console.print()

    if result.success:
        sys.exit(0)
    else:
        sys.exit(1)


@cli.command()
@click.option(
    "--context",
    "-c",
    required=True,
    help="Context/feature area (e.g., 'service-hub', 'kubernetes-dashboard/pod-logs')",
)
@click.option(
    "--description",
    "-d",
    default="",
    help="Description of the image(s) content",
)
@click.option(
    "--destination",
    "-D",
    default=None,
    help="Override destination folder (relative to content/assets/images/)",
)
@click.option(
    "--no-upload",
    is_flag=True,
    help="Skip uploading to R2 after processing",
)
@click.option(
    "--dry-run",
    is_flag=True,
    help="Show what would be done without making changes",
)
def process(
    context: str,
    description: str,
    destination: Optional[str],
    no_upload: bool,
    dry_run: bool,
):
    """
    Process images from the inbox.

    Compresses images, renames them following the naming convention,
    moves them to the appropriate destination, and uploads to R2.

    Examples:

        # Process with context
        python -m tools.image_processor process -c "service-hub" -d "deployment-pipeline"

        # Process with nested context
        python -m tools.image_processor process -c "service-hub/deployment-environments"

        # Dry run to see what would happen
        python -m tools.image_processor process -c "kubernetes-dashboard" --dry-run
    """
    print_header()

    # Get images from inbox
    images = get_inbox_images()

    if not images:
        console.print()
        console.print("[yellow]No images found in inbox.[/yellow]")
        console.print(f"[dim]Inbox location: {INBOX_DIR}[/dim]")
        console.print()
        sys.exit(0)

    console.print()
    console.print(f"[bold]Found {len(images)} image(s) in inbox[/bold]")
    console.print()

    # Run preflight checks (unless dry run)
    if not dry_run:
        if not run_preflight():
            console.print("[red]Preflight checks failed. Cannot process.[/red]")
            console.print()
            sys.exit(1)

    # Determine destination directory
    dest_context = destination if destination else context
    dest_dir = determine_destination(dest_context)

    console.print(f"[bold]Processing Images[/bold]")
    console.print()

    # Track statistics
    total_original = 0
    total_compressed = 0
    processed_files: list[Path] = []
    errors: list[str] = []

    for i, image_path in enumerate(images, 1):
        console.print(f"  [{i}/{len(images)}] [cyan]{image_path.name}[/cyan]")

        # Generate new filename
        # Use index as description suffix if multiple images
        img_desc = description
        if len(images) > 1 and description:
            img_desc = f"{description}-{i}"
        elif not description:
            # Use original filename stem as description fallback
            img_desc = image_path.stem

        new_filename = generate_filename(
            context=context,
            description=img_desc,
            original_suffix=image_path.suffix,
        )

        console.print(f"        → [green]{new_filename}[/green]")

        if dry_run:
            console.print(f"        [dim]Would move to: {dest_dir}[/dim]")
            console.print(f"        [dim](dry run - no changes made)[/dim]")
            console.print()
            continue

        # Compress to temp file first
        with tempfile.TemporaryDirectory() as tmpdir:
            temp_output = Path(tmpdir) / new_filename

            console.print("        Compressing...", end=" ")
            result = compress_image(image_path, temp_output)

            if result.success:
                total_original += result.original_size
                total_compressed += result.compressed_size

                if result.saved_percent > 0:
                    console.print(
                        f"[green]{format_size(result.original_size)} → "
                        f"{format_size(result.compressed_size)} "
                        f"({result.saved_percent:.0f}% saved)[/green]"
                    )
                else:
                    console.print(f"[dim]{format_size(result.original_size)} (no reduction)[/dim]")

                # Move to destination
                final_path = move_image(temp_output, dest_dir, new_filename)
                processed_files.append(final_path)

                console.print(f"        [dim]Destination: {dest_dir.relative_to(CONTENT_DIR.parent)}[/dim]")
                console.print(f"        [green]✓[/green] Moved")

                # Delete original from inbox
                image_path.unlink()

            else:
                errors.append(f"{image_path.name}: {result.error}")
                console.print(f"[red]Failed: {result.error}[/red]")

        console.print()

    # Print summary
    if not dry_run and processed_files:
        console.print("[bold]Summary[/bold]")
        console.print()

        saved_bytes = total_original - total_compressed
        saved_percent = (saved_bytes / total_original * 100) if total_original > 0 else 0

        console.print(f"  [green]✓[/green] {len(processed_files)} image(s) processed")
        console.print(
            f"  [green]✓[/green] {format_size(saved_bytes)} saved ({saved_percent:.0f}% reduction)"
        )

        if errors:
            console.print(f"  [red]✗[/red] {len(errors)} error(s)")
            for error in errors:
                console.print(f"      [red]{error}[/red]")

        console.print()

        # Upload to R2
        if not no_upload:
            result = sync_assets()
            display_upload_result(result)

            if result.success and processed_files:
                console.print()
                console.print("[bold]URLs[/bold]")
                console.print()
                for path in processed_files:
                    url = get_r2_url(path)
                    console.print(f"  [dim]{url}[/dim]")

    console.print()


@cli.command()
def inbox():
    """List images currently in the inbox."""
    print_header()

    images = get_inbox_images()

    if not images:
        console.print()
        console.print("[yellow]No images found in inbox.[/yellow]")
        console.print(f"[dim]Inbox location: {INBOX_DIR}[/dim]")
        console.print()
        return

    console.print()
    console.print(f"[bold]Images in Inbox ({len(images)})[/bold]")
    console.print()

    table = Table(show_header=True, header_style="bold")
    table.add_column("Filename")
    table.add_column("Size", justify="right")
    table.add_column("Type")

    for image in images:
        size = format_size(image.stat().st_size)
        suffix = image.suffix.lower().lstrip(".")
        table.add_row(image.name, size, suffix.upper())

    console.print(table)
    console.print()


# Allow running as module
if __name__ == "__main__":
    cli()

