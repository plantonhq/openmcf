# KubernetesCertificate Terraform module.
#
# Creates one cert-manager Certificate. The CR applies through
# kubectl_manifest (alekc/kubectl): no cluster connection at plan time, so a
# Certificate can be PLANNED before cert-manager's CRDs exist — single-run
# infra charts (cert-manager + issuer + certificate together) and offline
# plan proofs both depend on it.
#
# No wait_for block, deliberately: issuance time belongs to the issuer (an
# ACME order can take minutes; an unreachable CA would block forever) — the
# same never-block-on-a-controller posture as Ingress. Consumers that need
# the TLS Secret express the dependency through composition; the E2E lanes
# verify issuance by polling the live cluster. Pulumi equivalent: typed
# Certificate resource without await annotations.

resource "kubectl_manifest" "certificate" {
  yaml_body = yamlencode({
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = local.certificate_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.certificate_spec
  })

  server_side_apply = true
}
