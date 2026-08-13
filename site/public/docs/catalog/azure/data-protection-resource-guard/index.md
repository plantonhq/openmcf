---
title: "Data Protection Resource Guard"
description: "Data Protection Resource Guard deployment documentation"
icon: "package"
order: 100
componentName: "azuredataprotectionresourceguard"
---

# Azure Data Protection Resource Guard

Creates a Data Protection Resource Guard -- the approval gate behind Multi-User Authorization for Azure Backup vaults. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Resource Guard** -- the Multi-User Authorization gate, with its critical-operation exclusion list

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureResourceGroup** -- ideally one a DIFFERENT administrator controls than the vaults the guard will protect (that separation IS the security model).

### Azure Subscription

- **The guard is free** -- it is a pure governance object.
- **Scope separation is the whole point** -- a guard in the same scope as its vaults is a speed bump, not a control. Plan who owns the guard's resource group before deploying.
- **Vaults opt in by reference** -- creating the guard changes nothing until a vault references its ARM ID.

## Deploy

### Console

Open the deployment store, find **Azure Data Protection Resource Guard**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **MUA Guard** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f guard.yaml
```

## After Deploy

The guard's `resource_guard_id` output is what backup vaults reference to put their privileged operations under the guard's approval requirement.
