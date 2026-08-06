output "key_id" {
  description = "The generated key ID (UUID; mrk-... for multi-Region keys)."
  value       = aws_kms_key.this.key_id
}

output "key_arn" {
  description = "The key ARN -- the join key encryption-at-rest fields across the catalog reference."
  value       = aws_kms_key.this.arn
}

output "alias_names" {
  description = "The alias names attached to the key (each alias/...), in spec order."
  value       = var.spec.aliases
}
