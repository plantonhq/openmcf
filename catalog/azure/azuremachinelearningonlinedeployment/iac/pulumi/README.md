# AzureMachineLearningOnlineDeployment Pulumi Module

## Overview

Creates a managed online deployment on an Azure Machine Learning online endpoint -- a running copy of a model behind the endpoint's address, on Azure-managed VMs with health probes, request limits, and optional model data collection.

## Resources Created

- `machinelearningservices.OnlineDeployment` (pulumi-azure-native) -- the deployment, an ARM child of the endpoint, with `ManagedOnlineDeploymentArgs` as its properties (the managed compute type only -- the Kubernetes and AzureMLCompute variants are a recorded deferral)

**Why azure-native, not the classic provider:** the classic pulumi-azure SDK (bridged from azurerm) carries NO ML deployment resources, because azurerm itself has none (tracked at hashicorp/terraform-provider-azurerm#32011). The official azure-native SDK's typed resources are generated from the same ARM specification the Terraform module's raw-API shape pins, so both engines write the same ARM wire shapes. azure-native's `machinelearningservices` module targets ARM api-version 2025-09-01 -- a later GA line than the Terraform module's 2025-06-01 pin whose modeled surface is verified IDENTICAL property-by-property; the component's live verification reads ARM at 2025-06-01 either way, so one source of truth checks both engines.

## Key Implementation Details

- The shared `pulumiazurenativeprovider` builder resolves the credential mechanism (static client secret, keyless web identity, or ambient chain).
- azure-native addresses ARM children by ancestor NAMES; the module parses the resource-group, workspace, and endpoint names out of the spec's endpoint ARM ID (`parseEndpointId` in `module/locals.go`).
- `instance_count` rides the ARM SKU (`Name: "Default"`, `Capacity`) -- the service's autoscaling contract for managed deployments, and the one dial the service applies without rolling containers.
- Scale settings are deliberately absent: the managed variant's only legal mode is Default; TargetUtilization is Kubernetes-only.

## Outputs

- `online_deployment_id`, `online_deployment_name` -- identical to the Terraform module's outputs.
