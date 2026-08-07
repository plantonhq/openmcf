# Azure Virtual Network

Creates an Azure Virtual Network (VNet) -- the isolated, private IP address space every network-attached Azure workload lives inside -- with the full network surface: multi-CIDR and dual-stack address spaces, Azure Network Manager IPAM delegation, custom DNS, BGP community advertisement, DDoS Protection Plan attachment, VM-to-VM encryption, and flow-timeout tuning.

## What Gets Created

When you deploy an AzureVirtualNetwork resource, Planton provisions:

- **Virtual Network** — an `azurerm_virtual_network` carrying the address space and network-wide policy

The network is deliberately just the network. Subnets (`AzureSubnet`), outbound NAT (`AzureNatGateway`), and private DNS attachments (`AzurePrivateDnsZoneVirtualNetworkLink`) are their own composable resources referencing this network's outputs -- so topologies grow without touching the network itself.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the network in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/virtualNetworks/write` (Network Contributor, Contributor, or Owner)

## Quick Start

Create a file `network.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetwork
metadata:
  name: prod-network
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetwork.prod-network
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: prod-vnet
  addressSpaces:
    - "10.0.0.0/16"
```

Deploy:

```shell
planton apply -f network.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the network. A regional resource; changing it replaces the network. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Network name, unique within the resource group. | Required, 2-64 chars, Azure naming rules |

### Address Source (exactly one)

| Field | Type | Description |
|-------|------|-------------|
| `addressSpaces` | `list(string)` | Self-managed CIDR blocks. Multiple blocks are first-class -- grow the space or run dual-stack. |
| `ipAddressPools` | `list(object)` | Delegate allocation to Azure Network Manager IPAM pools (max 2: one per IP version). Each entry: `id` (pool ARM ID) + `numberOfIpAddresses` (positive number as string; can grow, never shrink). |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `dnsServers` | `list(string)` | Custom DNS server IPs. Empty means Azure's default resolver (168.63.129.16) -- required for private DNS zone resolution. VMs pick up changes on DHCP lease renewal. |
| `bgpCommunity` | `string` | BGP community advertised with the network's routes over ExpressRoute, `asn:community` notation (ASN segment is always 12076 today). |
| `ddosProtectionPlan` | `object` | Attach an existing DDoS Protection Plan: `id` (plan ARM ID) + `enable`. The plan is a separate, billed, shareable resource. |
| `encryption` | `enum` | VM-to-VM encryption enforcement: `ALLOW_UNENCRYPTED` or `DROP_UNENCRYPTED`. Unset = encryption off. Note ARM currently accepts only `ALLOW_UNENCRYPTED`. |
| `flowTimeoutInMinutes` | `int` | Connection-tracking timeout for intra-network flows, 4-30. Unset = Azure's 4-minute default. |
| `privateEndpointVnetPolicies` | `enum` | Network-wide private endpoint policy: `BASIC`. Unset = ARM's default (`Disabled`). |
| `edgeZone` | `string` | Deploy into an Azure Edge Zone (metro-local extension). Changing it replaces the network. |
| `tags` | `map(string)` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Growing Network with Custom DNS

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetwork
metadata:
  name: hub-network
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetwork.hub-network
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: hub-vnet
  addressSpaces:
    - "10.0.0.0/16"
    - "10.1.0.0/16"
  dnsServers:
    - "10.0.0.4"
    - "10.0.0.5"
  flowTimeoutInMinutes: 15
```

### DDoS-Protected Public-Facing Network

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetwork
metadata:
  name: edge-network
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetwork.edge-network
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: edge-rg
  name: edge-vnet
  addressSpaces:
    - "10.50.0.0/16"
  ddosProtectionPlan:
    id: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/security-rg/providers/Microsoft.Network/ddosProtectionPlans/org-ddos-plan
    enable: true
```

### Dual-Stack Network

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetwork
metadata:
  name: dualstack-network
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetwork.dualstack-network
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: dualstack-vnet
  addressSpaces:
    - "10.2.0.0/16"
    - "fd00:db8:deca::/48"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `virtual_network_id` | `string` | Full ARM ID of the network -- the join key for subnets, peerings, and DNS links |
| `virtual_network_name` | `string` | The network's name as deployed |
| `guid` | `string` | ARM's stable network GUID |
| `address_spaces` | `list(string)` | The ACTUAL ranges carried (IPAM-provisioned when pools delegate allocation) |

## Related Components

- [AzureSubnet](/docs/catalog/azure/subnet) — partitions the network's address space into workload segments
- [AzureNatGateway](/docs/catalog/azure/nat-gateway) — managed outbound connectivity for a subnet
- [AzurePrivateDnsZoneVirtualNetworkLink](/docs/catalog/azure/private-dns-zone-virtual-network-link) — makes a private DNS zone resolvable from this network
- [AzureNetworkSecurityGroup](/docs/catalog/azure/network-security-group) — traffic filtering for subnets and NICs
