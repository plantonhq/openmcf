# Overview

The **AzureTrafficManagerEndpoint** component deploys one destination of a Traffic Manager profile -- the place the DNS traffic director can send users. Exactly one variant per endpoint: a public Azure resource by ARM ID (`azure`), anything with a DNS name or IP (`external`), or another Traffic Manager profile (`nested`) composing routing methods into trees.

## Purpose

- **Destinations with independent lifecycles**: endpoints join, drain, and leave a profile without touching it or each other -- add a region, retire a datacenter, or shift weights one endpoint at a time.
- **Three target worlds, one grammar**: Azure resources (address tracked by Azure), external services (other clouds, on-premises), and child profiles (a Performance parent over regional Weighted children).
- **Routing inputs live with the destination**: weight, priority, geography and subnet claims, and per-endpoint probe headers ride the endpoint they describe.

## Key Features

- Full azurerm v5 surface across all three endpoint resources, folded into one union spec (shared routing fields at the root, per-type arguments in the variant blocks -- each mirroring its provider resource exactly).
- Typed references: the owning profile and nested child profiles wire by output; the azure variant's target takes any Azure resource's id output.
- Honest defaults: weight defaults to 1 (always sent); priority is left to Azure's creation-order assignment unless set -- exactly the provider's behavior.

## Use Cases

- **Regional deployments behind one name**: one azure endpoint per region's public IP or App Service, weighted or latency-routed by the profile.
- **Hybrid and multi-cloud fronts**: external endpoints pointing at on-premises datacenters or other clouds during migrations.
- **Routing trees**: nested endpoints composing a global Performance parent over per-region Weighted children for two-level traffic policy.

## Future Enhancements

- Geographic-hierarchy code validation (WORLD, GEO-EU, country codes) stays apply-time -- Azure validates claims against the live hierarchy, which cannot be introspected offline.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
