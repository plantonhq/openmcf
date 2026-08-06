# Terraform Module to Deploy AwsEcsCluster

This module provisions an AWS ECS cluster and its capacity: the cluster
itself (Container Insights, ECS Exec auditing, Fargate storage
encryption, Service Connect defaults), one `aws_ecs_capacity_provider`
per folded EC2 entry (each wrapping a referenced auto-scaling group),
and a single `aws_ecs_cluster_capacity_providers` association that PUTs
the union of Fargate built-ins and EC2 provider names onto the cluster
together with the default strategy.

Generated `variables.tf` reflects the proto schema for `AwsEcsCluster`
(generator-owned; never hand-edit).

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest hack/manifest.yaml
planton tofu plan --manifest hack/manifest.yaml
planton tofu apply --manifest hack/manifest.yaml --auto-approve
planton tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see [`examples.md`](./examples.md) and [`hack/manifest.yaml`](../hack/manifest.yaml).
