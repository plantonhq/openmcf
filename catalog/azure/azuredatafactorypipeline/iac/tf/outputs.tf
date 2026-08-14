output "pipeline_id" {
  description = "The Azure Resource Manager ID of the pipeline ({factory_id}/pipelines/{name})"
  value       = azurerm_data_factory_pipeline.main.id
}

output "pipeline_name" {
  description = "The pipeline's name -- what triggers and pipeline-run API calls reference"
  value       = azurerm_data_factory_pipeline.main.name
}
