# Create the batch deployment -- the job recipe (model, compute,
# batching behavior) behind its endpoint's address, as an ARM child of
# the endpoint (.../batchEndpoints/{endpoint}/deployments/{name}).
#
# Written at the pinned raw-ARM shape (no azurerm resource exists for ML
# deployments); the body mirrors the ARM specification exactly, and the
# spec's validation rules are the only pre-apply safety net. Everything
# in the recipe updates in place via full PUT (ARM flags nothing
# immutable on this surface -- unlike the online deployment's ForceNew
# instance type); name, region and endpoint replace the deployment.
# Nothing runs or bills at create time -- each endpoint invocation
# materializes a job from this recipe. The envelope's sku is
# deliberately absent: batch scale lives on resources.instanceCount
# per job, not on an autoscaling SKU (the online deployment's dial).
resource "azapi_resource" "main" {
  type      = "Microsoft.MachineLearningServices/workspaces/batchEndpoints/deployments@2025-06-01"
  name      = var.spec.name
  parent_id = var.spec.endpoint_id
  location  = var.spec.region

  body = {
    properties = local.deployment_properties
  }

  tags = local.final_tags

  # azapi validates the body against its embedded ARM schemas at plan
  # time -- the closest a raw-API resource gets to provider-side
  # validation; kept on deliberately.
  schema_validation_enabled = true
}
