#!/usr/bin/env bash
set -euo pipefail

# Guard: a kind's two IaC engines must declare the SAME resource-identity tag
# key set, in the canonical shape.
#
# The defect class this retires (caught live on AwsElasticIp, 2026-08): the
# same manifest produced differently-tagged AWS resources depending on which
# engine deployed it -- one engine's identity map carried keys the other's
# never sent (a missing Name, a metadata.labels merge on one side only). Tag
# content is live send surface: cost-allocation reports, orphan-cleanup
# queries, and console searches all diverge when the engines disagree.
#
# WHAT IT CHECKS (catalog/aws; other providers carry their own tag idioms)
# For each kind directory `catalog/aws/<kind>/`:
#   1. Extract the identity-tag key set from `iac/tf/locals.tf` (the literal
#      "Name" / "planton.ai/*" keys) and from `iac/pulumi/module/locals.go`
#      (the awstagkeys.* constants). Identity tags live ONLY in those files --
#      that is the prescribed shape, which is what makes this check possible.
#   2. A kind whose resources are untaggable declares NO identity keys on
#      EITHER engine (with a comment saying why) -- one-engine-only
#      declarations fail.
#   3. Both key sets must equal the canonical six-key identity map (Name +
#      the five planton.ai/* keys), or the kind's recorded exception below.
#   4. The Pulumi map must use the shared awstagkeys constants -- literal key
#      strings and Name entries hardcoded at resource send sites hide from
#      this extraction and from every future sweep.
#   5. `metadata.labels` never merges into identity tags (the canonical
#      decision: user labels are platform metadata, not cloud tags).
#   6. Canonical e2e manifests carry no legacy `pulumi.planton.dev/*`
#      annotations (they hijack preview/stack resolution).
#
# BOUNDARY: this guard compares DECLARED key sets, not rendered sends -- a
# module that declares the map and then drops it at a send site is the proof
# lane's catch. Declared-set equality is the class that shipped; this pins it.

repo_root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root_dir"

provider_root="catalog/aws"

# The canonical identity key set (session-recorded catalog decision).
canonical_keys="Name
planton.ai/environment
planton.ai/organization
planton.ai/resource
planton.ai/resource-id
planton.ai/resource-kind"

# Five-key convention set: kinds whose ONE primary resource carries its own
# explicit name argument drop the Name tag by recorded decision -- a shared
# Name tag would restate (or, on multi-resource appliances, mislabel) what
# the resource already declares.
no_name_keys="planton.ai/environment
planton.ai/organization
planton.ai/resource
planton.ai/resource-id
planton.ai/resource-kind"

# Recorded exceptions: "<kind>|<set>|<reason>" where <set> names a variable
# above. Add a row ONLY with a recorded justification in the kind's module
# comments; never to silence a finding.
exceptions=(
  "awsecrrepo|no_name_keys|repository carries its own name (recorded 2026-08)"
  "awsglobalaccelerator|no_name_keys|accelerator carries its own name (recorded 2026-08)"
  "awssagemakerdomain|no_name_keys|domain carries its own domain_name (recorded 2026-08)"
  "awsplantonrunner|no_name_keys|internal appliance; every resource explicitly named (recorded 2026-08)"
  "awscodepipeline|no_name_keys|pipeline carries its own name (recorded 2026-08)"
  "awscodebuildproject|no_name_keys|project carries its own name (recorded 2026-08)"
)

# Labels-merge exceptions: kinds where metadata.labels DELIBERATELY reach
# cloud tags, on BOTH engines, with the reason documented in both modules.
labels_exceptions=(
  "awssubnet|labels are the chart-side discovery-tag mechanism (karpenter.sh/discovery on subnets; documented in both modules)"
)

