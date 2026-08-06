#!/usr/bin/env bash
# =============================================================================
# Module change detection for auto-tagging.
#
# Reads changed file paths on stdin and emits the unique IaC module
# directories they belong to, one per line:
#
#   --pulumi     apis/dev/planton/provider/{p}/{k}/{version}/iac/pulumi
#   --terraform  apis/dev/planton/provider/{p}/{k}/{version}/iac/tf
#
# The version segment matches the full maturity grammar
# (v<N>[alpha<N>|beta<N>]), never a bare v<N> assumption: a pattern that
# matched only "v1" once made module auto-tagging silently stop firing after
# an api-version directory rename -- no error, releases just stopped. This
# script is that lesson made structural: auto-tag.yaml consumes it instead of
# carrying inline patterns, and the --self-test fixtures prove the detection
# fires for every channel the grammar allows.
#
# Usage:
#   git diff --name-only A B | detect_module_dirs.sh --pulumi
#   detect_module_dirs.sh --self-test    (used by auto-tag before detection)
# =============================================================================
set -euo pipefail

VERSION_ERE='v[0-9]+((alpha|beta)[0-9]+)?'

detect() {
  # detect <flavor-dir>  (pulumi|tf), paths on stdin
  local flavor="$1"
  # NOTE: the sed delimiter must not be '|' -- the grammar's alternation
  # pipes inside (alpha|beta) would split the expression.
  grep -E "^apis/dev/planton/provider/[^/]+/[^/]+/${VERSION_ERE}/iac/${flavor}/" \
    | sed -E "s#(apis/dev/planton/provider/[^/]+/[^/]+/${VERSION_ERE}/iac/${flavor})/.*#\1#" \
    | sort -u \
    || true
}

self_test() {
  local fail=0
  check() {
    # check <flavor> <input-path> <expected-output ('' for no match)>
    local got
    got=$(printf '%s\n' "$2" | detect "$1")
    if [ "$got" != "$3" ]; then
      echo "SELF-TEST FAIL: detect --${1} '${2}' = '${got}', want '${3}'"
      fail=1
    fi
  }
  # Every maturity channel must fire.
  check tf "apis/dev/planton/provider/aws/awsvpc/v1alpha1/iac/tf/main.tf" \
           "apis/dev/planton/provider/aws/awsvpc/v1alpha1/iac/tf"
  check tf "apis/dev/planton/provider/aws/awsvpc/v1beta1/iac/tf/main.tf" \
           "apis/dev/planton/provider/aws/awsvpc/v1beta1/iac/tf"
  check tf "apis/dev/planton/provider/aws/awsvpc/v1/iac/tf/main.tf" \
           "apis/dev/planton/provider/aws/awsvpc/v1/iac/tf"
  check pulumi "apis/dev/planton/provider/gcp/gcpgkecluster/v1alpha1/iac/pulumi/main.go" \
               "apis/dev/planton/provider/gcp/gcpgkecluster/v1alpha1/iac/pulumi"
  # Non-module changes must NOT fire.
  check tf "apis/dev/planton/provider/aws/awsvpc/v1alpha1/spec.proto" ""
  check tf "apis/dev/planton/provider/aws/awsvpc/v1alpha1/iac/pulumi/main.go" ""
  check pulumi "pkg/iac/pulumi/pulumimodule/module_directory.go" ""
  # Two files in one module collapse to one dir.
  local got
  got=$(printf '%s\n%s\n' \
    "apis/dev/planton/provider/aws/awsvpc/v1alpha1/iac/tf/main.tf" \
    "apis/dev/planton/provider/aws/awsvpc/v1alpha1/iac/tf/outputs.tf" | detect tf | wc -l | tr -d ' ')
  if [ "$got" != "1" ]; then
    echo "SELF-TEST FAIL: two files in one module produced ${got} dirs, want 1"
    fail=1
  fi
  if [ $fail -eq 0 ]; then
    echo "self-test: module change detection fires for every grammar channel"
  fi
  return $fail
}

case "${1:-}" in
  --pulumi) detect pulumi ;;
  --terraform) detect tf ;;
  --self-test) self_test ;;
  *)
    echo "usage: $0 --pulumi|--terraform|--self-test (paths on stdin)" >&2
    exit 2
    ;;
esac
