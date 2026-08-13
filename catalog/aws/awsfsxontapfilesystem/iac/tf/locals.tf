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
    "planton.ai/resource-kind" = "AwsFsxOntapFileSystem"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Whether this is a multi-AZ deployment — drives the preferred-subnet
  # derivation and the multi-AZ-only argument gating below.
  is_multi_az = contains(["MULTI_AZ_1", "MULTI_AZ_2"], var.spec.deployment_type)

  # The provider marks preferred_subnet_id Required for EVERY deployment type,
  # while the decision only exists for multi-AZ (single-AZ has one subnet).
  # The spec therefore requires it for multi-AZ and forbids it for single-AZ;
  # here the single-AZ value derives deterministically from the only subnet.
  preferred_subnet_id = local.is_multi_az ? var.spec.preferred_subnet_id : var.spec.subnet_ids[0]

  # Empty strings become null so unset stays indistinguishable from the AWS
  # defaults; empty-string arguments would otherwise fail provider validation
  # (ARN/CIDR format checks) or create phantom diffs.
  kms_key_id                = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  fsx_admin_password        = var.spec.fsx_admin_password != "" ? var.spec.fsx_admin_password : null
  endpoint_ip_address_range = var.spec.endpoint_ip_address_range != "" ? var.spec.endpoint_ip_address_range : null

  daily_automatic_backup_start_time = var.spec.daily_automatic_backup_start_time != "" ? var.spec.daily_automatic_backup_start_time : null
  weekly_maintenance_start_time     = var.spec.weekly_maintenance_start_time != "" ? var.spec.weekly_maintenance_start_time : null

  # Empty security_group_ids means "let AWS attach the VPC default SG" —
  # pass null so the provider omits the argument entirely (it is ForceNew;
  # an empty set and an omitted argument are different plans).
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # Empty route_table_ids means "let AWS manage routes in the VPC main route
  # table" — omit the argument rather than sending an empty set.
  route_table_ids = length(var.spec.route_table_ids) > 0 ? var.spec.route_table_ids : null
}
