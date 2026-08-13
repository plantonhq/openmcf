# Amazon SageMaker AI notebook instance with its folded lifecycle
# configuration (bootstrap scripts).
#
# Lifecycle facts the renders below depend on:
#   - most instance changes ride the provider's stop-update-start
#     choreography (SageMaker requires a Stopped instance for
#     UpdateNotebookInstance) - budget several minutes per change;
#   - growing volume_size updates in place, SHRINKING replaces the
#     instance (provider CustomizeDiff mirrors AWS's no-shrink rule);
#   - the lifecycle scripts are sent base64-encoded (the module encodes;
#     the spec carries plain shell) and run as root with a 5-minute
#     limit;
#   - clearing a script upstream does NOT clear it in AWS (the
#     provider's update omits empty fields) - replacing the text is the
#     reliable path, taught on the spec fields.

resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "this" {
  count = local.has_lifecycle ? 1 : 0

  name      = local.lifecycle_config_name
  on_create = var.spec.lifecycle_config.on_create != "" ? base64encode(var.spec.lifecycle_config.on_create) : null
  on_start  = var.spec.lifecycle_config.on_start != "" ? base64encode(var.spec.lifecycle_config.on_start) : null

  tags = local.aws_tags
}

resource "aws_sagemaker_notebook_instance" "this" {
  # The component's name IS the notebook instance name.
  name = local.notebook_name

  instance_type = var.spec.instance_type
  role_arn      = var.spec.role_arn

  volume_size = var.spec.volume_size_gb

  subnet_id       = var.spec.subnet_id != "" ? var.spec.subnet_id : null
  security_groups = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  kms_key_id      = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  # Changing either replaces the instance (provider-enforced).
  direct_internet_access = var.spec.direct_internet_access != "" ? var.spec.direct_internet_access : null
  platform_identifier    = var.spec.platform_identifier != "" ? var.spec.platform_identifier : null

  root_access = var.spec.root_access != "" ? var.spec.root_access : null

  default_code_repository      = var.spec.default_code_repository != "" ? var.spec.default_code_repository : null
  additional_code_repositories = length(var.spec.additional_code_repositories) > 0 ? var.spec.additional_code_repositories : null

  lifecycle_config_name = local.has_lifecycle ? aws_sagemaker_notebook_instance_lifecycle_configuration.this[0].name : null

  dynamic "instance_metadata_service_configuration" {
    for_each = var.spec.imds_minimum_version != "" ? [var.spec.imds_minimum_version] : []
    content {
      minimum_instance_metadata_service_version = instance_metadata_service_configuration.value
    }
  }

  tags = local.aws_tags
}
