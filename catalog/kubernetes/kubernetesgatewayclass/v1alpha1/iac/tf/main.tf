# KubernetesGatewayClass Terraform module.
#
# Applies one Gateway API GatewayClass custom resource (cluster-scoped). The CR spec
# arrives from the proto->tfvars converter already manifest-shaped (camelCase
# keys, null-pruned), so this module hands it to the engine verbatim — no
# snake->camel, null-prune, or oneof logic. The apiserver plus Planton
# protovalidate are the schema authority.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a GatewayClass can be PLANNED before the Gateway API
# CRDs exist, which is what lets an infra chart deploy the CRDs and the
# class in one run (and lets offline plan proofs work).
#
# No wait_for block, deliberately: the Accepted condition appears when the
# named controller reconciles the class, which is not part of applying it —
# the same never-block-on-a-controller posture as KubernetesIngress.

resource "kubectl_manifest" "gateway_class" {
  yaml_body = yamlencode({
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "GatewayClass"
    metadata = {
      name   = var.metadata.name
      labels = local.labels
    }
    spec = local.manifest_spec
  })

  server_side_apply = true
}
