# KubernetesOpenFga Terraform module — the Helm release.
#
# Installs OpenFGA from the official chart as a real Helm release. The
# typed spec renders into chart values (locals.typed_helm_values);
# declared pre-shared API keys materialize into a module-owned Secret
# before the release (authn_secret.tf) so key material never rides
# rendered values; the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.
#
# ESCAPE-HATCH REALITY (verified at chart 0.3.10): the chart ships a
# CLOSED values schema (values.schema.json additionalProperties: false)
# — a key the chart does not define fails the install outright, so
# helm_values can only override EXISTING chart values (extraEnvVars for
# the ~50 server flags without values paths, TLS file wiring, sidecars),
# never invent new ones.

resource "helm_release" "openfga" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the server to become Ready. Safe here ONLY because
  # migrations run as an init container (locals.datastore_block): the
  # chart's default hook-Job mode would deadlock this wait — the
  # Deployment's wait-for-migration init container waits on a
  # post-install hook Job that Helm only runs AFTER --wait returns.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). The
  # Service / `-migrate` Job / `-datastore-secret` names — and the
  # exported endpoints built from them — all derive from the fullname;
  # letting an override move it would break every output.
  values = concat(
    [yamlencode(local.typed_helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.openfga,
    kubernetes_secret_v1.authn_keys,
  ]

  lifecycle {
    # NAME BUDGET: the chart truncates its fullname at 63 characters
    # THEN derives `<fullname>-migrate` for the migration Job, whose pod
    # label value also caps at 63 — a name past 55 would truncate
    # silently or push the Job's label over the limit. Fail loudly
    # instead (Pulumi twin: the length check in main.go Resources).
    precondition {
      condition     = length(var.metadata.name) <= 55
      error_message = "metadata.name exceeds the OpenFGA chart's 55-character name budget (the chart truncates its fullname at 63 and appends \"-migrate\" for the migration Job)."
    }
  }
}
