# Regional Runner

This preset deploys the minimal runner appliance: an always-on worker in
one GCP region that receives deploy operations from the control plane
and executes them from inside your project. Two decisions -- a region
and a token reference -- and everything else self-configures. The
30-second decision for getting a runner standing in a project.

## When to Use

- First runner in a GCP project: targets are reachable over their public
  endpoints and no VPC placement is needed yet
- You want the leanest possible manifest -- default sizing (1 vCPU /
  512Mi), latest runner build, the platform-created service account

## Key Configuration Choices

- **Token as a managed-secret reference** -- the token authorizes the
  runner to JOIN, nothing more: on first boot the runner registers
  itself with the control plane and receives its own individually
  revocable identity, and revoking the token never touches runners it
  already admitted. On Planton the platform mints the token and writes
  it at exactly this reference before the infrastructure applies.
- **No mode or replica knobs** -- everything beyond the join (work
  queue, tunnel, API endpoints) arrives in the join response, so the
  runner self-configures on arrival; it runs as exactly one always-on
  instance by design
- **`region: us-central1`** -- pick the region hosting the endpoints the
  runner must reach; latency to targets matters more than latency to the
  control plane
- **No service account set** -- the deployment creates a dedicated
  permissionless account (never the Compute Engine default), so the
  identity seam exists and permissions can be granted later without
  replacing the runner

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |

The `runner-token` secret slug is yours to choose -- on Planton the
platform writes the token there automatically; elsewhere, create a token
with `planton runner token create` and store it under that slug.

## Related Presets

- `02-private-vpc-runner` -- routes private-range egress into a VPC to reach private endpoints
- `03-high-capacity` -- pinned version and larger sizing for heavy IaC workloads
