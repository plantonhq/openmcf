locals {
  # Honor the spec contract: an empty project_id falls back to the
  # provider's default project (null lets the google provider resolve its
  # own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # instance_name falls back to metadata.name — explicit conditional, so
  # both engines derive the identical cloud-side name.
  instance_name = var.spec.instance_name != "" ? var.spec.instance_name : var.metadata.name

  # Empty optional strings become null so the provider applies its own
  # defaults instead of receiving an empty string it would reject.
  hostname                    = var.spec.hostname != "" ? var.spec.hostname : null
  min_cpu_platform            = var.spec.min_cpu_platform != "" ? var.spec.min_cpu_platform : null
  desired_status              = var.spec.desired_status != "" ? var.spec.desired_status : null
  key_revocation_action_type  = var.spec.key_revocation_action_type != "" ? var.spec.key_revocation_action_type : null
  total_egress_bandwidth_tier = var.spec.total_egress_bandwidth_tier != "" ? var.spec.total_egress_bandwidth_tier : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.instance_name
    "planton-ai_kind"     = "gcpcomputeinstance"
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

  # User labels first so platform attribution labels win on key conflicts —
  # identical merge order to the Pulumi module.
  final_labels = merge(
    var.spec.labels,
    local.base_labels,
    local.org_label,
    local.env_label,
    local.id_label,
  )

  # SSH keys fold into the metadata "ssh-keys" key (newline-joined) —
  # byte-identical to the Pulumi module's fold. The startup script rides
  # the dedicated metadata_startup_script attribute, never plain metadata,
  # so it re-runs on every boot exactly as GCP documents.
  ssh_keys_metadata = (
    length(var.spec.ssh_keys) > 0
    ? { "ssh-keys" = join("\n", var.spec.ssh_keys) }
    : {}
  )

  final_metadata = merge(var.spec.metadata, local.ssh_keys_metadata)

  # Spot semantics: SPOT requires the API's legacy preemptible flag and
  # forbids automatic restart. Deriving both here (identically in the
  # Pulumi module) keeps the spec's single provisioning_model switch
  # honest.
  is_spot = var.spec.scheduling != null && var.spec.scheduling.provisioning_model == "SPOT"

  automatic_restart = (
    local.is_spot
    ? false
    : (var.spec.scheduling != null && var.spec.scheduling.automatic_restart != null ? var.spec.scheduling.automatic_restart : true)
  )
}
