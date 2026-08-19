# High Capacity (Production Hardened)

This preset is the production posture for heavy workloads: a sized-up
runner (4 vCPUs / 4Gi) with a pinned version. Nothing tracks `latest` --
the build is chosen deliberately, and rollback is re-pinning the
previous tag.

## When to Use

- Large stacks or high operation concurrency (memory pressure shows up
  as failed IaC operations mid-apply)
- Production environments with change control: versions are pinned and
  bumped deliberately

## Key Configuration Choices

- **`memory: 4Gi` before more cpu** -- IaC operations are memory-hungry
  (provider plugins, plans held in memory); memory is the dimension that
  fails first. Cloud Run requires at least 2Gi at 4 vCPUs, so 4Gi leaves
  real headroom.
- **`cpu: "4"`** -- Cloud Run admits 1, 2, 4, 6, or 8 vCPUs; 4 handles
  concurrent operations without queueing behind a single plan
- **`runnerVersion: v0.4.0`** -- a pinned tag (replace with the current
  release); new instances pull the tag on every (re)start, so bumping it
  is the whole upgrade
- **Token as a managed-secret reference** -- the token authorizes
  joining only; the runner registers itself on first boot and receives
  its own individually revocable identity, and instance replacement
  re-joins with the same token

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |

## Related Presets

- `01-regional-runner` -- the default-sized starting point
- `02-private-vpc-runner` -- routes private-range egress into a VPC to reach private endpoints
