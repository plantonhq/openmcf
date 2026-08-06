# KubernetesTemporal Terraform module.
#
# Installs Temporal from the official Helm chart as a real Helm release.
# The typed spec renders into chart values (locals.helm_values); database
# passwords ride the chart's existingSecret contract (a secretKeyRef the
# server and schema Jobs resolve at runtime — never rendered as values);
# the helm_values escape hatch is passed as a SECOND values document,
# which the provider merges over the first with Helm -f semantics — the
# exact semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "temporal" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "temporal" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the cluster to become Ready — a server that never starts
  # (empty schema, unreachable database, wrong credential Secret name)
  # should fail THIS apply, not the first workflow. The pre-install
  # schema Jobs run inside this budget too (Helm hooks execute before
  # the release resources).
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 900

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name (`<name>-frontend`, `<name>-web`, ...) — and the exported
  # outputs built from them — derive from the fullname; letting an
  # override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  lifecycle {
    # FAIL LOUDLY on names past the chart's fullname budget: child names
    # are `<fullname>-<component>` and the chart's componentname helper
    # truncates the FULLNAME to fit 63 characters — the longest
    # component ("internal-frontend", 17 chars) silently truncates the
    # fullname past 45 and breaks the naming contract the exported
    # outputs are built on. Twin: the Pulumi module's Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 45
      error_message = "The temporal chart derives child names from the resource name and silently truncates past 45 characters (62 minus its longest component suffix), which would break the naming contract — use a name of at most 45 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.temporal,
  ]
}
