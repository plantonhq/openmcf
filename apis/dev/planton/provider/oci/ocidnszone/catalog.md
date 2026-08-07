# DNS Zone on OCI

Deploys a managed authoritative DNS zone on Oracle Cloud Infrastructure supporting public (global) and private (VCN-scoped) resolution, primary and secondary zone types, zone transfers via external masters and downstreams, and DNSSEC signing. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Zone** -- a managed authoritative zone in the specified compartment. The zone name is derived from `metadata.name`. Supports global (publicly resolvable) and private (VCN-only) scopes, primary and secondary zone types, optional DNSSEC signing, and optional zone transfer configuration with external DNS servers.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the zone for tracking and governance

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the DNS zone in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A DNS view OCID -- required only for private-scoped zones. OCI DNS views are not modeled as Planton deployment components; provide the OCID as a literal value or via ValueFromRef from an external source.
- External master DNS server addresses -- required only for secondary zones that replicate from on-premises or third-party DNS servers.
- TSIG key OCIDs (optional) -- for authenticating zone transfers between this zone and external DNS servers.

## Deploy

### Console

Open the deployment store, find **DNS Zone on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public Primary** preset in the [Presets](#presets) tab to pre-populate a standard public DNS zone.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDnsZone
metadata:
  name: example.com
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  zoneType: primary
```

```shell
planton apply -f dns-zone.yaml
```

This creates a public primary DNS zone for `example.com`. No DNSSEC signing or zone transfers are configured. OCI assigns authoritative nameservers that you configure as NS records at your domain registrar.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DNS zone to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the DNS zone with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone type** -- Set `zoneType` to `primary` for zones where OCI is the authoritative source of truth (you manage records via OciDnsRecord). Set to `secondary` for zones that replicate from external master DNS servers. Secondary zones require at least one `externalMasters` entry. Zone type is immutable after creation.

**Resolution scope** -- Defaults to global (publicly resolvable). Set `scope` to `scope_private` for zones resolvable only within VCNs via a DNS view. Private zones require `viewId`. OCI does not support private secondary zones. Scope is immutable after creation.

**DNSSEC** -- Set `isDnssecEnabled: true` to enable DNSSEC signing. OCI generates KSK and ZSK key pairs and signs zone records automatically. Only meaningful for global-scoped zones. Can be toggled after creation.

**Zone transfers** -- Configure `externalDownstreams` on primary zones to push zone transfers to external DNS servers (hybrid DNS). Configure `externalMasters` on secondary zones to pull from on-premises or third-party DNS. Use `tsigKeyId` on external server entries for TSIG-authenticated transfers.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zoneId` | OCID of the DNS zone | OciDnsRecord zone reference, IAM policy scoping |
| `nameservers` | Comma-separated list of OCI-assigned authoritative nameserver hostnames | NS record configuration at domain registrar |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public primary zone** -- A standard globally resolvable DNS zone where OCI is the authoritative source. The most common pattern for hosting domain DNS on OCI. Start from the **Public Primary** preset.

**Private VCN zone** -- A private DNS zone resolvable only within VCNs via a DNS view. Used for internal service discovery and split-horizon DNS where internal names should not be resolvable from the public internet. Start from the **Private VCN** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DNS zone