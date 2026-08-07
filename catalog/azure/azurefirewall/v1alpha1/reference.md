# AzureFirewall

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFirewallSpec** defines an Azure Firewall instance: the managed,
stateful network firewall data plane that enforces an
AzureFirewallPolicy. The firewall carries WHERE enforcement runs -- the
dedicated subnet, public IPs, availability zones, and deployment model
-- while the attached policy carries WHAT is enforced (rules, threat
intelligence, TLS inspection, IDPS).

**Two deployment models**, selected by sku_name:
- **AZFW_VNET** (the default): the firewall deploys into a dedicated
  subnet of your virtual network. The subnet MUST be named exactly
  "AzureFirewallSubnet" with at least a /26 prefix -- ARM rejects any
  other name or a smaller subnet. Hub-spoke egress control is built on
  this shape: spoke route tables send 0.0.0.0/0 to the firewall's
  private IP (a VIRTUAL_APPLIANCE next hop).
- **AZFW_HUB**: the firewall deploys into a Virtual WAN hub
  (virtual_hub block); Azure manages its addressing and publishes the
  assigned IPs as outputs.

**Forced tunneling / management plane**: management_ip_configuration
gives the firewall a separate management path (subnet named exactly
"AzureFirewallManagementSubnet", /26 or larger, with its own public
IP). With it in place, the data-path ip_configuration may be fully
private. Azure requires it for BASIC-tier firewalls and for forced
tunneling (sending 0.0.0.0/0 on-premises). The block is fixed at
creation -- adding or removing it replaces the firewall.

**Tiers must match**: sku_tier must equal the attached policy's sku.

**ForceNew fields**: `name`, `region`, `resource_group`, `sku_name`,
`zones`, `management_ip_configuration`, and each ip_configuration's
`subnet_id`. Azure Firewall provisions and deletes SLOWLY (10-20+
minutes each way) -- design changes to avoid replacement.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFirewall
metadata:
  name: test-firewall
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hub-egress-fw
  skuName: AZFW_VNET
  skuTier: STANDARD
  ipConfigurations:
    - name: primary
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/AzureFirewallSubnet
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/fw-pip
    - name: extra-pip
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/fw-pip-2
  # No firewall-level DNS here: a policy-attached firewall takes DNS from
  # the policy (ARM rejects firewall-level DNS parameters otherwise).
  firewallPolicyId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/firewallPolicies/egress-baseline
  privateIpRanges:
    - "IANAPrivateRanges"
    - "100.64.0.0/10"
  zones:
    - "1"
    - "2"
    - "3"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.skuName` | `enum` |  |  |  |
| `spec.skuTier` | `enum` |  |  |  |
| `spec.ipConfigurations` | `[]AzureFirewallIpConfiguration` |  |  |  |
| `spec.ipConfigurations[].name` | `string` | yes |  |  |
| `spec.ipConfigurations[].subnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.ipConfigurations[].publicIpAddressId` | `string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.managementIpConfiguration` | `AzureFirewallManagementIpConfiguration` |  |  |  |
| `spec.managementIpConfiguration.name` | `string` | yes |  |  |
| `spec.managementIpConfiguration.subnetId` | `string \| valueFrom` | yes |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.managementIpConfiguration.publicIpAddressId` | `string \| valueFrom` | yes |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.firewallPolicyId` | `string \| valueFrom` |  |  | AzureFirewallPolicy (`status.outputs.firewall_policy_id`) |
| `spec.threatIntelMode` | `enum` |  |  |  |
| `spec.dnsServers` | `[]string` |  |  |  |
| `spec.dnsProxyEnabled` | `bool` |  |  |  |
| `spec.privateIpRanges` | `[]string` |  |  |  |
| `spec.virtualHub` | `AzureFirewallVirtualHub` |  |  |  |
| `spec.virtualHub.virtualHubId` | `string \| valueFrom` | yes |  |  |
| `spec.virtualHub.publicIpCount` | `int32` |  | `1` |  |
| `spec.zones` | `[]string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the firewall lives in, e.g. "eastus". Must match the
virtual network (or virtual hub) it deploys into. Changing the region
replaces the firewall.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the firewall is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The firewall's name, unique within the resource group. 1-80
characters; must begin with a letter or number, end with a letter,
number, or underscore, and may contain only letters, numbers,
underscores, periods, or hyphens. Changing the name replaces the
firewall.

- rule: Firewall names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.skuName

`enum`

