# Azure Traffic Manager Endpoint

Deploys one destination of a Traffic Manager profile: a public Azure resource by ARM ID (`azure`), a DNS name or IP address (`external`), or another profile composing routing trees (`nested`) -- exactly one variant per endpoint. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **One Traffic Manager endpoint** of the type your spec's variant declares, inside the referenced profile

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureTrafficManagerProfile** -- the endpoint is created inside a referenced profile (the preset wires it by reference, which also orders the deploy).

### Azure Subscription

- **Which fields matter depends on the PROFILE's routing method** -- weight (Weighted), priority (Priority), geo claims (Geographic: every code claimed by exactly one endpoint), subnet claims (Subnet: no overlaps) -- Azure evaluates them at apply time.
- **Azure endpoints need a PUBLIC address on the target** -- the referenced resource must hold a public IP (Standard tier for Public IPs; Basic is not steerable); external targets under Performance routing need an explicit `endpointLocation`.
- **Endpoints are free at rest** -- probes and queries bill on the profile.

## Deploy

### Console

Open the deployment store, find **Azure Traffic Manager Endpoint**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **External Endpoint** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f endpoint.yaml
```

## After Deploy

The profile starts probing the endpoint immediately; it enters DNS answers once probes pass. Drain it anytime with `enabled: false` (it leaves answers, configuration stays). Check the portal's endpoint monitor status when an answer surprises you.
