---
title: "DNAT: Publish an Internal Service"
description: "This preset creates a DNAT group that publishes an internal service through the firewall's public IP -- inbound HTTPS arriving at the firewall is translated to an internal frontend (typically an..."
type: "preset"
rank: "02"
presetSlug: "02-dnat-publish-service"
componentSlug: "firewall-policy-rule-collection-group"
componentTitle: "Firewall Policy Rule Collection Group"
provider: "azure"
icon: "package"
order: 2
---

# DNAT: Publish an Internal Service

This preset creates a DNAT group that publishes an internal service
through the firewall's public IP -- inbound HTTPS arriving at the
firewall is translated to an internal frontend (typically an internal
load balancer). A matching DNAT rule implicitly allows the translated
flow, so no companion network rule is needed.

## When to Use

- Exposing an internal load balancer or jumpbox without a public IP of
  its own
- Centralizing ALL inbound exposure at the firewall chokepoint, where
  threat intelligence and (on Premium) IDPS inspect it

## Key Configuration Choices

- **Group priority 100** -- DNAT collections are processed before every
  network/application collection anyway, but a low group priority keeps
  inbound rules evaluated ahead of other groups' DNAT entries
- **`destinationAddress`** must be one of the firewall's public IPs --
  compose it from the referenced `AzurePublicIp`'s `ip_address` output
- **One port entry** -- ARM caps a DNAT rule at a single destination-port
  entry (a port or a range); publish additional ports as additional rules

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<firewall-policy-name>` | The parent AzureFirewallPolicy | That policy's Planton resource name |
| `<firewall-public-ip>` | The firewall's public address | The referenced `AzurePublicIp`'s `status.outputs.ip_address` |
| `<internal-frontend-ip>` | The internal address to translate to | The internal load balancer's frontend IP |
