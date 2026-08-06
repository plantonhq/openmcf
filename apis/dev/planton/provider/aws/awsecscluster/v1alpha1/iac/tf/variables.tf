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
  description = "AwsEcsCluster specification"
  type = object({
    region = string
    container_insights = optional(string, "")
    capacity_providers = optional(list(string), [])
    ec2_capacity_providers = optional(list(object({
      name = string
      auto_scaling_group_arn = string
      managed_scaling = optional(object({
        status = optional(string, "")
        target_capacity = optional(number, 0)
        minimum_scaling_step_size = optional(number, 0)
        maximum_scaling_step_size = optional(number, 0)
        instance_warmup_period_seconds = optional(number, 0)
      }))
      managed_termination_protection = optional(string, "")
      managed_draining = optional(string, "")
    })), [])
    default_capacity_provider_strategy = optional(list(object({
      capacity_provider = string
      base = optional(number, 0)
      weight = optional(number, 0)
    })), [])
    execute_command_configuration = optional(object({
      logging = optional(string, "")
      log_configuration = optional(object({
        cloud_watch_log_group_name = optional(string, "")
        cloud_watch_encryption_enabled = optional(bool, false)
        s3_bucket_name = optional(string, "")
        s3_key_prefix = optional(string, "")
        s3_bucket_encryption_enabled = optional(bool, false)
      }))
      kms_key_id = optional(string, "")
    }))
    managed_storage_configuration = optional(object({
      fargate_ephemeral_storage_kms_key_id = optional(string, "")
      kms_key_id = optional(string, "")
    }))
    service_connect_namespace_arn = optional(string, "")
  })
}
