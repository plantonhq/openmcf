---
title: "ExpressRoute Port"
description: "ExpressRoute Port deployment documentation"
icon: "package"
order: 100
componentName: "azureexpressrouteport"
---

# Azure ExpressRoute Port

Deploys an ExpressRoute Port -- your own pair of physical ports on a Microsoft edge router at a colocation facility (ExpressRoute Direct). The port is the capacity object of the largest hybrid estates: you order 10 or 100 Gbps of dual links, hand the port's letter-of-authorization facts to the facility to complete the cross-connects, and carve ExpressRoute circuits from its bandwidth. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ExpressRoute Port** -- the ARM object owning the physical link pair, with its peering location, bandwidth, encapsulation, billing model, optional managed identity, and per-link admin/MACsec configuration
- **Port Authorizations** -- one per `authorizations` entry: ARM-generated keys that let circuits in OTHER subscriptions be built on this port's capacity
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with your `tags`, applied to the port

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Resource Group** where the port will be created.
- **ExpressRoute Direct availability**: the peering-location vocabulary is `az network express-route port location list` -- and some subscriptions need Microsoft enrollment for ExpressRoute Direct before ARM accepts the create.
- **Budget sign-off**: the port bills its full monthly rate (one of the most expensive single objects in Azure networking) from the moment it is created.

## Deploy

### Console

Open the deployment store, find **Azure ExpressRoute Port**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Port** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRoutePort
metadata:
  name: hq-port
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "network-rg"
  name: hq-port
  peeringLocation: "Equinix-Ashburn-DC2"
  bandwidthInGbps: 10
  encapsulation: DOT1Q
  link1:
    adminEnabled: true
  link2:
    adminEnabled: true
```

```shell
planton apply -f azure-express-route-port.yaml
```

The port object provisions its links as a pair; the physical cross-connects are ordered out-of-band with the facility using the per-link outputs (router, interface, patch panel, rack).

### InfraChart

In a hybrid-connectivity chart, the port anchors the ExpressRoute Direct chain: port → Direct-mode circuit → private peering → gateway connection, each wiring to the previous by reference.

## Key Configuration

These are the most important decisions when configuring a port. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The physical facts** -- `peeringLocation` (the ExpressRoute DIRECT vocabulary, narrower than circuit peering locations), `bandwidthInGbps` (10 or 100 at real locations), and `encapsulation` (DOT1Q for one VLAN tag per circuit -- the common choice; QINQ for stacked tags when customer VLAN ranges overlap) are all fixed at creation.

**Billing model** -- `billingType` defaults to METERED_DATA (port fee plus per-GB outbound on its circuits); UNLIMITED_DATA is a higher flat fee that wins at sustained high utilization.

**MACsec** -- layer-2 encryption per link: both Key Vault secret IDs (CKN + CAK) travel together, and the port needs a USER_ASSIGNED identity with Key Vault secret GET to read them -- both contracts are spec-enforced upfront.

**Authorizations** -- each named entry issues an ARM-generated key (surfaced, sensitive, in `authorization_keys`) that a circuit in another subscription redeems. Deleting an entry revokes it.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureUserAssignedIdentity** | `identity.identityIds` | `status.outputs.identity_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `express_route_port_id` | Azure Resource Manager ID of the port | AzureExpressRouteCircuit's `expressRoutePortId` (Direct mode) |
| `express_route_port_name` | Name of the port | Operational tooling |
| `guid` / `ethertype` / `mtu` | Port-level physical facts | Facility ordering, network design |
| `link1_*` / `link2_*` | Per-link router, interface, patch panel, rack, connector | The facility's cross-connect (LOA) order |
| `system_assigned_identity_principal_id` | The system-assigned identity's principal | Key Vault access policies |
| `authorization_keys` | Name-keyed issued keys (sensitive) | Circuits in other subscriptions |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard port** -- 10 Gbps Dot1Q metered with both links enabled. Start from the **Standard Port** preset.

**MACsec-encrypted port** -- layer-2 encryption keyed from your Key Vault via a user-assigned identity. Start from the **MACsec Port** preset.

## Works With

- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- provides the resource group the port is created in
- [**Azure ExpressRoute Circuit**](/cloud-catalog/azure-express-route-circuit) -- Direct-mode circuits carved from this port's bandwidth
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- the identity MACsec uses to read Key Vault secrets
