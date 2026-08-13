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
    "planton.ai/resource-kind" = "AwsFsxOpenzfsFileSystem"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Empty strings become null so unset stays indistinguishable from the AWS
  # defaults; empty-string arguments would otherwise fail provider validation
  # (ARN/CIDR format checks) or create phantom diffs.
  kms_key_id                = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  backup_id                 = var.spec.backup_id != "" ? var.spec.backup_id : null
  preferred_subnet_id       = var.spec.preferred_subnet_id != "" ? var.spec.preferred_subnet_id : null
  endpoint_ip_address_range = var.spec.endpoint_ip_address_range != "" ? var.spec.endpoint_ip_address_range : null

  daily_automatic_backup_start_time = var.spec.daily_automatic_backup_start_time != "" ? var.spec.daily_automatic_backup_start_time : null
  weekly_maintenance_start_time     = var.spec.weekly_maintenance_start_time != "" ? var.spec.weekly_maintenance_start_time : null

  # Empty collections become null so the provider omits the argument entirely:
  # an omitted security_group_ids lets AWS attach the VPC default SG (the
  # argument is ForceNew — an empty set and an omitted one are different
  # plans), and omitted route_table_ids let AWS use the VPC default table.
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  route_table_ids    = length(var.spec.route_table_ids) > 0 ? var.spec.route_table_ids : null
  delete_options     = length(var.spec.delete_options) > 0 ? var.spec.delete_options : null

  # final_backup_tags only matter when a final backup is actually taken.
  final_backup_tags = length(var.spec.final_backup_tags) > 0 ? var.spec.final_backup_tags : null
}