# Maps awstagkeys constant names to the literal keys they render as
# (pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys).
constant_to_key() {
  case "$1" in
    Name) echo "Name" ;;
    Resource) echo "planton.ai/resource" ;;
    Organization) echo "planton.ai/organization" ;;
    Environment) echo "planton.ai/environment" ;;
    ResourceKind) echo "planton.ai/resource-kind" ;;
    ResourceId) echo "planton.ai/resource-id" ;;
  esac
}

# TF identity keys: the quoted keys of the identity map in locals.tf. "Name"
# counts only when the file carries planton.ai keys too, so an unrelated Name
# usage in a tag-free module never trips the guard.
extract_tf_keys() {
  local f="$1"
  [[ -f "$f" ]] || return 0
  grep -qE '"planton\.ai/resource"' "$f" || return 0
  grep -hoE '"(Name|planton\.ai/[a-z-]+)"[[:space:]]*=' "$f" 2>/dev/null \
    | sed -E 's/^"//; s/"[[:space:]]*=$//' | sort -u || true
}

# Pulumi identity keys: the awstagkeys constants referenced in locals.go.
extract_pulumi_keys() {
  local f="$1"
  [[ -f "$f" ]] || return 0
  { grep -hoE 'awstagkeys\.(Name|Resource|Organization|Environment|ResourceKind|ResourceId)\b' "$f" 2>/dev/null || true; } \
    | sed -E 's/^awstagkeys\.//' | sort -u \
    | while IFS= read -r c; do constant_to_key "$c"; done | sort -u
}

set_violations=()
shape_violations=()
labels_violations=()
annotation_violations=()

while IFS= read -r kinddir; do
  kind="${kinddir#${provider_root}/}"
  tf_locals="${kinddir}/iac/tf/locals.tf"
  pulumi_locals="${kinddir}/iac/pulumi/module/locals.go"
  [[ -d "${kinddir}/iac" ]] || continue

  tf_keys="$(extract_tf_keys "$tf_locals")"
  pulumi_keys="$(extract_pulumi_keys "$pulumi_locals")"

  # Prescribed-shape checks on the Pulumi module (any module file).
  if compgen -G "${kinddir}/iac/pulumi/module/*.go" > /dev/null; then
    if grep -lqE '"planton\.ai/' "${kinddir}"/iac/pulumi/module/*.go 2>/dev/null; then
      shape_violations+=("${kind} -> literal planton.ai/* key strings in the Pulumi module; use the awstagkeys constants (locals.go)")
    fi
    if grep -lqE 'pulumi\.String\("Name"\)' "${kinddir}"/iac/pulumi/module/*.go 2>/dev/null; then
      shape_violations+=("${kind} -> Name tag hardcoded at a resource send site; declare it in the identity map in locals.go")
    fi
  fi

  # Both-or-neither, then set equality against canonical/exception.
  if [[ -z "$tf_keys" && -z "$pulumi_keys" ]]; then
    continue # untaggable-resource kind: no identity block on either engine
  fi
  if [[ -z "$tf_keys" || -z "$pulumi_keys" ]]; then
    set_violations+=("${kind} -> identity tags declared by ONE engine only (tf: $(echo $tf_keys | tr '\n' ' ' | sed 's/ $//'); pulumi: $(echo $pulumi_keys | tr '\n' ' ' | sed 's/ $//'))")
    continue
  fi

  expected="$canonical_keys"
  for rule in "${exceptions[@]}"; do
    rule_kind="${rule%%|*}"
    rest="${rule#*|}"
    rule_set="${rest%%|*}"
    if [[ "$kind" == "$rule_kind" ]]; then
      expected="${!rule_set}"
    fi
  done

  if [[ "$tf_keys" != "$expected" ]]; then
    set_violations+=("${kind} -> tf identity keys diverge from the expected set: [$(echo $tf_keys | tr '\n' ' ')] vs [$(echo $expected | tr '\n' ' ')]")
  fi
  if [[ "$pulumi_keys" != "$expected" ]]; then
    set_violations+=("${kind} -> pulumi identity keys diverge from the expected set: [$(echo $pulumi_keys | tr '\n' ' ')] vs [$(echo $expected | tr '\n' ' ')]")
  fi

  # The labels-merge class: user labels never reach cloud tags (except the
  # recorded rows above, which must hold on BOTH engines). Detection covers
  # the single-line and multi-line merge() forms and the Pulumi copy loop.
  labels_excepted=0
  for rule in "${labels_exceptions[@]}"; do
    [[ "$kind" == "${rule%%|*}" ]] && labels_excepted=1
  done
  tf_merges=0
  pulumi_merges=0
  if [[ -f "$tf_locals" ]] && grep -qE 'var\.metadata\.labels' "$tf_locals"; then
    tf_merges=1
  fi
  if [[ -f "$pulumi_locals" ]] && grep -qiE 'metadata\.Labels' "$pulumi_locals"; then
    pulumi_merges=1
  fi
  if [[ $labels_excepted -eq 1 ]]; then
    if [[ $tf_merges -ne 1 || $pulumi_merges -ne 1 ]]; then
      labels_violations+=("${kind} -> recorded labels-merge exception must hold on BOTH engines (tf: ${tf_merges}, pulumi: ${pulumi_merges})")
    fi
  else
    if [[ $tf_merges -eq 1 ]]; then
      labels_violations+=("${kind} -> metadata.labels merged into identity tags (iac/tf/locals.tf)")
    fi
    if [[ $pulumi_merges -eq 1 ]]; then
      labels_violations+=("${kind} -> metadata labels copied into the identity map (iac/pulumi/module/locals.go)")
    fi
  fi
done < <(find "$provider_root" -maxdepth 1 -mindepth 1 -type d | sort)

# Legacy annotation rot in canonical e2e manifests.
while IFS= read -r m; do
  annotation_violations+=("${m} -> legacy pulumi.planton.dev/* annotation (hijacks preview/stack resolution)")
done < <(grep -lE 'pulumi\.planton\.dev/' "$provider_root"/*/e2e/manifest.yaml 2>/dev/null || true)

