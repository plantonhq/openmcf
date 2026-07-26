# KubernetesTektonOperator Terraform module.
#
# Installs the Tekton Operator from its released single-file manifest —
# the operator's OFFICIAL distribution (the in-repo Helm chart is
# unpublished, version "devel"). The manifest applies per document: the
# namespace `tekton-operator`, the 14 operator.tekton.dev CRDs, the
# operator and webhook Deployments (with the spec's typed overrides
# patched on — see locals.tf), ConfigMaps, RBAC, Services and the webhook
# cert Secret.
#
# AUTO-INSTALL IS ALWAYS DISABLED: the tekton-config-defaults ConfigMap's
# AUTOINSTALL_COMPONENTS key is patched from the release's "true" to
# "false" so the operator never creates its own TektonConfig — the
# KubernetesTekton declaration kind is the single owner of the cluster's
# Tekton configuration. Installing the operator alone deploys no Tekton
# components.
#
# DESTROY SEMANTICS: every document deletes with the resource, INCLUDING
# the CRDs — which cascade-deletes any TektonConfig on the cluster.
# Always destroy the KubernetesTekton resource FIRST while the operator
# still runs (TektonInstallerSet finalizers are operator-processed; see
# the spec's destroy note).

# Fetch the released manifest from the pinned release tag.
data "http" "tekton_operator_manifest" {
  url = local.manifest_url

  request_headers = {
    Accept = "application/yaml"
  }
}

# Apply every document. Keyed by each document's COMPOSED IDENTITY
# (apiVersion//kind//name[//namespace]) — stable across manifest
# reorderings, and exactly the ID form the kubectl importer takes so the
# import map derives IDs blind from the address keys.
#
# server_side_apply keeps re-installs tolerant of the operator's own
# field management on its ConfigMaps and matches the Pulumi twin's apply
# mode.
resource "kubectl_manifest" "tekton_operator" {
  for_each = local.applied_documents

  yaml_body = yamlencode(each.value)

  server_side_apply = true
  force_conflicts   = true
}
