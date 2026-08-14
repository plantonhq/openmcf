# Microsoft Fabric Capacity

Deploys a Microsoft Fabric capacity -- the billing and compute anchor of Microsoft Fabric: workspaces assign themselves to it, and its F-SKU sets how much compute every workload on it shares. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Fabric capacity** -- the capacity itself: F-SKU, administrators, tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.

### Azure Subscription

- **A running capacity bills PER HOUR from the moment it exists** -- F2 is the smallest step; F2048 is a thousand times that. Deploy it when you are ready to use it.
- **At least one administrator is required** -- an Entra user principal name (e.g. `admin@contoso.com`) or a service principal's object ID. Administrators manage the capacity from the Fabric side.
- **Fabric is not available in every Azure region** -- check the Fabric region availability list before choosing.
- **The SKU tier is always "Fabric"** -- Azure defines no other value today, so the platform sends it for you.

## Deploy

### Console

Open the deployment store, find **Microsoft Fabric Capacity**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Starter Capacity** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f fabric-capacity.yaml
```

## After Deploy

Assign workspaces to the capacity from the Fabric portal (workspace settings -> License info -> select the capacity by its `fabric_capacity_name`). Watch capacity-unit utilization in the Fabric Capacity Metrics app and move the `sku_name` up or down in place as usage settles; the capacity bills per hour for as long as it exists, so delete environments that are not in use.
