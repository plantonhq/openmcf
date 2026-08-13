# AzureMachineLearningOnlineDeployment Terraform Module

## Overview

Creates a managed online deployment on an Azure Machine Learning online endpoint -- a running copy of a model behind the endpoint's address, on Azure-managed VMs with health probes, request limits, and optional model data collection.

## Resources Created

- `azapi_resource` -- the deployment, written at the pinned raw-ARM shape `Microsoft.MachineLearningServices/workspaces/onlineEndpoints/deployments@2025-06-01`, an ARM child of the endpoint

**Why azapi, not azurerm:** azurerm carries NO resource for ML online deployments (its endpoint draft is tracked at hashicorp/terraform-provider-azurerm#32011). The azapi provider is pinned EXACT (2.11.0) and the api-version is pinned in the resource type -- never `latest`. When azurerm ships native resources, this module migrates azapi → native in the next minor release (state move / re-import). The kind's spec carries the full validation burden: azapi has no provider-side schema, so every ARM contract is a manifest-time rule or a documented apply-time boundary.

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningOnlineDeploymentSpec fields; the endpoint reference arrives as a resolved literal ARM ID

## Outputs

- `online_deployment_id` -- the deployment's full ARM ID
- `online_deployment_name` -- the key the endpoint's traffic map routes by

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **`endpointComputeType` is PINNED to `Managed`** -- this kind models the managed compute type only; the Kubernetes and AzureMLCompute variants are a recorded deferral (they require attached compute whose supported story does not exist yet).
- **Scale settings are deliberately absent from the body** -- the managed variant's only legal mode is Default (fixed instance count); TargetUtilization is Kubernetes-only. `instance_count` rides the ARM SKU (`name: "Default"`, `capacity`) -- the service's autoscaling contract, and the one dial that updates without rolling containers.
- **Updates go through full PUT** -- the service rolls the deployment's instances on model/environment/code changes; ship model changes as NEW deployments and shift the endpoint's traffic map instead.
- **Instances bill while the deployment lives** -- there is no scale-to-zero; managed-endpoint VM quota (separate from regular compute quota) gates provisioning.
