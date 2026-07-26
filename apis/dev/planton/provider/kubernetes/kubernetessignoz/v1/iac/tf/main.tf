# KubernetesSignoz Terraform module.
#
# Installs SigNoz from the official Helm chart as a real Helm release.
# The typed spec renders into chart values (locals.helm_values); the
# telemetry store is a COMPOSED ClickHouse connection (the bundled
# subchart stays permanently off — nothing ClickHouse-related installs,
# and the release carries no operator and no CRDs, so uninstall is
# ordinary object deletion); the ClickHouse password reaches the server
# as a secretKeyRef (never a rendered value); the helm_values escape
# hatch is passed as a SECOND values document, with fullnameOverride
# re-pinned in a THIRD — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps + re-pin.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "signoz" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "signoz" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the whole platform to become Ready — a SigNoz whose schema
  # migration fails against the composed ClickHouse should fail THIS
  # apply, not the first trace query. The migrator's own `migrate ready`
  # init container blocks until the telemetry store answers, so the
  # budget also absorbs a ClickHouse that is still coming up.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 1200

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second, and the
  # fullname pin re-asserted LAST — the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name — and the exported outputs built from them — derives from
  # this fullname; letting an override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  # The collector Deployment (`<name>-otel-collector`) is the longest
  # fullname-derived child; its pod names must fit Kubernetes'
  # 63-character cap — an over-long resource name corrupts the naming
  # contract the outputs promise. Fail THIS plan loudly instead (twin:
  # the Pulumi module's MaxNameLength guard).
  lifecycle {
    precondition {
      condition     = local.name_within_budget
      error_message = "metadata.name '${local.release_name}' is ${length(local.release_name)} characters — the signoz chart's child-name budget allows at most ${local.max_name_length} (the collector composes names like <name>-otel-collector-<replicaset>-<pod> within Kubernetes' 63-character cap)."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.signoz,
  ]
}
