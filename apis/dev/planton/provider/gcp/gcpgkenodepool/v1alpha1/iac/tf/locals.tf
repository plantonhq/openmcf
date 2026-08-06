locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Pool name defaults to metadata.name (the spec-level contract), so a
  # manifest only sets node_pool_name when the cloud name must differ from
  # the Planton object name.
  node_pool_name = var.spec.node_pool_name != "" ? var.spec.node_pool_name : var.metadata.name

  # Node config is optional end to end: a spec that omits it gets GKE's
  # defaults. All node_config reads below go through this local so the
  # rest of the module never null-checks the object itself.
  nc = var.spec.node_config

  # Empty optional strings become null so the provider omits them from the
  # API payload instead of sending empty values it would reject or diff on.
  version          = var.spec.version != "" ? var.spec.version : null
  machine_type     = try(local.nc.machine_type, "e2-medium")
  disk_type        = try(local.nc.disk_type, "") != "" ? local.nc.disk_type : null
  image_type       = try(local.nc.image_type, "COS_CONTAINERD")
  service_account  = try(local.nc.service_account, "") != "" ? local.nc.service_account : null
  min_cpu_platform = try(local.nc.min_cpu_platform, "") != "" ? local.nc.min_cpu_platform : null
  boot_disk_kms_key = try(local.nc.boot_disk_kms_key, "") != "" ? local.nc.boot_disk_kms_key : null
  logging_variant  = try(local.nc.logging_variant, "") != "" ? local.nc.logging_variant : null
  max_run_duration = try(local.nc.max_run_duration, "") != "" ? local.nc.max_run_duration : null

  # The same planton-ai_* label set the Pulumi module applies, so a pool is
  # attributable to its Planton object regardless of the engine that created
  # it. The pool name (not metadata.name) keys the name label so the label
  # matches what is visible in the GCP console.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.node_pool_name
    "planton-ai_kind"     = "gcpgkenodepool"
  }

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  # User resource labels merge in first so the platform attribution labels
  # can never be clobbered by a spec label with the same key. These are GCE
  # resource labels on the node VMs — distinct from Kubernetes node labels,
  # which pass through node_config.labels untouched.
  final_resource_labels = merge(
    try(local.nc.resource_labels, {}),
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )

  # GKE requires disable-legacy-endpoints=true on every node; it is
  # enforced beneath any user metadata so a spec entry can never weaken it.
  node_metadata = merge(
    try(local.nc.metadata, {}),
    { "disable-legacy-endpoints" = "true" },
  )

  # Outputs-side sizing: the effective min/max mirror either the fixed
  # node_count or the autoscaling bounds (per-zone arm, else total arm) so
  # downstream consumers read one honest pair regardless of sizing mode.
  effective_min = (
    var.spec.autoscaling != null
    ? coalesce(var.spec.autoscaling.min_nodes, var.spec.autoscaling.total_min_nodes, 0)
    : coalesce(var.spec.node_count, 0)
  )
  effective_max = (
    var.spec.autoscaling != null
    ? coalesce(var.spec.autoscaling.max_nodes, var.spec.autoscaling.total_max_nodes, 0)
    : coalesce(var.spec.node_count, 0)
  )
}
