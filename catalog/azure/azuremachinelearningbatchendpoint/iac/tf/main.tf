# Create the Machine Learning batch endpoint -- the stable address
# batch scoring jobs are submitted to, as an ARM child of its workspace
# (.../workspaces/{ws}/batchEndpoints/{name}).
#
# Written at the pinned raw-ARM shape (no azurerm resource exists for ML
# endpoints); the body mirrors the ARM specification exactly, and the
# spec's validation rules are the only pre-apply safety net. The
# endpoint updates via full PUT: everything in the body updates in
# place (the default-deployment pointer is the routing dial); name,
# region and workspace replace the endpoint. Nothing runs or bills
# while no job is active.
resource "azapi_resource" "main" {
  type      = "Microsoft.MachineLearningServices/workspaces/batchEndpoints@2025-06-01"
  name      = var.spec.name
  parent_id = var.spec.workspace_id
  location  = var.spec.region

  # Optional here where the online endpoint requires one: batch jobs
  # run under the INVOKER's Entra token plus the COMPUTE pool's
  # managed identity, so the endpoint's own identity sits outside the
  # batch data path.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_wire[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # The ARM properties object, assembled in locals so unset optionals
  # are omitted and ARM's own defaults apply. There is NO keys arm
  # here (unlike the online endpoint): the batch service accepts only
  # AADToken auth -- it rejects Key mode outright -- so ARM's
  # create-time keys property is dead surface for this kind.
  body = {
    properties = local.endpoint_properties
  }

  tags = local.final_tags

  # The job-submission surface the outputs publish.
  response_export_values = [
    "properties.scoringUri",
    "properties.swaggerUri",
  ]

  # azapi validates the body against its embedded ARM schemas at plan
  # time -- the closest a raw-API resource gets to provider-side
  # validation; kept on deliberately.
  schema_validation_enabled = true
}
