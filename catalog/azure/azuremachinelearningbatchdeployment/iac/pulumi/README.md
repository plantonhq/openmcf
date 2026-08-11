# AzureMachineLearningBatchDeployment Pulumi Module

## Overview

Creates a batch deployment on an Azure Machine Learning batch endpoint -- the job recipe (model, compute, batching behavior) behind the endpoint's address. Nothing runs or bills at create time; each endpoint invocation materializes a job from the recipe.

## Resources Created

- `machinelearningservices.BatchDeployment` (pulumi-azure-native) -- the deployment, an ARM child of the endpoint

**Why azure-native, not the classic provider:** the classic pulumi-azure SDK (bridged from azurerm) carries NO ML deployment resources, because azurerm itself has none (tracked at hashicorp/terraform-provider-azurerm#32011). The official azure-native SDK's typed resources are generated from the same ARM specification the Terraform module's raw-API shape pins, so both engines write the same ARM wire shapes. azure-native's `machinelearningservices` module targets ARM api-version 2025-09-01 -- a later GA line than the Terraform module's 2025-06-01 pin whose modeled surface is verified IDENTICAL property-by-property (definitions, enums, defaults, and mutability all match); the component's live verification reads ARM at 2025-06-01 either way, so one source of truth checks both engines.

## Key Implementation Details

- The shared `pulumiazurenativeprovider` builder resolves the credential mechanism (static client secret, keyless web identity, or ambient chain).
- azure-native addresses ARM children by ancestor NAMES; the module parses the resource-group, workspace, and endpoint names out of the spec's endpoint ARM ID (`parseEndpointId` in `module/locals.go`).
- The model reference union maps to the SDK's typed asset-reference args (`IdAssetReferenceArgs` / `DataPathAssetReferenceArgs` / `OutputPathAssetReferenceArgs`) with the `referenceType` discriminator derived from which spec block is set.
- The pipeline-component arm maps to `BatchPipelineComponentDeploymentConfigurationArgs`; its `settings` bag is string-valued -- the SDK's own generated typing, matching the spec's map shape, so both engines send identical wire values.
- One SDK nuance, wire-equivalent by construction: the typed args' own `Defaults()` fills `errorThreshold` with -1 when unset, which IS ARM's default -- the Terraform module omits the property and ARM applies the same value. `miniBatchSize` is float64 in the SDK where ARM types int64 -- the same JSON number on the wire.

## Outputs

- `batch_deployment_id`, `batch_deployment_name` -- identical to the Terraform module's outputs.
