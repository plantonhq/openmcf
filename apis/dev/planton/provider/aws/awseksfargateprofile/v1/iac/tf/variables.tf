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
  description = "AwsEksFargateProfile specification"
  type = object({
    region = string
    cluster_name = string
    pod_execution_role_arn = string
    subnet_ids = list(string)
    selectors = list(object({
      namespace = string
      labels = optional(map(string), {})
    }))
  })
}
