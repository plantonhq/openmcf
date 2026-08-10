locals {
  # The cloud-side policy name defaults to metadata.name when the spec
  # leaves policy_name empty — the same naming basis every kind uses.
  policy_name = (
    var.spec.policy_name != null && var.spec.policy_name != ""
    ? var.spec.policy_name
    : var.metadata.name
  )

  # Which parent arm is set. All empty means "the provider's default
  # project" — made concrete via google_client_config below.
  is_folder_parent = var.spec.parent != null && var.spec.parent.folder_id != ""
  is_org_parent    = var.spec.parent != null && var.spec.parent.organization_id != ""
  parent_project_id = (
    var.spec.parent != null && var.spec.parent.project_id != ""
    ? var.spec.parent.project_id
    : (local.is_folder_parent || local.is_org_parent ? "" : data.google_client_config.current[0].project)
  )

  # GCP's API identifies the attach point by its URL-ENCODED full resource
  # name. The module renders it from the typed parent so manifests never
  # hand-assemble it — urlencode() only needs to escape "/" on this
  # charset, identically to the Pulumi module's encoding.
  parent = (
    local.is_folder_parent
    ? urlencode("cloudresourcemanager.googleapis.com/folders/${trimprefix(var.spec.parent.folder_id, "folders/")}")
    : local.is_org_parent
    ? urlencode("cloudresourcemanager.googleapis.com/organizations/${trimprefix(var.spec.parent.organization_id, "organizations/")}")
    : urlencode("cloudresourcemanager.googleapis.com/projects/${local.parent_project_id}")
  )
}

# The provider's own resolved configuration — the source of the default
# project when the parent is omitted. Count-gated on that one case so every
# plan that names its attach point runs credential-free.
data "google_client_config" "current" {
  count = (
    var.spec.parent == null || (var.spec.parent.project_id == "" && var.spec.parent.folder_id == "" && var.spec.parent.organization_id == "")
  ) ? 1 : 0
}
