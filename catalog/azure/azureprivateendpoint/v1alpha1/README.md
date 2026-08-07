# AzurePrivateEndpoint

## Overview

`AzurePrivateEndpoint` provisions an Azure Private Endpoint: a network
interface that gives a Private Link-enabled service a private IP address
inside your virtual network. The target service -- an Azure PaaS resource
(SQL, PostgreSQL, Storage, Key Vault, Cosmos DB, ...) or a custom Private
Link Service -- becomes reachable over a private IP on the Microsoft
backbone, never the public internet.

## Why It Matters

- **Private connectivity** -- traffic to the service stays on the Microsoft
  backbone; no public exposure, no service endpoints, no NAT
- **Data-exfiltration protection** -- each endpoint maps to one sub-resource
  ("blob", "vault", "postgresqlServer"), not the whole service, so a
  compromised client cannot pivot to other data
- **Correct DNS, atomically** -- the private DNS zone group is part of this
  resource; without it the service FQDN resolves to the PUBLIC IP, silently
  defeating the private link

## Key Features

- **Polymorphic target** -- `private_service_connection.private_connection_resource_id`
  references any Private Link-capable service, or `connection_alias`
  connects through a partner's Private Link Service alias
- **Private DNS registration** -- `private_dns_zone_ids` registers the
  endpoint IP as an A record in each zone (typically the service's
  `privatelink.*` zone)
- **Manual approval flow** -- `is_manual_connection` + `request_message` for
  cross-tenant/cross-subscription connections that need owner approval
- **Static IPs** -- `ip_configurations` pins sub-resources to fixed
  addresses when firewall allowlists or hard-coded DNS require it
- **ASG membership** -- `application_security_group_ids` puts the endpoint's
  NIC in workload groups so NSG rules govern traffic by group
- **Composable** -- subnet, resource group, DNS zones, and ASGs are all
  references that default to their respective Planton kinds

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | Azure region (must match the subnet); fixed at creation |
| `resource_group` | StringValueOrRef | Yes | Resource group (defaults to AzureResourceGroup) |
| `name` | string | Yes | Endpoint name, unique in the resource group; fixed at creation |
| `subnet_id` | StringValueOrRef | Yes | Subnet the IP is drawn from (defaults to AzureSubnet); fixed at creation |
| `private_service_connection` | message | Yes | The private link connection (target + sub-resource) |
| `private_dns_zone_ids` | repeated StringValueOrRef | No | Private DNS zones to register into (defaults to AzurePrivateDnsZone) |
| `ip_configurations` | repeated message | No | Static IP assignments; empty = dynamic allocation |
| `application_security_group_ids` | repeated StringValueOrRef | No | ASGs the NIC joins (defaults to AzureApplicationSecurityGroup) |
| `custom_network_interface_name` | string | No | Custom NIC name; fixed at creation |
| `tags` | map | No | User tags, merged over Planton-derived tags |

### `private_service_connection`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `private_connection_resource_id` | StringValueOrRef | one-of | Target service by ARM ID (polymorphic) |
| `connection_alias` | string | one-of | Private Link Service alias (partner cross-tenant) |
| `subresource_names` | repeated string | No | Sub-resource / group ID (e.g. "postgresqlServer") |
| `is_manual_connection` | bool | No | Requires owner approval (default false) |
| `request_message` | string | No | Approval message (required iff manual, 1-140 chars) |

## Outputs

| Output | Description |
|--------|-------------|
| `private_endpoint_id` | Full ARM ID of the endpoint |
| `private_endpoint_name` | The endpoint's name as deployed |
| `private_ip_address` | The private IP allocated from the subnet |
| `network_interface_id` | ARM ID of the auto-created NIC |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateEndpoint
metadata:
  name: pg-private-endpoint
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: pg-private-endpoint
  subnetId:
    valueFrom:
      name: pe-subnet
  privateServiceConnection:
    privateConnectionResourceId:
      valueFrom:
        kind: AzurePostgresqlFlexibleServer
        name: prod-postgres
        fieldPath: status.outputs.server_id
    subresourceNames:
      - postgresqlServer
  privateDnsZoneIds:
    - valueFrom:
        name: postgres-privatelink-zone
```

## Lifecycle Notes

- `region`, `name`, `subnet_id`, the connection block, `ip_configurations`,
  and `custom_network_interface_name` are **fixed at creation**
- `private_dns_zone_ids`, `application_security_group_ids`, and `tags` update
  in place
- A manual connection stays in `Pending` until the target owner approves it;
  the endpoint provisions but does not carry traffic until then

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
