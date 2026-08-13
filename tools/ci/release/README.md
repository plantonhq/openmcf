# Release Tooling

Scripts and utilities for managing Planton releases (IaC modules, content
distribution zips, the InfraChart catalog, and definitions).

## Quick Start

```bash
# Show what the next version would be (defaults to patch bump)
make next-version

# Show next minor version
make next-version bump=minor

# Create a release (triggers GitHub Actions workflow)
make release                   # patch bump: v0.0.0 -> v0.0.1
make release bump=minor        # minor bump: v0.0.0 -> v0.1.0
make release bump=major        # major bump: v0.0.0 -> v1.0.0
```

## Release Flow

When you run `make release`, here's what happens:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  Local (make release)                                                       │
│  ├── Calculate next version using tools/ci/release/next_version.py          │
│  ├── Create git tag (e.g., v1.0.0)                                          │
│  └── Push tag to origin                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  GitHub Actions (.github/workflows/release.yaml)                            │
│  ├── Create the GitHub Release with auto-generated notes                    │
│  ├── Pulumi module binaries    -> Cloudflare R2 (downloads.planton.dev)     │
│  ├── Terraform module zips     -> Cloudflare R2 (downloads.planton.dev)     │
│  ├── Content distribution zips -> Cloudflare R2 (downloads.planton.dev)     │
│  ├── Definitions               -> Cloudflare R2 (downloads.planton.dev)     │
│  └── InfraChart catalog zip    -> attached to the GitHub Release            │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Scripts

### next_version.py

Calculates the next semantic version based on existing git tags.

```bash
# Usage
python3 tools/ci/release/next_version.py [patch|minor|major]

# Examples
python3 tools/ci/release/next_version.py          # patch bump (default)
python3 tools/ci/release/next_version.py minor    # minor bump
python3 tools/ci/release/next_version.py major    # major bump
```

The script:
- Finds the latest tag matching strict `vX.Y.Z` pattern
- Defaults to `v0.0.0` if no tags exist
- Outputs the next version to stdout

### detect_module_dirs.sh

Single source of truth for which changed files count as Pulumi/Terraform
module changes (maturity-grammar-aware). Used by `auto-tag.yaml`, which runs
its `--self-test` before every detection pass.

### package_content.sh

Packages the content distribution zips (presets, IaC source, catalog pages,
proto source, the component reference pack) for the release's R2 upload.

## Required GitHub Secrets

The module/content/definitions release workflows upload to Cloudflare R2 and
need `CLOUDFLARE_R2_ACCESS_KEY_ID`, `CLOUDFLARE_R2_SECRET_ACCESS_KEY`, and
`CLOUDFLARE_R2_ENDPOINT`. The GitHub Release itself and the chart-catalog
attachment need only the automatically provided `GITHUB_TOKEN`.

## Troubleshooting

### Release workflow failed

Check the GitHub Actions logs at:
https://github.com/plantonhq/planton/actions

### Version not incrementing correctly

The script only recognizes strict `vX.Y.Z` tags. Tags with suffixes (like `v1.0.0-beta`) are ignored.

```bash
# Check what tag will be used as base
git tag --list 'v*' --sort=-v:refname | head -5
```
