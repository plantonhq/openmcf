# Domain-scoped companions, each its own provider resource keyed off the
# domain: SAML authentication for Dashboards and the grantor side of
# cross-account VPC endpoint access. Both attach to exactly this domain
# and are meaningless without it, which is why they fold here instead of
# being their own kinds.

# SAML sign-in for OpenSearch Dashboards. The block's presence is the
# enable switch: removing spec.saml_options destroys this resource, which
# disables SAML on the domain (the provider's delete calls the disable
# API). Rides on fine-grained access control (CEL-enforced).
resource "aws_opensearch_domain_saml_options" "this" {
  count = var.spec.saml_options != null ? 1 : 0

  domain_name = aws_opensearch_domain.this.domain_name

  saml_options {
    enabled = true

    idp {
      entity_id        = var.spec.saml_options.idp_entity_id
      metadata_content = var.spec.saml_options.idp_metadata_content
    }

    master_backend_role = var.spec.saml_options.master_backend_role != "" ? var.spec.saml_options.master_backend_role : null
    master_user_name    = var.spec.saml_options.master_user_name != "" ? var.spec.saml_options.master_user_name : null
    roles_key           = var.spec.saml_options.roles_key != "" ? var.spec.saml_options.roles_key : null
    subject_key         = var.spec.saml_options.subject_key != "" ? var.spec.saml_options.subject_key : null

    # 0 keeps the AWS default (60 minutes).
    session_timeout_minutes = var.spec.saml_options.session_timeout_minutes > 0 ? var.spec.saml_options.session_timeout_minutes : null
  }
}

# Cross-account private access, grantor side: each listed account is
# authorized to create OpenSearch-managed VPC endpoints against this
# domain. One resource per account, keyed by the account ID, so grants
# come and go independently.
resource "aws_opensearch_authorize_vpc_endpoint_access" "this" {
  for_each = toset(var.spec.authorized_vpc_endpoint_access_accounts)

  domain_name = aws_opensearch_domain.this.domain_name
  account     = each.value
}
