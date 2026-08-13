# AzureMachineLearningOnlineEndpoint Pulumi Module

## Overview

Creates a managed online endpoint on an Azure Machine Learning workspace -- the stable HTTPS address applications call to score against deployed models, with authentication and a traffic dial across the endpoint's deployments.

## Resources Created

- `machinelearningservices.OnlineEndpoint` (pulumi-azure-native) -- the endpoint, an ARM child of the workspace

**Why azure-native, not the classic provider:** the classic pulumi-azure SDK (bridged from azurerm) carries NO ML endpoint resources, because azurerm itself has none (tracked at hashicorp/terraform-provider-azurerm#32011). The official azure-native SDK's typed resources are generated from the same ARM specification the Terraform module's raw-API shape pins, so both engines write the same ARM wire shapes. azure-native's `machinelearningservices` module targets ARM api-version 2025-09-01 -- a later GA line than the Terraform module's 2025-06-01 pin whose modeled surface is verified IDENTICAL property-by-property (definitions, enums, defaults, and mutability all match); the component's live verification reads ARM at 2025-06-01 either way, so one source of truth checks both engines.

## Key Implementation Details

- The shared `pulumiazurenativeprovider` builder resolves the credential mechanism (static client secret, keyless web identity, or ambient chain) -- this kind is a consumer of the azure-native provider path, unlike the classic-provider siblings.
- azure-native addresses ARM children by ancestor NAMES; the module parses the resource-group and workspace names out of the spec's workspace ARM ID (`parseWorkspaceId` in `module/locals.go`).
- Bring-your-own auth keys are wrapped with `pulumi.ToSecret` so they never land readable in state; ARM never returns them on reads, and they are deliberately NOT outputs (retrieval is the listKeys action).
- The identity type is sent as ARM's own common-types literal (`SystemAssigned,UserAssigned` -- no space), the same ARM value the Terraform module's azapi form normalizes to.

## Outputs

- `online_endpoint_id`, `online_endpoint_name`, `scoring_uri`, `swagger_uri`, `system_assigned_identity_principal_id` -- identical to the Terraform module's outputs.
