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
  description = "AwsCodePipeline specification"
  type = object({
    region = string
    pipeline_type = optional(string)
    execution_mode = optional(string)
    role_arn = string
    artifact_stores = list(object({
      location = string
      region = optional(string, "")
      encryption_key_id = optional(string, "")
    }))
    stages = list(object({
      name = string
      actions = list(object({
        name = string
        category = string
        owner = string
        provider = string
        version = string
        configuration = optional(map(string), {})
        input_artifacts = optional(list(string), [])
        output_artifacts = optional(list(string), [])
        namespace = optional(string, "")
        region = optional(string, "")
        role_arn = optional(string, "")
        run_order = optional(number, 0)
        timeout_in_minutes = optional(number, 0)
        commands = optional(list(string), [])
        output_artifacts_for_compute_action = optional(list(object({
          name = string
          files = optional(list(string), [])
        })), [])
        output_variables = optional(list(string), [])
      }))
      before_entry = optional(object({
        result = optional(string, "")
        rules = list(object({
          name = string
          rule_type_id = object({
            provider = string
            category = optional(string)
            owner = optional(string)
            version = optional(string, "")
          })
          configuration = optional(map(string), {})
          commands = optional(list(string), [])
          input_artifacts = optional(list(string), [])
          region = optional(string, "")
          role_arn = optional(string, "")
          timeout_in_minutes = optional(number, 0)
        }))
      }))
      on_success = optional(object({
        result = optional(string, "")
        rules = list(object({
          name = string
          rule_type_id = object({
            provider = string
            category = optional(string)
            owner = optional(string)
            version = optional(string, "")
          })
          configuration = optional(map(string), {})
          commands = optional(list(string), [])
          input_artifacts = optional(list(string), [])
          region = optional(string, "")
          role_arn = optional(string, "")
          timeout_in_minutes = optional(number, 0)
        }))
      }))
      on_failure = optional(object({
        result = optional(string, "")
        retry_configuration = optional(object({
          retry_mode = string
        }))
        condition = optional(object({
          result = optional(string, "")
          rules = list(object({
            name = string
            rule_type_id = object({
              provider = string
              category = optional(string)
              owner = optional(string)
              version = optional(string, "")
            })
            configuration = optional(map(string), {})
            commands = optional(list(string), [])
            input_artifacts = optional(list(string), [])
            region = optional(string, "")
            role_arn = optional(string, "")
            timeout_in_minutes = optional(number, 0)
          }))
        }))
      }))
    }))
    triggers = optional(list(object({
      provider_type = string
      git_configuration = object({
        source_action_name = string
        push = optional(list(object({
          branches = optional(object({
            includes = optional(list(string), [])
            excludes = optional(list(string), [])
          }))
          file_paths = optional(object({
            includes = optional(list(string), [])
            excludes = optional(list(string), [])
          }))
          tags = optional(object({
            includes = optional(list(string), [])
            excludes = optional(list(string), [])
          }))
        })), [])
        pull_request = optional(list(object({
          branches = optional(object({
            includes = optional(list(string), [])
            excludes = optional(list(string), [])
          }))
          file_paths = optional(object({
            includes = optional(list(string), [])
            excludes = optional(list(string), [])
          }))
          events = optional(list(string), [])
        })), [])
      })
    })), [])
    variables = optional(list(object({
      name = string
      default_value = optional(string, "")
      description = optional(string, "")
    })), [])
  })
}