The deployment model: AZFW_VNET (dedicated subnet in your VNet -- the
default) or AZFW_HUB (Virtual WAN hub). Fixed at creation.
Unspecified deploys AZFW_VNET, the standard hub-spoke shape.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_sku_name_unspecified` -- Not specified -- deploys AZFW_VNET, the standard hub-spoke shape.
- `AZFW_VNET` -- Deploy into a dedicated "AzureFirewallSubnet" subnet of your virtual network. The hub-spoke egress-control shape.
- `AZFW_HUB` -- Deploy into a Virtual WAN hub (secured virtual hub). Azure manages the firewall's addressing.

### spec.skuTier

`enum`

The firewall's capability tier: BASIC (SMB-scale; requires a
management IP configuration), STANDARD (the production default), or
PREMIUM (TLS inspection + IDPS, paired with a Premium policy). Must
match the attached policy's tier. Unspecified deploys STANDARD.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_sku_tier_unspecified` -- Not specified -- deploys STANDARD, the production default.
- `BASIC` -- SMB-scale tier. Requires a management_ip_configuration and pairs only with a BASIC policy.
- `STANDARD` -- The production default: full rule engine and threat intelligence.
- `PREMIUM` -- Adds TLS inspection and IDPS -- pairs with a PREMIUM policy.

### spec.ipConfigurations

`[]AzureFirewallIpConfiguration`

The data-path IP configurations (AZFW_VNET model). EXACTLY ONE
configuration carries subnet_id (the firewall's one subnet); add
more configurations -- each with its own public IP, no subnet -- to
scale the firewall's public IP set (each extra public IP adds SNAT
ports and DNAT frontends). Azure requires at least one public IP on
the data path unless a management_ip_configuration provides the
management path.

### spec.ipConfigurations[].name

`string` · required

The configuration's name, unique on the firewall (and distinct from
the management configuration's name).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.ipConfigurations[].subnetId

`string | valueFrom`

The firewall's subnet -- set on EXACTLY ONE configuration. References
an AzureSubnet's ARM id. ARM requires the subnet's name to be exactly
"AzureFirewallSubnet" and its prefix at least /26; the subnet must
carry no other workloads. Fixed at creation (changing it replaces
the firewall).

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.ipConfigurations[].publicIpAddressId

`string | valueFrom`

The Standard-SKU static public IP this configuration fronts --
references an AzurePublicIp's ARM id. Each public IP adds SNAT ports
and a DNAT frontend. Optional on the subnet-bearing configuration
only when a management_ip_configuration provides the management
path; required on every additional configuration.

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.managementIpConfiguration

`AzureFirewallManagementIpConfiguration`

The management IP configuration: a separate management path on the
dedicated "AzureFirewallManagementSubnet" (/26+) with a required
public IP. Azure requires it for BASIC-tier firewalls and for forced
tunneling. FIXED AT CREATION -- adding or removing this block (or
changing its subnet) replaces the firewall.

### spec.managementIpConfiguration.name

`string` · required

The configuration's name (must differ from every data-path
configuration name).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.managementIpConfiguration.subnetId

`string | valueFrom` · required

The management subnet -- references an AzureSubnet's ARM id. ARM
requires the subnet's name to be exactly
"AzureFirewallManagementSubnet" with at least a /26 prefix. Fixed at
creation.

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.managementIpConfiguration.publicIpAddressId

`string | valueFrom` · required

The management path's public IP -- required by ARM. References an
AzurePublicIp's ARM id (Standard SKU, static).

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.firewallPolicyId

`string | valueFrom`

The firewall policy this instance enforces -- references an
AzureFirewallPolicy's ARM id. The policy carries rules, threat
intelligence, TLS inspection, and IDPS; its tier must match sku_tier.

