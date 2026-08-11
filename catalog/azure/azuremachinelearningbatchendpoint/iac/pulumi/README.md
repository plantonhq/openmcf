# AzureMachineLearningBatchEndpoint Pulumi Module

## Overview

Creates a batch endpoint on an Azure Machine Learning workspace -- the stable address batch scoring jobs are submitted to, with Microsoft Entra authentication and a default-deployment pointer that routes submissions.

## Resources Created

- `machinelearningservices.BatchEndpoint` (pulumi-azure-native) -- the endpoint, an ARM child of the workspace

**Why azure-native, not the classic provider:** the classic pulumi-azure SDK (bridged from azurerm) carries NO ML endpoint resources, because azurerm itself has none (tracked at hashicorp/terraform-provider-azurerm#32011). The official azure-native SDK's typed resources are generated from the same ARM specification the Terraform module's raw-API shape pins, so both engines write the same ARM wire shapes. azure-native's `machinelearningservices` module targets ARM api-version 2025-09-01 -- a later GA line than the Terraform module's 2025-06-01 pin whose modeled surface is verified IDENTICAL property-by-property (definitions, enums, defaults, and mutability all match); the component's live verification reads ARM at 2025-06-01 either way, so one source of truth checks both engines.

## Key Implementation Details

- The shared `pulumiazurenativeprovider` builder resolves the credential mechanism (static client secret, keyless web identity, or ambient chain).
- azure-native addresses ARM children by ancestor NAMES; the module parses the resource-group and workspace names out of the spec's workspace ARM ID (`parseWorkspaceId` in `module/locals.go`).
- authMode always sends `AADToken` -- ARM requires the property, and it is the only value the batch service accepts (the spec's vocabulary enforces it; the module fills the default when unset). There is no keys arm: with Key auth rejected by the service, ARM's create-time keys property is dead surface for this kind.
- Identity is OPTIONAL (nil when the spec omits the block), unlike the online endpoint sibling -- batch jobs run under the invoker's token plus the compute pool's identity. When present, the type is sent as ARM's own common-types literal (`SystemAssigned,UserAssigned` -- no space).

## Outputs

- `batch_endpoint_id`, `batch_endpoint_name`, `scoring_uri`, `swagger_uri`, `system_assigned_identity_principal_id` -- identical to the Terraform module's outputs.
