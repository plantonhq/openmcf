locals {
  resource_name = coalesce(try(var.metadata.name, null), "aws-iam-role")
  tags          = merge({
    "Name" = local.resource_name
  }, try(var.metadata.labels, {}))

  description = try(var.spec.description, null)
  path        = try(var.spec.path, "/")

  trust_policy_json = try(var.spec.trust_policy, null)

  managed_policy_arns = try(var.spec.managed_policy_arns, [])

  inline_policies_map = try(var.spec.inline_policies, {})
}



