# Computed values shared across the module.

locals {
  # Resource-identity labels: the kuberneteslabelkeys set, identical to what
  # the Pulumi module stamps for the same manifest. Stamped on the namespace
  # this module creates — never injected into the chart's own resources
  # (Helm owns those).
  identity_labels_base = {
    "planton.ai/resource" = "true"
    "planton.ai/name"     = var.metadata.name
    "planton.ai/kind"     = "KubernetesHelmRelease"
  }

  id_label = (
    var.metadata.id != null && try(var.metadata.id, "") != ""
  ) ? { "planton.ai/id" = var.metadata.id } : {}

  org_label = (
    var.metadata.org != null && try(var.metadata.org, "") != ""
  ) ? { "planton.ai/organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && try(var.metadata.env, "") != ""
  ) ? { "planton.ai/environment" = var.metadata.env } : {}

  labels = merge(
    local.identity_labels_base,
    local.id_label,
    local.org_label,
    local.env_label,
  )

  # The namespace the release installs into (the generator flattens the
  # spec's value-or-ref to a plain string) and whether this module creates
  # it.
  namespace_name   = var.spec.namespace
  create_namespace = try(var.spec.create_namespace, false)

  # The Helm release name: spec.release_name when set, otherwise the
  # resource's metadata.name — identical resolution in the Pulumi module.
  release_name = try(var.spec.release_name, "") != "" ? var.spec.release_name : var.metadata.name

  # Helm's own defaults resolved for the optional knobs so both engines send
  # identical values whether or not the spec set the fields (the Pulumi
  # module resolves the same defaults in its locals).
  timeout_seconds = try(var.spec.timeout_seconds, null) != null ? var.spec.timeout_seconds : 300
  max_history     = try(var.spec.max_history, null) != null ? var.spec.max_history : 10

  # skip_await inverts helm_release's wait flag (the Pulumi module passes
  # SkipAwait directly). Both engines default to awaiting readiness.
  wait = !try(var.spec.skip_await, false)
}
