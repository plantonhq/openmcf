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
  description = "AwsCertManagerCert specification"
  type = object({
    region = string
    primary_domain_name = optional(string, "")
    alternate_domain_names = optional(list(string), [])
    validation_method = optional(string, "")
    validation_options = optional(list(object({
      domain_name = string
      validation_domain = string
    })), [])
    key_algorithm = optional(string, "")
    route53_hosted_zone_id = optional(string, "")
    wait_for_validation = optional(bool)
    options = optional(object({
      certificate_transparency_logging_preference = optional(string, "")
      export = optional(string, "")
    }))
    imported = optional(object({
      certificate_body = string
      private_key = string
      certificate_chain = optional(string, "")
    }))
    certificate_authority_arn = optional(string, "")
    early_renewal_duration = optional(string, "")
  })
}
