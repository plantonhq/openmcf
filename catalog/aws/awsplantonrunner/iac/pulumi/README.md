# Pulumi Module to Deploy AwsPlantonRunner

This module provisions a standing Planton runner appliance on ECS
Fargate: the runner token in the `<name>-token` Secrets Manager secret
(zero-day recovery window -- the token is re-mintable credential
material -- injected into the container by the ECS agent at task
start), the execution role (image pull, logs, and a read grant scoped
to exactly the one secret), the runtime IAM role (the runner's own AWS
identity -- created permissionless, or the referenced `task_role`
passed through), an explicit-retention CloudWatch log group, an
outbound-only security group in the VPC derived from the first
referenced subnet, a dedicated ECS cluster, the task definition, and a
Fargate service holding exactly one runner. The appliance's name basis
is `metadata.name`, matching the Terraform module key-for-key.

The container's environment contract: `PLANTON_RUNNER_TOKEN` (resolved
from the secret by the ECS agent), `PLANTON_RUNNER_NAME` (the
registration name, `<env>-<metadata.name>` or `metadata.name` outside
an environment), `PLANTON_RUNNER_ENDPOINT` only when a control-plane
endpoint is declared (omitted, the runner's built-in hosted default
applies), `PORT` 50051, and `LOG_LEVEL` info. No execution-mode
variable exists: the runner self-configures its mode from the join
response.

The service deliberately does not gate on task steady state: the
runner's readiness contract is its work queue (operations wait there
until the worker polls), never ECS task liveness.

## CLI usage (Planton pulumi)

```bash
# Preview
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir .

# Update (apply)
planton pulumi update \
  --manifest ../../e2e/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir . \
  --yes

# Refresh
planton pulumi refresh \
  --manifest ../../e2e/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir .

# Destroy
planton pulumi destroy \
  --manifest ../../e2e/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir .
```

## Debugging

You can debug the Pulumi program with Delve by pointing
`runtime.options.binary` in `Pulumi.yaml` at a wrapper script:

```yaml
runtime:
  name: go
  options:
    binary: ./debug.sh
```

Then run your Pulumi commands as usual. For detailed steps, see
`docs/pages/docs/guide/debug-pulumi-modules.mdx`.
