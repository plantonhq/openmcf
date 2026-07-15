# Hub-Spoke Egress Firewall

This preset creates the production hub firewall: zone-redundant,
STANDARD tier, policy-attached, deployed into the hub VNet's
`AzureFirewallSubnet` with one Standard static public IP. It is the
chokepoint every spoke's default route points at.

## When to Use

- The hub of any hub-spoke network that centralizes egress control
- Landing zones where spokes must not reach the internet directly

## Key Configuration Choices

- **`zones: ["1","2","3"]`** -- zone redundancy is free on Azure
  Firewall and is the production posture; zones are fixed at creation
- **`skuTier: STANDARD`** -- matches the policy's tier (both must agree);
  go PREMIUM only when the policy carries TLS inspection or IDPS
- **One public IP to start** -- add ip_configurations (each with its own
  PIP, no subnet) later if SNAT port exhaustion appears; each extra IP
  adds ports in place

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group | Its `status.outputs.resource_group_name` |
| `<firewall-subnet-name>` | The AzureSubnet named exactly `AzureFirewallSubnet` (/26+) | That subnet's Planton resource name |
| `<firewall-pip-name>` | The Standard static AzurePublicIp | That address's Planton resource name |
| `<firewall-policy-name>` | The AzureFirewallPolicy to enforce | That policy's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Steer every spoke through the firewall:

```yaml
# AzureRouteTable (one per spoke, attached to the spoke's subnets)
spec:
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress:
        valueFrom:
          name: <the firewall resource's name>
  bgpRoutePropagationEnabled: false
```
