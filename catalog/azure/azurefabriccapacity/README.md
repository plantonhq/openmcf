# Overview

The **AzureFabricCapacity** component deploys a Microsoft Fabric capacity -- the billing and compute anchor of Microsoft Fabric. Workspaces assign themselves to a capacity, and the capacity's F-SKU sets how much compute every workload on it (lakehouses, warehouses, Power BI, real-time analytics) shares.

## Purpose

- **The one ARM-side Fabric resource**: azurerm's entire Fabric surface is this capacity; workspaces and the items inside them are managed through Microsoft's dedicated `fabric` Terraform provider, the Fabric portal, or its APIs.
- **Compute as a dial**: the F-SKU (F2 through F2048, doubling per step) scales up and down in place -- start small, grow with real usage.
- **Administration is explicit**: the capacity's administrators (Entra users or service principals) are declared on the resource, at least one at all times.

## Key Features

- Full azurerm v5 surface: the F-SKU vocabulary, administration members, and tags. The SKU tier has exactly one legal value ("Fabric") and is deliberately not part of the spec -- both engines send it explicitly.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup; the `fabric_capacity_id` and `fabric_capacity_name` outputs are what Fabric-side automation assigns workspaces against.
- Create and delete are polled long-running operations -- the module waits for completion.

## Use Cases

- **The analytics platform anchor**: one capacity per environment; teams assign Fabric workspaces to it from the Fabric portal.
- **Cost-controlled BI**: run an F2 for development, scale the production capacity with demand -- the SKU moves in place.
- **Copilot-enabled analytics**: F64 and above carry Copilot and unlimited Power BI sharing.

## Future Enhancements

- Fabric workspaces and items live in Microsoft's dedicated `fabric` Terraform provider -- a future per-provider catalog decision, deliberately outside this kind.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
