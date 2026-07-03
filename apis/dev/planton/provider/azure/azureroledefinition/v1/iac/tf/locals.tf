# Role definitions carry no ARM tags (Microsoft.Authorization resources do
# not support them), so the usual metadata-derived tag locals are
# intentionally absent from this module.
locals {
  # Normalize empty strings to null: azurerm validates these fields as
  # non-empty WHEN PRESENT, so an empty string must mean "unset", exactly as
  # the Pulumi module omits empty optionals.
  description        = var.spec.description == null || var.spec.description == "" ? null : var.spec.description
  role_definition_id = var.spec.role_definition_id == null || var.spec.role_definition_id == "" ? null : var.spec.role_definition_id

  # An empty assignable_scopes list is passed as null rather than [] so
  # azurerm applies its own defaulting (assignable_scopes = [scope]) -- the
  # same behavior the Pulumi module gets by omitting the argument. Passing []
  # explicitly would still default server-side, but null keeps the plan
  # honest about "not configured".
  assignable_scopes = length(var.spec.assignable_scopes) == 0 ? null : var.spec.assignable_scopes
}
