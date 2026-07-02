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
  description = "AwsEksNodeGroup specification"
  type = object({
    region = string
    cluster_name = string
    node_role_arn = string
    subnet_ids = list(string)
    launch_template = optional(object({
      launch_template_id = string
      version = optional(string, "")
    }))
    instance_types = optional(list(string), [])
    ami_type = optional(string, "")
    capacity_type = optional(string, "")
    disk_size_gb = optional(number, 0)
    scaling = object({
      min_size = optional(number, 0)
      max_size = optional(number, 0)
      desired_size = optional(number, 0)
    })
    remote_access = optional(object({
      ec2_ssh_key = optional(string, "")
      source_security_group_ids = optional(list(string), [])
    }))
    labels = optional(map(string), {})
    taints = optional(list(object({
      key = string
      value = optional(string, "")
      effect = string
    })), [])
    update_config = optional(object({
      max_unavailable = optional(number, 0)
      max_unavailable_percentage = optional(number, 0)
      update_strategy = optional(string, "")
    }))
    node_repair_config = optional(object({
      enabled = optional(bool, false)
      max_parallel_nodes_repaired_count = optional(number, 0)
      max_parallel_nodes_repaired_percentage = optional(number, 0)
      max_unhealthy_node_threshold_count = optional(number, 0)
      max_unhealthy_node_threshold_percentage = optional(number, 0)
      overrides = optional(list(object({
        min_repair_wait_time_mins = optional(number, 0)
        node_monitoring_condition = string
        node_unhealthy_reason = string
        repair_action = string
      })), [])
    }))
    version = optional(string, "")
    release_version = optional(string, "")
    force_update_version = optional(bool, false)
  })
}
