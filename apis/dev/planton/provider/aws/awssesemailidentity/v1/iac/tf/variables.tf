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
  description = "AwsSesEmailIdentity specification"
  type = object({
    region = string
    email_identity = string
    configuration_set = optional(string, "")
    dkim_signing = optional(object({
      next_signing_key_length = optional(string, "")
      domain_signing_private_key = optional(string, "")
      domain_signing_selector = optional(string, "")
    }))
    mail_from = optional(object({
      mail_from_domain = string
      behavior_on_mx_failure = optional(string)
    }))
    email_forwarding_enabled = optional(bool)
    policies = optional(any, [])
  })
}
