# ---------------------------------------------------------------------------
# AWS EFS File System
# ---------------------------------------------------------------------------
# Core file system with encryption, throughput, lifecycle policies, and
# replication-overwrite protection. ForceNew attributes: encrypted,
# kms_key_id, performance_mode, availability_zone_name -- changing any of
# them replaces the file system.
# ---------------------------------------------------------------------------

resource "aws_efs_file_system" "this" {
  # The creation token is the idempotency key AND the only human-controlled
  # physical identity EFS offers (the console name is the Name tag). Pinned to
  # metadata.name -- the same basis the Pulumi module uses.
  creation_token = local.resource_name

  # Encryption at rest (ForceNew -- cannot be enabled after creation).
  encrypted  = var.spec.encrypted
  kms_key_id = local.kms_key_id

  # Performance is create-time (ForceNew); throughput is mutable but AWS
  # enforces a 24-hour cooldown between throughput-mode changes.
  performance_mode                = local.performance_mode
  throughput_mode                 = local.throughput_mode
  provisioned_throughput_in_mibps = local.provisioned_throughput_in_mibps

  # One Zone storage (ForceNew). Null keeps the regional multi-AZ shape.
  availability_zone_name = local.availability_zone_name

  # Replication-overwrite protection. AWS defaults to ENABLED; only an
  # explicit value is sent so unset stays indistinguishable from the AWS
  # default. DISABLED is required before this file system can be targeted as
  # a replication destination.
  dynamic "protection" {
    for_each = local.replication_overwrite_protection != null ? [local.replication_overwrite_protection] : []
    content {
      replication_overwrite = protection.value
    }
  }

  # Lifecycle policies -- one block per transition rule (IA, Archive,
  # back-to-primary), the provider's own shape.
  dynamic "lifecycle_policy" {
    for_each = local.lifecycle_policies
    content {
      transition_to_ia                    = try(lifecycle_policy.value.transition_to_ia, null)
      transition_to_archive               = try(lifecycle_policy.value.transition_to_archive, null)
      transition_to_primary_storage_class = try(lifecycle_policy.value.transition_to_primary_storage_class, null)
    }
  }

  tags = local.aws_tags
}

# ---------------------------------------------------------------------------
# Mount Targets -- one per subnet (and thus per AZ)
# ---------------------------------------------------------------------------
# AWS allows at most one mount target per Availability Zone and returns the
# SAME mount-target ID to parallel same-AZ create calls, so the spec declares
# one subnet per AZ. Security groups must allow NFS (TCP 2049) from clients;
# when none are declared, AWS attaches the VPC's default security group.
# ---------------------------------------------------------------------------

resource "aws_efs_mount_target" "this" {
  for_each = local.mount_targets

  file_system_id  = aws_efs_file_system.this.id
  subnet_id       = each.key
  security_groups = local.security_groups

  # Static addresses and the address family are all ForceNew. Empty strings
  # become null so unset keeps the AWS defaults (auto-assigned IPv4).
  ip_address      = each.value.ip_address != "" ? each.value.ip_address : null
  ip_address_type = each.value.ip_address_type != "" ? each.value.ip_address_type : null
  ipv6_address    = each.value.ipv6_address != "" ? each.value.ipv6_address : null
}

# ---------------------------------------------------------------------------
# Backup Policy -- automatic daily backups
# ---------------------------------------------------------------------------
# The resource has no true delete (removal PUTs status DISABLED), so it is
# only materialized when backups are enabled; absent means AWS's default
# (disabled) and a toggle-off flows through the same resource lifecycle.
# ---------------------------------------------------------------------------

resource "aws_efs_backup_policy" "this" {
  count = var.spec.backup_enabled ? 1 : 0

  file_system_id = aws_efs_file_system.this.id

  backup_policy {
    status = "ENABLED"
  }
}

# ---------------------------------------------------------------------------
# File System Policy -- IAM resource policy
# ---------------------------------------------------------------------------
# The spec models the policy as a Struct, which arrives from the tfvars
# pipeline as a NESTED OBJECT (never a JSON string) -- jsonencode() turns it
# into the document the provider expects.
# ---------------------------------------------------------------------------

resource "aws_efs_file_system_policy" "this" {
  count = var.spec.policy != null ? 1 : 0

  file_system_id = aws_efs_file_system.this.id
  policy         = jsonencode(var.spec.policy)

  # Only set when the user deliberately opts out of AWS's lockout safety
  # check (a policy that denies the deploying principal future
  # PutFileSystemPolicy calls is otherwise rejected).
  bypass_policy_lockout_safety_check = var.spec.bypass_policy_lockout_safety_check
}

# ---------------------------------------------------------------------------
# Replication -- cross-region / cross-AZ disaster recovery
# ---------------------------------------------------------------------------
# One-per-file-system and create-time immutable (the provider has no Update;
# any change replaces the configuration). Deleting the replication stops
# syncing but leaves the destination file system in place -- modifying or
# deleting the destination then requires ITS replication_overwrite_protection
# to be DISABLED.
# ---------------------------------------------------------------------------

resource "aws_efs_replication_configuration" "this" {
  count = var.spec.replication != null ? 1 : 0

  source_file_system_id = aws_efs_file_system.this.id

  destination {
    # At least one of region / availability_zone_name is required (CEL-
    # enforced to match the provider's own constraint). An AZ destination
    # creates the replica as a One Zone file system -- the cheaper DR shape.
    region                 = try(var.spec.replication.destination_region, "") != "" ? var.spec.replication.destination_region : null
    availability_zone_name = try(var.spec.replication.destination_availability_zone_name, "") != "" ? var.spec.replication.destination_availability_zone_name : null

    # Replicas are always encrypted; null uses the destination region's
    # AWS-managed key aws/elasticfilesystem.
    kms_key_id = try(var.spec.replication.destination_kms_key_id, "") != "" ? var.spec.replication.destination_kms_key_id : null

    # Replicate into an existing file system instead of letting AWS mint one.
    file_system_id = try(var.spec.replication.destination_file_system_id, "") != "" ? var.spec.replication.destination_file_system_id : null
  }
}
