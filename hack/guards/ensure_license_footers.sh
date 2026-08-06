#!/usr/bin/env bash
set -euo pipefail

# Guard: every component README (apis/dev/planton/**/<component>/v1/README.md),
# every infra chart README (the README.md beside a chart's Chart.yaml, at any
# nesting depth), and every helm chart README (helm/<chart>/README.md) must end
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
  find apis/dev/planton -path '*/v*/README.md' 2>/dev/null
  # Helm chart READMEs: helm/<chart>/README.md — published distribution surfaces.
  find helm -mindepth 2 -maxdepth 2 -name README.md 2>/dev/null
)

# Infra chart READMEs: the README beside each Chart.yaml, wherever the chart
# nests in the tree (charts/<provider>/<name>/, charts/<cloud>/kubernetes/<name>/,
# ...). Discovery keys off Chart.yaml presence exactly like the structure guard
# and every chart loader, so nesting depth is never a blind spot. A chart with
# no README at all fails too: the README is part of the chart's anatomy and the
# footer is the attribution that travels with a copied chart directory.
while IFS= read -r chart_yaml; do
  [[ -z "$chart_yaml" ]] && continue
  readme="$(dirname "$chart_yaml")/README.md"
  if [[ ! -f "$readme" ]]; then
    missing+=("$readme (README.md missing entirely)")
    continue
  fi
  last_line="$(awk 'NF {line=$0} END {print line}' "$readme")"
  if [[ "$last_line" != "$footer" ]]; then
    missing+=("$readme")
  fi
done < <(find charts -type f -name 'Chart.yaml' -not -path '*/build/*' 2>/dev/null)

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
