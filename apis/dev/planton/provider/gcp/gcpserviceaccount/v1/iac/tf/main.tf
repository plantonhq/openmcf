# The service account identity. account_id and project are immutable (ForceNew
# in the provider): changing either destroys and recreates the account, which
# invalidates every IAM binding and workload identity referencing the old email.
# display_name, description, and disabled are all mutable in place.
resource "google_service_account" "main" {
  account_id   = var.spec.service_account_id
  display_name = local.display_name
  project      = local.project_id
  description  = var.spec.description

  # GCP flips disabled state via separate Enable/Disable API calls (not the
  # regular update mask); the provider handles that internally, so toggling
  # this field never recreates the account.
  disabled = local.disabled
}

# Optional user-managed JSON key. Created only when spec.create_key is true —
# keyless patterns (Workload Identity, impersonation, federation) are the
# recommended default. The private key is returned once at creation and is
# marked sensitive in state.
resource "google_service_account_key" "main" {
  count = local.create_key ? 1 : 0

  service_account_id = google_service_account.main.name
}

# Project-level role grants. google_project_iam_member is ADDITIVE: each grant
# manages exactly one (role, member) pair and never clobbers other members'
# bindings on the same role — safe to compose with grants made elsewhere.
resource "google_project_iam_member" "project_roles" {
  for_each = toset(local.project_iam_roles)

  # The grant must target the account's home project explicitly; deriving it
  # from the created account (rather than re-reading spec) keeps the grant
  # correct when the project fell back to the provider default.
  project = google_service_account.main.project
  role    = each.value
  member  = "serviceAccount:${local.service_account_email}"

  depends_on = [google_service_account.main]
}

# Organization-level role grants (additive, same member semantics as above).
# Org-scope roles affect every project under the organization.
resource "google_organization_iam_member" "org_roles" {
  for_each = toset(local.org_iam_roles)

  org_id = var.spec.org_id
  role   = each.value
  member = "serviceAccount:${local.service_account_email}"

  depends_on = [google_service_account.main]
}
