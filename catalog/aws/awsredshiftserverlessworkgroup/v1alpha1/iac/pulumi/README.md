# Pulumi Module to Deploy AwsRedshiftServerlessWorkgroup

This Pulumi Go program deploys an Amazon Redshift Serverless workgroup -- the
compute plane of the serverless warehouse -- using the Planton API and module.
The namespace it serves, the subnets it places compute in, and the security
groups on its endpoint all attach by reference. The module never creates a
resource that deserves to be its own node.

## Requirements

- Planton CLI built locally
- Valid AWS credential provided via the CLI stack input (not in `spec`)

## CLI commands

Preview:

```shell
planton pulumi preview \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```

Update (apply):

```shell
planton pulumi update \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

Destroy:

```shell
planton pulumi destroy \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

## Resources Created

1. **`redshiftserverless.Workgroup`** — the compute plane: RPU capacity (fixed
   baseline or price-performance target), VPC placement, reachability, config
   parameters, release track

## Outputs

| Key | Description |
|-----|-------------|
| `workgroup_name` | The handle the serverless APIs and credentials API address |
| `workgroup_id` | Unique identifier AWS assigns to the workgroup |
| `arn` | ARN of the workgroup |
| `endpoint_address` | DNS hostname SQL clients connect to |
| `port` | TCP port for connections |
