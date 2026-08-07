# The schema group, addressed by the parent namespace's ARM ID. This
# resource has NO update surface -- Azure exposes no mutable properties
# on a schema group, so every field is ForceNew and any change replaces
# the group (dropping the schemas registered inside it). The registry's
# tier contract (STANDARD or higher namespace) is enforced by Azure at
# apply time.
resource "azurerm_eventhub_namespace_schema_group" "main" {
  # ForceNew: renaming replaces the group and drops its registered
  # schemas.
  name         = var.spec.schema_group_name
  namespace_id = var.spec.namespace_id

  # Evolution policy and format, mapped from the spec enums to ARM's
  # wire values in locals. Both ForceNew.
  schema_compatibility = local.schema_compatibility
  schema_type          = local.schema_type
}
