output "resource_server_identifier" {
  description = "The resource server's identifier within its pool -- the scope prefix access tokens carry"
  value       = aws_cognito_resource_server.this.identifier
}

output "scope_identifiers" {
  description = "The fully-qualified scope identifiers ({identifier}/{scope_name}) app clients list in allowed_oauth_scopes"
  value       = local.scope_identifiers
}

output "user_pool_id" {
  description = "The user pool this resource server belongs to (resolved from the spec reference) -- AWS keys resource servers by the (pool id, identifier) pair"
  value       = aws_cognito_resource_server.this.user_pool_id
}
