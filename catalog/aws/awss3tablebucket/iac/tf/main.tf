# One S3 table bucket with its full contents: namespaces, Iceberg
# tables, resource policies, and replication - seven provider
# resources under one declarative owner.
#
# Lifecycle facts the render below depends on:
#   - namespaces key by name; tables key by "namespace//table" (stable
#     across list reorders); both are create-only at AWS (renames
#     replace);
#   - a table's iceberg_schema/properties are CREATE-ONLY input: the
#     provider never reads them back (schema evolution happens through
#     query engines), so they cannot drift and never round-trip on
#     import;
#   - the table format argument is module-pinned to ICEBERG - the
#     provider's enum holds exactly that one value;
#   - encryption_configuration and maintenance_configuration are
#     untyped object attributes at the provider (protocol v6 pending)
#     - values must be complete, correctly-shaped objects, which the
#     typed spec guarantees;
#   - policies are JSON-normalized by AWS (importIgnore at the
#     provider); replication carries a version_token
#     optimistic-concurrency handshake the provider manages.

locals {
  bucket_encryption = var.spec.encryption != null ? {
    sse_algorithm = var.spec.encryption.sse_algorithm
    kms_key_arn   = var.spec.encryption.kms_key_arn != "" ? var.spec.encryption.kms_key_arn : null
  } : null

  maintenance = var.spec.unreferenced_file_removal != null ? {
    iceberg_unreferenced_file_removal = {
      status = var.spec.unreferenced_file_removal.disabled ? "disabled" : "enabled"
      settings = {
        non_current_days  = var.spec.unreferenced_file_removal.non_current_days > 0 ? var.spec.unreferenced_file_removal.non_current_days : null
        unreferenced_days = var.spec.unreferenced_file_removal.unreferenced_days > 0 ? var.spec.unreferenced_file_removal.unreferenced_days : null
      }
    }
  } : null

  # Tables flattened out of their namespaces, keyed
  # "namespace//table" (the "//" separator is the import bridge's
  # address-key-segment convention; the output maps reuse the key).
  tables = merge([
    for namespace in var.spec.namespaces : {
      for table in namespace.tables :
      "${namespace.name}//${table.name}" => merge(table, { namespace = namespace.name })
    }
  ]...)
}

resource "aws_s3tables_table_bucket" "this" {
  name = var.metadata.name

  encryption_configuration  = local.bucket_encryption
  maintenance_configuration = local.maintenance
  force_destroy             = var.spec.force_destroy

  tags = local.aws_tags
}

resource "aws_s3tables_table_bucket_policy" "this" {
  count = var.spec.resource_policy != "" ? 1 : 0

  table_bucket_arn = aws_s3tables_table_bucket.this.arn
  resource_policy  = var.spec.resource_policy
}

resource "aws_s3tables_table_bucket_replication" "this" {
  count = var.spec.replication != null ? 1 : 0

  table_bucket_arn = aws_s3tables_table_bucket.this.arn
  role             = var.spec.replication.role

  rule {
    dynamic "destination" {
      for_each = var.spec.replication.destination_table_bucket_arns
      content {
        destination_table_bucket_arn = destination.value
      }
    }
  }
}

resource "aws_s3tables_namespace" "this" {
  for_each = { for namespace in var.spec.namespaces : namespace.name => namespace }

  namespace        = each.value.name
  table_bucket_arn = aws_s3tables_table_bucket.this.arn
}

resource "aws_s3tables_table" "this" {
  for_each = local.tables

  name             = each.value.name
  namespace        = aws_s3tables_namespace.this[each.value.namespace].namespace
  table_bucket_arn = aws_s3tables_table_bucket.this.arn
  # The provider's enum holds exactly this one value.
  format = "ICEBERG"

  encryption_configuration = each.value.encryption != null ? {
    sse_algorithm = each.value.encryption.sse_algorithm
    kms_key_arn   = each.value.encryption.kms_key_arn != "" ? each.value.encryption.kms_key_arn : null
  } : null

  maintenance_configuration = (each.value.compaction != null || each.value.snapshot_management != null) ? {
    iceberg_compaction = each.value.compaction != null ? {
      status = each.value.compaction.disabled ? "disabled" : "enabled"
      settings = {
        target_file_size_mb = each.value.compaction.target_file_size_mb > 0 ? each.value.compaction.target_file_size_mb : null
      }
    } : null
    iceberg_snapshot_management = each.value.snapshot_management != null ? {
      status = each.value.snapshot_management.disabled ? "disabled" : "enabled"
      settings = {
        max_snapshot_age_hours = each.value.snapshot_management.max_snapshot_age_hours > 0 ? each.value.snapshot_management.max_snapshot_age_hours : null
        min_snapshots_to_keep  = each.value.snapshot_management.min_snapshots_to_keep > 0 ? each.value.snapshot_management.min_snapshots_to_keep : null
      }
    } : null
  } : null

  dynamic "metadata" {
    for_each = each.value.iceberg_schema != null ? [each.value.iceberg_schema] : []
    content {
      iceberg {
        properties = length(each.value.properties) > 0 ? each.value.properties : null
        schema {
          dynamic "field" {
            for_each = metadata.value.fields
            content {
              name     = field.value.name
              type     = field.value.type
              required = field.value.required
            }
          }
        }
      }
    }
  }

  # The schema (metadata) is a birth certificate: CreateTable consumes
  # it, no read path ever returns it, and the provider marks every leaf
  # RequiresReplace. Without this ignore, importing an existing table
  # plans a DESTRUCTIVE delete+create against the empty imported state
  # (upstream's own import tests ImportStateVerifyIgnore the whole
  # block). Post-create schema evolution flows through query engines
  # (ALTER TABLE), never through this field - editing it after create
  # is deliberately inert. PARITY: the Pulumi module ignores the same
  # path via IgnoreChanges.
  lifecycle {
    ignore_changes = [metadata]
  }

  tags = local.aws_tags
}

resource "aws_s3tables_table_policy" "this" {
  for_each = { for key, table in local.tables : key => table if table.resource_policy != "" }

  name             = aws_s3tables_table.this[each.key].name
  namespace        = each.value.namespace
  table_bucket_arn = aws_s3tables_table_bucket.this.arn
  resource_policy  = each.value.resource_policy
}

resource "aws_s3tables_table_replication" "this" {
  for_each = { for key, table in local.tables : key => table if table.replication != null }

  table_arn = aws_s3tables_table.this[each.key].arn
  role      = each.value.replication.role

  rule {
    dynamic "destination" {
      for_each = each.value.replication.destination_table_bucket_arns
      content {
        destination_table_bucket_arn = destination.value
      }
    }
  }
}
