# Amazon Bedrock AgentCore identity-and-access bundle: workload
# identities, vaulted outbound credentials (API keys, OAuth2 clients),
# and a Cedar policy engine with its policies.
#
# Lifecycle facts the renders below depend on:
#   - AWS vaults the credential secrets in Secrets Manager under the
#     service's token vault -- consumers reference the provider ARN,
#     never the secret;
#   - the write-only credential argument variants (api_key_wo,
#     client_id_wo, client_secret_wo) are excluded by design: the spec's
#     sensitive fields arrive just-in-time resolved, and the plain
#     arguments let the provider detect rotation;
#   - a Cedar policy is a structural child of its engine (created after,
#     destroyed before).

# Named identities AgentCore workloads present when calling other
# services.
resource "aws_bedrockagentcore_workload_identity" "this" {
  for_each = local.workload_identities

  name = each.value.name

  allowed_resource_oauth2_return_urls = length(each.value.allowed_resource_oauth2_return_urls) > 0 ? each.value.allowed_resource_oauth2_return_urls : null
}

# Vaulted API keys for outbound calls.
resource "aws_bedrockagentcore_api_key_credential_provider" "this" {
  for_each = local.api_key_providers

  name = each.value.name

  # The spec value is sensitive end to end (JIT-resolved secret
  # reference); the plain argument lets the provider detect rotation.
  api_key = each.value.api_key

  tags = local.aws_tags
}

# Vaulted OAuth2 clients. The spec's vendor field selects which of the
# provider's six structurally-identical vendor blocks renders (see
# locals.oauth2_vendor_values).
resource "aws_bedrockagentcore_oauth2_credential_provider" "this" {
  for_each = local.oauth2_providers

  name                       = each.value.name
  credential_provider_vendor = local.oauth2_vendor_values[each.value.vendor]

  oauth2_provider_config {
    dynamic "custom_oauth2_provider_config" {
      for_each = each.value.vendor == "CUSTOM" ? [each.value] : []
      content {
        client_id     = custom_oauth2_provider_config.value.client_id
        client_secret = custom_oauth2_provider_config.value.client_secret

        # Required for CUSTOM vendors (spec-validated): exactly one of a
        # discovery URL or spelled-out endpoints.
        oauth_discovery {
          discovery_url = custom_oauth2_provider_config.value.oauth_discovery.discovery_url != "" ? custom_oauth2_provider_config.value.oauth_discovery.discovery_url : null

          dynamic "authorization_server_metadata" {
            for_each = custom_oauth2_provider_config.value.oauth_discovery.authorization_server_metadata != null ? [custom_oauth2_provider_config.value.oauth_discovery.authorization_server_metadata] : []
            content {
              issuer                 = authorization_server_metadata.value.issuer
              authorization_endpoint = authorization_server_metadata.value.authorization_endpoint
              token_endpoint         = authorization_server_metadata.value.token_endpoint
              response_types         = length(authorization_server_metadata.value.response_types) > 0 ? authorization_server_metadata.value.response_types : null
            }
          }
        }
      }
    }

    dynamic "github_oauth2_provider_config" {
      for_each = each.value.vendor == "GITHUB" ? [each.value] : []
      content {
        client_id     = github_oauth2_provider_config.value.client_id
        client_secret = github_oauth2_provider_config.value.client_secret
      }
    }

    dynamic "google_oauth2_provider_config" {
      for_each = each.value.vendor == "GOOGLE" ? [each.value] : []
      content {
        client_id     = google_oauth2_provider_config.value.client_id
        client_secret = google_oauth2_provider_config.value.client_secret
      }
    }

    dynamic "microsoft_oauth2_provider_config" {
      for_each = each.value.vendor == "MICROSOFT" ? [each.value] : []
      content {
        client_id     = microsoft_oauth2_provider_config.value.client_id
        client_secret = microsoft_oauth2_provider_config.value.client_secret
      }
    }

    dynamic "salesforce_oauth2_provider_config" {
      for_each = each.value.vendor == "SALESFORCE" ? [each.value] : []
      content {
        client_id     = salesforce_oauth2_provider_config.value.client_id
        client_secret = salesforce_oauth2_provider_config.value.client_secret
      }
    }

    dynamic "slack_oauth2_provider_config" {
      for_each = each.value.vendor == "SLACK" ? [each.value] : []
      content {
        client_id     = slack_oauth2_provider_config.value.client_id
        client_secret = slack_oauth2_provider_config.value.client_secret
      }
    }
  }

  tags = local.aws_tags
}

# The Cedar authorization engine.
resource "aws_bedrockagentcore_policy_engine" "this" {
  count = local.has_policy_engine ? 1 : 0

  name = var.spec.policy_engine.engine_name

  description = var.spec.policy_engine.description != "" ? var.spec.policy_engine.description : null

  # Changing the key replaces the engine (provider-enforced).
  encryption_key_arn = var.spec.policy_engine.encryption_key_arn != "" ? var.spec.policy_engine.encryption_key_arn : null

  tags = local.aws_tags
}

# The engine's Cedar policies -- structural children of the engine.
resource "aws_bedrockagentcore_policy" "this" {
  for_each = local.policies

  policy_engine_id = aws_bedrockagentcore_policy_engine.this[0].policy_engine_id
  name             = each.value.name

  description     = each.value.description != "" ? each.value.description : null
  validation_mode = each.value.validation_mode != "" ? each.value.validation_mode : null

  definition {
    cedar {
      statement = each.value.cedar_statement
    }
  }
}
