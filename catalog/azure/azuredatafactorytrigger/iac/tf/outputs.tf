# Exactly one of the four resources exists (the spec's variant block
# decides); all four share the {factory_id}/triggers/{name} ID shape.
output "trigger_id" {
  description = "The Azure Resource Manager ID of the trigger ({factory_id}/triggers/{name})"
  value = try(
    azurerm_data_factory_trigger_schedule.main[0].id,
    azurerm_data_factory_trigger_tumbling_window.main[0].id,
    azurerm_data_factory_trigger_blob_event.main[0].id,
    azurerm_data_factory_trigger_custom_event.main[0].id,
    null
  )
}

output "trigger_name" {
  description = "The trigger's name -- what tumbling window dependency entries resolve against"
  value = try(
    azurerm_data_factory_trigger_schedule.main[0].name,
    azurerm_data_factory_trigger_tumbling_window.main[0].name,
    azurerm_data_factory_trigger_blob_event.main[0].name,
    azurerm_data_factory_trigger_custom_event.main[0].name,
    null
  )
}
