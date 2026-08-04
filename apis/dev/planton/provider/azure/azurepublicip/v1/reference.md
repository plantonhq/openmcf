# AzurePublicIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzurePublicIpSpec** defines the configuration for creating an Azure
Public IP Address: a static, internet-routable address that load
balancers, application gateways, NAT gateways, firewalls, and VMs attach
for inbound or outbound connectivity.

The public IP is a foundational, first-class node: higher-level resources
reference it by ARM ID rather than creating their own, so one address can
move between consumers (e.g. re-pointing a frontend during a blue/green
cutover) without changing what the world has allowlisted.

Allocation is always STATIC and is deliberately not modeled: dynamic
allocation existed only for the Basic SKU, whose creation Azure
discontinued in 2025 (fully retired September 30, 2025), and every
current SKU requires static allocation. The address is assigned at
creation and persists for the resource's lifetime.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: test-pip
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-public-ip
  zones:
    - "1"
    - "2"
    - "3"
  domainNameLabel: planton-hack-test
  domainNameLabelScope: TENANT_REUSE
  idleTimeoutInMinutes: 10
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.skuTier` | `enum` |  |  |  |
| `spec.ipVersion` | `enum` |  |  |  |
| `spec.zones` | `[]string` |  |  |  |
| `spec.publicIpPrefixId` | `string \| valueFrom` |  |  | AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`) |
| `spec.domainNameLabel` | `string` |  |  |  |
| `spec.domainNameLabelScope` | `enum` |  |  |  |
| `spec.reverseFqdn` | `string` |  |  |  |
| `spec.idleTimeoutInMinutes` | `int32` |  | `4` |  |
| `spec.ipTags` | `map<string, string>` |  |  |  |
| `spec.ddosProtectionMode` | `enum` |  |  |  |
| `spec.ddosProtectionPlanId` | `string` |  |  |  |
| `spec.edgeZone` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the public IP will be created, e.g. "eastus",
"westeurope". Must match the region of the resource it will attach to.
Changing the region replaces the address.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the public IP will be created in. Can be a
literal resource-group name or a reference to an AzureResourceGroup's
name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the public IP, unique within the resource group. 1-80
characters (alphanumerics, underscores, periods, and hyphens; must
start with a letter or number and end with a letter, number, or
underscore). Changing the name replaces the resource -- and with it
the actual address -- so name it after the endpoint it represents
("prod-gateway-frontend", "prod-nat-egress").

- rule: Public IP names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.sku

`enum`

The SKU. Unspecified applies Azure's default (STANDARD) -- the
production tier every current architecture uses (the retired Basic SKU
is deliberately not modeled). STANDARD_V2 is Azure's next-generation
SKU, required to attach the address to a StandardV2 NAT gateway.
Fixed at creation; a GLOBAL-tier address must keep the STANDARD SKU.

Allowed values (use exactly as shown):

- `azure_public_ip_sku_unspecified` -- Not specified: Azure's default (Standard).
- `STANDARD` -- The production SKU every current architecture uses: static, secure by default (closed until a NSG admits traffic), zone-capable.
- `STANDARD_V2` -- Azure's next-generation SKU, required for StandardV2 NAT gateway attachment. Not valid with the GLOBAL tier.

### spec.skuTier

`enum`

The SKU tier. Unspecified applies Azure's default (REGIONAL) --
correct for virtually everything. GLOBAL exists solely for
cross-region load balancer frontends and requires the STANDARD SKU.
Fixed at creation.

Allowed values (use exactly as shown):

- `azure_public_ip_sku_tier_unspecified` -- Not specified: Azure's default (Regional).
- `REGIONAL` -- A regional address -- correct for virtually everything.
- `GLOBAL` -- A globally-anycast address for cross-region load balancer frontends. Requires the STANDARD SKU.

### spec.ipVersion

`enum`

The IP version. Unspecified applies Azure's default (IPv4). Fixed at
creation.

Allowed values (use exactly as shown):

- `azure_public_ip_ip_version_unspecified` -- Not specified: Azure's default (IPv4).
- `IPV4` -- An IPv4 address.
- `IPV6` -- An IPv6 address.

### spec.zones

`[]string`

Availability zones the address is anchored to. Multiple zones
("1","2","3") make it zone-redundant -- the production default; a
single zone pins it; empty leaves the address non-zonal (NOT
zone-redundant). Zone support varies by region. Fixed at creation.

- rule: {"repeated":{"items":{"string":{"in":["1","2","3"]}}}}

### spec.publicIpPrefixId

`string | valueFrom`

Allocate this address from a reserved Public IP Prefix instead of
Microsoft's general pool, by ARM ID. Addresses drawn from a prefix
come from one contiguous, allowlistable range. Fixed at creation.

- references: AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIpPrefix, name: <that resource's name>, fieldPath: status.outputs.public_ip_prefix_id}} -- a bare string does not parse

### spec.domainNameLabel

`string`

An optional DNS label that gives the address a stable Azure-managed
name: {label}.{region}.cloudapp.azure.com. Must be unique within the
region (see domain_name_label_scope for relaxing that to a hashed,
reuse-scoped name). 3-63 characters: lowercase letters, digits, and
hyphens; starts with a letter, ends with a letter or digit.

- rule: domain_name_label must start with a lowercase letter, end with a letter or digit, contain only lowercase letters, digits, and hyphens, and be 3-63 characters

