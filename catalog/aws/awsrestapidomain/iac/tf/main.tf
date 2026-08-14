# An API Gateway custom domain: TLS termination on your certificate,
# base-path mappings fanning paths out across REST APIs, and -- for
# PRIVATE domains -- VPC-endpoint access associations.
#
# Lifecycle facts the renders below depend on:
#   - the certificate reference fans in by endpoint type: EDGE and
#     PRIVATE domains take certificate_arn (the AWS SDK documents that
#     argument as "edge-optimized endpoint or private endpoint"; EDGE
#     certs live in us-east-1), REGIONAL domains take
#     regional_certificate_arn -- one spec field, wired to the right
#     provider argument here (the same fan-in in both engines);
#   - uploaded certificate material follows the same fan-in
#     (certificate_name for EDGE/PRIVATE, regional_certificate_name
#     for REGIONAL);
#   - domain create/update waits on DomainNameStatus AVAILABLE (up to
#     60 minutes upstream -- enhanced security policies trigger a
#     post-create update);
#   - base-path mapping creation retries briefly upstream to absorb
#     stage-propagation lag on freshly deployed stages;
#   - an access association is a Plugin Framework resource with ARN
#     identity and NO update -- every change replaces it, which is free.

resource "aws_api_gateway_domain_name" "this" {
  domain_name = var.spec.domain_name

  # Certificate fan-in by endpoint type (see the header comment):
  # certificate_arn serves EDGE and PRIVATE, regional_certificate_arn
  # serves REGIONAL only.
  certificate_arn          = var.spec.certificate_arn != "" && local.endpoint_type != "REGIONAL" ? var.spec.certificate_arn : null
  regional_certificate_arn = var.spec.certificate_arn != "" && local.endpoint_type == "REGIONAL" ? var.spec.certificate_arn : null

  certificate_name          = var.spec.uploaded_certificate != null && local.endpoint_type != "REGIONAL" ? var.spec.uploaded_certificate.name : null
  regional_certificate_name = var.spec.uploaded_certificate != null && local.endpoint_type == "REGIONAL" ? var.spec.uploaded_certificate.name : null
  certificate_body          = var.spec.uploaded_certificate != null ? var.spec.uploaded_certificate.body : null
  certificate_chain         = var.spec.uploaded_certificate != null && var.spec.uploaded_certificate.chain != "" ? var.spec.uploaded_certificate.chain : null
  certificate_private_key   = var.spec.uploaded_certificate != null ? var.spec.uploaded_certificate.private_key : null

  endpoint_configuration {
    types           = [local.endpoint_type]
    ip_address_type = var.spec.endpoint_configuration != null && try(var.spec.endpoint_configuration.ip_address_type, "") != "" ? var.spec.endpoint_configuration.ip_address_type : null
  }

  endpoint_access_mode = var.spec.endpoint_access_mode != "" ? var.spec.endpoint_access_mode : null
  security_policy      = var.spec.security_policy != "" ? var.spec.security_policy : null

  dynamic "mutual_tls_authentication" {
    for_each = var.spec.mutual_tls != null ? [var.spec.mutual_tls] : []
    content {
      truststore_uri     = mutual_tls_authentication.value.truststore_uri
      truststore_version = mutual_tls_authentication.value.truststore_version != "" ? mutual_tls_authentication.value.truststore_version : null
    }
  }

  ownership_verification_certificate_arn = var.spec.ownership_verification_certificate_arn != "" ? var.spec.ownership_verification_certificate_arn : null

  routing_mode = var.spec.routing_mode != "" ? var.spec.routing_mode : null

  # The domain resource policy (PRIVATE domains). The spec carries it
  # as structured YAML; jsonencode renders the exact JSON AWS stores.
  policy = var.spec.policy != null ? jsonencode(var.spec.policy) : null

  tags = local.aws_tags
}

# Base-path mappings fanning the domain's paths out across REST APIs.
# PRIVATE domains additionally pass domain_name_id: AWS permits many
# private domains sharing one hostname account-wide, so the name alone
# is ambiguous there.
resource "aws_api_gateway_base_path_mapping" "this" {
  for_each = local.base_path_mappings

  domain_name    = aws_api_gateway_domain_name.this.domain_name
  domain_name_id = local.endpoint_type == "PRIVATE" ? aws_api_gateway_domain_name.this.domain_name_id : null
  api_id         = each.value.rest_api_id
  base_path      = each.value.base_path != "" ? each.value.base_path : null
  stage_name     = each.value.stage_name != "" ? each.value.stage_name : null
}

# VPC endpoints granted access to a PRIVATE domain.
resource "aws_api_gateway_domain_name_access_association" "this" {
  for_each = local.access_associations

  access_association_source      = each.value.vpc_endpoint_id
  access_association_source_type = "VPCE"
  domain_name_arn                = aws_api_gateway_domain_name.this.arn

  tags = local.aws_tags
}
