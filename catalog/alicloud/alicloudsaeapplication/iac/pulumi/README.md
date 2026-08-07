# Pulumi Module to Deploy AliCloudSaeApplication

## CLI usage (Planton pulumi)

```bash
# Preview
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

# Update (apply)
planton pulumi update \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes

# Refresh
planton pulumi refresh \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

# Destroy
planton pulumi destroy \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```

## Debugging

This module includes a `debug.sh` helper. To enable debugging, edit `Pulumi.yaml` and uncomment the `runtime.options.binary` line so Pulumi runs the program via the script:

```yaml
name: planton-alicloud-module-test
runtime:
  name: go
#  options:
#    binary: ./debug.sh
```

Then make the script executable and run your command (e.g., `preview` or `update`). See `docs/pages/docs/guide/debug-pulumi-modules.mdx` for full instructions.

```bash
chmod +x debug.sh
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```

## Module Overview

This Pulumi module deploys an Alibaba Cloud SAE application using a single `sae.Application` resource. The module reads an `AliCloudSaeApplicationStackInput` protobuf message, initializes locals (tag merging, environment variable serialization), and provisions the application with the specified compute, networking, health check, update strategy, and logging configuration.

The module converts the `envs` map into the JSON array format that the SAE API expects, maps health check specs to the provider's `LivenessV2` and `ReadinessV2` types, and conditionally sets optional fields only when provided — avoiding zero-value overrides of provider defaults.

Both stack outputs (`app_id` and `app_name`) are exported for use by downstream components.

---

## Further Reading

- **[e2e/manifest.yaml](../../e2e/manifest.yaml)**: Minimal test manifest.
