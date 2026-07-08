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
  description = "AwsAppRunnerAutoScalingConfiguration specification"
  type = object({
    region = string
    max_concurrency = optional(number)
    max_size = optional(number)
    min_size = optional(number)
  })
}