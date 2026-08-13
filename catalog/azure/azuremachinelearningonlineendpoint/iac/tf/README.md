# AzureMachineLearningOnlineEndpoint Terraform Module

## Overview

Creates a managed online endpoint on an Azure Machine Learning workspace -- the stable HTTPS address applications call to score against deployed models, with authentication and a traffic dial across the endpoint's deployments.

## Resources Created

- `azapi_resource` -- the endpoint, written at the pinned raw-ARM shape `Microsoft.MachineLearningServices/workspaces/onlineEndpoints@2025-06-01`, an ARM child of the workspace

**Why azapi, not azurerm:** azurerm carries NO resource for ML online endpoints (its draft is tracked at hashicorp/terraform-provider-azurerm#32011). The azapi provider is pinned EXACT (2.11.0) and the api-version is pinned in the resource type -- never `latest`. When azurerm ships native resources, this module migrates azapi → native in the next minor release (state move / re-import). The kind's spec carries the full validation burden: azapi has no provider-side schema, so every ARM contract is a manifest-time rule or a documented apply-time boundary.

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningOnlineEndpointSpec fields; the workspace and identity references arrive as resolved literal ARM IDs

## Outputs

- `online_endpoint_id` -- the endpoint's full ARM ID (what deployments wire to)
- `online_endpoint_name` -- the traffic map's routing key space
- `scoring_uri` / `swagger_uri` -- the scoring address and its OpenAPI document
- `system_assigned_identity_principal_id` -- for registry / storage grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Update surface**: the body updates via full PUT -- traffic, mirror traffic, description, public network access, and auth mode all update in place; name, region, and workspace replace the endpoint. Endpoint names are reserved REGION-WIDE per subscription.
- **Bring-your-own keys ride `sensitive_body`**: azapi's write-only overlay merges the keys into the ARM request without ever storing them in Terraform state -- exactly right for a property ARM treats as create-time input and never returns on any read (retrieval is the separate listKeys action). Because ARM never echoes them, keys are also deliberately NOT outputs.
- **`ignore_missing_property` (azapi default true) is load-bearing here**: create-only properties like the initial keys never come back on reads, and the default keeps that from registering as drift.
- **Identity wire values**: azapi accepts the azurerm-style `"SystemAssigned, UserAssigned"` and normalizes it to ARM's own `SystemAssigned,UserAssigned` -- the same ARM identity type the Pulumi module sends directly.
