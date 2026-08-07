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
  description = "AwsEcrRepo specification"
  type = object({
    region = string
    repository_name = string
    image_tag_mutability = optional(string)
    image_tag_mutability_exclusion_filters = optional(list(string), [])
    encryption_type = optional(string)
    kms_key_id = optional(string, "")
    scan_on_push = optional(bool)
    force_delete = optional(bool, false)
    lifecycle_rules = optional(list(object({
      rule_priority = number
      description = optional(string, "")
      tag_status = string
      tag_prefixes = optional(list(string), [])
      tag_patterns = optional(list(string), [])
      count_type = string
      count_number = number
    })), [])
    repository_policy = optional(any)
  })
}
