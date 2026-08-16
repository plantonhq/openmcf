# An EventBridge API destination and/or its connection, as two
# independently deployable arms (the spec's CELs guarantee at least
# one arm, and exactly one connection source for the destination).
#
# Lifecycle facts the render below depends on:
#   - both names are fixed for life (replace-on-change);
#   - the connection's authorization_type is DERIVED from whichever
#     auth block the spec sets - the two can never disagree;
#   - AWS stores the credential values in a Secrets Manager secret it
#     creates and owns (the secret_arn output); DescribeConnection
#     never returns them, so the provider reads secrets back from
#     prior state - imports cannot recover them (declared
#     write-normalized in the import catalog);
#   - connection creates/updates wait through an auth state machine
#     (CREATING/AUTHORIZING -> AUTHORIZED, up to 20 minutes);
#   - the destination's rate limit defaults to 300/s at AWS.

# --- connection arm --------------------------------------------------------

resource "aws_cloudwatch_event_connection" "this" {
  count = var.spec.connection != null ? 1 : 0

  name        = var.spec.connection.name
  description = var.spec.connection.description != "" ? var.spec.connection.description : null

  # Derived: exactly one auth block is set (spec CEL).
  authorization_type = var.spec.connection.api_key != null ? "API_KEY" : (var.spec.connection.basic != null ? "BASIC" : "OAUTH_CLIENT_CREDENTIALS")

  auth_parameters {
    dynamic "api_key" {
      for_each = var.spec.connection.api_key != null ? [var.spec.connection.api_key] : []
      content {
        key   = api_key.value.key
        value = api_key.value.value
      }
    }

    dynamic "basic" {
      for_each = var.spec.connection.basic != null ? [var.spec.connection.basic] : []
      content {
        username = basic.value.username
        password = basic.value.password
      }
    }

    dynamic "oauth" {
      for_each = var.spec.connection.oauth != null ? [var.spec.connection.oauth] : []
      content {
        authorization_endpoint = oauth.value.authorization_endpoint
        http_method            = oauth.value.http_method

        client_parameters {
          client_id     = oauth.value.client_id
          client_secret = oauth.value.client_secret
        }

        oauth_http_parameters {
          dynamic "body" {
            for_each = oauth.value.oauth_http_parameters.body != null ? oauth.value.oauth_http_parameters.body : []
            content {
              key             = body.value.key
              value           = body.value.value
              is_value_secret = body.value.is_value_secret
            }
          }
          dynamic "header" {
            for_each = oauth.value.oauth_http_parameters.header != null ? oauth.value.oauth_http_parameters.header : []
            content {
              key             = header.value.key
              value           = header.value.value
              is_value_secret = header.value.is_value_secret
            }
          }
          dynamic "query_string" {
            for_each = oauth.value.oauth_http_parameters.query_string != null ? oauth.value.oauth_http_parameters.query_string : []
            content {
              key             = query_string.value.key
              value           = query_string.value.value
              is_value_secret = query_string.value.is_value_secret
            }
          }
        }
      }
    }

    dynamic "invocation_http_parameters" {
      for_each = var.spec.connection.invocation_http_parameters != null ? [var.spec.connection.invocation_http_parameters] : []
      content {
        dynamic "body" {
          for_each = invocation_http_parameters.value.body != null ? invocation_http_parameters.value.body : []
          content {
            key             = body.value.key
            value           = body.value.value
            is_value_secret = body.value.is_value_secret
          }
        }
        dynamic "header" {
          for_each = invocation_http_parameters.value.header != null ? invocation_http_parameters.value.header : []
          content {
            key             = header.value.key
            value           = header.value.value
            is_value_secret = header.value.is_value_secret
          }
        }
        dynamic "query_string" {
          for_each = invocation_http_parameters.value.query_string != null ? invocation_http_parameters.value.query_string : []
          content {
            key             = query_string.value.key
            value           = query_string.value.value
            is_value_secret = query_string.value.is_value_secret
          }
        }
      }
    }

    # A private OAuth authorization endpoint, reached through VPC
    # Lattice.
    dynamic "connectivity_parameters" {
      for_each = var.spec.connection.private_authorization_endpoint != null ? [var.spec.connection.private_authorization_endpoint] : []
      content {
        resource_parameters {
          resource_configuration_arn = connectivity_parameters.value.resource_configuration_arn
        }
      }
    }
  }

  # A private invocation endpoint, reached through VPC Lattice.
  dynamic "invocation_connectivity_parameters" {
    for_each = var.spec.connection.private_invocation_endpoint != null ? [var.spec.connection.private_invocation_endpoint] : []
    content {
      resource_parameters {
        resource_configuration_arn = invocation_connectivity_parameters.value.resource_configuration_arn
      }
    }
  }

  kms_key_identifier = var.spec.connection.kms_key_identifier != "" ? var.spec.connection.kms_key_identifier : null
}

# --- destination arm -------------------------------------------------------

resource "aws_cloudwatch_event_api_destination" "this" {
  count = var.spec.destination != null ? 1 : 0

  name        = var.spec.destination.name
  description = var.spec.destination.description != "" ? var.spec.destination.description : null

  # The owned connection when the instance has one, else the external
  # ARN (spec CEL: exactly one source).
  connection_arn = var.spec.connection != null ? aws_cloudwatch_event_connection.this[0].arn : var.spec.destination.connection_arn

  invocation_endpoint              = var.spec.destination.invocation_endpoint
  http_method                      = var.spec.destination.http_method
  invocation_rate_limit_per_second = var.spec.destination.invocation_rate_limit_per_second != null ? var.spec.destination.invocation_rate_limit_per_second : null
}
