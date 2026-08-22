# AzureMachineLearningComputeCluster Terraform Module

## Overview

Creates an auto-scaling compute cluster on an Azure Machine Learning workspace -- the pool of VMs that training jobs and pipelines run on, growing and shrinking between its configured node bounds.

## Resources Created

- `azurerm_machine_learning_compute_cluster` -- the cluster, an ARM child of the workspace (`.../workspaces/{ws}/computes/{name}`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningComputeClusterSpec fields; the workspace and subnet references arrive as resolved literal ARM IDs

## Outputs

- `machine_learning_compute_cluster_id` -- the cluster's full ARM ID
- `machine_learning_compute_cluster_name` -- what jobs and pipelines reference as their compute target
- `system_assigned_identity_principal_id` -- for storage / Key Vault / ACR grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Update surface**: only `identity`, `scale_settings`, and `tags` update in place -- every other argument is ForceNew on the provider. Replacement is routine for clusters (they hold no data), but running jobs fail when it happens.
- **The region split**: `location` is the NODES' region and may differ from the workspace's (the only ML compute with this ability); the provider writes the cluster ENVELOPE at the workspace's region, so ARM reads the envelope back there.
- **Subnet is Optional+Computed**: unset lets Azure network the nodes -- a workspace managed network assigns one and the value is read back after apply.
- **SSH credentials**: the admin password is sensitive (reference a secret); `key_value` carries the SSH PUBLIC key -- public material. At least one is required when the block is present (spec CEL mirrors the provider's AtLeastOneOf).

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege action manifest the deploying principal needs.