- references: AzureFirewallPolicy (`status.outputs.firewall_policy_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewallPolicy, name: <that resource's name>, fieldPath: status.outputs.firewall_policy_id}} -- a bare string does not parse

### spec.threatIntelMode

`enum`

How the firewall acts on traffic matching Microsoft's threat
intelligence feed when NO policy is attached (policy-attached
firewalls take their threat-intelligence posture from the policy).
Unspecified lets Azure apply its default (Alert).

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_threat_intel_mode_unspecified` -- Not specified -- the field is not sent; Azure applies its default (Alert).
- `ALERT` -- Log matches; let traffic flow.
- `DENY` -- Log and block matches.
- `OFF` -- Disable threat-intelligence filtering.

### spec.dnsServers

`[]string`

Custom upstream DNS servers the firewall resolves against. Setting
servers implicitly turns the DNS proxy ON (Azure couples them on the
wire) regardless of dns_proxy_enabled. ONLY legal on a firewall
WITHOUT an attached policy -- ARM rejects firewall-level DNS
parameters when firewall_policy_id is set (DNS is managed by the
policy's dns block instead).

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dnsProxyEnabled

`bool`

Run the firewall as a DNS proxy for clients pointing their DNS at
its private IP. Required for FQDN-based network rules to resolve
deterministically. Setting dns_servers forces this on regardless.
ONLY legal on a firewall WITHOUT an attached policy -- configure the
policy's dns block instead when firewall_policy_id is set.

### spec.privateIpRanges

`[]string`

The address ranges the firewall treats as PRIVATE (never SNATed).
Each entry is a CIDR or the literal token "IANAPrivateRanges" (the
RFC 1918 set). When unset, Azure defaults to the IANA private
ranges. Policy-attached firewalls normally configure SNAT ranges on
the policy instead.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.virtualHub

`AzureFirewallVirtualHub`

The Virtual WAN hub to deploy into (AZFW_HUB model only).

### spec.virtualHub.virtualHubId

`string | valueFrom` · required

The Virtual WAN hub's ARM id. A bare reference: Planton has no
Virtual WAN kind yet, so supply the id as a literal or an explicit
kind/fieldPath reference.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.virtualHub.publicIpCount

`int32` · optional (explicit presence)

How many public IPs Azure assigns the hub firewall (Azure manages
the addresses; they surface as outputs). At least 1; defaults to 1.

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.zones

`[]string`

The availability zones the firewall spans, e.g. ["1", "2", "3"].
Zone redundancy is free on Azure Firewall (you pay only the normal
deployment cost) and is the production posture. Fixed at creation.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.tags

`map<string, string>`

Free-form tags applied to the firewall, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins.

## Validation Rules

- `hub_sku_requires_virtual_hub`: An AZFW_HUB firewall deploys into a Virtual WAN hub -- set virtual_hub (and remove ip_configurations); an AZFW_VNET firewall uses ip_configurations instead of virtual_hub
- `virtual_hub_excludes_ip_configurations`: A hub-deployed firewall gets its addressing from the Virtual WAN hub -- remove ip_configurations
- `vnet_firewall_requires_ip_configuration`: A VNet-deployed firewall needs at least one ip_configuration (the block carrying the AzureFirewallSubnet subnet_id)
- `exactly_one_subnet_bearing_ip_configuration`: Exactly one ip_configuration must carry subnet_id -- the firewall lives in one subnet; additional configurations only add public IPs
- `public_ip_required_without_management_path`: At least one ip_configuration needs a public IP unless a management_ip_configuration provides the management path (Azure's create contract)
- `management_name_distinct_from_ip_configurations`: The management_ip_configuration's name must differ from every ip_configuration name
- `policy_owns_dns_configuration`: A policy-attached firewall takes its DNS configuration from the policy -- move dns_servers/dns_proxy_enabled to the firewall policy's dns block (ARM rejects firewall-level DNS parameters when a policy is attached)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFirewall, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_id` | `string` | The Azure Resource Manager ID of the firewall. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/azureFirewalls/{name} |
| `status.outputs.firewall_name` | `string` | The name of the firewall resource. |
| `status.outputs.private_ip_address` | `string` | The firewall's private IP on its data-path subnet -- THE hub-spoke seam: spoke route tables send 0.0.0.0/0 (or any egress-controlled prefix) here via a VIRTUAL_APPLIANCE next hop (AzureRouteTable.routes[].next_hop_in_ip_address). Empty for hub-deployed firewalls (see virtual_hub_private_ip_address). |
| `status.outputs.management_private_ip_address` | `string` | The management path's private IP, when a management_ip_configuration is deployed. |
| `status.outputs.virtual_hub_public_ip_addresses` | `[]string` | The public IPs Azure assigned to a hub-deployed firewall (AZFW_HUB). Empty for VNet firewalls -- their public IPs are the referenced AzurePublicIp resources' own outputs. |
| `status.outputs.virtual_hub_private_ip_address` | `string` | The private IP Azure assigned to a hub-deployed firewall. Empty for VNet firewalls (see private_ip_address). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.ipConfigurations[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.ipConfigurations[].publicIpAddressId` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.managementIpConfiguration.subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.managementIpConfiguration.publicIpAddressId` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.firewallPolicyId` | AzureFirewallPolicy | `status.outputs.firewall_policy_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureRouteTable | `spec.routes[].nextHopInIpAddress` | `status.outputs.private_ip_address` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
