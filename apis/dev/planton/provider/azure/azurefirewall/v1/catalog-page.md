# Azure Firewall

Creates an Azure Firewall -- the managed, stateful firewall data plane that enforces an Azure Firewall Policy from a dedicated subnet (or Virtual WAN hub). Spoke route tables send egress to its private IP; the attached policy decides what gets through.

## What Gets Created

When you deploy an AzureFirewall resource, Planton provisions:

- **Azure Firewall** — an `azurerm_firewall` in the specified region and resource group, deployed into the referenced `AzureFirewallSubnet` (or Virtual WAN hub), fronted by the referenced Standard static public IPs, enforcing the referenced policy

Rules and inspection live on the `AzureFirewallPolicy`; the firewall is the enforcement point.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A subnet named exactly `AzureFirewallSubnet`** (at least /26) in the hub VNet — an `AzureSubnet`
- **A Standard-SKU static public IP** — an `AzurePublicIp`
- **A firewall policy** to enforce — an `AzureFirewallPolicy` (same tier)
- **Network write rights**: `Microsoft.Network/azureFirewalls/write`

## Quick Start

Create a file `firewall.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewall
metadata:
  name: hub-egress-fw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFirewall.hub-egress-fw
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: hub-egress-fw
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: firewall-subnet
      publicIpAddressId:
        valueFrom:
          name: firewall-pip
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  zones: ["1", "2", "3"]
```

Deploy (expect 10-20 minutes):

```shell
planton apply -f firewall.yaml
```

After deployment, read `status.outputs.private_ip_address` — the address spoke route tables send egress to.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region; must match the VNet/hub. | Required, ForceNew |
| `resourceGroup` | `StringValueOrRef` | Resource group name. | Required |
| `name` | `string` | Firewall name, unique within the resource group. | Required, 1-80 chars, ForceNew |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `skuName` | `enum` | `AZFW_VNET` | `AZFW_VNET` (your subnet) or `AZFW_HUB` (Virtual WAN). ForceNew. |
| `skuTier` | `enum` | `STANDARD` | `BASIC`, `STANDARD`, `PREMIUM`. Must match the policy's tier. |
| `ipConfigurations` | `list` | -- | Exactly one carries `subnetId` (the AzureFirewallSubnet); extra entries add public IPs (more SNAT ports, more DNAT frontends). |
| `managementIpConfiguration` | `object` | -- | Dedicated `AzureFirewallManagementSubnet` + required public IP. Required for BASIC tier and forced tunneling. ForceNew. |
| `firewallPolicyId` | `StringValueOrRef` | -- | The `AzureFirewallPolicy` to enforce. |
| `threatIntelMode` | `enum` | Azure's `Alert` | Only meaningful without a policy. |
| `dnsServers` | `list(string)` | -- | Custom upstreams; setting them implicitly enables the DNS proxy. Policy-less firewalls only -- a policy-attached firewall takes DNS from the policy (validated). |
| `dnsProxyEnabled` | `bool` | `false` | Run the firewall as a DNS proxy (needed for FQDN network rules). Policy-less firewalls only. |
| `privateIpRanges` | `list(string)` | IANA private | SNAT-exempt ranges; CIDRs or the literal `IANAPrivateRanges`. |
| `virtualHub` | `object` | -- | `AZFW_HUB` model: hub ARM id + Azure-assigned public IP count. |
| `zones` | `list(string)` | -- | Availability zones (free on Azure Firewall). ForceNew. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags. |

## Examples

### Hub-Spoke Egress Control (the flagship shape)

The firewall in the hub, spokes default-routing through it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: spoke-routes
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: spoke-routes
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      # Resolves to the firewall's private_ip_address output -- the route
      # follows the firewall instead of pinning a hand-copied IP.
      nextHopInIpAddress:
        valueFrom:
          name: hub-firewall
  bgpRoutePropagationEnabled: false
```

### Forced Tunneling (private data path)

```yaml
spec:
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: firewall-subnet
  managementIpConfiguration:
    name: management
    subnetId:
      valueFrom:
        name: firewall-mgmt-subnet
    publicIpAddressId:
      valueFrom:
        name: firewall-mgmt-pip
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `firewall_id` | `string` | Full ARM ID of the firewall |
| `firewall_name` | `string` | The firewall's name as deployed |
| `private_ip_address` | `string` | The data-path private IP — spoke route tables' next hop |
| `management_private_ip_address` | `string` | The management path's private IP (when configured) |
| `virtual_hub_public_ip_addresses` | `list` | Azure-assigned public IPs (hub model) |
| `virtual_hub_private_ip_address` | `string` | Azure-assigned private IP (hub model) |

## Related Components

- [AzureFirewallPolicy](/docs/catalog/azure/firewall-policy) — WHAT the firewall enforces
- [AzureFirewallPolicyRuleCollectionGroup](/docs/catalog/azure/firewall-policy-rule-collection-group) — the rules
- [AzureSubnet](/docs/catalog/azure/subnet) — the `AzureFirewallSubnet` home
- [AzurePublicIp](/docs/catalog/azure/public-ip) — the firewall's frontends
- [AzureRouteTable](/docs/catalog/azure/route-table) — steers spokes through the firewall's private IP
