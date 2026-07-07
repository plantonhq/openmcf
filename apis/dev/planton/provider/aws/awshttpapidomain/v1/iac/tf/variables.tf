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
  description = "AwsHttpApiDomain specification"
  type = object({
    region = string
    domain_name = string
    certificate_arn = string
    ip_address_type = optional(string, "")
    mutual_tls = optional(object({
      truststore_uri = string
      truststore_version = optional(string, "")
    }))
    api_mappings = optional(list(object({
      api_id = string
      stage = string
      api_mapping_key = optional(string, "")
    })), [])
  })
}