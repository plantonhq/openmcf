# Create the Machine Learning online endpoint -- the stable HTTPS
# address applications call to score against deployed models, as an ARM
# child of its workspace (.../workspaces/{ws}/onlineEndpoints/{name}).
#
# Written at the pinned raw-ARM shape (no azurerm resource exists for ML
# endpoints); the body mirrors the ARM specification exactly, and the
# spec's validation rules are the only pre-apply safety net. The
# endpoint updates via full PUT: everything in the body updates in
# place (traffic is the blue/green dial); name, region and workspace
# replace the endpoint. The endpoint's NAME is reserved region-wide per
# subscription.
resource "azapi_resource" "main" {
  type      = "Microsoft.MachineLearningServices/workspaces/onlineEndpoints@2025-06-01"
  name      = var.spec.name
  parent_id = var.spec.workspace_id
  location  = var.spec.region

  # Required by the spec (a recorded tightening of ARM's optional): an
  # endpoint without an identity cannot pull images or models and every
  # deployment on it fails at provisioning.
  identity {
    type         = local.identity_type_wire[var.spec.identity.type]
    identity_ids = length(var.spec.identity.identity_ids) > 0 ? var.spec.identity.identity_ids : null
  }

  # The ARM properties object, assembled in locals so unset optionals
  # are omitted and ARM's own defaults apply.
  body = {
    properties = local.endpoint_properties
  }

  # Bring-your-own auth keys (Key mode): write-only overlay -- merged
  # into the ARM request, never stored in state, never read back (ARM
  # returns keys only through the listKeys action).
  sensitive_body = length(local.initial_keys) > 0 ? {
    properties = {
      keys = local.initial_keys
    }
  } : null

  tags = local.final_tags

  # The scoring surface the outputs publish. ARM never returns key
  # values on reads, so keys are deliberately not exported.
  response_export_values = [
    "properties.scoringUri",
    "properties.swaggerUri",
  ]

  # azapi validates the body against its embedded ARM schemas at plan
  # time -- the closest a raw-API resource gets to provider-side
  # validation; kept on deliberately.
  schema_validation_enabled = true
}
