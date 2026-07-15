# Cognito resource server -- the OAuth 2.0 resource (an API) a user pool
# mints custom access-token scopes for. The referenced pool ID arrives
# pre-resolved as a plain string.

resource "aws_cognito_resource_server" "this" {
  # The identifier is the resource server's identity within the pool
  # (ForceNew) and the prefix of every scope it mints.
  identifier   = var.spec.identifier
  name         = var.spec.name
  user_pool_id = var.spec.user_pool_id

  # Scopes update in place: removing one invalidates it for future tokens;
  # already-issued tokens carry it until they expire.
  dynamic "scope" {
    for_each = var.spec.scopes
    content {
      scope_name        = scope.value.scope_name
      scope_description = scope.value.scope_description
    }
  }
}
