# Create the standard web test. The request block is required; optional
# fields are sent only when set (null) so an unspecified spec and Azure's
# defaults (GET, follow-redirects, 5-minute frequency, 200 status) deploy
# identically on both engines. geo_locations are the Azure web-test location
# ids the test runs FROM.
resource "azurerm_application_insights_standard_web_test" "main" {
  name                    = var.spec.name
  resource_group_name     = var.spec.resource_group
  application_insights_id = var.spec.application_insights_id
  location                = var.spec.region
  geo_locations           = var.spec.geo_locations

  frequency     = var.spec.frequency
  timeout       = var.spec.timeout
  enabled       = var.spec.enabled
  retry_enabled = var.spec.retry_enabled
  description   = var.spec.description != "" ? var.spec.description : null

  request {
    url                              = var.spec.request.url
    http_verb                        = var.spec.request.http_verb != "" ? var.spec.request.http_verb : null
    body                             = var.spec.request.body != "" ? var.spec.request.body : null
    follow_redirects_enabled         = var.spec.request.follow_redirects_enabled
    parse_dependent_requests_enabled = var.spec.request.parse_dependent_requests_enabled

    dynamic "header" {
      for_each = var.spec.request.headers
      content {
        name  = header.value.name
        value = header.value.value
      }
    }
  }

  dynamic "validation_rules" {
    for_each = var.spec.validation_rules != null ? [var.spec.validation_rules] : []
    content {
      expected_status_code        = validation_rules.value.expected_status_code
      ssl_cert_remaining_lifetime = validation_rules.value.ssl_cert_remaining_lifetime
      ssl_check_enabled           = validation_rules.value.ssl_check_enabled

      dynamic "content" {
        for_each = validation_rules.value.content != null ? [validation_rules.value.content] : []
        content {
          content_match      = content.value.content_match
          ignore_case        = content.value.ignore_case
          pass_if_text_found = content.value.pass_if_text_found
        }
      }
    }
  }

  tags = local.final_tags
}
