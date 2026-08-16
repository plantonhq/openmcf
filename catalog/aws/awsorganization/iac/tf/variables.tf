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
  description = "AwsOrganization specification"
  type = object({
    region                        = string
    feature_set                   = optional(string, "")
    aws_service_access_principals = optional(list(string), [])
    enabled_policy_types          = optional(list(string), [])
    delegated_administrators = optional(list(object({
      account_id        = optional(string, "")
      service_principal = string
    })), [])
    resource_policy = optional(any)
  })
}