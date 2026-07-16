---
title: "Production Locked-Down Account"
description: "This preset creates a geo-zone-redundant account with the production security posture: a DENY firewall admitting only declared subnets and trusted Microsoft services, anonymous access made..."
type: "preset"
rank: "02"
presetSlug: "02-production-locked-down"
componentSlug: "storage-account"
componentTitle: "Storage Account"
provider: "azure"
icon: "package"
order: 2
---

# Production Locked-Down Account

This preset creates a geo-zone-redundant account with the production
security posture: a DENY firewall admitting only declared subnets and
trusted Microsoft services, anonymous access made unrepresentable, SAS
lifetimes policed, full blob data protection, and a lifecycle schedule
that tiers stale data down automatically.

## When to Use

- Production application data that must survive zone AND regional loss
- Accounts subject to a security baseline (no anonymous access, no
  unbounded SAS tokens, no open network path)
- Data with a natural hot-to-cold aging curve (the lifecycle rule pays
  for itself)

## Key Configuration Choices

- **GZRS replication** -- ZRS at home plus async geo-replication
- **`defaultAction: DENY` + subnet references** -- the data plane admits
  only declared networks; ARM management operations are unaffected
- **`allowNestedItemsToBePublic: false`** -- no container can ever opt
  into anonymous access
- **SAS policy `7.00:00:00` / BLOCK** -- tokens older than a week are
  rejected, not just logged
- **Lifecycle rule** -- cool after 30 days without a read (bouncing back
  to hot on access), archive after 180 days without modification,
  delete after 2 years

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your resource-group composition |
| `myorgprodstorage` | 3-24 lowercase letters/digits, globally unique | Your naming convention (no hyphens!) |
| `<app-subnet-resource-name>` | The AzureSubnet admitted through the firewall (needs the Microsoft.Storage service endpoint) | Your network composition |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Remix Ideas

- Add `customerManagedKey` + `identity` for bring-your-own-key
  encryption (see the account catalog page's production example)
- Set `publicNetworkAccessEnabled: false` and front the account with
  private endpoints to remove the public endpoint entirely
- Set `sharedAccessKeyEnabled: false` once every consumer authenticates
  with Microsoft Entra
