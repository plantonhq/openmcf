---
title: "Shared Zone with Public DNS Fallback"
description: "This preset links a zone with `NX_DOMAIN_REDIRECT` resolution: names the private zone cannot answer are retried against public DNS instead of failing. The fallback pattern for privatelink zones..."
type: "preset"
rank: "03"
presetSlug: "03-public-fallback"
componentSlug: "private-dns-zone-virtual-network-link"
componentTitle: "Private DNS Zone Virtual Network Link"
provider: "azure"
icon: "package"
order: 3
---

# Shared Zone with Public DNS Fallback

This preset links a zone with `NX_DOMAIN_REDIRECT` resolution: names the
private zone cannot answer are retried against public DNS instead of
failing. The fallback pattern for privatelink zones shared across
environments -- where some service instances have private endpoints (and
private records) while others are still reached publicly.

Without the redirect, a linked network treats the private zone as
authoritative: any name missing from it returns NXDOMAIN even when a
public record exists. With it, private records win when present and
public resolution covers the rest.

## When to Use

- A shared privatelink zone serving several environments, where only some
  service instances have private endpoints
- Migration windows: services move to private endpoints one at a time
  while the rest stay publicly resolvable

## Key Configuration Choices

- **`NX_DOMAIN_REDIRECT`** -- private-first, public-fallback; switch to
  `DEFAULT` (strict private) once every instance has its private record
- **Unset is also valid** -- Azure applies its own per-zone-type default;
  set the policy only when you need the behavior pinned

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `shared-vnet-link` | A short name identifying the linked network | Your naming convention |
| `<private-dns-zone-arm-id>` | The shared zone's full ARM ID | The zone's `status.outputs.zone_id` |
| `<virtual-network-arm-id>` | The network's full ARM ID | The network's `status.outputs.virtual_network_id` |
