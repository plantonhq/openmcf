# AzureMachineLearningComputeInstance Terraform Module

## Overview

Creates a compute instance on an Azure Machine Learning workspace -- a single always-on VM serving as one data scientist's cloud workstation for notebooks, interactive debugging, and small jobs.

## Resources Created

- `azurerm_machine_learning_compute_instance` -- the instance, an ARM child of the workspace (`.../workspaces/{ws}/computes/{name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningComputeInstanceSpec fields; the workspace and subnet references arrive as resolved literal ARM IDs

## Outputs

- `machine_learning_compute_instance_id` -- the instance's full ARM ID
- `machine_learning_compute_instance_name` -- what its owner selects as their compute
- `system_assigned_identity_principal_id` -- for storage / Key Vault grants, when a system identity is enabled
- `ssh_username`, `ssh_port` -- assigned by the service, populated only when the ssh block is configured

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **NO update path**: every argument is ForceNew on the provider, tags included -- any change replaces the instance, and its OS disk and local files go with it. Keep work in datastores and git.
- **Workspace-region only**: there is no location argument -- the instance always runs in its workspace's region (the service's own rule; clusters are the ML compute that can run elsewhere).
- **Region-wide name reservation**: instance names are unique per Azure region per subscription, not per workspace.
- **The subnet/public-IP contract lives at apply time**: with `node_public_ip_enabled` false, the provider requires `subnet_resource_id` UNLESS the workspace runs a managed network -- it inspects the live workspace's isolation mode, so the rule cannot be validated at manifest time (recorded on the spec fields).

## Required Permissions

The deploying principal needs `Microsoft.MachineLearningServices/workspaces/computes/*` on the workspace's resource group (Contributor covers it), plus regional vCPU quota for the chosen VM size.
