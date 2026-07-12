# Azure DNS Zone

Creates an Azure public DNS zone -- an internet-facing, authoritative DNS zone hosted on Azure's global anycast name-server fleet.

## What Gets Created

When you deploy an AzureDnsZone resource, Planton provisions:

- **Public DNS zone** -- an `azurerm_dns_zone` in the referenced resource group, with optional Start of Authority customization and governance tags

Records are separate resources: declare them with [Azure DNS Record](/docs/catalog/azure/azurednsrecord), one per record set, referencing this zone.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** for the zone's ARM lifecycle (referenced through `resourceGroup`)
- **A registered domain** (at any registrar) if the zone should answer the internet -- creating the zone alone does not make it authoritative

## Quick Start

Create a file `zone.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureDnsZone
metadata:
  name: example-com
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureDnsZone.example-com
spec:
  zoneName: example.com
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  tags:
    team: platform
```

Deploy:

```shell
planton apply -f zone.yaml
```

After deployment, update your domain registrar to use the four name servers in `status.outputs.name_servers` -- that delegation is what makes the zone authoritative. For a subdomain zone (team.example.com), publish those name servers as an NS record in the parent zone instead.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `zone_id` | The zone's ARM ID -- referenced by Front Door custom domains and AKS web-app routing |
| `zone_name` | What Azure DNS Record resources reference to address record sets |
| `resource_group_name` | Pairs with `zone_name` for management-plane addressing |
| `name_servers` | The four name servers to configure at the registrar |
| `max_number_of_record_sets` | The zone's record-set capacity limit |

## Related Resources

- [Azure DNS Record](/docs/catalog/azure/azurednsrecord) -- the record sets declared in the zone
- [Azure Private DNS Zone](/docs/catalog/azure/azureprivatednszone) -- name resolution inside virtual networks
- [Azure Front Door Custom Domain](/docs/catalog/azure/azurefrontdoorcustomdomain) -- watches the zone for its validation records
