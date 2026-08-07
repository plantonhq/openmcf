# ---------------------------------------------------------------------------
# AWS FSx data repository association
# ---------------------------------------------------------------------------
# One aws_fsx_data_repository_association resource links a Lustre directory to
# an S3 bucket/prefix. The association's identity is the (file system, path,
# bucket) triple — file_system_id, file_system_path, and data_repository_path
# are all ForceNew. The sync policies and imported_file_chunk_size update in
# place; batch_import_meta_data_on_create and delete_data_in_filesystem are
# create-/delete-time behaviors, not AWS state.
# ---------------------------------------------------------------------------

resource "aws_fsx_data_repository_association" "this" {
  # The (file system, path, bucket) identity triple. References arrive
  # pre-resolved as plain strings.
  file_system_id       = var.spec.file_system_id
  file_system_path     = var.spec.file_system_path
  data_repository_path = var.spec.data_repository_path

  # Bidirectional sync policies. The provider wraps them in an `s3` block;
  # the spec exposes the two event lists directly (the wrapper carries no
  # information of its own).
  dynamic "s3" {
    for_each = local.has_s3_policies ? [1] : []
    content {
      dynamic "auto_import_policy" {
        for_each = length(var.spec.auto_import_events) > 0 ? [1] : []
        content {
          events = var.spec.auto_import_events
        }
      }
      dynamic "auto_export_policy" {
        for_each = length(var.spec.auto_export_events) > 0 ? [1] : []
        content {
          events = var.spec.auto_export_events
        }
      }
    }
  }

  # Stripe size for imported files (null keeps the AWS default of 1024 MiB).
  imported_file_chunk_size = var.spec.imported_file_chunk_size

  # Create-time batch import of the existing S3 metadata — without it, only
  # objects changing AFTER creation appear in the namespace.
  batch_import_meta_data_on_create = var.spec.batch_import_meta_data_on_create

  # Delete-time cascade: remove the linked files from the file system when
  # the association is deleted (default keeps them).
  delete_data_in_filesystem = var.spec.delete_data_in_filesystem

  tags = local.aws_tags
}
