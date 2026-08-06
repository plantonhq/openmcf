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
# An optional zip name limits the run to that single artifact -- the release
# workflow builds each zip in its own job so every artifact is visible by
# name in the Actions UI. Without a target, all zips are built (local/dev
# behavior). Every guard (_test refusal, empty-selection failure) applies
# per zip either way.
#
# Usage:
#   bash tools/ci/release/package_content.sh <version-tag> [zip-name] [--dry-run]
#   bash tools/ci/release/package_content.sh v0.3.50
#   bash tools/ci/release/package_content.sh v0.3.50 --dry-run
#   bash tools/ci/release/package_content.sh v0.3.50 reference-pack
# =============================================================================

set -euo pipefail

VERSION="${1:?Usage: package_content.sh <version-tag> [zip-name] [--dry-run]}"
shift

TARGET=""
DRY_RUN=""
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN="--dry-run" ;;
    presets|iac-source|catalog-pages|proto-source|reference-pack) TARGET="$arg" ;;
    *)
      echo "ERROR: Unknown argument: $arg"
      echo "       Valid zip names: presets, iac-source, catalog-pages, proto-source, reference-pack"
      exit 1
      ;;
  esac
done

# True when the named zip should be built in this run.
wants() { [ -z "$TARGET" ] || [ "$TARGET" = "$1" ]; }

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

  # The _test provider is permanent internal test infrastructure ("never
  # shipped to users"). Refusing it HERE -- not only in the selectors --
  # means a future selector edit cannot quietly start shipping it.
  if grep -q "/_test/" "$tmp_list"; then
    echo "  ERROR: ${zip_name} selected _test provider content, which must"
    echo "         never ship to users:"
    grep "/_test/" "$tmp_list" | head -5 | sed 's/^/           /'
    rm -f "$tmp_list"
    exit 1
  fi

  if [ "$count" -eq 0 ]; then
    # An empty selection means a selector pattern no longer matches the tree
    # (e.g. an api-version directory rename) -- shipping a release with a
    # silently empty content zip is worse than failing the release.
    echo "  ERROR: No files found for ${description} (${zip_name}). A selector"
    echo "         pattern no longer matches the repository layout; fix the"
    echo "         pattern before releasing."
    rm -f "$tmp_list"
    exit 1
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

# ─── Presets ─────────────────────────────────────────────────────────────────
if wants presets; then
  echo "Presets..."
  {
    find "$PROVIDER_BASE" \( -path '*/v1alpha1/presets/*.yaml' -o -path '*/v1alpha1/presets/*.md' \) ! -path '*/_test/*'
    echo "apis/dev/planton/shared/cloudresourcekind/cloud_resource_kind.proto"
  } | create_zip "presets.zip" "presets"
fi

# ─── IaC Source ──────────────────────────────────────────────────────────────
# Mirrors the ALLOWED_EXTENSIONS in iac-bundler.ts: .go, .tf, .md, .yaml
# Excludes hidden dirs, vendor, and node_modules (same as iac-bundler.ts).
if wants iac-source; then
  echo "IaC source..."
  find "$PROVIDER_BASE" -path '*/v1alpha1/iac/*' ! -path '*/_test/*' \
      \( -name '*.go' -o -name '*.tf' -o -name '*.md' -o -name '*.yaml' \) \
      ! -path '*/vendor/*' \
      ! -path '*/node_modules/*' \
      ! -path '*/.terraform/*' \
      ! -path '*/.*' \
    | create_zip "iac-source.zip" "IaC source"
fi

# ─── Catalog Pages ───────────────────────────────────────────────────────────
if wants catalog-pages; then
  echo "Catalog pages..."
  find "$PROVIDER_BASE" -path '*/v1alpha1/catalog-page.md' ! -path '*/_test/*' \
    | create_zip "catalog-pages.zip" "catalog pages"
fi

# ─── Proto Source ────────────────────────────────────────────────────────────
if wants proto-source; then
  echo "Proto source..."
  find "$PROVIDER_BASE" \( \
      -path '*/v1alpha1/spec.proto' \
      -o -path '*/v1alpha1/api.proto' \
      -o -path '*/v1alpha1/stack_input.proto' \
      -o -path '*/v1alpha1/stack_outputs.proto' \
    \) ! -path '*/_test/*' | create_zip "proto-source.zip" "proto source"
fi

# ─── Reference Pack ──────────────────────────────────────────────────────────
# The pack is selected by file NAME, never by version-segment path: these
# names are the pack's frozen public contract (a contributor-edited file is
# the same file an agent reads -- no renames between source and artifact),
# pinned by the reference generator's contract tests. Name-based selection
# also survives api-version directory renames, which path patterns would not.
if wants reference-pack; then
  echo "Reference pack..."
  find "$PROVIDER_BASE" \( \
      -name 'reference.md' \
      -o -name 'GUIDE.md' \
      -o -name 'reference-index.md' \
      -o -name 'reference-graph.yaml' \
      -o -name 'reference-commons.md' \
      -o -path "$PROVIDER_BASE/patterns/*.md" \
    \) ! -path '*/_test/*' | create_zip "reference-pack.zip" "reference pack"
fi

echo ""
echo "=== Done ==="
