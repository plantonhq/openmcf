# The flink.apache.org/v1beta1 FlinkDeployment CR — the single
# declaration the Flink Kubernetes Operator (a KubernetesFlinkOperator
# install, the prerequisite) reconciles into the JobManager, its
# TaskManagers, the `<name>-rest` Service and (in application mode) the
# job they run.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a FlinkDeployment can be PLANNED before the
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its clusters in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, job submission, TaskManager registration) that
# is not part of applying the resource — the verifier owns readiness,
# the same never-block-on-a-controller posture as the sibling
# operator-CR modules.
resource "kubectl_manifest" "flinkdeployment" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = "FlinkDeployment"
    metadata = {
      name      = local.flinkdeployment_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.flinkdeployment_spec
  })

  server_side_apply = true
  force_conflicts   = true

  # BACKGROUND deletion, explicitly: the OPERATOR owns the
  # FlinkDeployment's cascade — its finalizer tears down the JobManager,
  # TaskManagers and Services. Foreground propagation would block the
  # delete on children the operator keeps reconciling. Pulumi twin: the
  # pulumi.com/deletionPropagationPolicy annotation.
  delete_cascade = "Background"

  lifecycle {
    # FAIL LOUDLY on names past the operator's naming budget: the
    # operator derives `<name>-rest` and `<name>-taskmanager-N-M` child
    # names, and Kubernetes object names cap at 63 characters — past 45
    # the derived names silently break the contract the exported
    # outputs are built on. Twin: the Pulumi module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 45
      error_message = "The Flink operator derives child names from the resource name (`<name>-rest`, `<name>-taskmanager-N-M`) and Kubernetes names cap at 63 characters — use a name of at most 45 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
