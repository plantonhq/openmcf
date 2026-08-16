# An AWS Backup vault: exactly one of the two AWS vault types (the
# spec's exactly-one union mirrors AWS's own VaultType discriminator),
# plus the standard arm's attachable satellites.
#
# Lifecycle facts the render below depends on:
#   - the three satellites (lock, policy, notifications) attach by
#     VAULT NAME and only to STANDARD vaults - the provider's readers
#     reject other vault types, so the union gates them structurally;
#   - force_destroy is deploy-side delete behavior, never reported
#     back by AWS - invisible to imports, asserted only at destroy;
#   - lock MODE is decided by changeable_for_days alone: unset =
#     governance (removable), set = compliance (immutable once the
#     cooling-off window passes) - and AWS never reports the window
#     back;
#   - an air-gapped vault is immutable apart from tags: every argument
#     forces replacement, and its recovery points cannot be manually
#     deleted (they age out by retention).

resource "aws_backup_vault" "this" {
  count = var.spec.standard != null ? 1 : 0

  # metadata.name is the vault name on both engines (AWS: 2-50
  # characters of letters, digits, hyphens, underscores).
  name = var.metadata.name

  # Rendered only on an explicit choice so the module never fights the
  # provider default (the AWS Backup service key).
  kms_key_arn = var.spec.standard.kms_key_arn != "" ? var.spec.standard.kms_key_arn : null

  # Deploy-side recovery-point drain at destroy (see the header note).
  force_destroy = var.spec.standard.force_destroy

  tags = local.aws_tags
}

resource "aws_backup_logically_air_gapped_vault" "this" {
  count = var.spec.air_gapped != null ? 1 : 0

  name = var.metadata.name

  # Both retention bounds are REQUIRED by AWS on this vault type (min
  # floor 7 days) and both force replacement.
  min_retention_days = var.spec.air_gapped.min_retention_days
  max_retention_days = var.spec.air_gapped.max_retention_days

  # Rendered only on an explicit choice so the module never fights the
  # provider default (the AWS-owned key).
  encryption_key_arn = var.spec.air_gapped.encryption_key_arn != "" ? var.spec.air_gapped.encryption_key_arn : null

  tags = local.aws_tags
}

resource "aws_backup_vault_lock_configuration" "this" {
  count = var.spec.standard != null ? (var.spec.standard.lock != null ? 1 : 0) : 0

  backup_vault_name = aws_backup_vault.this[0].name

  # Present = COMPLIANCE mode (immutable after the window); absent =
  # governance mode. Write-only at AWS - never read back.
  changeable_for_days = var.spec.standard.lock.changeable_for_days

  min_retention_days = var.spec.standard.lock.min_retention_days
  max_retention_days = var.spec.standard.lock.max_retention_days
}

resource "aws_backup_vault_policy" "this" {
  count = var.spec.standard != null ? (var.spec.standard.policy != null ? 1 : 0) : 0

  backup_vault_name = aws_backup_vault.this[0].name
  policy            = jsonencode(var.spec.standard.policy)
}

resource "aws_backup_vault_notifications" "this" {
  count = var.spec.standard != null ? (var.spec.standard.notifications != null ? 1 : 0) : 0

  backup_vault_name   = aws_backup_vault.this[0].name
  sns_topic_arn       = var.spec.standard.notifications.sns_topic_arn
  backup_vault_events = var.spec.standard.notifications.events
}
