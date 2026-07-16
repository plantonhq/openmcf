---
title: "Geo-DR Pairing"
description: "This preset pairs two PREMIUM namespaces for geo-disaster recovery with a scoped alias credential: metadata replicates continuously, and clients hold alias connection strings that survive a failover..."
type: "preset"
rank: "01"
presetSlug: "01-geo-dr-pairing"
componentSlug: "service-bus-disaster-recovery-config"
componentTitle: "Service Bus Disaster Recovery Config"
provider: "azure"
icon: "package"
order: 1
---

# Geo-DR Pairing

This preset pairs two PREMIUM namespaces for geo-disaster recovery with
a scoped alias credential: metadata replicates continuously, and
clients hold alias connection strings that survive a failover
unchanged.

## When to Use

- Business-critical messaging that must survive a regional outage with
  its topology intact
- Compliance regimes mandating a tested DR posture for messaging

## Key Configuration Choices

- **`aliasName`** -- point ALL clients at the alias, never at either
  namespace directly; that is what makes failover transparent
- **`aliasAuthorizationRuleId`** -- a namespace-scoped rule on the
  primary gives least-privilege alias credentials; unset falls back to
  the root rule
- **The partner must be empty at pairing time** -- pair immediately
  after creating the standby namespace, before any entity lands on it

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-app-bus-alias` | The stable DNS identity | Your naming convention |
| `my-bus-eastus` | The active PREMIUM namespace | Your messaging composition |
| `my-bus-westus` | The standby PREMIUM namespace (different region, empty) | Your DR composition |
| `my-dr-clients-rule` | A namespace-scoped AzureServiceBusAuthorizationRule on the primary | Your credential composition |

## Operating It

Failover is an operational action on the SECONDARY (portal/CLI/SDK)
during an incident -- it promotes the standby and BREAKS the pairing.
Re-pair to a new partner afterwards by updating `partnerNamespaceId`.
