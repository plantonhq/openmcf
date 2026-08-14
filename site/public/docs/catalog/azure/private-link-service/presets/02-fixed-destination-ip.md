---
title: "Fixed Destination IP"
description: "This preset publishes a single-instance service at one fixed private IP -- Private Link without a load balancer in front. Consumer traffic is NATed straight to the destination address."
type: "preset"
rank: "02"
presetSlug: "02-fixed-destination-ip"
componentSlug: "private-link-service"
componentTitle: "Private Link Service"
provider: "azure"
icon: "package"
order: 2
---

# Fixed Destination IP

This preset publishes a single-instance service at one fixed private IP -- Private Link without a load balancer in front. Consumer traffic is NATed straight to the destination address.

## When to Use

- Single-instance services (an appliance, a legacy host) that need private cross-network consumption without deploying a load balancer
- Simple point publications where load-balancer semantics (probes, rules, pools) add nothing

## Key Configuration Choices

- **Exactly one destination form** -- this preset uses `destinationIpAddress`; the load-balancer list must stay empty (the spec enforces the exactly-one-of)
- **The destination must be reachable from the NAT subnet** -- the NAT addresses are where consumer traffic originates inside your network
- **Same subnet policy flag** -- the NAT subnet still needs `privateLinkServiceNetworkPoliciesEnabled: false`
- **Prefer the load-balancer shape for anything that scales** -- a fixed IP publication has no health awareness and no failover

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-pls-subnet-arm-id>` | The policies-disabled subnet's ARM ID | `AzureSubnet` status outputs (`subnet_id`) |
| `<your-service-private-ip>` | The service instance's private IPv4 | Your service's NIC configuration |
