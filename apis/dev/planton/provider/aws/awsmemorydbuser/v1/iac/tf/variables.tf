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
  description = "AwsMemorydbUser specification"
  type = object({
    region = string
    access_string = string
    authentication_mode = object({
      type = string
      passwords = optional(list(string), [])
    })
  })
}