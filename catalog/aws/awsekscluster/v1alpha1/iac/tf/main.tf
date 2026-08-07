# The cluster is only the CONTROL PLANE: the managed API server, etcd, and
# cluster-level posture. Compute attaches as separate composable nodes
# (AwsEksNodeGroup, or EKS Auto Mode enabled below), and the cluster role
# is a referenced AwsIamRole that carries its own AmazonEKSClusterPolicy --
# this module never modifies a role it merely references.
#
# Create-only at AWS: the name, the cluster role, ip_family,
# service_ipv4_cidr, bootstrap flags -- and two one-way doors handled by
# AWS as forced replacement: disabling envelope encryption after enabling
# it, and reverting CUSTOMER_ROUTED egress back to AWS_MANAGED.
resource "aws_eks_cluster" "this" {
  name     = local.cluster_name
  role_arn = var.spec.cluster_role_arn

  # An empty version lets AWS pick its current default; EKS upgrades one
  # minor at a time and never downgrades.
  version              = var.spec.version != "" ? var.spec.version : null
  force_update_version = var.spec.force_update_version ? true : null

  vpc_config {
    subnet_ids         = var.spec.subnet_ids
    security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

    # Endpoint exposure is a pair of independent toggles at AWS (public
    # defaults true, private defaults false). endpoint_public_access is
    # proto-optional so an explicit false ("private-only cluster") is
    # distinguishable from unset: the contract delivers null when unset and
    # the value passes straight through.
    endpoint_public_access  = var.spec.endpoint_public_access
    endpoint_private_access = var.spec.endpoint_private_access

    public_access_cidrs = length(var.spec.public_access_cidrs) > 0 ? var.spec.public_access_cidrs : null

    # Sent only when set so AWS's default (AWS_MANAGED) keeps applying.
    # Reverting CUSTOMER_ROUTED to AWS_MANAGED forces cluster replacement.
    # PARITY-EXCEPTION: pulumi-aws (v7.35.0) does not yet model
    # control_plane_egress_mode, so only this engine implements it. Revisit
    # on the next pulumi-aws upgrade. Stack outputs are unaffected.
    control_plane_egress_mode = var.spec.control_plane_egress_mode != "" ? var.spec.control_plane_egress_mode : null
  }

  # Address-family settings are create-only; the elastic_load_balancing
  # capability inside the same block belongs to Auto Mode and must move in
  # lockstep with compute_config/storage_config below.
  dynamic "kubernetes_network_config" {
    for_each = (var.spec.ip_family != "" || var.spec.service_ipv4_cidr != "" || local.auto_mode_enabled) ? [1] : []
    content {
      ip_family         = var.spec.ip_family != "" ? var.spec.ip_family : null
      service_ipv4_cidr = var.spec.service_ipv4_cidr != "" ? var.spec.service_ipv4_cidr : null

      dynamic "elastic_load_balancing" {
        for_each = local.auto_mode_enabled ? [1] : []
        content {
          enabled = true
        }
      }
    }
  }

  # Envelope encryption is a one-way door: AWS only ever associates a key
  # (never dissociates or re-keys in place), so the block is sent only when
  # a key is configured. "secrets" is the only resource type the EKS API
  # accepts here, which is why the spec folds it away.
  dynamic "encryption_config" {
    for_each = var.spec.kms_key_arn != "" ? [1] : []
    content {
      provider {
        key_arn = var.spec.kms_key_arn
      }
      resources = ["secrets"]
    }
  }

  dynamic "access_config" {
    for_each = var.spec.access_config != null ? [var.spec.access_config] : []
    content {
      authentication_mode = access_config.value.authentication_mode != "" ? access_config.value.authentication_mode : null
      # Proto-optional: explicit false ("no creator admin") must be
      # distinguishable from unset (AWS defaults to true). Create-only.
      bootstrap_cluster_creator_admin_permissions = access_config.value.bootstrap_cluster_creator_admin_permissions
    }
  }

  # EKS Auto Mode compute: AWS provisions and scales EC2 capacity itself.
  # The node role is required by AWS when built-in node pools are enabled
  # (spec validation enforces it).
  dynamic "compute_config" {
    for_each = local.auto_mode_enabled ? [var.spec.auto_mode] : []
    content {
      enabled       = true
      node_pools    = length(compute_config.value.node_pools) > 0 ? compute_config.value.node_pools : null
      node_role_arn = compute_config.value.node_role_arn != "" ? compute_config.value.node_role_arn : null
    }
  }

  # Auto Mode block storage -- the third leg of the all-or-nothing trio.
  dynamic "storage_config" {
    for_each = local.auto_mode_enabled ? [1] : []
    content {
      block_storage {
        enabled = true
      }
    }
  }

  # Control-plane logs stream to CloudWatch; an empty set keeps logging
  # off exactly as AWS defaults it.
  enabled_cluster_log_types = length(var.spec.enabled_cluster_log_types) > 0 ? var.spec.enabled_cluster_log_types : null

  dynamic "upgrade_policy" {
    for_each = var.spec.upgrade_support_type != "" ? [var.spec.upgrade_support_type] : []
    content {
      support_type = upgrade_policy.value
    }
  }

  dynamic "zonal_shift_config" {
    for_each = var.spec.zonal_shift_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  # Always send the explicit boolean: the provider attribute is
  # Optional+Computed, so a null config means "keep the cluster's current
  # value" -- omitting false would make protection impossible to turn off
  # once enabled (and destroys would stay blocked forever).
  deletion_protection = var.spec.deletion_protection

  # Proto-optional: explicit false ("bring your own add-ons") must be
  # distinguishable from unset (AWS defaults to true). Create-only.
  bootstrap_self_managed_addons = var.spec.bootstrap_self_managed_addons

  tags = local.aws_tags
}
