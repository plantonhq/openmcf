locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  bucket_name = var.spec.bucket_name
  location    = var.spec.location

  # Empty optional strings become null so the provider applies its own
  # defaults (STANDARD class, inherited public-access prevention, DEFAULT
  # rpo) instead of receiving an empty string it would reject.
  storage_class            = var.spec.storage_class != "" ? var.spec.storage_class : null
  public_access_prevention = var.spec.public_access_prevention != "" ? var.spec.public_access_prevention : null
  rpo                      = var.spec.rpo != "" ? var.spec.rpo : null
  kms_key_name             = var.spec.kms_key_name != "" ? var.spec.kms_key_name : null

  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.spec.bucket_name
    "planton-ai_kind"     = "gcpgcsbucket"
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

  # Additive IAM grants keyed by (role, member) — the grant's identity.
  # A condition never changes the key: the same (role, member) pair with a
  # different condition is a replacement of that one grant.
  iam_members = {
    for m in var.spec.iam_members :
    "${m.role}/${m.member}" => m
  }

  # Folders split by nesting depth (1-5, the spec's cap). The Storage API
  # never auto-creates missing parents, and instances of one for_each
  # cannot depend on each other — so each depth is its own resource
  # group, chained with depends_on: parents create first, children
  # destroy first. Keys are the folder paths — the same identity the
  # Pulumi module names its folder resources with.
  folders_by_depth = {
    for depth in range(1, 6) :
    depth => {
      for f in var.spec.folders :
      f.name => f if length(regexall("/", f.name)) == depth
    }
  }

  # Managed folders keyed by path — independent prefix anchors, no
  # ordering concerns.
  managed_folders = {
    for f in var.spec.managed_folders :
    f.name => f
  }

  # Notifications keyed by list index: the resource is immutable end to
  # end (every change is a replace), so index-keyed churn on reorder is
  # the resource's only change mode anyway. Matches the Pulumi module's
  # notification-<index> names.
  notifications = {
    for idx, n in var.spec.notifications :
    tostring(idx) => n
  }
}
