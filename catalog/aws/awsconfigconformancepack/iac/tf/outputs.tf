output "pack_name" {
  description = "The pack's name (also the provider's import ID, at either scope)"
  value       = var.metadata.name
}

output "pack_arn" {
  description = "The pack's ARN (account- or organization-scope, whichever this instance deployed)"
  value = coalesce(
    one(aws_config_conformance_pack.this[*].arn),
    one(aws_config_organization_conformance_pack.this[*].arn),
  )
}

output "region" {
  description = "The AWS region the pack lives in (packs are addressed by region + name)"
  value = coalesce(
    one(aws_config_conformance_pack.this[*].region),
    one(aws_config_organization_conformance_pack.this[*].region),
  )
}
