#!/usr/bin/env bash
# =============================================================================
# Package Planton content for distribution via Cloudflare R2.
#
# Creates five zip files, each scoped to a single concern:
#
#   presets.zip        -- Preset YAML + MD files, kind enum proto
#   iac-source.zip     -- IaC source (.go, .tf, .md, .yaml under iac/)
#   catalog-pages.zip  -- Per-component catalog-page.md files
#   proto-source.zip   -- Raw proto source (spec, api, stack_input, stack_outputs)
#   reference-pack.zip -- The component reference pack: generated reference
#                         pages, catalog indexes, the cross-reference graph,
#                         the commons page, and the authored GUIDE.md /
#                         patterns wisdom layer
#
# All zips preserve repo-relative paths so they can be extracted into a single
# directory and overlay into a virtual Planton root. Consumers like the Planton
# upgrade scripts use this merged directory as --planton-path or PLANTON_ROOT.
#
# The version tag is accepted as an argument for logging purposes only; zip
# filenames are version-free because the version is encoded in the R2 path
# (releases/{tag}/content/{name}.zip).
#
# Usage:
#   bash tools/ci/release/package_content.sh v0.3.50
#   bash tools/ci/release/package_content.sh v0.3.50 --dry-run
# =============================================================================

set -euo pipefail

VERSION="${1:?Usage: package_content.sh <version-tag> [--dry-run]}"
DRY_RUN="${2:-}"

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"

PROVIDER_BASE="apis/dev/planton/provider"

if [ ! -d "$PROVIDER_BASE" ]; then
  echo "ERROR: Provider base directory not found: $PROVIDER_BASE"
  exit 1
fi

echo "=== Packaging Planton content for ${VERSION} ==="
echo ""

create_zip() {
  local zip_name="$1"
  local description="$2"
  shift 2

  local tmp_list
  tmp_list=$(mktemp)

  # Read file paths from stdin into a sorted temp file.
  sort > "$tmp_list"

  local count
  count=$(wc -l < "$tmp_list" | tr -d ' ')

  if [ "$count" -eq 0 ]; then
    echo "  WARNING: No files found for ${description}. Skipping ${zip_name}."
    rm -f "$tmp_list"
    return
  fi

  if [ "$DRY_RUN" = "--dry-run" ]; then
    echo "  [dry-run] ${zip_name}: ${count} files"
    rm -f "$tmp_list"
    return
  fi

  zip -q -@ "$zip_name" < "$tmp_list"
  rm -f "$tmp_list"

  local size
  size=$(du -h "$zip_name" | cut -f1)
  printf "  %-30s %6s  (%s files)\n" "$zip_name" "$size" "$count"
}

# ─── 1. Presets ───────────────────────────────────────────────────────────────
echo "1/5  Presets..."
{
  find "$PROVIDER_BASE" \( -path '*/v1/presets/*.yaml' -o -path '*/v1/presets/*.md' \)
  echo "apis/dev/planton/shared/cloudresourcekind/cloud_resource_kind.proto"
} | create_zip "presets.zip" "presets"

# ─── 2. IaC Source ────────────────────────────────────────────────────────────
# Mirrors the ALLOWED_EXTENSIONS in iac-bundler.ts: .go, .tf, .md, .yaml
# Excludes hidden dirs, vendor, and node_modules (same as iac-bundler.ts).
echo "2/5  IaC source..."
find "$PROVIDER_BASE" -path '*/v1/iac/*' \
    \( -name '*.go' -o -name '*.tf' -o -name '*.md' -o -name '*.yaml' \) \
    ! -path '*/vendor/*' \
    ! -path '*/node_modules/*' \
    ! -path '*/.terraform/*' \
    ! -path '*/.*' \
  | create_zip "iac-source.zip" "IaC source"

# ─── 3. Catalog Pages ────────────────────────────────────────────────────────
echo "3/5  Catalog pages..."
find "$PROVIDER_BASE" -path '*/v1/catalog-page.md' \
  | create_zip "catalog-pages.zip" "catalog pages"

# ─── 4. Proto Source ──────────────────────────────────────────────────────────
echo "4/5  Proto source..."
find "$PROVIDER_BASE" \( \
    -path '*/v1/spec.proto' \
    -o -path '*/v1/api.proto' \
    -o -path '*/v1/stack_input.proto' \
    -o -path '*/v1/stack_outputs.proto' \
  \) | create_zip "proto-source.zip" "proto source"

# ─── 5. Reference Pack ────────────────────────────────────────────────────────
# The pack is selected by file NAME, never by version-segment path: these
# names are the pack's frozen public contract (a contributor-edited file is
# the same file an agent reads -- no renames between source and artifact),
# pinned by the reference generator's contract tests. Name-based selection
# also survives api-version directory renames, which path patterns would not.
echo "5/5  Reference pack..."
find "$PROVIDER_BASE" \( \
    -name 'reference.md' \
    -o -name 'GUIDE.md' \
    -o -name 'reference-index.md' \
    -o -name 'reference-graph.yaml' \
    -o -name 'reference-commons.md' \
    -o -path "$PROVIDER_BASE/patterns/*.md" \
  \) | create_zip "reference-pack.zip" "reference pack"

echo ""
echo "=== Done ==="
