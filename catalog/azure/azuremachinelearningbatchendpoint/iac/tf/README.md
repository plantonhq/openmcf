# AzureMachineLearningBatchEndpoint Terraform Module

## Overview

Creates a batch endpoint on an Azure Machine Learning workspace -- the stable address batch scoring jobs are submitted to, with Microsoft Entra authentication and a default-deployment pointer that routes submissions.

## Resources Created

- `azapi_resource` -- the endpoint, written at the pinned raw-ARM shape `Microsoft.MachineLearningServices/workspaces/batchEndpoints@2025-06-01`, an ARM child of the workspace

**Why azapi, not azurerm:** azurerm carries NO resource for ML batch endpoints (its endpoint draft is tracked at hashicorp/terraform-provider-azurerm#32011). The azapi provider is pinned EXACT (2.11.0) and the api-version is pinned in the resource type -- never `latest`. When azurerm ships native resources, this module migrates azapi → native in the next minor release (state move / re-import). The kind's spec carries the full validation burden: azapi has no provider-side schema, so every ARM contract is a manifest-time rule or a documented apply-time boundary.

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningBatchEndpointSpec fields; the workspace and identity references arrive as resolved literal ARM IDs

## Outputs

- `batch_endpoint_id` -- the endpoint's full ARM ID (what deployments wire to)
- `batch_endpoint_name` -- the default-deployment pointer's routing key space
- `scoring_uri` / `swagger_uri` -- the job-submission address and its OpenAPI document
- `system_assigned_identity_principal_id` -- for role grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Update surface**: the body updates via full PUT -- the default-deployment pointer, description, and properties all update in place; name, region, and workspace replace the endpoint.
- **auth mode always sends `AADToken`**: ARM requires authMode, and AADToken is the only value the batch service accepts -- it rejects `Key` and `AMLToken` with "AuthMode must be 'AADToken'" even though the shared ARM enum advertises them. The spec's vocabulary enforces this at manifest time; the module fills the default when the field is unset.
- **No keys arm, no `sensitive_body`**: with Key auth impossible, ARM's create-time `keys` property is dead surface for this kind (deliberately not modeled) -- the one structural difference from the online endpoint sibling's module.
- **No `publicNetworkAccess`**: the batch surface does not carry the property (reachability follows the workspace's network settings) -- another deliberate difference from the online sibling.
- **Identity is OPTIONAL** (a `dynamic` block): batch jobs run under the invoker's Entra token plus the compute pool's managed identity, so the endpoint's identity sits outside the batch data path. Wire values follow the family pattern: azapi accepts the azurerm-style `"SystemAssigned, UserAssigned"` and normalizes it to ARM's own `SystemAssigned,UserAssigned`.
