# GcpCloudComposerUserWorkloadsConfigMap - Terraform Module

This Terraform module provisions a Kubernetes ConfigMap in a Cloud Composer environment's workloads namespace (`google_composer_user_workloads_config_map`). It is the Terraform-side implementation of the Planton `GcpCloudComposerUserWorkloadsConfigMap` resource kind and has feature parity with the Pulumi module.

## Overview

The ConfigMap's data updates in place; name, environment, region, and project are immutable (ForceNew). Data is plain configuration — use the user workloads Secret kind for secret material. No API enablement here: the Composer API is enabled by the environment this ConfigMap is delivered into. An empty `project_id` falls back to the provider's default project.

## Usage with Planton CLI

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
planton tofu destroy --manifest ../hack/manifest.yaml --auto-approve
```

Credentials are provided via stack input (by the CLI), not in the manifest `spec`. Manifest file: `../hack/manifest.yaml`.

## Direct Terraform Usage

```bash
cd apis/dev/planton/provider/gcp/gcpcloudcomposeruserworkloadsconfigmap/v1/iac/tf
terraform init
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `metadata` | Resource metadata (name, labels, etc.) | — |
| `spec` | GcpCloudComposerUserWorkloadsConfigMap spec | — |

The `spec` object includes: `project_id` (optional — provider default when empty), `region`, `environment` (resolved ref to the environment name), `config_map_name`, and `data` (min one entry).

## Outputs

| Name | Description |
|------|-------------|
| `name` | Fully qualified resource name (`projects/{p}/locations/{r}/environments/{e}/userWorkloadsConfigMaps/{n}`) |
| `config_map_name` | The Kubernetes ConfigMap name — what DAGs reference |
