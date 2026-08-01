# The ray.io/v1 RayCluster CR — the single declaration the KubeRay
# operator (a KubernetesKubeRayOperator install, the prerequisite)
# reconciles into the head pod, the worker group pods, the
# `<name>-head-svc` Service every exported endpoint rides, and — in
# token auth mode without a bring-your-own Secret — the generated
# bearer-token Secret named exactly after this resource.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a RayCluster can be PLANNED before the
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its clusters in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls — Ray images are multi-GB — autoscaler sidecar
# injection, GCS startup) that is not part of applying the resource —
# the verifier owns readiness, the same never-block-on-a-controller
# posture as the sibling operator-CR modules.
resource "kubectl_manifest" "ray_cluster" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = "RayCluster"
    metadata = {
      name      = local.cluster_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.raycluster_spec
  })

  server_side_apply = true
  force_conflicts   = true

  # BACKGROUND deletion, explicitly: the OPERATOR owns the RayCluster's
  # cascade — its finalizer tears down the head and worker pods, the
  # head Service, and the generated token Secret. Foreground propagation
  # would block the delete on children the operator keeps reconciling.
  # Pulumi twin: the pulumi.com/deletionPropagationPolicy annotation.
  delete_cascade = "Background"

  lifecycle {
    # FAIL LOUDLY on names past the operator's naming budget: the
    # operator derives `<name>-head-svc` (9-character suffix) and
    # per-group worker pod names (`<name>-<group>-worker-<random>`) —
    # Kubernetes names cap at 63 characters, and 40 keeps every derived
    # name inside the budget with short group names. Twin: the Pulumi
    # module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 40
      error_message = "The KubeRay operator derives `<name>-head-svc` and per-group worker pod names (`<name>-<group>-worker-…`) from the resource name, and Kubernetes names cap at 63 characters — use a name of at most 40 characters (and keep worker group names short)."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
