# Exactly one of the two resources exists (the spec's flowlet flag
# decides); both share the {factory_id}/dataflows/{name} ID shape.
output "data_flow_id" {
  description = "The Azure Resource Manager ID of the data flow ({factory_id}/dataflows/{name})"
  value = try(
    azurerm_data_factory_data_flow.main[0].id,
    azurerm_data_factory_flowlet_data_flow.main[0].id,
    null
  )
}

output "data_flow_name" {
  description = "The data flow's name -- what other data flows' flowlet references resolve against"
  value = try(
    azurerm_data_factory_data_flow.main[0].name,
    azurerm_data_factory_flowlet_data_flow.main[0].name,
    null
  )
}
