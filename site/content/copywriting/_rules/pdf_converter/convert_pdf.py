#!/usr/bin/env python3
"""
PDF to Markdown Converter for Planton.ai Copywriting Workflow

This script converts PDF files in the _workspace directory to Markdown format
using pdftotext (from Poppler) and pandoc.

Requirements:
- Poppler (provides pdftotext)
- Pandoc (converts text to markdown)

Install via Homebrew:
  brew install poppler pandoc
"""

import os
import sys
import subprocess
import shutil
from pathlib import Path
from typing import Tuple, List


def check_command_exists(command: str) -> bool:
    """Check if a command exists in PATH."""
    return shutil.which(command) is not None


def check_dependencies() -> Tuple[bool, List[str]]:
    """
    Check if required dependencies (pdftotext, pandoc) are installed.
    
    Returns:
        Tuple of (all_installed: bool, missing_commands: List[str])
    """
    missing = []
    
    if not check_command_exists('pdftotext'):
        missing.append('pdftotext (from Poppler)')
    
    if not check_command_exists('pandoc'):
        missing.append('pandoc')
    
    return (len(missing) == 0, missing)


def print_error_message(missing_commands: List[str]):
    """Print a nicely formatted error message for missing dependencies."""
    print("\n" + "="*70)
    print("  ❌ MISSING DEPENDENCIES")
    print("="*70)
    print("\nThe following tools are required but not installed:\n")
    
    for cmd in missing_commands:
        print(f"  • {cmd}")
    
    print("\n" + "-"*70)
    print("  📦 INSTALLATION")
    print("-"*70)
    print("\nInstall the missing dependencies using Homebrew:\n")
    
    brew_packages = []
    if 'pdftotext' in ' '.join(missing_commands):
        brew_packages.append('poppler')
    if 'pandoc' in ' '.join(missing_commands):
        brew_packages.append('pandoc')
    
    if brew_packages:
        print(f"  brew install {' '.join(brew_packages)}")
    
    print("\n" + "-"*70)
    print("  📖 MORE INFO")
    print("-"*70)
    print("\n  Poppler: https://poppler.freedesktop.org/")
    print("  Pandoc:  https://pandoc.org/")
    print("\n" + "="*70 + "\n")


def convert_pdf_to_markdown(pdf_path: Path, output_path: Path) -> bool:
    """
    Convert a PDF file to Markdown using pdftotext and pandoc.
    
    Args:
        pdf_path: Path to input PDF file
        output_path: Path to output Markdown file
    
    Returns:
        True if conversion successful, False otherwise
    """
    try:
        # Step 1: Extract text from PDF using pdftotext
        print(f"  📄 Extracting text from: {pdf_path.name}")
        pdftotext_result = subprocess.run(
            ['pdftotext', str(pdf_path), '-'],
            capture_output=True,
            text=True,
            check=True
        )
        
        if not pdftotext_result.stdout.strip():
            print(f"  ⚠️  Warning: No text extracted from {pdf_path.name}")
            print(f"      This might be a scanned PDF requiring OCR")
            return False
        
        # Step 2: Convert text to Markdown using pandoc
        print(f"  🔄 Converting to Markdown...")
        pandoc_result = subprocess.run(
            ['pandoc', '-t', 'markdown', '--wrap=none'],
            input=pdftotext_result.stdout,
            capture_output=True,
            text=True,
            check=True
        )
        
        # Step 3: Write output
        output_path.write_text(pandoc_result.stdout, encoding='utf-8')
        print(f"  ✅ Created: {output_path.name}")
        
        return True
        
    except subprocess.CalledProcessError as e:
        print(f"  ❌ Error converting {pdf_path.name}:")
        print(f"     {e}")
        if e.stderr:
            print(f"     {e.stderr}")
        return False
    except Exception as e:
        print(f"  ❌ Unexpected error: {e}")
        return False


def find_pdf_files(workspace_dir: Path) -> List[Path]:
    """Find all PDF files in the workspace directory."""
    return list(workspace_dir.glob('*.pdf'))


def main():
    """Main entry point for PDF to Markdown conversion."""
    
    # Print header
    print("\n" + "="*70)
    print("  📑 PDF to Markdown Converter")
    print("  Planton.ai Copywriting Workflow")
    print("="*70 + "\n")
    
    # Check dependencies
    print("🔍 Checking dependencies...")
    all_installed, missing = check_dependencies()
    
    if not all_installed:
        print_error_message(missing)
        sys.exit(1)
    
    print("  ✅ All dependencies installed\n")
    
    # Determine workspace directory
    script_dir = Path(__file__).parent
    workspace_dir = script_dir.parent.parent / '_workspace'
    
    if not workspace_dir.exists():
        print(f"❌ Error: Workspace directory not found")
        print(f"   Expected: {workspace_dir}")
        sys.exit(1)
    
    print(f"📂 Workspace: {workspace_dir.relative_to(Path.cwd())}\n")
    
    # Find PDF files
    pdf_files = find_pdf_files(workspace_dir)
    
    if not pdf_files:
        print("ℹ️  No PDF files found in workspace")
        print("   Nothing to convert\n")
        sys.exit(0)
    
    print(f"📋 Found {len(pdf_files)} PDF file(s)\n")
    
    # Convert each PDF
    converted_count = 0
    failed_count = 0
    
    for pdf_path in sorted(pdf_files):
        print(f"Processing: {pdf_path.name}")
        print("-" * 70)
        
        # Output path: same name but .md extension
        output_path = pdf_path.with_suffix('.md')
        
        if convert_pdf_to_markdown(pdf_path, output_path):
            converted_count += 1
            
            # Optionally remove the original PDF
            # For now, we keep it for reference
            # pdf_path.unlink()
            
        else:
            failed_count += 1
        
        print()
    
    # Print summary
    print("="*70)
    print("  📊 CONVERSION SUMMARY")
    print("="*70)
    print(f"\n  Total PDF files:     {len(pdf_files)}")
    print(f"  Successfully converted:  {converted_count} ✅")
    if failed_count > 0:
        print(f"  Failed:              {failed_count} ❌")
    print("\n" + "="*70 + "\n")
    
    # Exit with appropriate code
    if failed_count > 0:
        sys.exit(1)
    else:
        sys.exit(0)


if __name__ == '__main__':
    main()

