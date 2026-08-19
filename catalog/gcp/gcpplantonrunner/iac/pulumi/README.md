# Pulumi Module to Deploy GcpPlantonRunner

This module provisions a standing Planton runner appliance on Cloud Run.
The pieces, in dependency order: the Cloud Run Admin API enablement
(`run.googleapis.com`, left enabled on destroy), the runtime service
account (created permissionless when no `service_account` is referenced
-- never the Compute Engine default), the token secret in Secret Manager
(`<name>-token`) with its value and an accessor grant scoped to exactly
that secret and exactly the runtime account, and finally the Cloud Run v2
service holding exactly one always-on runner (gen2, always-allocated
CPU, optional Direct VPC egress). The project, VPC, and any referenced
service account are referenced resources -- the module never creates or
mutates them. Every derivation matches the Terraform module key-for-key.

ENROLLMENT IS TOKEN-FIRST: the service ships the runner TOKEN (via
Secret Manager), never an identity. The runner joins the control plane
on first boot, registers itself, and receives its own individually
revocable identity; instance replacement re-joins with the same token
(its lineage re-admits the runner it originally admitted).

## Environment contract

The container's environment, exactly as the module wires it:

| Variable | Value |
|----------|-------|
| `PLANTON_RUNNER_TOKEN` | Secret Manager reference to `<name>-token`, version `latest`, resolved by Cloud Run at instance start -- never plaintext in the service definition |
| `PLANTON_RUNNER_NAME` | The registration name: `<env>-<metadata.name>`, or `metadata.name` outside an environment |
| `PLANTON_RUNNER_ENDPOINT` | Only when `control_plane_endpoint` is set; omitted, the runner's built-in hosted default applies |
| `PORT` | `50051` (the gRPC/CloudOps server; the container port Cloud Run routes to, named `h2c`) |
| `LOG_LEVEL` | `info` |

The startup probe reads the runner's health server at `/healthz` on port
8093; it answers independently of control-plane reachability, so a
runner whose control plane is momentarily unreachable still starts.

## The singleton law

The service is pinned to exactly one instance (min = max = 1). A
runner's identity is minted for ONE live instance -- a second instance
joining under the same name would revoke the first's key (token lineage:
re-admission re-mints and revokes). Never enable scaling here without
redesigning enrollment for fleets.

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
