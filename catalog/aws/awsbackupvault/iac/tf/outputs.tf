output "vault_arn" {
  description = "The vault's ARN (either arm) - what copy actions and air-gapped targeting reference"
  value       = var.spec.standard != null ? aws_backup_vault.this[0].arn : aws_backup_logically_air_gapped_vault.this[0].arn
}

output "vault_name" {
  description = "The vault's name (also the provider's import ID on either arm) - what backup plan rules target"
  value       = var.spec.standard != null ? aws_backup_vault.this[0].name : aws_backup_logically_air_gapped_vault.this[0].name
}
