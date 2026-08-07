# One Google Cloud project — the Layer-0 container every other GCP
# resource lives in. IAM grants are deliberately NOT bundled here; model
# them as first-class GcpProjectIamMember resources.
resource "google_project" "this" {
  name       = local.display_name
  project_id = var.spec.project_id

  billing_account = local.billing_account_id
  org_id          = local.parent_org_id
  folder_id       = local.parent_folder_id

  labels = length(local.final_labels) > 0 ? local.final_labels : null

  # Resource Manager tags bind at create time only; changing them
  # afterwards forces recreation (bind tag values out-of-band instead).
  tags = length(var.spec.tags) > 0 ? var.spec.tags : null

  # False by default: deleting the auto-created "default" network is a
  # standard hardening step, and explicit GcpVpcNetwork resources are the
  # composable path.
  auto_create_network = var.spec.auto_create_network

  deletion_policy = local.deletion_policy
}

# Pre-enable the requested Cloud APIs. disable_on_destroy is false:
# removing an entry (or the project resource) must never disable a service
# that other tooling in the project depends on.
resource "google_project_service" "enabled_apis" {
  for_each = toset(var.spec.enabled_apis)

  project = google_project.this.project_id
  service = each.value

  disable_dependent_services = true
  disable_on_destroy         = false

  depends_on = [google_project.this]
}
