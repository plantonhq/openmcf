---
title: "Shared Capacity Port"
description: "This preset creates the multi-tenant shape: a 100 Gbps QinQ port whose capacity is shared across subscriptions through issued authorizations. Each named authorization generates a key (surfaced,..."
type: "preset"
rank: "03"
presetSlug: "03-shared-capacity-port"
componentSlug: "expressroute-port"
componentTitle: "ExpressRoute Port"
provider: "azure"
icon: "package"
order: 3
---

# Shared Capacity Port

This preset creates the multi-tenant shape: a 100 Gbps QinQ port whose capacity is shared across subscriptions through issued authorizations. Each named authorization generates a key (surfaced, sensitive, in the `authorization_keys` output) that a circuit in another subscription redeems to be carved from this port.

## When to Use

- One organization-owned port serving circuits for multiple teams or subscriptions
- Platform teams selling internal connectivity capacity to product teams

## Key Configuration Choices

- **QINQ makes tenant VLANs coexist** -- Azure manages an outer S-Tag per circuit, so consuming teams' VLAN plans never collide
- **UNLIMITED_DATA suits aggregation** -- a shared 100 Gbps port typically runs at the sustained utilization where flat-rate beats metering
- **One authorization per consumer, named for the consumer** -- deleting the entry revokes that consumer's access; treat the keys as credentials
- **Oversubscription is designed-for** -- aggregate circuit bandwidth may exceed the port (up to 2x); plan for burst sharing, not sustained simultaneous saturation

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-peering-location>` | The ExpressRoute Direct facility | `az network express-route port location list` |
