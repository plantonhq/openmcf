# KubernetesClusterSecretStore Terraform module.
#
# Creates one External Secrets Operator ClusterSecretStore plus the
# credential Secret its backend needs — in the spec's secrets namespace,
# because cluster-scoped stores read their referenced Secrets from an
# EXPLICIT namespace.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a ClusterSecretStore can be PLANNED before the
# External Secrets Operator's CRDs exist, which is what lets an infra
# chart deploy the operator and its stores in one run (and lets offline
# plan proofs work).
#
# No wait_for block, deliberately: store readiness depends on external
# reachability (the cloud secrets API, Vault) that is not part of applying
# the resource — the same never-block-on-a-controller posture as the
# cert-manager issuers. Pulumi equivalent: CustomResource without await
# annotations.

# Credential Secret first; the CR depends on it so ESO never observes a
# store whose secretRefs dangle.
resource "kubernetes_secret_v1" "credentials" {
  for_each = local.credential_secrets

  metadata {
    name      = each.key
    namespace = local.secrets_namespace
    labels    = local.labels
  }

  # The provider's `data` argument takes plaintext values (marked sensitive
  # in state) and handles the base64 encoding — the Secret lands identical
  # to the Pulumi module's stringData.
  data = each.value
}

resource "kubectl_manifest" "cluster_secret_store" {
  yaml_body = yamlencode({
    apiVersion = "external-secrets.io/v1"
    kind       = "ClusterSecretStore"
    metadata = {
      name   = local.store_name
      labels = local.labels
    }
    spec = local.store_spec
  })

  server_side_apply = true

  depends_on = [kubernetes_secret_v1.credentials]
}
