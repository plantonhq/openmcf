# All three scope variants expose identical attribute faces; the locals
# coalesce per attribute across whichever ARM type materialized.

# AzureServiceBusDisasterRecoveryConfig's alias_authorization_rule_id
# consumes this with zero translation.
output "authorization_rule_id" {
  description = "The Azure Resource Manager ID of the authorization rule"
  value       = local.rule_id
}

output "rule_name" {
  description = "The rule's name (the SharedAccessKeyName clients present)"
  value       = local.rule_name
}

output "primary_key" {
  description = "The primary key"
  value       = local.primary_key
  sensitive   = true
}

output "secondary_key" {
  description = "The secondary key (rotation partner)"
  value       = local.secondary_key
  sensitive   = true
}

output "primary_connection_string" {
  description = "The ready-to-use primary connection string"
  value       = local.primary_connection_string
  sensitive   = true
}

output "secondary_connection_string" {
  description = "The secondary connection string (rotation partner)"
  value       = local.secondary_connection_string
  sensitive   = true
}

# Alias faces are only populated when the namespace carries a geo-DR
# pairing (AzureServiceBusDisasterRecoveryConfig); empty otherwise.
output "primary_connection_string_alias" {
  description = "The primary connection string addressing the geo-DR alias"
  value       = local.primary_connection_string_alias
  sensitive   = true
}

output "secondary_connection_string_alias" {
  description = "The secondary alias connection string (rotation partner)"
  value       = local.secondary_connection_string_alias
  sensitive   = true
}
