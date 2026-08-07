# Black-Hole Guardrails

This preset drops traffic to a protected prefix at the routing layer
(`nextHopType: NONE`) -- a coarse, cheap guardrail that stops lateral
movement toward a sensitive tier (a data subnet, a management range)
from every subnet the table is attached to, regardless of NSG state.

Routing black-holes complement NSGs rather than replace them: an NSG
filters per-flow with rules and logging; a black-hole route removes the
path entirely and cannot be "allowed around" by a later NSG change.

## When to Use

- Preventing web-tier subnets from ever addressing the data tier
  directly (force traffic through an app tier or firewall)
- Cheap defense-in-depth under NSG rules for high-value ranges
- Temporarily cutting off a compromised prefix during incident response
  (routes update in place, immediately, for all attached subnets)

## Key Configuration Choices

- **`NONE` next hop** -- traffic is dropped silently (no ICMP
  unreachable); factor that into debugging expectations
- **Most-specific-prefix wins** -- a black-hole for a /24 coexists with
  a broader default route; add one route per protected range

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the table in | The resource group's `status.outputs.resource_group_name` |
| `<blocked-cidr>` | The protected prefix to drop traffic toward (e.g. `10.0.20.0/24`) | Your network plan |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
