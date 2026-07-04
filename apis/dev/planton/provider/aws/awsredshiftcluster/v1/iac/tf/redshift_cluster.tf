# The cluster composes onto its neighbors instead of embedding them:
# subnets, security groups, IAM roles, KMS keys, and the Elastic IP
# attach by reference, and warehouse ingress rules live on the
# referenced AwsSecurityGroup nodes -- this module never creates or
# mutates a resource that deserves to be its own node.
#
# Create-only in AWS: the identifier, the subnet group, the master
# username, and the snapshot restore sources. Node type and node count
# are NOT create-only -- changing them triggers an in-place (but
# access-interrupting) elastic/classic resize, never a replacement.
resource "aws_redshift_cluster" "this" {
  cluster_identifier = local.cluster_identifier

  # Compute topology. number_of_nodes 0 keeps the AWS default (1,
  # single-node); the provider derives cluster_type (single-node vs
  # multi-node) from the count, so it is never set here.
  node_type       = var.spec.node_type
  number_of_nodes = var.spec.number_of_nodes != 0 ? var.spec.number_of_nodes : null

  # Empty pins nothing: "1.0" is the only version family Redshift has
  # ever shipped; actual engine patches ride the maintenance track.
  cluster_version = var.spec.cluster_version != "" ? var.spec.cluster_version : null

  # Empty keeps the AWS default first database ("dev").
  database_name = var.spec.database_name != "" ? var.spec.database_name : null

  # Networking: the subnet group managed here (or referenced), the VPC
  # default SG when no groups are given (AWS's own default), and
  # AWS-picked AZ placement unless explicitly pinned.
  cluster_subnet_group_name = local.manage_subnet_group ? aws_redshift_subnet_group.this[0].name : (var.spec.cluster_subnet_group_name != "" ? var.spec.cluster_subnet_group_name : null)
  vpc_security_group_ids    = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
  availability_zone         = var.spec.availability_zone != "" ? var.spec.availability_zone : null

  # Zone relocation moves the single cluster between AZs on outage or
  # demand (RA3-only; the port must sit in 5431-5455 or 8191-8215 --
  # the default 5439 qualifies). Mutually exclusive with multi_az
  # (CEL-enforced): relocation moves, Multi-AZ fails over to a standby.
  availability_zone_relocation_enabled = var.spec.availability_zone_relocation_enabled
  multi_az                             = var.spec.multi_az

  # 0 keeps the AWS default (5439).
  port = var.spec.port != 0 ? var.spec.port : null

  # Public reachability is opt-in; the static leader-node address only
  # exists on a public cluster (CEL ties elastic_ip to
  # publicly_accessible so the constraint fails at validate, not at
  # deploy).
  publicly_accessible = var.spec.publicly_accessible
  elastic_ip          = var.spec.elastic_ip != "" ? var.spec.elastic_ip : null

  # Enhanced VPC routing forces COPY/UNLOAD data movement through the
  # VPC where flow logs and endpoints can see and govern it.
  enhanced_vpc_routing = var.spec.enhanced_vpc_routing

  master_username = var.spec.master_username != "" ? var.spec.master_username : null

  # The password contract (CEL enforces exactly one strategy): the
  # AWS-managed Secrets Manager secret (recommended -- no secret in
  # manifest or state) or a directly supplied password.
  # manage_master_password is forwarded ONLY when true: an explicit
  # false conflicts with master_password in the provider's
  # ConflictsWith machinery.
  manage_master_password            = var.spec.manage_master_password ? true : null
  master_password                   = var.spec.master_password != "" ? var.spec.master_password : null
  master_password_secret_kms_key_id = var.spec.master_password_secret_kms_key_id != "" ? var.spec.master_password_secret_kms_key_id : null

  # Encryption at rest. AWS defaults new clusters to encrypted and the
  # spec keeps that default; toggling later is an in-place but
  # long-running migration. The provider models encrypted as a nullable
  # bool (string under the hood) -- Terraform converts the bool cleanly.
  encrypted  = var.spec.encrypted
  kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # IAM roles the warehouse assumes for COPY/UNLOAD/Spectrum. The
  # default role must also be in iam_roles (an AWS requirement the
  # error message makes obvious enough to leave to the API).
  iam_roles            = length(var.spec.iam_roles) > 0 ? var.spec.iam_roles : null
  default_iam_role_arn = var.spec.default_iam_role_arn != "" ? var.spec.default_iam_role_arn : null

  # Snapshots: automated retention bounds recovery; manual retention 0
  # keeps the AWS default (-1, indefinite).
  automated_snapshot_retention_period = var.spec.automated_snapshot_retention_period
  manual_snapshot_retention_period    = var.spec.manual_snapshot_retention_period != 0 ? var.spec.manual_snapshot_retention_period : null

  # Maintenance: empty window lets AWS assign one; empty track keeps
  # "current" (the AWS default).
  preferred_maintenance_window = var.spec.preferred_maintenance_window != "" ? var.spec.preferred_maintenance_window : null
  maintenance_track_name       = var.spec.maintenance_track_name != "" ? var.spec.maintenance_track_name : null
  allow_version_upgrade        = var.spec.allow_version_upgrade
  apply_immediately            = var.spec.apply_immediately

  # Deletion safety: the CEL contract requires a final-snapshot name
  # unless skipping is explicit, so a delete can never fail late on a
  # missing snapshot identifier.
  skip_final_snapshot       = var.spec.skip_final_snapshot
  final_snapshot_identifier = var.spec.final_snapshot_identifier != "" ? var.spec.final_snapshot_identifier : null

  # Create-time restore sources (mutually exclusive, CEL-enforced):
  # a snapshot by name (with optional source-cluster disambiguation) or
  # by ARN (the cross-account/cross-region sharing shape). owner_account
  # covers snapshots shared by another AWS account. A restored cluster
  # inherits the snapshot's credentials, so master_username stays empty.
  snapshot_identifier         = var.spec.snapshot_identifier != "" ? var.spec.snapshot_identifier : null
  snapshot_arn                = var.spec.snapshot_arn != "" ? var.spec.snapshot_arn : null
  snapshot_cluster_identifier = var.spec.snapshot_cluster_identifier != "" ? var.spec.snapshot_cluster_identifier : null
  owner_account               = var.spec.owner_account != "" ? var.spec.owner_account : null

  # Parameter groups: the managed inline group, an existing referenced
  # group, or the Redshift default.
  cluster_parameter_group_name = local.effective_parameter_group

  tags = local.aws_tags
}
