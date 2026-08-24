# Azure Bastion Host

Deploys Azure Bastion -- the managed jump service for RDP/SSH sessions to virtual machines over their private addresses, with no public IPs on the machines themselves. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Bastion host** -- the jump service, on dedicated infrastructure (Basic/Standard/Premium: the `AzureBastionSubnet` + a Standard static public IP) or Azure-shared infrastructure (Developer: attached to a virtual network)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureSubnet named EXACTLY `AzureBastionSubnet`** (/26 or larger) in the target virtual network -- dedicated-infrastructure SKUs deploy only there.
- **An AzurePublicIp** (Standard SKU, static) the host binds EXCLUSIVELY -- required on Basic/Standard; Premium may omit it for a private-only host.
- Developer SKU instead references the **AzureVirtualNetwork** directly -- no subnet, no public IP.

### Azure Subscription

- **Billing starts at provisioning**: hourly per SKU plus per-scale-unit on Standard/Premium. Developer is free.
- **Creates take ~10 minutes**; deletes are similar. SKU upgrades are in-place, downgrades replace the host.
- **Developer SKU is region-limited** -- check availability before choosing it.

## Deploy

### Console

Open the deployment store, find **Azure Bastion Host**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Host** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f bastion.yaml
```

## After Deploy

Sessions open from the Azure portal's Connect blade (or, with tunneling on Standard/Premium, from `az network bastion ssh/rdp/tunnel`). The `dns_name` output is the session endpoint; `private_only_enabled` reports whether the host is reachable only from connected networks.
