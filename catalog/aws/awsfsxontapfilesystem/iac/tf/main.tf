# ---------------------------------------------------------------------------
# AWS FSx for NetApp ONTAP file system
# ---------------------------------------------------------------------------
# One aws_fsx_ontap_file_system resource carries the whole spec. ForceNew
# attributes (replace the file system when changed): deployment_type,
# storage_type, subnet_ids, preferred_subnet_id, security_group_ids,
# kms_key_id, and endpoint_ip_address_range. On the first generation the
# scaling knobs are also destructive: increasing
# throughput_capacity_per_ha_pair on SINGLE_AZ_1/MULTI_AZ_1, or ha_pairs on
# SINGLE_AZ_1, replaces the file system (the provider models this as a
# conditional ForceNew). Storage capacity, IOPS, passwords, route tables, and
# the maintenance/backup windows update in place.
#
# Note there are NO backup-skip or tag-copy arguments here: ONTAP backups are
# volume-scoped — those decisions live on aws_fsx_ontap_volume.
# ---------------------------------------------------------------------------

resource "aws_fsx_ontap_file_system" "this" {
  # Core shape. deployment_type is always sent — the spec default
  # (SINGLE_AZ_2) is materialized before the module runs.
  deployment_type = var.spec.deployment_type

  # ONTAP file systems are SSD-only (the spec enforces it); sent explicitly
  # so the plan states the storage class rather than relying on provider
  # defaulting.
  storage_type = var.spec.storage_type

  # Capacity scales with HA pairs (1024–524288 GiB per pair; the spec's CEL
  # rules carry the formula). Grows in place, never shrinks.
  storage_capacity = var.spec.storage_capacity_gib

  # Throughput is sized through exactly one arm (spec-enforced XOR):
  # whole-file-system throughput_capacity (the first-generation sizing) or
  # per-HA-pair throughput_capacity_per_ha_pair (required for scale-out and
  # the second generation's tiers). The unset arm passes null and the
  # provider omits it — sending both is an AWS error.
  throughput_capacity             = var.spec.throughput_capacity
  throughput_capacity_per_ha_pair = var.spec.throughput_capacity_per_ha_pair

  # Scale-out: >1 only on SINGLE_AZ_2 (spec-enforced). AWS requires the
  # per-HA-pair throughput arm to be re-sent whenever HA pairs change; both
  # values flowing from the spec keeps that invariant automatically.
  ha_pairs = var.spec.ha_pairs

  # Networking (ForceNew). Single-AZ: one subnet; multi-AZ: two subnets with
  # the active file server in preferred_subnet_id. The provider requires
  # preferred_subnet_id for every deployment type, so single-AZ derives it
  # from the only subnet (see locals).
  subnet_ids          = var.spec.subnet_ids
  preferred_subnet_id = local.preferred_subnet_id
  security_group_ids  = local.security_group_ids

  # Multi-AZ floating endpoint range + the route tables AWS manages for
  # failover. Both omitted for single-AZ (and CEL-rejected in the spec).
  endpoint_ip_address_range = local.endpoint_ip_address_range
  route_table_ids           = local.route_table_ids

  # Encryption at rest: null falls back to the AWS-managed FSx key.
  kms_key_id = local.kms_key_id

  # ONTAP CLI / REST API access for the fsxadmin user. Updatable in place.
  fsx_admin_password = local.fsx_admin_password

  # SSD IOPS: AUTOMATIC (3 IOPS/GiB) unless USER_PROVISIONED sets an exact
  # figure. Updatable in place.
  dynamic "disk_iops_configuration" {
    for_each = var.spec.disk_iops_configuration != null ? [var.spec.disk_iops_configuration] : []
    content {
      mode = disk_iops_configuration.value.mode
      iops = disk_iops_configuration.value.iops > 0 ? disk_iops_configuration.value.iops : null
    }
  }

  # Automatic backups. Zero is a real value ("no automatic backups") — the
  # resolved default flows through as-is.
  automatic_backup_retention_days   = var.spec.automatic_backup_retention_days
  daily_automatic_backup_start_time = local.daily_automatic_backup_start_time

  # Maintenance window ("d:HH:MM", UTC).
  weekly_maintenance_start_time = local.weekly_maintenance_start_time

  tags = local.aws_tags
}
