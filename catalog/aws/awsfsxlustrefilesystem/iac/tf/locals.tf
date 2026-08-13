locals {
  # FSx file systems have no cloud name argument — the console name is the
  # Name tag. metadata.name is the same basis the Pulumi module pins, keeping
  # the two engines' physical identity converged.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsFsxLustreFileSystem"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Empty strings become null so unset stays indistinguishable from the AWS
  # defaults; empty-string arguments would otherwise fail provider validation
  # (ARN/URI format checks) or create phantom diffs.
  kms_key_id               = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  file_system_type_version = var.spec.file_system_type_version != "" ? var.spec.file_system_type_version : null
  backup_id                = var.spec.backup_id != "" ? var.spec.backup_id : null
  import_path              = var.spec.import_path != "" ? var.spec.import_path : null
  export_path              = var.spec.export_path != "" ? var.spec.export_path : null
  auto_import_policy       = var.spec.auto_import_policy != "" ? var.spec.auto_import_policy : null
  drive_cache_type         = var.spec.drive_cache_type != "" ? var.spec.drive_cache_type : null

  daily_automatic_backup_start_time = var.spec.daily_automatic_backup_start_time != "" ? var.spec.daily_automatic_backup_start_time : null
  weekly_maintenance_start_time     = var.spec.weekly_maintenance_start_time != "" ? var.spec.weekly_maintenance_start_time : null

  # Empty security_group_ids means "let AWS attach the VPC default SG" —
  # pass null so the provider omits the argument entirely (it is ForceNew;
  # an empty set and an omitted argument are different plans).
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # final_backup_tags only matter when a final backup is actually taken.
  final_backup_tags = length(var.spec.final_backup_tags) > 0 ? var.spec.final_backup_tags : null
}
