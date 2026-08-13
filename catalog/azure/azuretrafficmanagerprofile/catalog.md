# Azure Traffic Manager Profile

Deploys a Traffic Manager profile -- Azure's DNS-based traffic director, which answers lookups on its `{relative-name}.trafficmanager.net` name with the address of one of its endpoints, chosen by routing method and endpoint health. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Traffic Manager profile** -- the global routing object with its DNS identity and health-probe configuration (endpoints are separate AzureTrafficManagerEndpoint resources referencing it)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **A resource group** -- the profile's metadata record lives in a referenced resource group (the profile itself is global).

### Azure Subscription

- **The DNS relative name is globally unique across ALL of Azure** -- the trafficmanager.net namespace is shared; Azure rejects a taken name at apply time, so prefix with your organization.
- **Traffic Manager is a GLOBAL service** -- the spec carries no region; nothing about the profile is regional.
- **Billing is per million DNS queries plus per-endpoint health probes** -- fast-interval probes and Traffic View bill extra; the profile object itself is cheap at rest.

## Deploy

### Console

Open the deployment store, find **Azure Traffic Manager Profile**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Performance Routing** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f profile.yaml
```

## After Deploy

Add destinations with AzureTrafficManagerEndpoint resources referencing the `traffic_manager_profile_id` output, then point your own domain at the profile with a CNAME to the `fqdn` output. The profile answers only with healthy endpoints -- check endpoint monitor status in the portal when an answer surprises you.
