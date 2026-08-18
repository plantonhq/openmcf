# An IAM SAML identity provider: the account's federation trust
# anchor, created from the IdP's metadata XML.
#
# Lifecycle facts the render below depends on:
#   - the provider's name comes from metadata.name and is WRITE-ONCE at
#     AWS - a rename replaces the provider and invalidates every role
#     trust policy naming its ARN;
#   - the metadata document updates IN PLACE - certificate rotations
#     are ordinary updates, and valid_until (exported below) is the
#     date to rotate by;
#   - IAM is global: the provider exists account-wide regardless of the
#     endpoint region the stack ran against.
resource "aws_iam_saml_provider" "this" {
  name                   = var.metadata.name
  saml_metadata_document = var.spec.saml_metadata_document

  tags = local.aws_tags
}
