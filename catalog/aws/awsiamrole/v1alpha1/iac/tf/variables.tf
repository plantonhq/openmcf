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
  description = "AwsIamRole specification"
  type = object({
    region = string
    description = optional(string, "")
    path = optional(string, "")
    trust_policy = optional(any)
    oidc_trust = optional(object({
      provider_arn = string
      provider_url = string
      subjects = optional(list(string), [])
      wildcard_subjects = optional(list(string), [])
      audiences = optional(list(string), [])
    }))
    managed_policy_arns = optional(list(string), [])
    inline_policies = optional(any, {})
    max_session_duration = optional(number, 0)
    permissions_boundary = optional(string, "")
    force_detach_policies = optional(bool, false)
  })
}
