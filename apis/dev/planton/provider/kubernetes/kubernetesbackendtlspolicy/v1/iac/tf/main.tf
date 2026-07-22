# KubernetesBackendTlsPolicy Terraform module.
#
# Applies one Gateway API BackendTLSPolicy custom resource. The CR spec arrives
# from the proto->tfvars converter already manifest-shaped (camelCase keys,
# null-pruned, StringValueOrRef foreign keys resolved to literal strings), so
# this module hands it to the engine verbatim — no snake->camel, null-prune, or
# oneof logic. The apiserver plus Planton protovalidate are the schema
# authority.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a BackendTLSPolicy can be PLANNED before the
# Gateway API CRDs exist, which is what lets an infra chart deploy the CRDs, a
# Gateway, and its policies in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: the per-ancestor Accepted/ResolvedRefs
# conditions appear when a Gateway controller reconciles the policy, which is
# not part of applying it — the same never-block-on-a-controller posture as
# KubernetesIngress. Pulumi equivalent: CustomResource without await
# annotations.

resource "kubectl_manifest" "backend_tls_policy" {
  yaml_body = yamlencode({
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "BackendTLSPolicy"
    metadata = {
      name      = var.metadata.name
      namespace = var.spec.namespace
      labels    = local.labels
    }
    spec = local.manifest_spec
  })

  server_side_apply = true
}
