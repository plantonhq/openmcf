---
title: "Allowlisted Address from a Prefix"
description: "This preset allocates the public IP from a reserved `AzurePublicIpPrefix` instead of Microsoft's general pool. Every address drawn from the prefix falls inside one contiguous, pre-communicated CIDR..."
type: "preset"
rank: "03"
presetSlug: "03-allowlisted-from-prefix"
componentSlug: "public-ip"
componentTitle: "Public IP"
provider: "azure"
icon: "package"
order: 3
---

# Allowlisted Address from a Prefix

This preset allocates the public IP from a reserved `AzurePublicIpPrefix` instead of Microsoft's general pool. Every address drawn from the prefix falls inside one contiguous, pre-communicated CIDR -- partners and firewalls allowlist the prefix once, and any address you later allocate from it is already admitted.

## When to Use

- Endpoints whose address must fall inside a range a partner has already allowlisted
- Architectures that pre-reserve egress/ingress ranges before individual addresses exist
- Avoiding per-address firewall-change tickets as your endpoint count grows

## Key Configuration Choices

- **`publicIpPrefixId` reference** -- resolves to the prefix's `public_ip_prefix_id` output; the prefix must be in the same region and zones as this address
- **Zones must match the prefix** -- an address inherits its zonal shape from the range it is drawn from; this preset assumes a zone-redundant prefix
- **Fixed at creation** -- an address cannot move into or out of a prefix later

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the prefix's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-public-ip-prefix-resource-name>` | Planton metadata name of the `AzurePublicIpPrefix` | Your prefix resource |
