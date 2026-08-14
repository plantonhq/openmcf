---
title: "DNS Resolver Virtual Network Link"
description: "DNS Resolver Virtual Network Link deployment documentation"
icon: "package"
order: 100
componentName: "azureprivatednsresolvervirtualnetworklink"
---

# Azure DNS Resolver Virtual Network Link

Deploys the virtual network link that attaches ONE virtual network to a DNS forwarding ruleset -- the switch that turns the rule book on for that network. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Virtual network link** -- a child of the forwarding ruleset, one per ruleset-network pair

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzurePrivateDnsResolverForwardingRuleset** -- the rule book this link turns on (wired by reference, which also orders the deploy).
- **An AzureVirtualNetwork** to attach -- in the ruleset's region; peering to the resolver's network is NOT required, and cross-subscription networks are allowed.

### Azure Subscription

- **Links are free at rest** -- only the resolver's endpoint hours and query volume bill.
- **Everything except metadata is create-only** -- changing the network or ruleset replaces the link (a brief forwarding gap for that network).

## Deploy

### Console

Open the deployment store, find **Azure DNS Resolver Virtual Network Link**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Spoke Link** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f link.yaml
```

## After Deploy

Resources in the linked network start forwarding immediately -- queries for the ruleset's domains flow to their target servers; everything else resolves inside Azure as before. Verify with a lookup of a forwarded domain from a VM in the linked network.
