# Azure ExpressRoute Circuit

Deploys an ExpressRoute circuit -- the dedicated PRIVATE connection between your infrastructure and Microsoft, bought through a connectivity provider (Equinix, Megaport, ...) or carved from your own ExpressRoute Direct port. Creating the circuit issues the service key your provider uses to provision the physical cross-connect; peerings (AzureExpressRouteCircuitPeering) then make routes flow through it. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ExpressRoute Circuit** -- the ARM billing/identity object with its SKU, provisioning mode, and generated service key
- **Circuit Authorizations** -- one per `authorizations` entry: ARM-generated keys that let virtual network gateways in OTHER subscriptions connect to this circuit
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with your `tags`, applied to the circuit

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Resource Group** where the circuit will be created.
- **The connectivity facts**: your provider's exact name and peering location plus the bandwidth step you bought -- or your ExpressRoute Direct port's ARM ID. `az network express-route list-service-providers` shows the provider vocabulary.

## Deploy

### Console

Open the deployment store, find **Azure ExpressRoute Circuit**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Provider Circuit** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRouteCircuit
metadata:
  name: hq-circuit
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "network-rg"
  name: hq-circuit
  skuTier: STANDARD
  skuFamily: METERED_DATA
  serviceProviderName: "Equinix"
  peeringLocation: "Washington DC"
  bandwidthInMbps: 1000
```

```shell
planton apply -f azure-express-route-circuit.yaml
```

The circuit provisions in minutes and sits in `NotProvisioned` until you hand the `service_key` output to your provider. **Billing starts at creation**, even before the provider completes the cross-connect.

### InfraChart

In a hybrid-connectivity chart, the circuit anchors the chain: circuit → private peering → ExpressRoute-type virtual network gateway → gateway connection, each wiring to the previous by reference.

## Key Configuration

These are the most important decisions when configuring a circuit. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Provisioning mode** -- exactly one (spec-enforced): the service-provider trio (`serviceProviderName` + `peeringLocation` + `bandwidthInMbps`, all fixed at creation except bandwidth, which can only GROW) or the ExpressRoute Direct pair (`expressRoutePortId` + `bandwidthInGbps`).

**SKU** -- `skuTier` sizes reach (LOCAL: the circuit's metro only, no egress fees; STANDARD: the geopolitical area -- the common choice; PREMIUM: global reach and higher route limits) and `skuFamily` sizes billing (METERED_DATA per outbound GB, UNLIMITED_DATA flat -- economical above roughly two-thirds sustained utilization).

**Authorizations** -- each named entry issues an ARM-generated key (surfaced, sensitive, in `authorization_keys`) that a gateway in another subscription redeems. Deleting an entry revokes it.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `express_route_circuit_id` | Azure Resource Manager ID of the circuit | Operational tooling |
| `express_route_circuit_name` | Name of the circuit | AzureExpressRouteCircuitPeering's `expressRouteCircuitName` |
| `service_key` | The provisioning credential (sensitive) | Handed to the connectivity provider |
| `service_provider_provisioning_state` | NotProvisioned / Provisioning / Provisioned / Deprovisioning | Gating peering deployment |
| `authorization_keys` | Name-keyed issued keys (sensitive) | Far-side gateway connections in other subscriptions |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Provider circuit** -- The common shape: a carrier-provisioned circuit at a named peering location. Start from the **Provider Circuit** preset.

**Metro-local circuit** -- LOCAL tier for egress-free connectivity in one metro. Start from the **Local Metro Circuit** preset.

## Works With

- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- provides the resource group the circuit is created in
- [**Azure ExpressRoute Circuit Peering**](/cloud-catalog/azure-express-route-circuit-peering) -- the BGP configuration that makes routes flow
- [**Azure Virtual Network Gateway**](/cloud-catalog/azure-virtual-network-gateway) -- the EXPRESS_ROUTE-type gateway that connects a VNet to the circuit's private peering
