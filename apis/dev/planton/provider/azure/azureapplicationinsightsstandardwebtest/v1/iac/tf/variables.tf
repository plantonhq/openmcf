variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Application Insights Standard Web Test specification"
  type = object({
    resource_group          = string
    name                    = string
    region                  = string
    application_insights_id = string

    frequency     = optional(number)
    timeout       = optional(number)
    enabled       = optional(bool)
    retry_enabled = optional(bool)

    request = object({
      url                              = string
      http_verb                        = optional(string, "")
      body                             = optional(string, "")
      follow_redirects_enabled         = optional(bool)
      parse_dependent_requests_enabled = optional(bool)
      headers = optional(list(object({
        name  = string
        value = string
      })), [])
    })

    validation_rules = optional(object({
      expected_status_code        = optional(number)
      ssl_cert_remaining_lifetime = optional(number)
      ssl_check_enabled           = optional(bool)
      content = optional(object({
        content_match      = string
        ignore_case        = optional(bool)
        pass_if_text_found = optional(bool)
      }))
    }))

    # Azure web-test location ids the test runs FROM (at least one).
    geo_locations = list(string)

    description = optional(string, "")

    tags = optional(map(string), {})
  })
}
