variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsEcsCluster specification"
  type = object({
    region             = string
    container_insights = optional(string, "")
    capacity_providers = optional(list(string), [])
    ec2_capacity_providers = optional(list(object({
      name                   = string
      auto_scaling_group_arn = string
      managed_scaling = optional(object({
        status                         = optional(string, "")
        target_capacity                = optional(number, 0)
        minimum_scaling_step_size      = optional(number, 0)
        maximum_scaling_step_size      = optional(number, 0)
        instance_warmup_period_seconds = optional(number, 0)
      }))
      managed_termination_protection = optional(string, "")
      managed_draining               = optional(string, "")
    })), [])
    managed_instances_capacity_providers = optional(list(object({
      name                    = string
      infrastructure_role_arn = string
      instance_launch_template = object({
        ec2_instance_profile_arn = string
        network_configuration = object({
          subnets         = list(string)
          security_groups = optional(list(string), [])
        })
        capacity_option_type = optional(string, "")
        capacity_reservations = optional(object({
          reservation_preference = optional(string, "")
          reservation_group_arn  = optional(string, "")
        }))
        instance_requirements = optional(object({
          memory_mib = object({
            min = number
            max = optional(number, 0)
          })
          vcpu_count = object({
            min = number
            max = optional(number, 0)
          })
          allowed_instance_types                                  = optional(list(string), [])
          excluded_instance_types                                 = optional(list(string), [])
          instance_generations                                    = optional(list(string), [])
          cpu_manufacturers                                       = optional(list(string), [])
          bare_metal                                              = optional(string, "")
          burstable_performance                                   = optional(string, "")
          require_hibernate_support                               = optional(bool, false)
          spot_max_price_percentage_over_lowest_price             = optional(number, 0)
          max_spot_price_as_percentage_of_optimal_on_demand_price = optional(number, 0)
          on_demand_max_price_percentage_over_lowest_price        = optional(number, 0)
          local_storage                                           = optional(string, "")
          local_storage_types                                     = optional(list(string), [])
          total_local_storage_gb = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          memory_gib_per_vcpu = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          network_interface_count = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          network_bandwidth_gbps = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          baseline_ebs_bandwidth_mbps = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          accelerator_count = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
          accelerator_manufacturers = optional(list(string), [])
          accelerator_names         = optional(list(string), [])
          accelerator_types         = optional(list(string), [])
          accelerator_total_memory_mib = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
        }))
        use_local_storage = optional(bool)
        monitoring        = optional(string, "")
        storage_size_gib  = optional(number, 0)
      })
      scale_in_after_seconds = optional(number)
      propagate_tags         = optional(string, "")
    })), [])
    default_capacity_provider_strategy = optional(list(object({
      capacity_provider = string
      base              = optional(number, 0)
      weight            = optional(number, 0)
    })), [])
    execute_command_configuration = optional(object({
      logging = optional(string, "")
      log_configuration = optional(object({
        cloud_watch_log_group_name     = optional(string, "")
        cloud_watch_encryption_enabled = optional(bool, false)
        s3_bucket_name                 = optional(string, "")
        s3_key_prefix                  = optional(string, "")
        s3_bucket_encryption_enabled   = optional(bool, false)
      }))
      kms_key_id = optional(string, "")
    }))
    managed_storage_configuration = optional(object({
      fargate_ephemeral_storage_kms_key_id = optional(string, "")
      kms_key_id                           = optional(string, "")
    }))
    service_connect_namespace_arn = optional(string, "")
  })
}
