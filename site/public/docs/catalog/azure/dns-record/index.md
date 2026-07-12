---
title: "DNS Record"
description: "DNS Record deployment documentation"
icon: "package"
order: 100
componentName: "azurednsrecord"
---

# Azure DNS Record

Declares one record set in an Azure public DNS zone -- A, AAAA, CNAME, MX, SRV, CAA, TXT, NS, or PTR, with the record type determined by which typed payload the spec carries.

## What Gets Created

When you deploy an AzureDnsRecord resource, Planton provisions:

- **One DNS record set** -- the matching `azurerm_dns_*_record` in the referenced zone, carrying every value the zone answers for that (name, type) pair

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureDnsZone** (or the name of a zone managed outside Planton) to create the record in

## Quick Start

Create a file `record.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureDnsRecord
metadata:
  name: www-a-record
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureDnsRecord.www-a-record
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  zoneName:
    valueFrom:
      kind: AzureDnsZone
      name: example-com
      fieldPath: status.outputs.zone_name
  name: www
  ttlSeconds: 300
  a:
    addresses:
      - 203.0.113.10
```

Deploy:

```shell
planton apply -f record.yaml
```

Set exactly one payload -- it determines the record type. An MX record set carries preference+exchange pairs; an alias A record replaces `addresses` with a `targetResourceId` reference to an Azure resource (a Public IP, Traffic Manager profile, or CDN/Front Door endpoint) whose address Azure then tracks automatically:

```yaml
  a:
    targetResourceId:
      valueFrom:
        kind: AzurePublicIp
        name: frontend-ip
        fieldPath: status.outputs.public_ip_id
```

## Key Outputs

| Output | Purpose |
|--------|---------|
| `record_id` | The record set's ARM ID |
| `fqdn` | The record's fully qualified name (trailing dot) |

## Related Resources

- [Azure DNS Zone](/docs/catalog/azure/dns-zone) -- the zone the record lives in
- [Azure Public IP](/docs/catalog/azure/public-ip) -- a common alias-record target
- [Azure Container App Custom Domain](/docs/catalog/azure/container-app-custom-domain) -- consumes the `asuid.` TXT + CNAME records this kind declares
