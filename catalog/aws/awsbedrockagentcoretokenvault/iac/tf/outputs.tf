output "token_vault_id" {
  description = "The vault the setting was applied to (also the provider's import ID)"
  value       = aws_bedrockagentcore_token_vault_cmk.this.token_vault_id
}

output "key_type" {
  description = "The key ownership in effect: CustomerManagedKey or ServiceManagedKey"
  value       = try(aws_bedrockagentcore_token_vault_cmk.this.kms_configuration[0].key_type, "")
}

output "kms_key_arn" {
  description = "The customer-managed KMS key ARN in effect (empty under ServiceManagedKey)"
  value       = try(aws_bedrockagentcore_token_vault_cmk.this.kms_configuration[0].kms_key_arn, "")
}
