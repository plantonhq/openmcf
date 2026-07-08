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
  description = "AwsAppRunnerService specification"
  type = object({
    region = string
    image_source = optional(object({
      image_identifier = string
      image_repository_type = string
      access_role_arn = optional(string, "")
    }))
    code_source = optional(object({
      repository_url = string
      branch = string
      source_directory = optional(string, "")
      connection_arn = string
      configuration_source = string
      runtime = optional(string, "")
      build_command = optional(string, "")
    }))
    port = optional(string)
    start_command = optional(string, "")
    environment_variables = optional(map(string), {})
    environment_secrets = optional(map(string), {})
    cpu = optional(string)
    memory = optional(string)
    instance_role_arn = optional(string, "")
    health_check = optional(object({
      protocol = optional(string)
      path = optional(string)
      interval = optional(number)
      timeout = optional(number)
      healthy_threshold = optional(number)
      unhealthy_threshold = optional(number)
    }))
    auto_scaling_configuration_arn = optional(string, "")
    vpc_connector_arn = optional(string, "")
    observability_configuration_arn = optional(string, "")
    is_publicly_accessible = optional(bool)
    ip_address_type = optional(string)
    kms_key_arn = optional(string, "")
    auto_deployments_enabled = optional(bool, false)
    custom_domains = optional(list(object({
      domain_name = string
      enable_www_subdomain = optional(bool)
    })), [])
    web_acl_arn = optional(string, "")
  })
}