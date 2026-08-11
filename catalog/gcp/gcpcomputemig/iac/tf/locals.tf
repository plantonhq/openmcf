locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The cloud-side group name defaults to metadata.name when the spec
  # leaves mig_name empty — the same naming basis every kind uses.
  mig_name = (
    var.spec.mig_name != null && var.spec.mig_name != ""
    ? var.spec.mig_name
    : var.metadata.name
  )

  # Instances are named "<base_instance_name>-<suffix>"; defaults to the
  # group name.
  base_instance_name = (
    var.spec.base_instance_name != null && var.spec.base_instance_name != ""
    ? var.spec.base_instance_name
    : local.mig_name
  )

  autoscaler_name = (
    var.spec.autoscaler != null && var.spec.autoscaler.autoscaler_name != ""
    ? var.spec.autoscaler.autoscaler_name
    : local.mig_name
  )

  # Templates are immutable, so every template change ROTATES the
  # template through name_prefix + create_before_destroy: a fresh
  # "<mig-name>-<timestamp>" template is created first, the group
  # manager repoints, then the old template is deleted. The prefix is
  # capped at 37 characters so the provider uses its readable timestamp
  # naming (beyond 37 it falls back to a collision-prone shortened
  # UUID) — identical truncation in the Pulumi module.
  template_name_prefix = (
    length(local.mig_name) > 36
    ? "${substr(local.mig_name, 0, 36)}-"
    : "${local.mig_name}-"
  )

  # Empty region means a ZONAL group (zone set); a set region selects
  # the regional resource family. Both branches of every resource pair
  # in main.tf key off this. The spec's CEL guarantees exactly one.
  is_regional = var.spec.region != null && var.spec.region != ""

  # The group's location — exported for downstream scope-compatibility
  # checks (a regional backend service needs backends in its region).
  location = local.is_regional ? var.spec.region : var.spec.zone

  # User labels first so platform attribution labels win on key
  # conflicts — identical merge order to the Pulumi module. These are
  # the VM labels the template stamps on every instance (the ONLY
  # template surface GCP allows to change in place).
  final_labels = merge(
    var.spec.template.labels,
    {
      "planton-resource"      = "true"
      "planton-resource-name" = local.mig_name
      "planton-resource-kind" = "gcpcomputemig"
    },
    var.metadata.id != null && var.metadata.id != "" ? {
      "planton-resource-id" = var.metadata.id
    } : {},
    var.metadata.org != null && var.metadata.org != "" ? {
      "planton-organization" = var.metadata.org
    } : {},
    var.metadata.env != null && var.metadata.env != "" ? {
      "planton-environment" = var.metadata.env
    } : {}
  )

  # Versions: an empty list runs one default version on this kind's own
  # template; entries with an empty template_self_link also resolve to
  # the own template (the canary escape hatch pins external URLs).
  versions = (
    length(var.spec.versions) > 0
    ? var.spec.versions
    : [{
      version_name        = ""
      template_self_link  = ""
      target_size_fixed   = null
      target_size_percent = null
    }]
  )

  # Spot semantics: SPOT requires the API's legacy preemptible flag and
  # forbids automatic restart (FLEX_START likewise never auto-restarts).
  # Deriving both here — identically to the Pulumi module — keeps the
  # spec's single provisioning_model switch honest.
  scheduling_is_reclaimable = (
    var.spec.template.scheduling != null
    && contains(["SPOT", "FLEX_START"], var.spec.template.scheduling.provisioning_model)
  )

  # Keyed per-instance-config and resize-request maps so for_each
  # addresses entries by name (stable plans on reorder).
  per_instance_configs = { for config in var.spec.per_instance_configs : config.config_name => config }
  resize_requests      = { for request in var.spec.resize_requests : request.request_name => request }
}
