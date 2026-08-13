# The managed node group composes onto its neighbors instead of embedding
# them: the cluster attaches by reference, the node IAM role is a
# referenced AwsIamRole that carries its own worker policies (this module
# never modifies a role it merely references), and launch mechanics come
# either from the inline knobs or from a referenced AwsLaunchTemplate --
# the spec's CEL rules enforce AWS's mutual exclusions between the two
# styles before anything reaches the API.
#
# Create-only at AWS: the name, cluster, role, subnets, ami_type,
# capacity_type, disk_size, instance_types, and the whole remote_access
# block. Labels, taints, scaling, update/repair config update in place;
# version/release_version/launch_template changes roll the group node by
# node.
resource "aws_eks_node_group" "this" {
  node_group_name = local.node_group_name
  cluster_name    = var.spec.cluster_name
  node_role_arn   = var.spec.node_role_arn
  subnet_ids      = var.spec.subnet_ids

  # Inline launch mechanics: sent only when set so AWS's defaults keep
  # applying (t3.medium; 20 GiB Linux / 50 GiB Windows root disk; AMI
  # family inferred from the instance types).
  instance_types = length(var.spec.instance_types) > 0 ? var.spec.instance_types : null
  disk_size      = var.spec.disk_size_gb > 0 ? var.spec.disk_size_gb : null
  ami_type       = var.spec.ami_type != "" ? var.spec.ami_type : null
  capacity_type  = local.capacity_type

  scaling_config {
    min_size     = var.spec.scaling.min_size
    max_size     = var.spec.scaling.max_size
    desired_size = var.spec.scaling.desired_size
  }

  # Launch-template style: the template owns AMI, instance type, storage,
  # and network posture.
  dynamic "launch_template" {
    for_each = var.spec.launch_template != null ? [var.spec.launch_template] : []
    content {
      id = launch_template.value.launch_template_id
      # The provider requires a version. "$Default" follows the template's
      # default version -- the setup that lets promoting a template version
      # roll the fleet. Pin a numeric version for fully drift-free plans
      # (AWS reads "$Default" back as the resolved number).
      version = launch_template.value.version != "" ? launch_template.value.version : "$Default"
    }
  }

  # Immutable at AWS: any change here replaces the group. SSH without
  # explicit source groups makes AWS open port 22 to the internet -- the
  # created group is surfaced as the remote_access_sg_id output.
  dynamic "remote_access" {
    for_each = var.spec.remote_access != null ? [var.spec.remote_access] : []
    content {
      ec2_ssh_key               = remote_access.value.ec2_ssh_key != "" ? remote_access.value.ec2_ssh_key : null
      source_security_group_ids = length(remote_access.value.source_security_group_ids) > 0 ? remote_access.value.source_security_group_ids : null
    }
  }

  labels = length(var.spec.labels) > 0 ? var.spec.labels : null

  # Taints keep ordinary pods off dedicated capacity until they tolerate
  # them; updates diff in place (never a group replacement).
  dynamic "taint" {
    for_each = var.spec.taints
    content {
      key    = taint.value.key
      value  = taint.value.value != "" ? taint.value.value : null
      effect = taint.value.effect
    }
  }

  # Exactly one unavailability form is set (spec validation enforces it),
  # mirroring AWS's ExactlyOneOf on the same fields.
  dynamic "update_config" {
    for_each = var.spec.update_config != null ? [var.spec.update_config] : []
    content {
      max_unavailable            = update_config.value.max_unavailable > 0 ? update_config.value.max_unavailable : null
      max_unavailable_percentage = update_config.value.max_unavailable_percentage > 0 ? update_config.value.max_unavailable_percentage : null
      update_strategy            = update_config.value.update_strategy != "" ? update_config.value.update_strategy : null
    }
  }

  # Managed node auto-repair: EKS replaces or reboots nodes the cluster
  # reports unhealthy, within these parallelism and threshold bounds.
  dynamic "node_repair_config" {
    for_each = var.spec.node_repair_config != null ? [var.spec.node_repair_config] : []
    content {
      enabled = node_repair_config.value.enabled

      max_parallel_nodes_repaired_count       = node_repair_config.value.max_parallel_nodes_repaired_count > 0 ? node_repair_config.value.max_parallel_nodes_repaired_count : null
      max_parallel_nodes_repaired_percentage  = node_repair_config.value.max_parallel_nodes_repaired_percentage > 0 ? node_repair_config.value.max_parallel_nodes_repaired_percentage : null
      max_unhealthy_node_threshold_count      = node_repair_config.value.max_unhealthy_node_threshold_count > 0 ? node_repair_config.value.max_unhealthy_node_threshold_count : null
      max_unhealthy_node_threshold_percentage = node_repair_config.value.max_unhealthy_node_threshold_percentage > 0 ? node_repair_config.value.max_unhealthy_node_threshold_percentage : null

      dynamic "node_repair_config_overrides" {
        for_each = node_repair_config.value.overrides
        content {
          min_repair_wait_time_mins = node_repair_config_overrides.value.min_repair_wait_time_mins
          node_monitoring_condition = node_repair_config_overrides.value.node_monitoring_condition
          node_unhealthy_reason     = node_repair_config_overrides.value.node_unhealthy_reason
          repair_action             = node_repair_config_overrides.value.repair_action
        }
      }
    }
  }

  # A pool of pre-initialized nodes that cuts scale-out from minutes to
  # seconds; updates in place. Stopped pools cost only EBS storage.
  dynamic "warm_pool_config" {
    for_each = var.spec.warm_pool_config != null ? [var.spec.warm_pool_config] : []
    content {
      pool_state = warm_pool_config.value.pool_state != "" ? warm_pool_config.value.pool_state : null
      min_size   = warm_pool_config.value.min_size > 0 ? warm_pool_config.value.min_size : null
      # Explicit 0 is meaningful (no prepared capacity beyond min_size), so
      # presence (null vs a value) decides what is sent; null keeps AWS's
      # default ceiling (max_size minus desired capacity).
      max_group_prepared_capacity = warm_pool_config.value.max_group_prepared_capacity
      reuse_on_scale_in           = warm_pool_config.value.reuse_on_scale_in ? true : null
    }
  }

  # Version changes roll the group node by node; release_version pins the
  # exact AMI release within the minor.
  version              = var.spec.version != "" ? var.spec.version : null
  release_version      = var.spec.release_version != "" ? var.spec.release_version : null
  force_update_version = var.spec.force_update_version ? true : null

  tags = local.aws_tags
}
