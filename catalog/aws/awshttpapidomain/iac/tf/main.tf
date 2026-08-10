# API Gateway v2 custom domain name + its API mappings.
#
# The domain binds an owned DNS name to an ACM certificate; mappings then
# publish APIs under the domain (optionally namespaced by a path key). DNS is
# composed downstream: the exported target_domain_name / hosted_zone_id feed a
# Route 53 alias record.

resource "aws_apigatewayv2_domain_name" "this" {
  domain_name = var.spec.domain_name

  # API Gateway v2 domains accept exactly one endpoint type (REGIONAL) and
  # one security policy (TLS_1_2) -- AWS validates against single-value enums
  # -- so neither is a spec field. Modeling them would be two decorative
  # knobs a user could never meaningfully turn.
  domain_name_configuration {
    certificate_arn = var.spec.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"

    # ipv4 or dualstack; AWS applies its default when unset.
    ip_address_type = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null

    # AWS-issued public certificate proving domain ownership -- required by
    # AWS when the TLS certificate is Private-CA-issued, or when mTLS is
    # combined with an ACM-imported certificate.
    ownership_verification_certificate_arn = var.spec.ownership_verification_certificate_arn != "" ? var.spec.ownership_verification_certificate_arn : null
  }

  # How requests route: static api_mappings only (AWS default), routing
  # rules only, or rules first with mapping fallback. Spec CEL guarantees
  # the mode and the routing_rules list agree.
  routing_mode = var.spec.routing_mode != "" ? var.spec.routing_mode : null

  # Mutual TLS: clients must present a certificate chaining to a CA in the
  # S3-hosted truststore. Pinning truststore_version makes CA rotation an
  # explicit change instead of a silent side effect of overwriting the
  # object.
  dynamic "mutual_tls_authentication" {
    for_each = var.spec.mutual_tls != null ? [var.spec.mutual_tls] : []
    content {
      truststore_uri     = mutual_tls_authentication.value.truststore_uri
      truststore_version = mutual_tls_authentication.value.truststore_version != "" ? mutual_tls_authentication.value.truststore_version : null
    }
  }

  tags = local.aws_tags
}

# One mapping resource per entry, addressed by path key (the domain root maps
# under the "(root)" alias -- see locals.tf). Referenced API IDs arrive
# pre-resolved as plain strings.
resource "aws_apigatewayv2_api_mapping" "this" {
  for_each = local.api_mapping_map

  api_id      = each.value.api_id
  domain_name = aws_apigatewayv2_domain_name.this.id
  stage       = each.value.stage

  # Empty means the domain root; AWS stores the absence of a key, not "".
  api_mapping_key = each.value.api_mapping_key != "" ? each.value.api_mapping_key : null
}

# One routing rule per entry, addressed by priority (unique per spec CEL,
# mirroring AWS's own uniqueness rule). Rules match on base path or header
# and invoke one API stage; the provider nests the target under
# action/invoke_api -- the spec's flat api_id/stage/strip_base_path fields
# render into that wrapper here.
resource "aws_apigatewayv2_routing_rule" "this" {
  for_each = local.routing_rule_map

  domain_name = aws_apigatewayv2_domain_name.this.id
  priority    = each.value.priority

  action {
    invoke_api {
      api_id          = each.value.api_id
      stage           = each.value.stage
      strip_base_path = each.value.strip_base_path
    }
  }

  # Each condition tests exactly one dimension (spec CEL mirrors the
  # provider's exactly-one-of); all conditions on a rule must match.
  dynamic "condition" {
    for_each = each.value.conditions
    content {
      dynamic "match_base_paths" {
        for_each = length(condition.value.base_paths) > 0 ? [condition.value.base_paths] : []
        content {
          any_of = match_base_paths.value
        }
      }

      dynamic "match_headers" {
        for_each = condition.value.header != null ? [condition.value.header] : []
        content {
          any_of {
            header     = match_headers.value.name
            value_glob = match_headers.value.value_glob
          }
        }
      }
    }
  }
}
