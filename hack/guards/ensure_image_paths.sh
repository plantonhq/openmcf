#!/usr/bin/env bash
set -euo pipefail

# Guard: every Planton container image is named by one rule, and retired
# addresses never come back.
#
# WHY: Planton's images live at ghcr.io/plantonhq/planton/<slug>, where the
# slug names what the image is (control-plane, runner, operator, mcp,
# runner-tunnel, client-apps/web) and Planton OS images carry the os/ prefix
# (os/mcp, os/client-apps/web). Two earlier schemes -- repository-mirrored
# paths under planton/product/... and the mcp-os spelling -- once coexisted
# with this one, and the documents disagreed with the workflows about which
# image used which. A default the operator compiles in, a chart value, an
# example manifest, or a page that names a retired address sends an adopter to
# an image that will never be updated again. This guard holds every tracked
# file to the current addresses; dated records (_changelog/) are the one place
# the old ones may remain, because they describe the past.
#
# The platform repository holds the same table for its scripts and workflows
# (product/tools/ci/registry.sh) and runs the same check.

repo_root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root_dir"

# Retired address forms, matched literally. Each entry: pattern|replacement.
retired=(
  'ghcr.io/plantonhq/planton/product/client-apps/planton/web|ghcr.io/plantonhq/planton/client-apps/web'
  'ghcr.io/plantonhq/planton/product/client-apps/planton-os/web|ghcr.io/plantonhq/planton/os/client-apps/web'
  'ghcr.io/plantonhq/planton/product/|ghcr.io/plantonhq/planton/<slug> (no product/ prefix)'
  'ghcr.io/plantonhq/planton/mcp-os|ghcr.io/plantonhq/planton/os/mcp'
)

status=0
for entry in "${retired[@]}"; do
  pattern="${entry%%|*}"
  replacement="${entry#*|}"
  # Tracked files (what a commit would carry); dated records and this guard's
  # own table are excluded by path. git grep exits 1 when nothing matches and
  # >1 on a real error.
  set +e
  hits="$(git grep -n -F -- "$pattern" -- ':(exclude)_changelog' ':(exclude)hack/guards/ensure_image_paths.sh')"
  rc=$?
  set -e
  if [[ $rc -gt 1 ]]; then
    echo "ERROR: git grep failed (exit $rc) while checking '${pattern}'" >&2
    exit "$rc"
  fi
  if [[ -n "$hits" ]]; then
    status=1
    echo "ERROR: retired image address '${pattern}' is still referenced:" >&2
    echo "$hits" | sed 's/^/  /' >&2
    echo "  Use ${replacement} instead." >&2
    echo >&2
  fi
done

if [[ $status -ne 0 ]]; then
  echo "Image path guard FAILED. Every Planton image is ghcr.io/plantonhq/planton/<slug>; Planton OS images carry the os/ prefix." >&2
  exit 1
fi

echo "OK: no retired image address is referenced outside dated records."
