# AzureFirewall

## Overview

`AzureFirewall` provisions an Azure Firewall instance: the managed,
stateful network firewall data plane that enforces an
`AzureFirewallPolicy`. The firewall carries WHERE enforcement runs -- the
dedicated subnet, public IPs, availability zones, and deployment model --
while the attached policy carries WHAT is enforced (rules, threat
intelligence, TLS inspection, IDPS).

The classic hub-spoke shape: the firewall sits in the hub VNet's
`AzureFirewallSubnet`, and every spoke's route table sends `0.0.0.0/0` to
the firewall's private IP (a `VIRTUAL_APPLIANCE` next hop on
`AzureRouteTable`).

## Why a First-Class Resource?

- **The data plane is regional infrastructure** -- it deploys per hub, is
  zoned, and provisions in tens of minutes; the policy it enforces is a
  fast control-plane document authored elsewhere
- **The hub-spoke seam** -- `private_ip_address` is the output every spoke
  route table consumes
- **Composable addressing** -- subnets and public IPs are referenced
  `AzureSubnet`/`AzurePublicIp` resources, never bundled

## Deployment Models

| Model | Selected by | Addressing |
|-------|-------------|------------|
| VNet (default) | `sku_name: AZFW_VNET` | Your subnet (named exactly `AzureFirewallSubnet`, /26+) + your Standard static public IPs |
| Virtual WAN hub | `sku_name: AZFW_HUB` | `virtual_hub` block; Azure assigns the IPs (surfaced as outputs) |

The pairing is validated at authoring time: a VNet firewall carries
`ip_configurations`, a hub firewall carries `virtual_hub` -- never both.

## Spec Fields (top level)

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region; ForceNew |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group (defaults to AzureResourceGroup reference) |
| `name` | string | Yes | -- | Firewall name, 1-80 chars; ForceNew |
| `sku_name` | enum | No | `AZFW_VNET` | Deployment model; ForceNew |
| `sku_tier` | enum | No | `STANDARD` | BASIC / STANDARD / PREMIUM; must match the policy's tier |
| `ip_configurations` | list | VNet model | -- | Exactly one carries `subnet_id`; extras add public IPs |
| `management_ip_configuration` | message | BASIC tier / forced tunneling | -- | Dedicated `AzureFirewallManagementSubnet` + required public IP; ForceNew |
| `firewall_policy_id` | StringValueOrRef | No | -- | The policy enforced (defaults to AzureFirewallPolicy reference) |
| `threat_intel_mode` | enum | No | Azure's Alert | Only meaningful without a policy |
| `dns_servers` / `dns_proxy_enabled` | list / bool | No | -- | Policy-less firewalls only (a policy owns DNS -- validated); servers implicitly force the proxy ON |
| `private_ip_ranges` | list | No | IANA private | CIDRs or the literal `IANAPrivateRanges` |
| `virtual_hub` | message | Hub model | -- | Hub id + Azure-assigned public IP count |
| `zones` | list | No | -- | Availability zones; ForceNew; free on Azure Firewall |
| `tags` | map | No | -- | User tags (merged over Planton-derived tags) |

## Outputs

| Output | Description |
|--------|-------------|
| `firewall_id` | Full ARM ID of the firewall |
| `firewall_name` | The firewall's name as deployed |
| `private_ip_address` | **The hub-spoke seam** -- spoke route tables' `next_hop_in_ip_address` |
| `management_private_ip_address` | The management path's private IP (when configured) |
| `virtual_hub_public_ip_addresses` | Azure-assigned public IPs (hub model) |
| `virtual_hub_private_ip_address` | Azure-assigned private IP (hub model) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFirewall
metadata:
  name: hub-egress-fw
  org: mycompany
  env: production
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

Steer a spoke through it:

```yaml
# AzureRouteTable
spec:
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress:
        valueFrom:
          name: <the firewall resource's name>
```

## Lifecycle Notes

- **Slow resource**: provisioning AND deletion each run 10-20+ minutes;
  the ForceNew surface (`name`, `region`, `sku_name`, `zones`,
  `management_ip_configuration`, each configuration's `subnet_id`) is
  expensive -- design to avoid replacement
- **Subnet contracts are ARM's**: the data subnet must be named exactly
  `AzureFirewallSubnet` (management: `AzureFirewallManagementSubnet`),
  each at least /26, carrying no other workloads
- **Public IP rule**: required on the data path unless a
  `management_ip_configuration` provides the management path (validated
  at authoring time); Azure requires the public IPs to be Standard SKU
  with Static allocation
- **BASIC tier** requires the management configuration and pairs only
  with a BASIC policy
- **Classic inline rule collections are deliberately not modeled** --
  policy-based management is Azure's direction; ARM rejects mixing them

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
