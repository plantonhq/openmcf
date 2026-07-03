---
title: "Zone-Redundant StandardV2 NAT Gateway"
description: "This preset creates Azure's next-generation StandardV2 NAT Gateway: zone-redundant automatically, with no zone pinning (`zones` must be left empty -- the spec enforces it). Where a Standard gateway..."
type: "preset"
rank: "03"
presetSlug: "03-zone-redundant-v2"
componentSlug: "nat-gateway"
componentTitle: "NAT Gateway"
provider: "azure"
icon: "package"
order: 3
---

# Zone-Redundant StandardV2 NAT Gateway

This preset creates Azure's next-generation StandardV2 NAT Gateway: zone-redundant automatically, with no zone pinning (`zones` must be left empty -- the spec enforces it). Where a Standard gateway needs one deployment per zone for resilience, one StandardV2 gateway survives any single zone's failure by itself.

## When to Use

- Zone-resilient egress without deploying and operating one gateway per zone
- New architectures in regions where StandardV2 is available
- Consolidating multi-zone egress behind a single gateway and address set

## Key Configuration Choices

- **`skuName: STANDARD_V2`** -- fixed at creation; a gateway cannot change SKU later
- **No `zones`** -- StandardV2 is zone-redundant by itself; setting zones is rejected
- **StandardV2 addresses required** -- the referenced `AzurePublicIp` (or prefix) must itself be `sku: STANDARD_V2`; a Standard address will not attach

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (StandardV2 availability varies by region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-standardv2-public-ip-resource-name>` | Planton metadata name of a `STANDARD_V2` `AzurePublicIp` | Your public IP resource |
