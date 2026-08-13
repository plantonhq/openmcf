# AWS Secrets Manager secret: a named, versioned, KMS-encrypted container
# for credential material with optional automatic rotation and cross-region
# replication.
#
# The module owns up to four resources, all discoverable by the secret's
# name: the secret itself, its resource policy (rendered through the
# standalone policy resource for its typed block_public_policy guard), the
# managed version (when a value arm is set), and the rotation
# configuration.

resource "aws_secretsmanager_secret" "this" {
  # Create-time immutable; doubles as the Name tag. metadata.name on both
  # engines -- never provider auto-naming (name_prefix stays unused).
  name = local.secret_name

  # Description is ALWAYS sent explicitly so the two engines never inject
  # differing defaults into state.
  description = var.spec.description

  # Customer-managed key when referenced; AWS-managed aws/secretsmanager
  # otherwise. Updates in place -- new versions encrypt under the new key.
  kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Consumed only at delete time: 0 forces immediate permanent deletion,
  # 7-30 keeps the soft-delete recovery window (default 30, materialized by
  # the manifest loader).
  recovery_window_in_days = var.spec.recovery_window_in_days

  # Always sent: false fails replication loudly on a name collision in a
  # replica region (the safe posture); true overwrites deliberately.
  force_overwrite_replica_secret = var.spec.force_overwrite_replica_secret

  # Cross-region replicas -- each encrypts under its own region's key (the
  # referenced customer key, or that region's AWS-managed key).
  # Two delete-time truths both engines inherit from the provider
  # (live-verified 2026-08-13): AWS deletes replica secrets ASYNCHRONOUSLY
  # after RemoveRegionsFromReplication and the provider does not wait for
  # it, so destroying with recovery_window 0 can strand a replica as a
  # live standalone secret (a recovery window lets the async deletion
  # complete); and replication is never waited on at create either, so a
  # Failed replication (e.g. against a stranded same-name ex-replica,
  # which force_overwrite does NOT clear) is silent at apply.
  dynamic "replica" {
    for_each = var.spec.replica_regions
    content {
      region     = replica.value.region
      kms_key_id = replica.value.kms_key_id != "" ? replica.value.kms_key_id : null
    }
  }

  # Managed external secret partner identifier (ForceNew). Omitted entirely
  # for ordinary self-managed secrets -- AWS treats an absent Type and an
  # empty Type differently in error paths, and omission is the documented
  # shape for self-managed secrets.
  type = var.spec.type != "" ? var.spec.type : null

  tags = local.aws_tags
}

# Resource policy -- rendered through the standalone policy resource (not
# the secret's inline policy argument) because only the standalone resource
# carries block_public_policy, the PutResourcePolicy guard that rejects
# policies granting anonymous access.
resource "aws_secretsmanager_secret_policy" "this" {
  count = local.create_policy ? 1 : 0

  secret_arn = aws_secretsmanager_secret.this.arn
  policy     = jsonencode(var.spec.policy)

  # Default true (materialized by the manifest loader): reject public
  # policies unless the manifest deliberately opts out.
  block_public_policy = var.spec.block_public_policy
}

# The managed version -- created only when a value arm is set. version_id
# is exported in every arm (empty for a shell secret) so both engines emit
# the same output set.
resource "aws_secretsmanager_secret_version" "this" {
  count = local.create_version ? 1 : 0

  secret_id     = aws_secretsmanager_secret.this.arn
  secret_string = var.spec.string_value != "" ? var.spec.string_value : null
  # The provider expects base64 in secret_binary and decodes it before
  # calling PutSecretValue (CEL already guaranteed the encoding at manifest
  # time).
  secret_binary = var.spec.binary_value != "" ? var.spec.binary_value : null

  # Custom staging labels ride ALONGSIDE AWSCURRENT: providing
  # version_stages REPLACES the automatic AWSCURRENT assignment, so the
  # module always includes it -- dropping it would leave the secret with no
  # current version.
  version_stages = length(var.spec.version_stages) > 0 ? concat(["AWSCURRENT"], var.spec.version_stages) : null
}

# Rotation. Ordered after the version: with rotate_immediately (the
# default) AWS invokes the rotation mechanism as soon as rotation is
# configured, and the rotation function reads the current value -- so the
# value must exist first.
resource "aws_secretsmanager_secret_rotation" "this" {
  count = local.create_rotation ? 1 : 0

  secret_id = aws_secretsmanager_secret.this.arn

  # Exactly one mechanism (CEL-enforced): the self-managed rotation Lambda,
  # or the partner-managed external rotation role (pairs with spec.type).
  rotation_lambda_arn               = var.spec.rotation.rotation_lambda_arn != "" ? var.spec.rotation.rotation_lambda_arn : null
  external_secret_rotation_role_arn = var.spec.rotation.external_rotation_role_arn != "" ? var.spec.rotation.external_rotation_role_arn : null

  dynamic "external_secret_rotation_metadata" {
    for_each = var.spec.rotation.external_rotation_metadata
    content {
      key   = external_secret_rotation_metadata.value.key
      value = external_secret_rotation_metadata.value.value
    }
  }

  # Default true (materialized): rotate once as soon as rotation is
  # configured. Explicit false only tests the configuration.
  rotate_immediately = var.spec.rotation.rotate_immediately

  rotation_rules {
    # Exactly one cadence arm (CEL-enforced). automatically_after_days and
    # schedule_expression are ExactlyOneOf at the provider too.
    automatically_after_days = var.spec.rotation.automatically_after_days
    schedule_expression      = var.spec.rotation.schedule_expression != "" ? var.spec.rotation.schedule_expression : null
    duration                 = var.spec.rotation.duration != "" ? var.spec.rotation.duration : null
  }

  depends_on = [aws_secretsmanager_secret_version.this]
}
