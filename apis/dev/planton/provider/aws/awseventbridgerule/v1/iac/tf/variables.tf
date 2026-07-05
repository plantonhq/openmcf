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
  description = "AwsEventBridgeRule specification"
  type = object({
    region = string
    event_bus_name = optional(string, "")
    description = optional(string, "")
    event_pattern = optional(any)
    schedule_expression = optional(string, "")
    state = optional(string, "")
    role_arn = optional(string, "")
    force_destroy = optional(bool, false)
    targets = optional(list(object({
      name = string
      arn = string
      role_arn = optional(string, "")
      input = optional(string, "")
      input_path = optional(string, "")
      input_transformer = optional(object({
        input_paths = optional(map(string), {})
        input_template = string
      }))
      dead_letter_config = optional(object({
        arn = string
      }))
      retry_policy = optional(object({
        maximum_event_age_in_seconds = optional(number)
        maximum_retry_attempts = optional(number)
      }))
      sqs_target = optional(object({
        message_group_id = optional(string, "")
      }))
      kinesis_target = optional(object({
        partition_key_path = optional(string, "")
      }))
      http_target = optional(object({
        path_parameter_values = optional(list(string), [])
        query_string_parameters = optional(map(string), {})
        header_parameters = optional(map(string), {})
      }))
      batch_target = optional(object({
        job_definition = string
        job_name = string
        array_size = optional(number, 0)
        job_attempts = optional(number, 0)
      }))
      ecs_target = optional(object({
        task_definition_arn = string
        task_count = optional(number, 0)
        launch_type = optional(string, "")
        platform_version = optional(string, "")
        group = optional(string, "")
        capacity_provider_strategy = optional(list(object({
          capacity_provider = string
          base = optional(number, 0)
          weight = optional(number, 0)
        })), [])
        network_configuration = optional(object({
          subnets = list(string)
          security_groups = optional(list(string), [])
          assign_public_ip = optional(bool, false)
        }))
        ordered_placement_strategy = optional(list(object({
          type = string
          field = optional(string, "")
        })), [])
        placement_constraints = optional(list(object({
          type = string
          expression = optional(string, "")
        })), [])
        propagate_tags = optional(string, "")
        enable_ecs_managed_tags = optional(bool, false)
        enable_execute_command = optional(bool, false)
      }))
    })), [])
  })
}
