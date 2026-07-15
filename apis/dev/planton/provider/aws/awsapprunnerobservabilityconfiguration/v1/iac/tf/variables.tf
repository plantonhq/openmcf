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
  description = "AwsAppRunnerObservabilityConfiguration specification"
  type = object({
    region = string
    trace_configuration = optional(object({
      vendor = optional(string)
    }))
  })
}
