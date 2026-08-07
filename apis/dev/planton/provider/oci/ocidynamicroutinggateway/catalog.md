# Dynamic Routing Gateway on OCI

Deploys an Oracle Cloud Infrastructure Dynamic Routing Gateway (DRG) with bundled VCN attachments, custom route tables, route distributions, distribution statements, and static route rules. A DRG is a virtual router that enables connectivity between a VCN and networks outside its region -- other VCNs (peering), on-premises networks (Site-to-Site VPN via IPSec, FastConnect), and cross-region VCNs (remote peering). Sub-resources reference each other by `displayName` rather than OCID, keeping the YAML self-contained. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DRG** -- the virtual router in the specified compartment
- **DRG Attachments** -- one per entry in `attachments`, connecting VCNs, IPSec tunnels, virtual circuits, remote peering connections, or loopbacks to the DRG
- **DRG Route Tables** -- one per entry in `routeTables`, with optional import distribution references and ECMP settings. OCI also auto-creates default route tables per network type.
- **Static Route Rules** -- one per entry in each route table's `staticRouteRules`, directing traffic for a CIDR to a specific DRG attachment
- **DRG Route Distributions** -- one per entry in `routeDistributions`, controlling which routes are advertised to route tables (import) or to attachments (export)
- **Distribution Statements** -- one per entry in each distribution's `statements`, with priority-based match criteria (match all, attachment type, or specific attachment)
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the DRG, attachments, route tables, and distributions

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the DRG in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- One or more VCNs to attach (for peering). Each VCN can be attached to only one DRG at a time. Provide VCN OCIDs directly or through separate OciVcn Cloud Resources.
- For on-premises connectivity: IPSec connections or FastConnect virtual circuits configured in OCI. These are attached to the DRG as additional attachment types.

## Deploy

### Console

Open the deployment store, find **Dynamic Routing Gateway on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single VCN Attachment** preset in the [Presets](#presets) tab to pre-populate a DRG with one VCN attachment.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDynamicRoutingGateway
metadata:
  name: hub-drg
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  displayName: hub-drg
  attachments:
    - displayName: app-vcn
      networkDetails:
        type: vcn
        id:
          value: "ocid1.vcn.oc1..example"
```

```shell
planton apply -f drg.yaml
```

This creates a DRG with a single VCN attachment using OCI's auto-generated default route tables. No custom route tables, distributions, or static rules are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DRG to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the DRG with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a DRG. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Attachment types** -- Each attachment connects a network to the DRG via `networkDetails.type`: `vcn` (VCN peering), `ipsec_tunnel` (Site-to-Site VPN), `virtual_circuit` (FastConnect), `remote_peering_connection` (cross-region), or `loopback`. The `networkDetails.id` field references the network resource OCID.

**Custom route tables** -- Define entries in `routeTables` to control traffic forwarding between attachments. Reference an import distribution via `importDrgRouteDistributionName` to automatically populate routes from matching attachments. Add `staticRouteRules` for explicit CIDR-to-attachment mappings. Enable `isEcmpEnabled` for equal-cost multi-path routing across redundant IPSec tunnels or virtual circuits.

**Route distributions** -- Define entries in `routeDistributions` to control which routes are advertised. Each distribution contains prioritized `statements` with match criteria (`matchAll`, `drgAttachmentType`, or `drgAttachmentId`). Distributions are referenced by name from route tables (import) and attachments (export).

**Hub-and-spoke routing** -- For multi-VCN environments, create a shared custom route table with an import distribution that matches all VCN attachments. Assign each spoke VCN attachment to this shared route table so spokes automatically learn each other's CIDRs.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `drg_id` | OCID of the DRG | OciSubnet route rules targeting the DRG, network auditing |
| `default_export_drg_route_distribution_id` | OCID of the auto-created default export distribution | Configuring external DRG attachments (IPSec, FastConnect) managed outside this component |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single VCN attachment** -- A DRG with one VCN attachment using OCI's auto-generated default route tables. The starting point for IPSec VPN, FastConnect, or future VCN peering. Start from the **Single VCN Attachment** preset.

**Hub-and-spoke** -- A DRG configured as a multi-VCN hub with two spoke VCN attachments sharing a custom route table that imports routes from all VCN attachments via an import distribution, enabling spoke-to-spoke communication. Start from the **Hub-and-Spoke** preset.

## Works With

- [**OCI Compartment**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DRG