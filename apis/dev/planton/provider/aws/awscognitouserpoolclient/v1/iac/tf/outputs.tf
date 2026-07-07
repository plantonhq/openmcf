output "client_id" {
  description = "The app client ID -- the public identifier applications present at sign-in and the 'aud' claim JWT authorizers validate"
  value       = aws_cognito_user_pool_client.this.id
}

output "client_secret" {
  description = "The app client secret (only minted when generate_secret is true). Sensitive -- treat as a credential"
  value       = aws_cognito_user_pool_client.this.client_secret
  sensitive   = true
}

output "user_pool_id" {
  description = "The user pool this client belongs to (resolved from the spec reference) -- application configs typically need the (pool id, client id) pair together"
  value       = aws_cognito_user_pool_client.this.user_pool_id
}
