# Pulumi Module to Deploy AwsEc2Instance

This module provisions exactly one AWS EC2 virtual machine: launch source
(inline AMI/type or a referenced launch template with per-instance
overrides), networking, storage, IMDSv2 posture, purchase options, and
lifecycle protections. Everything the instance composes with -- subnet,
security groups, IAM instance profile (by NAME), launch template, KMS
keys -- attaches by reference; the module creates no side resources. The
instance's display identity is the `Name` tag carried from
`metadata.name`, matching the Terraform module key-for-key.

## CLI usage (Planton pulumi)

```bash
# Preview
planton pulumi preview \
  --manifest ../hack/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir .

# Update (apply)
planton pulumi update \
  --manifest ../hack/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir . \
  --yes

# Refresh
planton pulumi refresh \
  --manifest ../hack/manifest.yaml \
  --stack <org>/<project>/<stack> \
  --module-dir .

# Destroy
planton pulumi destroy \
  --manifest ../hack/manifest.yaml \
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
