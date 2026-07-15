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
  description = "AwsEksAccessEntry specification"
  type = object({
    region = string
    cluster_name = string
    principal_arn = string
    type = optional(string, "")
    kubernetes_groups = optional(list(string), [])
    user_name = optional(string, "")
    policy_associations = optional(list(object({
      policy_arn = string
      access_scope = object({
        type = string
        namespaces = optional(list(string), [])
      })
    })), [])
  })
}
