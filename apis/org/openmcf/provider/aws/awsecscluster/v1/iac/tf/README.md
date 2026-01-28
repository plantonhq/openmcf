# Terraform Module to Deploy AwsEcsCluster

This module provisions an AWS ECS (Elastic Container Service) cluster with support for Fargate and Fargate Spot capacity providers.
It includes optional CloudWatch Container Insights and ECS Exec capabilities for monitoring and debugging.

Generated `variables.tf` reflects the proto schema for `AwsEcsCluster`.

## Usage

Use the OpenMCF CLI (tofu) with the default local backend:

```shell
openmcf tofu init --manifest hack/manifest.yaml
openmcf tofu plan --manifest hack/manifest.yaml
openmcf tofu apply --manifest hack/manifest.yaml --auto-approve
openmcf tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see [`examples.md`](./examples.md) and [`hack/manifest.yaml`](../hack/manifest.yaml).
