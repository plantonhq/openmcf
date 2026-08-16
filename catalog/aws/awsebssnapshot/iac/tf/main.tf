# One EBS snapshot - snapshot a volume XOR copy a snapshot XOR import
# a disk image - with fast snapshot restore and cross-account
# createVolumePermission grants managed in-line.
#
# Lifecycle facts the render below depends on:
#   - the three arms are three provider resources (aws_ebs_snapshot /
#     aws_ebs_snapshot_copy / aws_ebs_snapshot_import); exactly one
#     exists per the spec's union CEL, and all three expose the same
#     downstream surface (id/arn/owner/size);
#   - only the volume arm is importable at the provider - the copy and
#     import resources ship no importer (declared honestly in the
#     import catalog);
#   - storage_tier / permanent_restore / temporary_restore_days update
#     in place on all arms; every source field replaces;
#   - fast snapshot restore bills per zone-hour while enabled - each
#     zone is its own resource keyed by zone name;
#   - createVolumePermission grants are per-account resources keyed by
#     account id; encrypted snapshots additionally need the KMS key
#     shared out-of-band.

locals {
  is_copy   = var.spec.copy_from != null
  is_import = var.spec.import_from != null
  is_volume = !local.is_copy && !local.is_import

  storage_tier           = var.spec.storage_tier != "" ? var.spec.storage_tier : null
  permanent_restore      = var.spec.permanent_restore ? true : null
  temporary_restore_days = var.spec.temporary_restore_days > 0 ? var.spec.temporary_restore_days : null
  description            = var.spec.description != "" ? var.spec.description : null
}

# The VOLUME arm: snapshot a live volume.
resource "aws_ebs_snapshot" "this" {
  count = local.is_volume ? 1 : 0

  volume_id   = var.spec.volume_id
  description = local.description

  storage_tier           = local.storage_tier
  permanent_restore      = local.permanent_restore
  temporary_restore_days = local.temporary_restore_days

  tags = local.aws_tags
}

# The COPY arm: copy an existing snapshot (same- or cross-region),
# optionally re-encrypting.
resource "aws_ebs_snapshot_copy" "this" {
  count = local.is_copy ? 1 : 0

  source_snapshot_id = var.spec.copy_from.source_snapshot_id
  source_region      = var.spec.copy_from.source_region
  description        = local.description

  encrypted  = var.spec.copy_from.encrypted ? true : null
  kms_key_id = var.spec.copy_from.kms_key_id != "" ? var.spec.copy_from.kms_key_id : null

  completion_duration_minutes = var.spec.copy_from.completion_duration_minutes > 0 ? var.spec.copy_from.completion_duration_minutes : null

  storage_tier           = local.storage_tier
  permanent_restore      = local.permanent_restore
  temporary_restore_days = local.temporary_restore_days

  tags = local.aws_tags
}

# The IMPORT arm: build the snapshot from a disk image via VM
# Import/Export.
resource "aws_ebs_snapshot_import" "this" {
  count = local.is_import ? 1 : 0

  description = local.description
  role_name   = var.spec.import_from.role_name != "" ? var.spec.import_from.role_name : null

  encrypted  = var.spec.import_from.encrypted ? true : null
  kms_key_id = var.spec.import_from.kms_key_id != "" ? var.spec.import_from.kms_key_id : null

  disk_container {
    format      = var.spec.import_from.disk_container.format
    description = var.spec.import_from.disk_container.description != "" ? var.spec.import_from.disk_container.description : null
    url         = var.spec.import_from.disk_container.url != "" ? var.spec.import_from.disk_container.url : null

    dynamic "user_bucket" {
      for_each = var.spec.import_from.disk_container.s3_bucket != "" ? [1] : []
      content {
        s3_bucket = var.spec.import_from.disk_container.s3_bucket
        s3_key    = var.spec.import_from.disk_container.s3_key
      }
    }
  }

  storage_tier           = local.storage_tier
  permanent_restore      = local.permanent_restore
  temporary_restore_days = local.temporary_restore_days

  tags = local.aws_tags
}

locals {
  # The one snapshot whichever arm produced - the downstream surface
  # is identical.
  snapshot_id = local.is_volume ? aws_ebs_snapshot.this[0].id : (local.is_copy ? aws_ebs_snapshot_copy.this[0].id : aws_ebs_snapshot_import.this[0].id)
  snapshot_arn = local.is_volume ? aws_ebs_snapshot.this[0].arn : (local.is_copy ? aws_ebs_snapshot_copy.this[0].arn : aws_ebs_snapshot_import.this[0].arn)
  snapshot_owner_id = local.is_volume ? aws_ebs_snapshot.this[0].owner_id : (local.is_copy ? aws_ebs_snapshot_copy.this[0].owner_id : aws_ebs_snapshot_import.this[0].owner_id)
  snapshot_volume_size = local.is_volume ? aws_ebs_snapshot.this[0].volume_size : (local.is_copy ? aws_ebs_snapshot_copy.this[0].volume_size : aws_ebs_snapshot_import.this[0].volume_size)
}

# Fast snapshot restore, one resource per availability zone. Billed
# per zone-hour while enabled.
resource "aws_ebs_fast_snapshot_restore" "this" {
  for_each = toset(var.spec.fast_restore_availability_zones)

  availability_zone = each.value
  snapshot_id       = local.snapshot_id
}

# createVolumePermission grants, one resource per account.
resource "aws_snapshot_create_volume_permission" "this" {
  for_each = toset(var.spec.share_with_account_ids)

  snapshot_id = local.snapshot_id
  account_id  = each.value
}