failed=0

if [[ ${#set_violations[@]} -gt 0 ]]; then
  echo "ERROR: ${#set_violations[@]} kind(s) with divergent cross-engine identity-tag key sets." >&2
  echo "The same manifest must leave the same tag fingerprint through either engine; converge the" >&2
  echo "identity maps (or record a justified exception in this guard's table -- never silently):" >&2
  printf '  - %s\n' "${set_violations[@]}" >&2
  echo >&2
  failed=1
fi

if [[ ${#shape_violations[@]} -gt 0 ]]; then
  echo "ERROR: ${#shape_violations[@]} Pulumi module(s) declare identity tags outside the prescribed shape." >&2
  echo "Identity tags live in locals.go via the awstagkeys constants -- literal keys and send-site" >&2
  echo "Name entries hide from this guard's extraction and from every future sweep:" >&2
  printf '  - %s\n' "${shape_violations[@]}" >&2
  echo >&2
  failed=1
fi

if [[ ${#labels_violations[@]} -gt 0 ]]; then
  echo "ERROR: ${#labels_violations[@]} module(s) merge metadata labels into cloud tags." >&2
  echo "The canonical identity map is exactly six keys; user labels are platform metadata:" >&2
  printf '  - %s\n' "${labels_violations[@]}" >&2
  echo >&2
  failed=1
fi

if [[ ${#annotation_violations[@]} -gt 0 ]]; then
  echo "ERROR: ${#annotation_violations[@]} canonical e2e manifest(s) carry legacy pulumi.planton.dev annotations." >&2
  printf '  - %s\n' "${annotation_violations[@]}" >&2
  echo >&2
  failed=1
fi

if [[ $failed -ne 0 ]]; then
  echo "Cross-engine tag-key guard FAILED." >&2
  exit 1
fi

echo "Cross-engine tag-key guard passed: every kind's engines declare the same identity-tag key set in the prescribed shape."
