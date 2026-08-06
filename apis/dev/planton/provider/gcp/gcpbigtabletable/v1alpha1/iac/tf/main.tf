# Enable the Bigtable Admin API — table and GC-policy management run
# through it. disable_on_destroy is false: tearing down one table must
# never disable the API for everything else in the project.
resource "google_project_service" "bigtableadmin_api" {
  project = local.project_id
  service = "bigtableadmin.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Bigtable table: column families, pre-splits, change streams,
# automated backups, and deletion protection. Data (rows and cells) is
# application territory — this resource owns the structure applications
# write into.
#
# table name, instance, and split_keys are immutable (ForceNew); changing
# split_keys REPLACES the table and its data. Column families are mutable
# (append as the application grows). deletion_protection is the API-side
# guard (spec default PROTECTED): deletion by ANY client fails until it
# is set UNPROTECTED — stronger than an IaC-side guard for a
# data-bearing resource.
#
# Column families are created here with no GC policy; per-family
# retention lives in the google_bigtable_gc_policy resources below, so a
# policy change never touches the table object itself.
resource "google_bigtable_table" "this" {
  name          = local.table_name
  instance_name = var.spec.instance
  project       = local.project_id

  split_keys              = length(var.spec.split_keys) > 0 ? var.spec.split_keys : null
  change_stream_retention = var.spec.change_stream_retention != "" ? var.spec.change_stream_retention : null
  deletion_protection     = var.spec.deletion_protection
  row_key_schema          = var.spec.row_key_schema != "" ? var.spec.row_key_schema : null

  dynamic "column_family" {
    for_each = var.spec.column_families
    content {
      family = column_family.value.family
      type   = column_family.value.type != "" ? column_family.value.type : null
    }
  }

  dynamic "automated_backup_policy" {
    for_each = var.spec.automated_backup_policy != null ? [var.spec.automated_backup_policy] : []
    content {
      retention_period = automated_backup_policy.value.retention_period
      frequency        = automated_backup_policy.value.frequency
    }
  }

  depends_on = [google_project_service.bigtableadmin_api]
}

# One GC policy per column family that declares one — the API's own
# granularity. Bigtable never deletes old cell versions without a GC
# policy, so an unbounded family accumulates every write forever.
# Policies are mutable in place; deleting one resets the family to
# "no GC" rather than deleting data.
resource "google_bigtable_gc_policy" "this" {
  for_each = local.gc_policies

  instance_name = var.spec.instance
  project       = local.project_id
  table         = google_bigtable_table.this.name
  column_family = each.key

  # The raw JSON tree and the typed fields are mutually exclusive
  # (enforced by spec CEL) — exactly one shape reaches the provider.
  gc_rules = each.value.gc_rules != "" ? each.value.gc_rules : null
  mode     = each.value.mode != "" ? each.value.mode : null

  dynamic "max_age" {
    for_each = each.value.max_age != "" ? [each.value.max_age] : []
    content {
      duration = max_age.value
    }
  }

  dynamic "max_version" {
    for_each = each.value.max_versions > 0 ? [each.value.max_versions] : []
    content {
      number = max_version.value
    }
  }

  # Allows EXPANDING what is eligible for collection on a replicated
  # instance — Bigtable otherwise rejects the change as a data-loss
  # safety measure.
  ignore_warnings = each.value.ignore_warnings
}
