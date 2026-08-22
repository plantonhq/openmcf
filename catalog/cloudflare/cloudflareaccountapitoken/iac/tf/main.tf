# An account-owned API token. Each policy's `resources` travels to
# Cloudflare as ONE raw JSON object; the spec types it (whole-resource
# grant or nested sub-resource scoping per entry) and this module
# serializes each entry back to the API's shape with jsonencode.
#
# The token's secret value is returned by Cloudflare exactly once, on
# create -- see outputs. Cloudflare canonically re-orders policies and
# permission groups server-side; treat their order as insignificant.
resource "cloudflare_account_token" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name

  policies = [
    for policy in var.spec.policies : {
      effect = policy.effect
      permission_groups = [
        for id in policy.permission_group_ids : { id = id }
      ]
      # Whole-resource grants carry a string value, nested scopings a map --
      # two type-homogeneous comprehensions merged into one object, because a
      # single conditional cannot return both types.
      resources = jsonencode(merge(
        {
          for key, scope in policy.resources : key => scope.permission
          if try(scope.permission, "") != ""
        },
        {
          for key, scope in policy.resources : key => scope.subresources
          if try(scope.permission, "") == ""
        }
      ))
    }
  ]

  expires_on = try(var.spec.expires_on, "") != "" ? var.spec.expires_on : null
  not_before = try(var.spec.not_before, "") != "" ? var.spec.not_before : null
  status     = try(var.spec.status, "") != "" ? var.spec.status : null

  condition = try(var.spec.condition, null) != null ? {
    request_ip = try(var.spec.condition.request_ip, null) != null ? {
      in     = length(try(var.spec.condition.request_ip.in_cidrs, [])) > 0 ? var.spec.condition.request_ip.in_cidrs : null
      not_in = length(try(var.spec.condition.request_ip.not_in_cidrs, [])) > 0 ? var.spec.condition.request_ip.not_in_cidrs : null
    } : null
  } : null
}
