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
  description = "AwsEksAddon specification"
  type = object({
    region = string
    cluster_name = string
    addon_name = string
    addon_version = optional(string, "")
    resolve_conflicts_on_create = optional(string, "")
    resolve_conflicts_on_update = optional(string, "")
    configuration_values = optional(string, "")
    service_account_role_arn = optional(string, "")
    pod_identity_associations = optional(list(object({
      role_arn = string
      service_account = string
    })), [])
    preserve = optional(bool, false)
    namespace_config = optional(object({
      namespace = string
    }))
  })
}
