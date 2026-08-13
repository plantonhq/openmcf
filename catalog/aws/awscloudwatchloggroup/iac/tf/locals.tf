locals {
  # The log group's cloud name is the resource's metadata.name — the same
  # basis the Pulumi module uses, and what services that address the group by
  # name (ECS awslogs, ElastiCache log delivery) will see.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCloudwatchLogGroup"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Retention 0 means "never expire" — pass null so the provider leaves the
  # retention policy unmanaged instead of writing an explicit zero.
  retention_in_days = var.spec.retention_in_days != 0 ? var.spec.retention_in_days : null

  # KMS key arrives pre-resolved to a plain ARN string by the orchestrator
  # (the generator flattens StringValueOrRef fields).
  kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Null when unset so AWS defaults the class to STANDARD.
  log_group_class = var.spec.log_group_class != "" ? var.spec.log_group_class : null

  # Tri-state pass-through: unset (null) leaves the provider's Computed
  # attribute alone (new groups get AWS's unprotected default; existing groups
  # keep their state), true enables protection, and an EXPLICIT false is the
  # only way to disable protection once enabled — the provider keeps the
  # existing value whenever the argument is omitted.
  deletion_protection_enabled = var.spec.deletion_protection_enabled
}
