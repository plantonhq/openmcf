# Pulumi Module to Deploy AwsEcsCluster

This module provisions an AWS ECS cluster and its capacity: the cluster
itself (Container Insights, ECS Exec auditing, Fargate storage
encryption, Service Connect defaults), one capacity provider per folded
EC2 entry (each wrapping a referenced auto-scaling group with
ECS-managed scaling and draining), and a single association that PUTs
the union of Fargate built-ins and EC2 provider names onto the cluster
together with the default strategy. The cluster name is `metadata.name`,
matching the Terraform module key-for-key.

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

You can debug the Pulumi program with Delve. A `debug.sh` helper is
provided. To enable it, uncomment the `runtime.options.binary` line in
`Pulumi.yaml`:

```yaml
runtime:
  name: go
  options:
    binary: ./debug.sh
```

Then run your Pulumi commands as usual. For detailed steps, see
`docs/pages/docs/guide/debug-pulumi-modules.mdx`.
