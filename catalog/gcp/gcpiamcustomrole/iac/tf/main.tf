# A project-scoped IAM custom role — a named, least-privilege permission
# bundle. role_id and project are immutable (ForceNew): changing either
# destroys and recreates the role, breaking every grant that references the
# old projects/<project>/roles/<role_id> name. title, description, stage, and
# permissions all update in place, and permission edits propagate immediately
# to every existing grant of the role.
#
# GCP soft-deletes custom roles: after destroy, the role_id stays reserved for
# up to 14 days. Re-creating a role with a soft-deleted ID undeletes it and
# patches it to this configuration — the provider handles that flow natively,
# so a destroy/recreate cycle within the window converges rather than failing.
resource "google_project_iam_custom_role" "this" {
  role_id     = var.spec.role_id
  project     = local.project_id
  title       = var.spec.title
  description = var.spec.description
  permissions = local.permissions
  stage       = local.stage

  # DELETE (provider default) soft-deletes the role on destroy; PREVENT
  # fails the destroy; ABANDON leaves the role active.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
}
