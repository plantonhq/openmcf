# CLI Installation Documentation

**Date**: January 6, 2026
**Type**: Feature
**Components**: Documentation, CLI, Content Management

## Summary

Added comprehensive CLI installation documentation covering macOS, Windows, and Linux platforms. Removed deprecated install automation files (`install.sh` and `cli/latest`) that pointed to old infrastructure.

## Problem Statement

The Planton CLI installation path was unclear for non-macOS users. While Homebrew made macOS installation straightforward, Windows and Linux users lacked proper documentation. Additionally, legacy install automation files were pointing to deprecated GCS bucket URLs instead of the new Cloudflare R2 distribution.

### Pain Points

- No documented installation path for Windows users
- Linux installation not clearly documented
- Legacy `install.sh` pointed to old `storage.googleapis.com/planton-cli` bucket
- Stale `cli/latest` file showed outdated version (`v0.2.62` vs actual `v0.0.3`)
- GitHub Pages doesn't support custom headers, making the install script approach less viable

## Solution

Created a new documentation page at `/docs/cli` with:
- Clear installation instructions for all platforms
- macOS: Homebrew (recommended) + direct download
- Windows: WSL option (curl) + native PowerShell
- Linux: Direct download for amd64 and arm64
- Authentication and post-install steps
- Troubleshooting section

Removed deprecated automation files:
- `public/install.sh` - Shell script with hardcoded GCS URLs
- `public/cli/latest` - Stale version pointer

## Implementation Details

### New Documentation Page

**File**: `public/docs/cli.md`

Key sections with auto-generated anchors:
- `#install-on-macos` - Homebrew + direct download
- `#install-on-windows` - WSL and PowerShell options
- `#install-on-linux` - amd64 and arm64 downloads
- `#verify-installation` - Confirmation steps
- `#authentication` - Login and context setup
- `#upgrading` - Update instructions

### Windows Support

Added two installation paths for Windows users:

1. **WSL users**: Same curl commands as Linux
2. **Native Windows**: PowerShell with `Invoke-WebRequest`, PATH configuration, and checksum verification

### Fixes Applied

| Issue | Fix |
|-------|-----|
| Anchor syntax `{#id}` showing literally | Removed - anchors auto-generate from headings |
| Duplicate React key error | Renamed checksum sections to be unique per platform |
| Section order | Reordered to macOS → Windows → Linux |
| Private GitHub releases link | Removed entirely |

### Files Deleted

- `public/install.sh` (103 lines) - Deprecated shell installer
- `public/cli/latest` (1 line) - Stale version file

## Benefits

### For Users
- Clear installation path for every platform
- Copy-paste commands that work immediately
- WSL support for Windows developers who prefer Unix tools
- Architecture detection tips (`uname -m`, `$env:PROCESSOR_ARCHITECTURE`)

### For Maintenance
- No more version file to keep updated
- No dependency on custom headers (GitHub Pages compatible)
- Single source of truth in documentation
- Release workflow doesn't need to update planton.ai

## Impact

**Users affected**: All new CLI users, especially Windows and Linux users

**Pages updated**:
- New: `/docs/cli` - Complete installation guide

**Files removed**:
- `public/install.sh`
- `public/cli/latest`

## Related Work

- Unified Release Pipeline (`planton/_changelog/2025-12/2025-12-29-075919-unified-release-pipeline.md`)
- CLI binaries now hosted at `downloads.planton.ai/cli/`
- Homebrew formula at `plantonhq/homebrew-tap`

---

**Status**: ✅ Live
**Timeline**: ~30 minutes implementation

