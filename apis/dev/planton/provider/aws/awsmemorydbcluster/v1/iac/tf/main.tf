# Amazon MemoryDB: a durable, Multi-AZ, Redis/Valkey-compatible in-memory
# database. Topology is always sharded; authentication is ACL-based (the
# cluster's only auth model -- the acl_name below is the attachment point).
#
# The module owns a subnet group and a parameter group only when the folded
# arms (subnet_ids / parameters) are used; the bring-your-own name arms pass
# through untouched -- this module never mutates a resource it merely
# references.

resource "aws_memorydb_subnet_group" "this" {
  count = local.create_subnet_group ? 1 : 0

  # The group carries the cluster's own name: everything the module owns is
  # discoverable by one identifier, on both engines.
  name        = local.cluster_name
  description = "MemoryDB subnet group for ${local.cluster_name}"
  subnet_ids  = var.spec.subnet_ids
  tags        = local.aws_tags
}

resource "aws_memorydb_parameter_group" "this" {
  count = local.create_parameter_group ? 1 : 0

  # "-params" keeps the group distinct from the cluster while remaining
  # discoverable by the cluster's name, on both engines.
  name        = "${local.cluster_name}-params"
  family      = var.spec.parameter_group_family
  description = "Custom parameter group for ${local.cluster_name}"
  tags        = local.aws_tags

  # Removing an entry resets that parameter to its family default (AWS
  # semantics -- the provider issues a ResetParameterGroup for dropped
  # names).
  dynamic "parameter" {
    for_each = var.spec.parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }
}

resource "aws_memorydb_cluster" "this" {
  # Create-time immutable; doubles as the Name tag. metadata.name on both
  # engines -- never provider auto-naming.
  name = local.cluster_name

  # The ACL ref arrives pre-resolved (a literal like "open-access" or a
  # flattened reference to an AwsMemorydbAcl's exported name). Updates in
  # place. AWS couples it to TLS at create: tls_enabled=false only accepts
  # "open-access".
  acl_name  = var.spec.acl_name
  node_type = var.spec.node_type

  # Description is ALWAYS sent explicitly -- left to their own defaults the
  # two providers inject differing "Managed by ..." strings and the
  # engines' state permanently diverges.
  description = var.spec.description

  # Engine (redis/valkey; AWS supports redis -> valkey in place, never the
  # reverse). An empty version lets AWS pick the engine default.
  engine         = var.spec.engine
  engine_version = var.spec.engine_version != "" ? var.spec.engine_version : null

  # ForceNew.
  port = var.spec.port

  # Topology -- both dials scale in place (resharding redistributes slots
  # online; replicas roll one at a time).
  num_shards             = var.spec.num_shards
  num_replicas_per_shard = var.spec.num_replicas_per_shard

  # Networking: exactly one subnet-group source (module-owned or
  # bring-your-own; CEL keeps the arms exclusive). Neither set falls back
  # to AWS's account "default" subnet group.
  subnet_group_name  = local.effective_subnet_group_name
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # Dual-stack networking. network_type is ForceNew; ip_discovery (which
  # stack CLUSTER SLOTS/SHARDS report to clients) updates in place.
  network_type = var.spec.network_type != "" ? var.spec.network_type : null
  ip_discovery = var.spec.ip_discovery != "" ? var.spec.ip_discovery : null

  # Encryption. TLS default true (ForceNew). At-rest encryption is always
  # on; the KMS ref merely substitutes a customer-managed key (ForceNew).
  tls_enabled = var.spec.tls_enabled
  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  # Maintenance and snapshots (all update in place; final_snapshot_name is
  # consumed only at delete).
  maintenance_window       = var.spec.maintenance_window != "" ? var.spec.maintenance_window : null
  snapshot_retention_limit = var.spec.snapshot_retention_limit
  snapshot_window          = var.spec.snapshot_window != "" ? var.spec.snapshot_window : null
  final_snapshot_name      = var.spec.final_snapshot_name != "" ? var.spec.final_snapshot_name : null

  # Create-time restore sources (mutually exclusive by CEL; ForceNew).
  snapshot_arns = length(var.spec.snapshot_arns) > 0 ? var.spec.snapshot_arns : null
  snapshot_name = var.spec.snapshot_name != "" ? var.spec.snapshot_name : null

  # Parameter group: exactly one source (module-owned or bring-your-own).
  parameter_group_name = local.effective_parameter_group_name

  # Multi-region active-active membership (ForceNew): the multi-region
  # cluster is created outside this resource; this regional cluster joins
  # it by name.
  multi_region_cluster_name = var.spec.multi_region_cluster_name != "" ? var.spec.multi_region_cluster_name : null

  # Notifications (in place; clearing disables).
  sns_topic_arn = var.spec.sns_topic_arn != "" ? var.spec.sns_topic_arn : null

  # Advanced (both ForceNew).
  auto_minor_version_upgrade = var.spec.auto_minor_version_upgrade
  data_tiering               = var.spec.data_tiering

  tags = local.aws_tags
}
