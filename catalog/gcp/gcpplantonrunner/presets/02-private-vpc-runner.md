# Private VPC Runner

This preset routes the runner's egress into a VPC through Direct VPC
egress, which is what lets it reach private endpoints -- a private GKE
control plane, a private-IP Cloud SQL instance -- from inside the
network. Only traffic to private ranges rides the VPC; the runner's
dial-out to the control plane keeps its normal internet path, so joining
never depends on the VPC's routing.

## When to Use

- A private GKE control plane needs to be deployed to and operated
- Any target reachable only over private IPs (Cloud SQL private IP,
  internal services)
- You want zero inbound exposure: the runner only dials out, and the
  deployed service accepts no meaningful inbound traffic

## Key Configuration Choices

- **`vpcAccess` with network + subnetwork** -- Direct VPC egress needs
  no connector infrastructure; the runner draws IPs from the subnetwork,
  which must be in the runner's region. The refs here point at
  `GcpVpcNetwork`/`GcpSubnetwork` resources -- replace the whole
  `valueFrom` block with a literal `value:` if the VPC was built outside
  the platform.
- **Split egress by design** -- only private-range traffic enters the
  VPC; the control-plane dial-out stays on the internet path, so a
  misconfigured VPC route can break target reach but never the join
- **A network tag (`planton-runner`)** -- how VPC firewall rules select
  the runner's egress, e.g. a rule admitting it to the private GKE
  control plane by tag
- **Token as a managed-secret reference** -- the token authorizes
  joining only; the runner registers itself on first boot and receives
  its own individually revocable identity, and revoking the token never
  touches runners it already admitted

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |
| `my-vpc` | Name of your GcpVpcNetwork resource | Your VPC manifest's `metadata.name` |
| `my-subnet` | Name of your GcpSubnetwork resource (same region as the runner) | Your subnetwork manifest's `metadata.name` |

## Related Presets

- `01-regional-runner` -- the minimal shape when no private endpoints are involved
- `03-high-capacity` -- pinned version and larger sizing for heavy IaC workloads
