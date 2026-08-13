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
  description = "AwsRoute53Zone specification"
  type = object({
    region = string
    comment = optional(string, "")
    is_private = optional(bool, false)
    vpc_associations = optional(list(object({
      vpc_id = string
      vpc_region = optional(string, "")
    })), [])
    delegation_set_id = optional(string, "")
    force_destroy = optional(bool, false)
    enable_accelerated_recovery = optional(bool)
    query_logging = optional(object({
      cloudwatch_log_group_arn = string
    }))
    dnssec = optional(object({
      kms_key_arn = string
      key_signing_key_name = optional(string, "")
      key_signing_key_status = optional(string, "")
    }))
  })
}
