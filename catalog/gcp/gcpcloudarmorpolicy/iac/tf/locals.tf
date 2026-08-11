locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # policy_name falls back to metadata.name — explicit conditional, so both
  # engines derive the identical cloud-side name.
  policy_name = var.spec.policy_name != "" ? var.spec.policy_name : var.metadata.name

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it. Conditional labels appear under the same conditions on both
  # sides.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = local.policy_name
    "planton-ai_kind"     = "gcpcloudarmorpolicy"
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

  # User labels first: the platform labels win on key conflicts.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)
}
