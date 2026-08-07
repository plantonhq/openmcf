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
  description = "AwsWafIpSet specification"
  type = object({
    region = string
    scope = string
    ip_address_version = string
    addresses = optional(list(string), [])
    description = optional(string, "")
  })
}
