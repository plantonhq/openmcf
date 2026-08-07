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
  description = "AwsEcsService specification"
  type = object({
    region = string
    cluster_arn = string
    task_definition = string
    desired_count = optional(number)
    launch_type = optional(string, "")
    capacity_provider_strategy = optional(list(object({
      capacity_provider = string
      base = optional(number, 0)
      weight = optional(number, 0)
    })), [])
    platform_version = optional(string, "")
    scheduling_strategy = optional(string, "")
    network = optional(object({
      subnets = list(string)
      security_groups = optional(list(string), [])
      assign_public_ip = optional(bool, false)
    }))
    load_balancers = optional(list(object({
      target_group_arn = string
      container_name = string
      container_port = optional(number, 0)
      advanced_configuration = optional(object({
        alternate_target_group_arn = string
        production_listener_rule = string
        test_listener_rule = optional(string, "")
        role_arn = string
      }))
    })), [])
    health_check_grace_period_seconds = optional(number)
    deployment_maximum_percent = optional(number)
    deployment_minimum_healthy_percent = optional(number)
    deployment_circuit_breaker = optional(object({
      enable = optional(bool, false)
      rollback = optional(bool, false)
    }))
    alarms = optional(object({
      alarm_names = list(string)
      enable = optional(bool, false)
      rollback = optional(bool, false)
    }))
    deployment_configuration = optional(object({
      strategy = optional(string, "")
      bake_time_in_minutes = optional(number)
      canary_configuration = optional(object({
        canary_percent = optional(number, 0)
        canary_bake_time_in_minutes = optional(number)
      }))
      linear_configuration = optional(object({
        step_percent = optional(number, 0)
        step_bake_time_in_minutes = optional(number)
      }))
      lifecycle_hooks = optional(list(object({
        hook_target_arn = string
        role_arn = string
        lifecycle_stages = list(string)
        hook_details = optional(string, "")
      })), [])
    }))
    deployment_controller = optional(string, "")
    service_connect = optional(object({
      enabled = optional(bool, false)
      namespace = optional(string, "")
      services = optional(list(object({
        port_name = string
        discovery_name = optional(string, "")
        client_alias = optional(object({
          port = optional(number, 0)
          dns_name = optional(string, "")
        }))
        ingress_port_override = optional(number)
        timeout = optional(object({
          idle_timeout_seconds = optional(number, 0)
          per_request_timeout_seconds = optional(number, 0)
        }))
        tls = optional(object({
          aws_pca_authority_arn = string
          kms_key = optional(string, "")
          role_arn = optional(string, "")
        }))
      })), [])
      log_configuration = optional(object({
        log_driver = string
        options = optional(map(string), {})
        secret_options = optional(map(string), {})
      }))
    }))
    service_registries = optional(object({
      registry_arn = string
      container_name = optional(string, "")
      container_port = optional(number)
      port = optional(number)
    }))
    volume_configuration = optional(object({
      name = string
      managed_ebs_volume = object({
        role_arn = string
        size_in_gb = optional(number, 0)
        volume_type = optional(string, "")
        iops = optional(number, 0)
        throughput = optional(number, 0)
        encrypted = optional(bool)
        kms_key_id = optional(string, "")
        snapshot_id = optional(string, "")
        file_system_type = optional(string, "")
      })
    }))
    ordered_placement_strategy = optional(list(object({
      type = optional(string, "")
      field = optional(string, "")
    })), [])
    placement_constraints = optional(list(object({
      type = optional(string, "")
      expression = optional(string, "")
    })), [])
    availability_zone_rebalancing = optional(string, "")
    propagate_tags = optional(string, "")
    enable_ecs_managed_tags = optional(bool, false)
    enable_execute_command = optional(bool, false)
    force_delete = optional(bool, false)
    autoscaling = optional(object({
      min_tasks = optional(number, 0)
      max_tasks = optional(number, 0)
      cpu = optional(object({
        target_percent = optional(number, 0)
        scale_in_cooldown_seconds = optional(number)
        scale_out_cooldown_seconds = optional(number)
        disable_scale_in = optional(bool, false)
      }))
      memory = optional(object({
        target_percent = optional(number, 0)
        scale_in_cooldown_seconds = optional(number)
        scale_out_cooldown_seconds = optional(number)
        disable_scale_in = optional(bool, false)
      }))
      requests_per_target = optional(object({
        target_requests_per_target = optional(number, 0)
        load_balancer_arn_suffix = string
        target_group_arn_suffix = string
        scale_in_cooldown_seconds = optional(number)
        scale_out_cooldown_seconds = optional(number)
        disable_scale_in = optional(bool, false)
      }))
    }))
  })
}
