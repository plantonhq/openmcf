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

output "grant_ids" {
  description = "The AWS-generated grant IDs, keyed by the grant's position in spec.grants (the module's for_each key) -- the handles RetireGrant/RevokeGrant and state import take."
  value       = { for k, grant in aws_kms_grant.this : k => grant.grant_id }
}
