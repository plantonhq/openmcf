# Internal Load Balancer

This preset creates an internal (private VNet) Azure Load Balancer: a zone-redundant frontend with a pinned static private address in your subnet, an `app` backend pool, a TCP health probe, and a TCP rule on port 8080 with a raised idle timeout for long-lived east-west connections.

## When to Use

- Load balancing internal service tiers (APIs, databases, middleware) that must never be internet-reachable
- A stable, DNS-addressable private entry point for a service backed by VMs or a scale set
- Multi-tier architectures where a public LB fronts the web tier and this internal LB fronts the app tier

## Key Configuration Choices

- **Internal frontend with a pinned address** (`subnetId` + `privateIpAddress`) -- DNS records, firewall rules, and service discovery stay stable across redeployments
- **Zone-redundant frontend** (`zones: ["1","2","3"]`) -- the private address survives a zone outage; changing zones later replaces the frontend, so pick the posture up front
- **TCP probe** -- port-open health checking; switch to `PROBE_HTTP` with a `requestPath` when the workload exposes a health endpoint
- **Raised idle timeout with TCP reset** -- long-lived connections aren't silently dropped, and both ends learn immediately when they are

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match backend resources) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-lb-name>` | Name for the load balancer (unique within resource group) | Your naming convention |
| `<subnet-resource-id>` | Full ARM resource ID of the frontend subnet | `AzureSubnet` status outputs |
| `<static-private-ip>` | An unassigned address inside the subnet's range | Your IP allocation plan |

## Related Presets

- **01-public** -- internet-facing load balancing
- **03-outbound-and-nat** -- explicit SNAT egress and admin port forwarding
