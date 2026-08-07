# KubernetesManifest Terraform module.
#
# Applies raw multi-document YAML through kubectl_manifest (one resource per
# document), with an optional anchor namespace created first.
#
# NAMESPACE SEMANTICS (parity contract with the Pulumi module): documents
# that declare no metadata.namespace get the anchor namespace via
# override_namespace; documents with an explicit namespace pass through
# untouched. Cluster-scoped documents without a namespace also receive the
# override — the API server ignores metadata.namespace on cluster-scoped
# objects, so the outcome matches the Pulumi provider's scope-aware
# defaulting (which skips them client-side instead).

# The optional anchor namespace. Created before any document applies;
# documents and namespace are deleted together on destroy.
resource "kubernetes_namespace_v1" "manifest_namespace" {
  count = local.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace_name
    labels = local.labels
  }
}

# One kubectl_manifest per document. kubectl_manifest applies server-side
# (same apply mechanism as the Pulumi provider, whose helper enables
# EnableServerSideApply) and needs no cluster connection at plan time, so
# offline plan proofs work and CRD+CR manifests apply in one pass.
resource "kubectl_manifest" "documents" {
  for_each = local.manifest_docs_by_identity

  yaml_body = each.value

  # Anchor-namespace defaulting: only for documents that declare none.
  override_namespace = try(yamldecode(each.value).metadata.namespace, "") == "" ? local.namespace_name : null

  server_side_apply = true
  force_conflicts   = true

  # skip_await parity: wait/wait_for_rollout here, SkipAwait on the Pulumi
  # ConfigGroup. Both engines default to awaiting: wait_for_rollout blocks
  # apply until workload rollouts complete, and wait blocks destroy until
  # each document is actually gone (foreground propagation) — matching the
  # Pulumi engine's apply/deletion awaits, so verify-clean phases never
  # race. (Await BREADTH differs benignly: the Pulumi engine also readiness-
  # checks non-workload kinds like Services; kubectl awaits workload
  # rollouts only. Behavioral difference documented, not a parity
  # exception — the applied objects are identical.)
  wait             = !try(var.spec.skip_await, false)
  wait_for_rollout = !try(var.spec.skip_await, false)

  depends_on = [kubernetes_namespace_v1.manifest_namespace]
}
