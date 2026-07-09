locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # policy_name falls back to metadata.name — explicit conditional, so both
  # engines derive the identical cloud-side name.
  policy_name = var.spec.policy_name != "" ? var.spec.policy_name : var.metadata.name

  # Cloud Armor security policies carry no labels on the released provider
  # line — attribution is deliberately not attempted on either engine.
}
