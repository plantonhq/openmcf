#!/usr/bin/env bash
set -euo pipefail

# Guard: every component README (apis/dev/planton/**/<component>/v1/README.md)
# and every infra chart README (charts/<provider>/<name>/README.md) must end
# with the canonical license footer.
#
# WHY THIS EXISTS
# The realistic unit of copying for this catalog is one component directory,
# not the whole repository. The repo-root LICENSE does not travel with a copied
# directory -- the component's README does. The footer is the attribution that
# travels: it states the license and links to it with an absolute URL so the
# pointer survives copies, zips, and forks. Components are forged continuously;
# this guard makes it impossible for a new one to ship without its attribution.
#
# SCOPE
# A last-line PRESENCE check (mirroring the sibling grep-based guards, no
# toolchain beyond find + awk): README files that are the DIRECT child of a
# component's v1/ folder, plus infra chart READMEs. Deeper READMEs (v1/docs/,
# v1/iac/**) are inner documentation whose shipped artifacts carry LICENSE and
# NOTICE files instead of footers.

repo_root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root_dir"

footer='© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).'

missing=()

while IFS= read -r readme; do
  [[ -z "$readme" ]] && continue
  # The last non-empty line must be exactly the canonical footer.
  last_line="$(awk 'NF {line=$0} END {print line}' "$readme")"
  if [[ "$last_line" != "$footer" ]]; then
    missing+=("$readme")
  fi
done < <(
  # Component READMEs: direct child of a v1/ folder anywhere under the APIs
  # tree (path-shaped, so future domains are covered without editing this
  # guard). The pattern cannot match v1/docs/README.md or v1/iac/**/README.md.
  find apis/dev/planton -path '*/v1/README.md' 2>/dev/null
  # Infra chart READMEs: charts/<provider>/<name>/README.md exactly.
  find charts -mindepth 3 -maxdepth 3 -name README.md 2>/dev/null
)

if (( ${#missing[@]} > 0 )); then
  {
    echo "License-footer guard FAILED: ${#missing[@]} README(s) missing the canonical footer."
    printf '  %s\n' "${missing[@]}"
    echo
    echo "Fix: end each file with a separator and the canonical footer:"
    echo
    echo "---"
    echo
    echo "${footer}"
  } >&2
  exit 1
fi

echo "License-footer guard passed: every component and chart README carries the canonical footer."
