# Internal NLB

This preset creates the simplest useful Network Load Balancer: an internal
NLB with one node in one private subnet. It carries no listeners by design —
attach `AwsLbListener` resources (TCP/UDP/TCP_UDP/TLS, forward-only) against
its `load_balancer_arn` output, with `AwsLbTargetGroup` resources as the
destinations. Listener rules do not apply at Layer 4; routing is purely by
port and protocol.

## When to Use

- Layer-4 entry for service-to-service traffic inside a VPC (gRPC, databases,
  Redis, custom TCP protocols)
- The load-balancer half of a PrivateLink endpoint service (an internal NLB
  is the required target of a VPC endpoint service)
- Development and staging environments where one AZ is an accepted trade-off

## Key Configuration Choices

- **Internal scheme** (`internal: true`) — the NLB gets a private DNS name
  and is unreachable from the internet; the scheme is immutable, so this is
  a create-time decision
- **One subnet mapping** — the spec minimum. Subnet mappings are add-only in
  AWS: you can add AZs later but never remove one, so starting small is the
  reversible direction
- **No security groups** — the NLB admits all traffic on its listener ports,
  relying on VPC isolation. Deliberate: attaching security groups is a
  one-way door (the last one can never be removed), so the plain preset
  leaves that commitment unmade
- **AWS defaults everywhere else** — cross-zone distribution stays off (its
  AWS default), IP addressing stays `ipv4`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<nlb-name>` | Unique name for the NLB (AWS caps it at 32 characters) | Choose a descriptive name (e.g., `internal-grpc`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<private-subnet-id>` | Private subnet for the NLB node | AWS VPC console or `AwsSubnet` status outputs |

## Common Additions

- Add a second subnet mapping in another Availability Zone for zonal
  redundancy (recommended for anything production-bound)
- Set `privateIpv4Address` on a mapping to pin the node to a fixed address
  when downstream systems reference the NLB by IP
- Add `crossZoneLoadBalancingEnabled: true` if targets are unevenly
  distributed across AZs — and budget for the inter-AZ transfer it bills
- Use a `valueFrom` reference to an `AwsSubnet` resource instead of a
  literal subnet ID

## Related Presets

- **02-static-ip-internet-facing** — the public variant with Elastic-IP
  static addresses
- **03-private-link-hardened** — this scheme plus security groups,
  PrivateLink enforcement, and access logs
