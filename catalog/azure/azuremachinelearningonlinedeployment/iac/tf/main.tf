# Create the managed online deployment -- a running copy of a model
# behind its endpoint's address, as an ARM child of the endpoint
# (.../onlineEndpoints/{endpoint}/deployments/{name}).
#
# Written at the pinned raw-ARM shape (no azurerm resource exists for ML
# deployments); the body mirrors the ARM specification exactly, and the
# spec's validation rules are the only pre-apply safety net. Updates go
# through full PUT (the service rolls the deployment's containers);
# instance_count rides the SKU capacity -- the one dial the service
# scales without touching containers. Name, region and endpoint replace
# the deployment.
resource "azapi_resource" "main" {
  type      = "Microsoft.MachineLearningServices/workspaces/onlineEndpoints/deployments@2025-06-01"
  name      = var.spec.name
  parent_id = var.spec.endpoint_id
  location  = var.spec.region

  body = {
    # The service's ARM contract for autoscaling: managed deployments
    # carry SKU name "Default" with capacity as the instance count.
    sku = {
      name     = "Default"
      capacity = var.spec.instance_count != null ? var.spec.instance_count : 1
    }
    properties = local.deployment_properties
  }

  tags = local.final_tags

  # azapi validates the body against its embedded ARM schemas at plan
  # time -- the closest a raw-API resource gets to provider-side
  # validation; kept on deliberately.
  schema_validation_enabled = true
}
