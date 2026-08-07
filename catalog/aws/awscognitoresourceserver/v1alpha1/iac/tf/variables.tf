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
  description = "AwsCognitoResourceServer specification"
  type = object({
    region = string
    user_pool_id = string
    identifier = string
    name = string
    scopes = optional(list(object({
      scope_name = string
      scope_description = string
    })), [])
  })
}
