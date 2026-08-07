# KubernetesIssuer Terraform module.
#
# Creates one cert-manager Issuer (namespace-scoped) plus the credential
# Secrets its configuration needs — in the Issuer's OWN namespace, the only
# namespace a namespace-scoped issuer reads Secrets from.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — an Issuer can be PLANNED before cert-manager's
# CRDs exist, which is what lets an infra chart deploy cert-manager and its
# issuers in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: issuer readiness depends on external
# reachability (the ACME server, Vault, DNS) that is not part of applying
# the resource — the same never-block-on-a-controller posture as Ingress.
# Pulumi equivalent: CustomResource without await annotations.

# Credential Secrets first; the CR depends on them so cert-manager never
# observes an issuer whose secretRefs dangle.
resource "kubernetes_secret_v1" "credentials" {
  for_each = local.credential_secrets

  metadata {
    name      = each.key
    namespace = local.namespace
    labels    = local.labels
  }

  # The provider's `data` argument takes plaintext values (marked sensitive
  # in state) and handles the base64 encoding — the Secret lands identical
  # to the Pulumi module's stringData.
  data = each.value
}

resource "kubectl_manifest" "issuer" {
  yaml_body = yamlencode({
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = local.issuer_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.issuer_spec
  })

  server_side_apply = true

  depends_on = [kubernetes_secret_v1.credentials]
}