### spec.domainNameLabelScope

`enum`

Scope-based reuse policy for domain_name_label: Azure hashes the label
with the chosen scope so the same label can safely recur across
tenants, subscriptions, or resource groups (defense against dangling-
DNS subdomain takeover). Unspecified keeps the classic region-unique
label behavior. Requires domain_name_label; changing an existing
scope replaces the resource.

Allowed values (use exactly as shown):

- `azure_public_ip_domain_name_label_scope_unspecified` -- Not specified: the label is used as-is and must be unique within the region.
- `TENANT_REUSE` -- The hashed label is reusable across different tenants.
- `SUBSCRIPTION_REUSE` -- The hashed label is reusable across different subscriptions.
- `RESOURCE_GROUP_REUSE` -- The hashed label is reusable across different resource groups.
- `NO_REUSE` -- The hashed label is never reusable -- the strictest takeover defense.

### spec.reverseFqdn

`string`

A fully qualified domain name that resolves TO this address, recorded
as its reverse-DNS (PTR) name -- what mail servers and forward-
confirmed-reverse-DNS checks see. The FQDN must already resolve to the
address (create the forward record first, then set this). Updatable in
place.

### spec.idleTimeoutInMinutes

`int32` · optional (explicit presence)

How long Azure keeps an idle TCP connection open before reclaiming it,
in minutes (4-30). Unset uses Azure's default (4). Raise it for
long-lived idle connections (WebSockets, database sessions) that must
survive between keepalives. Updatable in place.

- default: `4`
- rule: {"int32":{"lte":30,"gte":4}}

### spec.ipTags

`map<string, string>`

Azure IP tags -- routing metadata attached to the address itself (e.g.
"RoutingPreference": "Internet" for cold-potato vs hot-potato transit),
NOT governance tags (use tags for those). Only specific tag/value
pairs are permitted by Azure, some requiring subscription enablement.
Fixed at creation.

### spec.ddosProtectionMode

`enum`

DDoS protection stance for this address. Unspecified applies Azure's
default (inherit from the virtual network's protection plan, if any).
ENABLED attaches dedicated IP-level protection -- pair it with
ddos_protection_plan_id; DISABLED opts the address out even when its
network is protected. Updatable in place.

Allowed values (use exactly as shown):

- `azure_public_ip_ddos_protection_mode_unspecified` -- Not specified: Azure's default -- inherit protection from the virtual network's DDoS plan, if any.
- `DISABLED` -- Opt the address out of DDoS protection even when its network is protected.
- `ENABLED` -- Dedicated IP-level protection; pair with ddos_protection_plan_id.

### spec.ddosProtectionPlanId

`string`

The ARM ID of the DDoS protection plan backing IP-level protection.
Only valid (and required in practice) when ddos_protection_mode is
ENABLED. Plain ARM ID: DDoS plans are shared, rarely-created
governance resources not yet modeled as a Planton kind.
Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/ddosProtectionPlans/{name}

### spec.edgeZone

`string`

Deploy the address into an Azure Edge Zone (a metro-local Azure
extension) instead of the main region, e.g. "losangeles". Leave unset
for the standard region -- the overwhelmingly common case. Fixed at
creation.

### spec.tags

`map<string, string>`

Free-form tags applied to the public IP, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Validation Rules

- `global_tier_requires_standard_sku`: A GLOBAL-tier public IP must keep the STANDARD SKU (ARM rejects StandardV2 with the Global tier)
- `ddos_plan_requires_enabled_mode`: ddos_protection_plan_id can only be set when ddos_protection_mode is ENABLED
- `label_scope_requires_label`: domain_name_label_scope requires domain_name_label to be set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePublicIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.public_ip_id` | `string` | The Azure Resource Manager ID of the public IP. This is the primary output: AzureApplicationGateway, AzureLoadBalancer, and AzureNatGateway reference it to attach the address. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/publicIPAddresses/{name} |
| `status.outputs.ip_address` | `string` | The allocated address itself. Static for the resource's lifetime -- the value that lands in DNS records and partner allowlists. |
| `status.outputs.fqdn` | `string` | The Azure-managed FQDN ({label}.{region}.cloudapp.azure.com), populated only when domain_name_label is set; empty otherwise. |
| `status.outputs.public_ip_name` | `string` | The name of the public IP resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.publicIpPrefixId` | AzurePublicIpPrefix | `status.outputs.public_ip_prefix_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAksCluster | `spec.networkProfile.loadBalancerProfile.outboundIpAddressIds` | `status.outputs.public_ip_id` |
| AzureApplicationGateway | `spec.frontendIpConfigurations[].publicIpAddressId` | `status.outputs.public_ip_id` |
| AzureFirewall | `spec.ipConfigurations[].publicIpAddressId` | `status.outputs.public_ip_id` |
| AzureFirewall | `spec.managementIpConfiguration.publicIpAddressId` | `status.outputs.public_ip_id` |
| AzureLoadBalancer | `spec.frontendIpConfigurations[].publicIpAddressId` | `status.outputs.public_ip_id` |
| AzureNatGateway | `spec.publicIpIds` | `status.outputs.public_ip_id` |
| AzureNetworkInterface | `spec.ipConfigurations[].publicIpAddressId` | `status.outputs.public_ip_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
