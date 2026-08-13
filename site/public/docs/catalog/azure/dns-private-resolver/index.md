---
title: "DNS Private Resolver"
description: "DNS Private Resolver deployment documentation"
icon: "package"
order: 100
componentName: "azureprivatednsresolver"
---

# Azure DNS Private Resolver

Deploys Azure DNS Private Resolver -- the managed DNS proxy that resolves names across the hybrid boundary: on-premises queries INTO Azure through inbound endpoints, Azure queries OUT to on-premises DNS through outbound endpoints. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Private Resolver** -- the managed proxy, anchored to one virtual network
- **Inbound endpoints** (optional, up to 5) -- private IPs that answer DNS queries sent to them, one per dedicated delegated subnet
- **Outbound endpoints** (optional, up to 5) -- egress points for queries leaving Azure, one per dedicated delegated subnet

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureVirtualNetwork** the resolver anchors to -- Azure allows at most ONE resolver per network.
- **One dedicated AzureSubnet per endpoint**, delegated to `Microsoft.Network/dnsResolvers`, sized /28 to /24, hosting nothing else.

### Azure Subscription

- **Endpoints bill hourly** from the moment they provision; the resolver object itself is free.
- **Everything except tags is create-only** -- plan the endpoint layout up front; edits replace the affected endpoint.

## Deploy

### Console

Open the deployment store, find **Azure DNS Private Resolver**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Hybrid Resolver** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f resolver.yaml
```

## After Deploy

Point on-premises DNS conditional forwarders at the `inbound_endpoint_ip` output -- that is the address that answers with Azure's private view. The `outbound_endpoint_id` output is what an AzurePrivateDnsResolverForwardingRuleset binds to steer outbound queries.
