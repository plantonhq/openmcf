# High Capacity (Production Hardened)

This preset is the production posture for heavy workloads: the largest
Consumption-plan sizing (2 vCPUs / 4Gi) with a pinned version. Nothing
tracks `latest` -- the build is chosen deliberately, and rollback is
re-pinning the previous tag.

## When to Use

- Large stacks or high operation concurrency (memory pressure shows up
  as failed IaC operations mid-apply)
- Production environments with change control: versions are pinned and
  bumped deliberately

## Key Configuration Choices

- **`cpu: 2` / `memory: 4Gi`** -- on the Consumption plan the pairing is
  fixed at memory = cpu x 2 (0.5/1Gi, 1/2Gi, 2/4Gi ...), so sizing
  memory up means sizing cpu with it; 2/4Gi is the largest pairing and
  the spec validates the combination up front instead of letting Azure
  reject it mid-deploy
- **`runnerVersion: v0.4.0`** -- a pinned tag (replace with the current
  release); new replicas pull the tag on every (re)start, so bumping it
  is the whole upgrade
- **Token as a managed-secret reference** -- the token authorizes
  joining only; the runner registers itself on first boot and receives
  its own individually revocable identity, and replica replacement
  re-joins with the same token

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |
| `<resource-group-resource-name>` | Name of the AzureResourceGroup resource | Your resource group manifest's `metadata.name` |
| `<container-app-environment-resource-name>` | Name of the AzureContainerAppEnvironment resource | Your environment manifest's `metadata.name` |

## Related Presets

- `01-environment-runner` -- the default-sized starting point
