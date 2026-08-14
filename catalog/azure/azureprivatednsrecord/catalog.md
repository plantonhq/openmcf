# Azure Private DNS Record

Deploys one DNS record set (A, AAAA, CNAME, MX, PTR, SRV, or TXT) in an Azure Private DNS zone -- name resolution visible only inside the virtual networks linked to the zone. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **One record set** of the type your spec's payload declares -- all values for the (name, type) pair in one resource

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzurePrivateDnsZone** -- the record is created inside a referenced zone (the preset wires it by reference, which also orders the deploy; pass a literal zone ARM id for a zone managed outside Planton).

### Azure Subscription

- **One record set per (name, type)** -- a second record with the same name and type in the same zone conflicts rather than merges; declare all of a name's values in one resource.
- **Resolution requires zone links** -- records answer only inside virtual networks linked to the zone (AzurePrivateDnsZoneVirtualNetworkLink); an unlinked zone's records are inert.
- **Free at rest** -- private DNS bills per zone and per million queries, never per record.

## Deploy

### Console

Open the deployment store, find **Azure Private DNS Record**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f record.yaml
```

## After Deploy

The `fqdn` output carries the record's full name (with DNS's trailing dot); resolution works from any network linked to the zone -- test with `nslookup <fqdn>` from a VM in a linked network.
