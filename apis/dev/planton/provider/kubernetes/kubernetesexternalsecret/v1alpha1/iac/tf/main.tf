# KubernetesExternalSecret Terraform module.
#
# Creates one External Secrets Operator ExternalSecret — the declaration
# that syncs entries from a store's backend into a materialized Kubernetes
# Secret in this namespace. No credential Secrets: authentication belongs
# to the store (KubernetesSecretStore / KubernetesClusterSecretStore),
# this resource only picks WHAT to sync.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — an ExternalSecret can be PLANNED before the
# External Secrets Operator's CRDs exist, which is what lets an infra
# chart deploy the operator, its stores, and their secrets in one run (and
# lets offline plan proofs work).
#
# No wait_for block, deliberately: the materialized Secret appears when
# the operator reaches the backend, which is not part of applying the
# resource — the same never-block-on-a-controller posture as the store
# kinds. Pulumi equivalent: CustomResource without await annotations.

resource "kubectl_manifest" "external_secret" {
  yaml_body = yamlencode({
    apiVersion = "external-secrets.io/v1"
    kind       = "ExternalSecret"
    metadata = {
      name      = local.external_secret_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.external_secret_spec
  })

  server_side_apply = true
}
