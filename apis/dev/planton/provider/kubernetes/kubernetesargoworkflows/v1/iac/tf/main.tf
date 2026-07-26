# KubernetesArgoWorkflows Terraform module.
#
# Installs Argo Workflows from the official Helm chart as a real Helm
# release. The typed spec renders into chart values (locals.helm_values);
# artifact-store and archive-database credentials ride the chart's own
# secret SELECTORS (resolved by the workloads at runtime — never rendered
# as values); the helm_values escape hatch is passed as a SECOND values
# document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "argo_workflows" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "argo_workflows" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the engine to become Ready — a controller that never starts
  # (bad archive credentials Secret name, unreachable database) should
  # fail THIS apply, not the first workflow submission. The budget covers
  # the full-schema CRD hook Job's download-and-apply on the default arm.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name (`<name>-server`, `<name>-workflow-controller`, ...) — and
  # the exported outputs built from them — derive from the fullname;
  # letting an override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  lifecycle {
    # FAIL LOUDLY on names past the chart's fullname budget: every child
    # name is `<fullname>-<component>` truncated at 63 characters — the
    # longest component suffix ("-workflow-controller", 20 chars) would
    # truncate SILENTLY past 43, breaking the naming contract the
    # exported outputs are built on. Twin: the Pulumi module's
    # Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 43
      error_message = "The argo-workflows chart derives child names from the resource name and silently truncates past 43 characters (63 minus its longest component suffix), which would break the naming contract — use a name of at most 43 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.argo_workflows,
  ]
}
