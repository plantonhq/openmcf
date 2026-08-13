# Azure DNS Forwarding Ruleset

Deploys a DNS forwarding ruleset for Azure DNS Private Resolver -- the declarative rule book that sends queries for chosen domains (say, your corporate Active Directory namespace) to the DNS servers that own them, through the resolver's outbound endpoint. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS forwarding ruleset** -- bound to the resolver's outbound endpoint(s)
- **Forwarding rules** (optional, up to 1,000) -- one per captured domain, each with its ordered target DNS servers

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzurePrivateDnsResolver with an outbound endpoint** -- the ruleset binds its ARM id (the preset wires it by reference, which also orders the deploy).

### Azure Subscription

- **A ruleset binds at most 2 outbound endpoints, both from the SAME resolver** -- Azure enforces it at deploy time.
- **Rules take effect only in LINKED networks** -- deploy AzurePrivateDnsResolverVirtualNetworkLink per network afterwards; the ruleset alone steers nothing.
- **Free at rest** -- only the resolver's endpoint hours and query volume bill.

## Deploy

### Console

Open the deployment store, find **Azure DNS Forwarding Ruleset**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **On-Premises Domain** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f ruleset.yaml
```

## After Deploy

Link each virtual network that should use the rules (AzurePrivateDnsResolverVirtualNetworkLink referencing the `dns_forwarding_ruleset_id` output) -- linked networks start forwarding immediately, no per-network DNS settings needed.
