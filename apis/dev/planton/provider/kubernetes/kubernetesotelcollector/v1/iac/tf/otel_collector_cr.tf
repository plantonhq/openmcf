# The opentelemetry.io/v1beta1 OpenTelemetryCollector CR — the single
# declaration the OpenTelemetry Operator (a KubernetesOtelOperator
# install, the prerequisite) reconciles into the collector workload
# (Deployment/DaemonSet/StatefulSet per mode, or sidecar registration),
# the `<name>-collector` Service with receiver-derived ports, the
# headless and monitoring Services, and the rendered config ConfigMap.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a collector can be PLANNED before the
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its collectors in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: collector readiness depends on the
# operator (webhook admission, image injection, workload rollout) that is
# not part of applying the resource — the verifier owns readiness, the
# same never-block-on-a-controller posture as the sibling operator-CR
# modules.
resource "kubectl_manifest" "otel_collector" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = local.kind
    metadata = {
      name      = local.resource_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.collector_spec
  })

  server_side_apply = true
  force_conflicts   = true

  # BACKGROUND deletion, explicitly: the OPERATOR owns the collector CR's
  # cascade — its ownership references tear down the workload, Services
  # and ConfigMap. Foreground propagation would block the delete on
  # children the operator keeps reconciling (WA: deleting an
  # operator-owned CR does not cascade cleanly under foreground). Pulumi
  # twin: the pulumi.com/deletionPropagationPolicy annotation.
  delete_cascade = "Background"

  lifecycle {
    # FAIL LOUDLY on names past the operator's naming budget: the
    # operator derives child names by suffixing —
    # "-collector-networkpolicy" is the longest at 25 characters
    # (feature-gated, but an operator-side gate flip must never break
    # existing collector names) — verified in the operator's naming
    # source at the pin; Kubernetes caps names at 63. Twin: the Pulumi
    # module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 38
      error_message = "metadata.name must be 38 characters or fewer: the operator derives \"<name>-collector-networkpolicy\" (25-char suffix) and Kubernetes caps names at 63."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
