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
  description = "AwsAutoScalingGroup specification"
  type = object({
    region = string
    subnets = list(string)
    launch_template = optional(object({
      launch_template_id = string
      version = optional(string, "")
    }))
    mixed_instances_policy = optional(object({
      launch_template = object({
        launch_template_id = string
        version = optional(string, "")
      })
      overrides = optional(list(object({
        instance_type = optional(string, "")
        weighted_capacity = optional(number, 0)
        launch_template = optional(object({
          launch_template_id = string
          version = optional(string, "")
        }))
        instance_requirements = optional(object({
          memory_mib = object({
            min = optional(number, 0)
            max = optional(number, 0)
          })
          vcpu_count = object({
            min = optional(number, 0)
            max = optional(number, 0)
          })
          allowed_instance_types = optional(list(string), [])
          excluded_instance_types = optional(list(string), [])
          instance_generations = optional(list(string), [])
          cpu_manufacturers = optional(list(string), [])
          bare_metal = optional(string, "")
          burstable_performance = optional(string, "")
          require_hibernate_support = optional(bool, false)
          spot_max_price_percentage_over_lowest_price = optional(number, 0)
          max_spot_price_as_percentage_of_optimal_on_demand_price = optional(number, 0)
          on_demand_max_price_percentage_over_lowest_price = optional(number, 0)
          local_storage = optional(string, "")
          local_storage_types = optional(list(string), [])
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
          accelerator_names = optional(list(string), [])
          accelerator_types = optional(list(string), [])
          accelerator_total_memory_mib = optional(object({
            min = optional(number, 0)
            max = optional(number, 0)
          }))
        }))
      })), [])
      instances_distribution = optional(object({
        on_demand_allocation_strategy = optional(string, "")
        on_demand_base_capacity = optional(number, 0)
        on_demand_percentage_above_base_capacity = optional(number)
        spot_allocation_strategy = optional(string, "")
        spot_instance_pools = optional(number, 0)
        spot_max_price = optional(string, "")
      }))
    }))
    min_size = optional(number, 0)
    max_size = optional(number, 0)
    desired_capacity = optional(number, 0)
    desired_capacity_type = optional(string, "")
    capacity_rebalance = optional(bool, false)
    default_cooldown_seconds = optional(number, 0)
    default_instance_warmup_seconds = optional(number, 0)
    health_check_type = optional(string, "")
    health_check_grace_period_seconds = optional(number, 0)
    target_groups = optional(list(string), [])
    termination_policies = optional(list(string), [])
    max_instance_lifetime_seconds = optional(number, 0)
    protect_from_scale_in = optional(bool, false)
    placement_group = optional(string, "")
    service_linked_role_arn = optional(string, "")
    enabled_metrics = optional(list(string), [])
    suspended_processes = optional(list(string), [])
    instance_refresh = optional(object({
      strategy = string
      triggers = optional(list(string), [])
      preferences = optional(object({
        min_healthy_percentage = optional(number)
        max_healthy_percentage = optional(number, 0)
        instance_warmup_seconds = optional(number, 0)
        checkpoint_percentages = optional(list(number), [])
        checkpoint_delay_seconds = optional(number, 0)
        auto_rollback = optional(bool, false)
        alarms = optional(list(string), [])
        scale_in_protected_instances = optional(string, "")
        standby_instances = optional(string, "")
        skip_matching = optional(bool, false)
      }))
    }))
    warm_pool = optional(object({
      pool_state = optional(string, "")
      min_size = optional(number, 0)
      max_group_prepared_capacity = optional(number)
      reuse_on_scale_in = optional(bool, false)
    }))
    instance_maintenance_policy = optional(object({
      min_healthy_percentage = optional(number, 0)
      max_healthy_percentage = optional(number, 0)
    }))
    capacity_distribution_strategy = optional(string, "")
    force_delete = optional(bool, false)
    wait_for_capacity_timeout = optional(string, "")
    scaling_policies = optional(list(object({
      name = string
      policy_type = string
      estimated_instance_warmup_seconds = optional(number, 0)
      target_tracking = optional(object({
        target_value = number
        predefined_metric_type = optional(string, "")
        resource_label = optional(string, "")
        customized_metric = optional(object({
          metric_name = optional(string, "")
          namespace = optional(string, "")
          statistic = optional(string, "")
          unit = optional(string, "")
          dimensions = optional(list(object({
            name = string
            value = string
          })), [])
          period_seconds = optional(number, 0)
          metrics = optional(list(object({
            id = string
            expression = optional(string, "")
            metric_stat = optional(object({
              metric_name = string
              namespace = string
              stat = string
              unit = optional(string, "")
              dimensions = optional(list(object({
                name = string
                value = string
              })), [])
              period_seconds = optional(number, 0)
            }))
            label = optional(string, "")
            return_data = optional(bool)
          })), [])
        }))
        disable_scale_in = optional(bool, false)
      }))
      step_scaling = optional(object({
        adjustment_type = string
        metric_aggregation_type = optional(string, "")
        min_adjustment_magnitude = optional(number, 0)
        step_adjustments = list(object({
          scaling_adjustment = optional(number, 0)
          metric_interval_lower_bound = optional(string, "")
          metric_interval_upper_bound = optional(string, "")
        }))
      }))
      simple_scaling = optional(object({
        adjustment_type = string
        scaling_adjustment = optional(number, 0)
        cooldown_seconds = optional(number, 0)
        min_adjustment_magnitude = optional(number, 0)
      }))
      predictive_scaling = optional(object({
        target_value = number
        predefined_metric_pair_type = optional(string, "")
        resource_label = optional(string, "")
        mode = optional(string, "")
        scheduling_buffer_time_seconds = optional(number, 0)
        max_capacity_breach_behavior = optional(string, "")
        max_capacity_buffer = optional(number, 0)
        predefined_load_metric = optional(object({
          metric_type = string
          resource_label = optional(string, "")
        }))
        predefined_scaling_metric = optional(object({
          metric_type = string
          resource_label = optional(string, "")
        }))
        customized_load_metric_queries = optional(list(object({
          id = string
          expression = optional(string, "")
          metric_stat = optional(object({
            metric_name = string
            namespace = string
            stat = string
            unit = optional(string, "")
            dimensions = optional(list(object({
              name = string
              value = string
            })), [])
            period_seconds = optional(number, 0)
          }))
          label = optional(string, "")
          return_data = optional(bool)
        })), [])
        customized_scaling_metric_queries = optional(list(object({
          id = string
          expression = optional(string, "")
          metric_stat = optional(object({
            metric_name = string
            namespace = string
            stat = string
            unit = optional(string, "")
            dimensions = optional(list(object({
              name = string
              value = string
            })), [])
            period_seconds = optional(number, 0)
          }))
          label = optional(string, "")
          return_data = optional(bool)
        })), [])
        customized_capacity_metric_queries = optional(list(object({
          id = string
          expression = optional(string, "")
          metric_stat = optional(object({
            metric_name = string
            namespace = string
            stat = string
            unit = optional(string, "")
            dimensions = optional(list(object({
              name = string
              value = string
            })), [])
            period_seconds = optional(number, 0)
          }))
          label = optional(string, "")
          return_data = optional(bool)
        })), [])
      }))
      disabled = optional(bool, false)
    })), [])
    scheduled_actions = optional(list(object({
      name = string
      recurrence = optional(string, "")
      time_zone = optional(string, "")
      start_time = optional(string, "")
      end_time = optional(string, "")
      min_size = optional(number)
      max_size = optional(number)
      desired_capacity = optional(number)
    })), [])
    lifecycle_hooks = optional(list(object({
      name = string
      lifecycle_transition = string
      default_result = optional(string, "")
      heartbeat_timeout_seconds = optional(number, 0)
      notification_target_arn = optional(string, "")
      role_arn = optional(string, "")
      notification_metadata = optional(string, "")
      apply_at_launch = optional(bool, false)
    })), [])
    notifications = optional(object({
      topic = string
      event_types = list(string)
    }))
    capacity_reservation = optional(object({
      preference = optional(string, "")
      capacity_reservation_ids = optional(list(string), [])
      capacity_reservation_resource_group_arns = optional(list(string), [])
    }))
    traffic_sources = optional(list(object({
      identifier = string
      type = optional(string, "")
    })), [])
    instance_lifecycle_policy = optional(object({
      terminate_hook_abandon = optional(string, "")
    }))
    ignore_failed_scaling_activities = optional(bool, false)
    force_delete_warm_pool = optional(bool, false)
    min_elb_capacity = optional(number, 0)
    wait_for_elb_capacity = optional(number, 0)
  })
}
