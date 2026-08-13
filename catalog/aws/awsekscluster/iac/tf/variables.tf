variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsEksCluster specification"
  type = object({
    region = string
    subnet_ids = list(string)
    cluster_role_arn = string
    version = optional(string, "")
    security_group_ids = optional(list(string), [])
    endpoint_public_access = optional(bool)
    endpoint_private_access = optional(bool, false)
    public_access_cidrs = optional(list(string), [])
    control_plane_egress_mode = optional(string, "")
    ip_family = optional(string, "")
    service_ipv4_cidr = optional(string, "")
    enabled_cluster_log_types = optional(list(string), [])
    kms_key_arn = optional(string, "")
    access_config = optional(object({
      authentication_mode = optional(string, "")
      bootstrap_cluster_creator_admin_permissions = optional(bool)
    }))
    auto_mode = optional(object({
      enabled = optional(bool, false)
      node_pools = optional(list(string), [])
      node_role_arn = optional(string, "")
    }))
    upgrade_support_type = optional(string, "")
    zonal_shift_enabled = optional(bool, false)
    deletion_protection = optional(bool, false)
    bootstrap_self_managed_addons = optional(bool)
    force_update_version = optional(bool, false)
    control_plane_scaling_tier = optional(string, "")
    remote_networks = optional(object({
      node_cidrs = optional(list(string), [])
      pod_cidrs = optional(list(string), [])
    }))
  })
}
