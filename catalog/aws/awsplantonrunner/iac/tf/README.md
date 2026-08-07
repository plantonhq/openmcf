# Terraform Module to Deploy AwsPlantonRunner

This module provisions a standing Planton runner appliance on ECS
Fargate: the credentials document in Secrets Manager (zero-day recovery
window, injected into the container by the ECS agent at task start), the
execution role (image pull, logs, and a read grant scoped to exactly the
one secret), the runtime IAM role (the runner's own AWS identity --
created permissionless, or the referenced `task_role` passed through),
an explicit-retention CloudWatch log group, an outbound-only security
group in the VPC derived from the first referenced subnet, a dedicated
ECS cluster, the task definition, and a Fargate service holding exactly
one runner.

The service deliberately does not gate on task steady state: the
runner's readiness contract is its work queue (operations wait there
until the worker polls), never ECS task liveness.

`variables.tf` reflects the proto schema for `AwsPlantonRunner`.

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest hack/manifest.yaml
planton tofu plan --manifest hack/manifest.yaml
planton tofu apply --manifest hack/manifest.yaml --auto-approve
planton tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

See [`hack/manifest.yaml`](../../e2e/manifest.yaml) for a minimal test manifest.
