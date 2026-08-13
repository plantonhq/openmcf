# Azure Availability Set

Deploys an availability set -- the classic placement grouping that spreads VMs across separate fault and update domains so one hardware failure or maintenance window cannot take them all down. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Availability set** -- the placement grouping: fault/update domain counts, managed-disk alignment, optional proximity placement group, tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name. VMs joining the set must live in the same group and region.

### Azure Subscription

- **Prefer zones where they exist** -- availability ZONES are the modern resilience unit; the set remains the right tool in regions without zones and for classic topologies. A VM uses a zone OR an availability set, never both.
- **Everything is fixed at creation** -- only tags update in place. VMs join the set when THEY are created; membership cannot change afterward.
- **Some regions support fewer than 3 fault domains** -- Azure rejects a count the region cannot provide.

## Deploy

### Console

Open the deployment store, find **Azure Availability Set**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Classic Web Tier** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f availability-set.yaml
```

## After Deploy

Create the tier's VMs with **AzureVirtualMachine**, referencing the set through `availability.availability_set_id` (it defaults to this kind's `availability_set_id` output). Two or more VMs in the set carry Azure's classic multi-VM SLA.
