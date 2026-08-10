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
  description = "AwsLambda specification"
  type = object({
    region      = string
    description = optional(string, "")
    role_arn    = string
    s3 = optional(object({
      bucket         = string
      key            = string
      object_version = optional(string, "")
    }))
    image_uri                   = optional(string, "")
    source_code_hash            = optional(string, "")
    code_sha256                 = optional(string, "")
    source_kms_key_arn          = optional(string, "")
    runtime                     = optional(string, "")
    handler                     = optional(string, "")
    architecture                = optional(string, "")
    memory_size_mb              = optional(number, 0)
    timeout_seconds             = optional(number, 0)
    ephemeral_storage_mb        = optional(number, 0)
    environment                 = optional(map(string), {})
    kms_key_arn                 = optional(string, "")
    subnet_ids                  = optional(list(string), [])
    security_group_ids          = optional(list(string), [])
    ipv6_allowed_for_dual_stack = optional(bool, false)
    dead_letter_target_arn      = optional(string, "")
    tracing_mode                = optional(string, "")
    file_system_config = optional(object({
      access_point_arn = string
      local_mount_path = string
    }))
    image_config = optional(object({
      entry_point       = optional(list(string), [])
      command           = optional(list(string), [])
      working_directory = optional(string, "")
    }))
    layer_arns                     = optional(list(string), [])
    publish                        = optional(bool, false)
    publish_to                     = optional(string, "")
    reserved_concurrent_executions = optional(number)
    snap_start                     = optional(bool, false)
    managed_instances = optional(object({
      capacity_provider_arn           = string
      memory_gib_per_vcpu             = optional(number, 0)
      max_concurrency_per_environment = optional(number, 0)
    }))
    durable_config = optional(object({
      execution_timeout_seconds = number
      retention_period_days     = optional(number, 0)
    }))
    tenant_isolation_mode = optional(string, "")
    logging_config = optional(object({
      log_format            = optional(string, "")
      application_log_level = optional(string, "")
      system_log_level      = optional(string, "")
      log_group             = optional(string, "")
    }))
    code_signing_config_arn = optional(string, "")
    aliases = optional(list(object({
      name                               = string
      description                        = optional(string, "")
      function_version                   = string
      routing_additional_version_weights = optional(map(number), {})
      provisioned_concurrent_executions  = optional(number)
    })), [])
    function_url = optional(object({
      authorization_type = string
      invoke_mode        = optional(string, "")
      cors = optional(object({
        allow_credentials = optional(bool, false)
        allow_origins     = optional(list(string), [])
        allow_methods     = optional(list(string), [])
        allow_headers     = optional(list(string), [])
        expose_headers    = optional(list(string), [])
        max_age_seconds   = optional(number, 0)
      }))
      qualifier = optional(string, "")
    }))
    invoke_permissions = optional(list(object({
      statement_id             = string
      principal                = string
      action                   = optional(string, "")
      source_arn               = optional(string, "")
      source_account           = optional(string, "")
      principal_org_id         = optional(string, "")
      function_url_auth_type   = optional(string, "")
      qualifier                = optional(string, "")
      event_source_token       = optional(string, "")
      invoked_via_function_url = optional(bool, false)
    })), [])
    async_invoke_config = optional(object({
      maximum_retry_attempts     = optional(number)
      maximum_event_age_seconds  = optional(number, 0)
      on_success_destination_arn = optional(string, "")
      on_failure_destination_arn = optional(string, "")
      qualifier                  = optional(string, "")
    }))
    recursive_loop = optional(string, "")
    runtime_management = optional(object({
      update_runtime_on   = string
      runtime_version_arn = optional(string, "")
      qualifier           = optional(string, "")
    }))
    scaling_configs = optional(list(object({
      qualifier                  = string
      min_execution_environments = optional(number)
      max_execution_environments = optional(number)
    })), [])
  })
}
