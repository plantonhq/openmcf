# Dual Mode (Deploys + Live CloudOps)

This preset runs the runner in `dual` mode: everything the private VPC
worker does, plus the real-time CloudOps channel -- live browsing of the
resources behind the runner (pods, services, cluster state) from the
console, through the runner, without any inbound network path.

## When to Use

- Teams that operate (not just deploy) private clusters from the console
  and want live resource views through the runner
- Day-2 troubleshooting of in-network targets without a bastion or VPN

## Key Configuration Choices

- **`executionMode: dual`** -- the deploy worker AND the CloudOps
  channel; the channel rides an outbound-initiated tunnel, so no inbound
  rule is ever needed
- **Tunnel material comes from the credentials document** -- generate the
  credentials for a tunneled runner (the default registration shape);
  a temporal-only credentials document will fail this mode's startup
  validation by design
- **Everything else matches the private VPC worker** -- same placement,
  same secret handling, same sizing default

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Match the runner registration's name |
| `<aws-region>` | AWS region code | The region hosting the private targets |
| `<private-subnet-a/b-resource-name>` | Names of the private AwsSubnet resources | Your subnet manifests' `metadata.name` |
| `<runner-credentials-secret-slug>` | The managed secret holding the credentials JSON | Created from `planton runner generate-credentials <runner-name>` |

## Related Presets

- `01-private-vpc-worker` -- deploys only, the leanest posture
- `03-high-capacity` -- pinned version and larger sizing for heavy IaC workloads
