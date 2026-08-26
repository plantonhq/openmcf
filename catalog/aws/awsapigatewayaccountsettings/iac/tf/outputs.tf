# The resource's own `id` attribute is deprecated at the provider
# (the pin's plans warn on it); the caller-identity data source is the
# canonical account-id read and returns the same 12-digit value.
data "aws_caller_identity" "this" {}

output "account_id" {
  description = "The 12-digit AWS account ID the settings belong to (also the provider's import ID)"
  value       = data.aws_caller_identity.this.account_id
}

output "api_key_version" {
  description = "The API key version the account is on (AWS-managed)"
  value       = aws_api_gateway_account.this.api_key_version
}

output "features" {
  description = "Feature flags AWS reports enabled on the account"
  value       = aws_api_gateway_account.this.features
}

# throttle_settings is a computed one-element list on the provider;
# try() keeps the outputs clean if AWS ever reports it empty.
output "throttle_burst_limit" {
  description = "Account-level maximum burst of concurrent requests"
  value       = try(aws_api_gateway_account.this.throttle_settings[0].burst_limit, 0)
}

output "throttle_rate_limit" {
  description = "Account-level steady-state request rate ceiling (requests per second)"
  value       = try(aws_api_gateway_account.this.throttle_settings[0].rate_limit, 0)
}
