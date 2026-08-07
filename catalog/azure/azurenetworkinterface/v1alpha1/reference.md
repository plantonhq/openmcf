# AzureNetworkInterface

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureNetworkInterfaceSpec** defines the configuration for creating an
Azure Network Interface (NIC): the attachment point that gives a virtual
machine its presence in a subnet.

The NIC is a first-class resource in Azure's own model -- a VM does not
contain its network configuration, it REFERENCES one or more NICs -- and
the catalog mirrors that honestly:
- an AzureVirtualMachine consumes this NIC's network_interface_id output
  (a VM can carry several NICs: management + data planes, appliance arms),
- each ip_configuration deploys into a referenced AzureSubnet and may
  front a referenced AzurePublicIp,
- a network security group attaches here (network_security_group_id) to
  filter this NIC's traffic specifically -- the per-workload complement
  to the subnet-level NSG attachment,
- load-balancer membership is expressed HERE, from the member side
  (Azure's own model): each ip_configuration lists the backend pools it
  joins and the inbound NAT rules it completes, referencing the load
  balancer's exported per-name IDs.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNetworkInterface
metadata:
  name: test-nic
  labels:
    environment: production
    team: platform
spec:
  region: eastus

  # The resource group the NIC lives in (literal value here; a manifest
  # can also reference an AzureResourceGroup's name output via valueFrom).
  resourceGroup:
    value: test-rg

  name: app-nic

  # A single IPv4 configuration: dynamic private address in a referenced
  # subnet, fronted by a referenced public IP, joining a load-balancer
  # pool, completing a single-target NAT rule, and joining an App Gateway
  # pool (each membership realized as its own association resource).
  ipConfigurations:
    - name: primary
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/app
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/test-pip
      loadBalancerBackendAddressPoolIds:
        - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/loadBalancers/test-lb/backendAddressPools/web
      loadBalancerInboundNatRuleIds:
        - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/loadBalancers/test-lb/inboundNatRules/ssh-admin
      applicationGatewayBackendAddressPoolIds:
        - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/applicationGateways/test-agw/backendAddressPools/web

  # Production NICs on supported VM sizes should enable SR-IOV.
  acceleratedNetworkingEnabled: true

  # NIC-level NSG, realized as an association resource.
  networkSecurityGroupId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/networkSecurityGroups/test-nsg

  # User tags merged over the metadata-derived tags.
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.ipConfigurations` | `[]AzureNetworkInterfaceIpConfiguration` | yes |  |  |
| `spec.ipConfigurations[].name` | `string` | yes |  |  |
| `spec.ipConfigurations[].subnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.ipConfigurations[].privateIpAllocation` | `enum` |  |  |  |
| `spec.ipConfigurations[].privateIpAddress` | `string` |  |  |  |
| `spec.ipConfigurations[].privateIpVersion` | `enum` |  |  |  |
| `spec.ipConfigurations[].publicIpAddressId` | `string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.ipConfigurations[].primary` | `bool` |  |  |  |
| `spec.ipConfigurations[].gatewayLoadBalancerFrontendIpConfigurationId` | `string` |  |  |  |
| `spec.ipConfigurations[].loadBalancerBackendAddressPoolIds` | `[]string \| valueFrom` |  |  | AzureLoadBalancer (`status.outputs.backend_pool_ids`) |
| `spec.ipConfigurations[].loadBalancerInboundNatRuleIds` | `[]string \| valueFrom` |  |  | AzureLoadBalancer (`status.outputs.nat_rule_ids`) |
| `spec.ipConfigurations[].applicationGatewayBackendAddressPoolIds` | `[]string \| valueFrom` |  |  | AzureApplicationGateway (`status.outputs.backend_address_pool_ids`) |
| `spec.dnsServers` | `[]string` |  |  |  |
| `spec.internalDnsNameLabel` | `string` |  |  |  |
| `spec.acceleratedNetworkingEnabled` | `bool` |  |  |  |
| `spec.ipForwardingEnabled` | `bool` |  |  |  |
| `spec.auxiliaryMode` | `enum` |  |  |  |
| `spec.auxiliarySku` | `enum` |  |  |  |
| `spec.edgeZone` | `string` |  |  |  |
| `spec.networkSecurityGroupId` | `string \| valueFrom` |  |  | AzureNetworkSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.applicationSecurityGroupIds` | `[]string \| valueFrom` |  |  | AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the NIC lives in, e.g. "eastus". Must match the
region of the virtual network it deploys into and of the VM that will
attach it. Changing the region replaces the NIC.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the NIC will be created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the NIC, unique within the resource group. 1-80 characters
(alphanumerics, underscores, periods, and hyphens; must start with a
letter or number and end with a letter, number, or underscore).
Changing the name replaces the NIC -- and detaches it from any VM.

- rule: Network interface names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.ipConfigurations

`[]AzureNetworkInterfaceIpConfiguration` · required

The NIC's IP configurations -- at least one. Most NICs carry exactly
one; multiple configurations serve dual-stack (an IPv4 and an IPv6
configuration side by side) and multi-IP scenarios (many addresses on
one NIC for per-site TLS or NAT pools). When more than one is
declared, the FIRST must be marked primary (ARM's contract).

- rule: {"repeated":{"minItems":"1"}}
- rule: STATIC private_ip_allocation requires private_ip_address (and DYNAMIC forbids it)
- rule: an IPv4 ip_configuration requires subnet_id (IPv6 configurations inherit the NIC's subnet placement)

### spec.ipConfigurations[].name

`string` · required

A label for this configuration, unique within the NIC. A
configuration-level name (visible in the portal and IaC state), not
an Azure resource name -- "primary" or "ipv6" reads well.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.ipConfigurations[].subnetId

`string | valueFrom`

The subnet this configuration's private address lives in, by ARM ID.
Required for IPv4 configurations (ARM's contract); an IPv6
configuration inherits the NIC's subnet placement. All of a NIC's
configurations must address subnets of the same virtual network.

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.ipConfigurations[].privateIpAllocation

`enum`

How the private address is assigned. Unspecified applies DYNAMIC --
Azure picks a free address from the subnet, which is right for
virtually all workloads (the address is stable for the NIC's
lifetime). STATIC pins a specific address (set private_ip_address)
for appliances and servers whose IP is configuration elsewhere.

Allowed values (use exactly as shown):

- `azure_network_interface_private_ip_allocation_unspecified` -- Not specified: DYNAMIC.
- `DYNAMIC` -- Azure assigns a free address from the subnet at creation; it stays stable for the NIC's lifetime.
- `STATIC` -- A specific address is pinned (set private_ip_address).

### spec.ipConfigurations[].privateIpAddress

`string`

For STATIC allocation: the exact private address to pin, which must
fall inside the subnet's range and be unassigned. Leave unset for
DYNAMIC (the assigned address surfaces in the outputs).

### spec.ipConfigurations[].privateIpVersion

`enum`

The address family of this configuration. Unspecified applies Azure's
default (IPV4). A dual-stack NIC carries an IPv4 configuration and an
IPV6 one side by side (in a dual-stack subnet).

Allowed values (use exactly as shown):

- `azure_network_interface_private_ip_version_unspecified` -- Not specified: IPv4.
- `IPV4` -- An IPv4 private address.
- `IPV6` -- An IPv6 private address (dual-stack NICs pair this with an IPv4 configuration).

### spec.ipConfigurations[].publicIpAddressId

`string | valueFrom`

The public IP fronting this configuration, by ARM ID. References a
first-class AzurePublicIp resource so the address is visible in the
resource graph, allowlistable, and reusable. Omit for private-only
NICs -- the production shape behind a load balancer or NAT gateway.

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.ipConfigurations[].primary

`bool`

Whether this is the NIC's primary configuration. Exactly one
configuration is primary; with a single configuration ARM treats it
as primary automatically, and with multiple the FIRST must be marked
(spec-level validation enforces both).

### spec.ipConfigurations[].gatewayLoadBalancerFrontendIpConfigurationId

`string`

The frontend IP configuration of a Gateway-SKU load balancer that
chains this NIC into a gateway appliance path, by ARM ID
(referenceable via the gateway load balancer's
frontend_ip_configuration_ids output). A niche service-chaining seam.

### spec.ipConfigurations[].loadBalancerBackendAddressPoolIds

`[]string | valueFrom`

Load-balancer backend pools this configuration joins, by pool ARM
ID -- membership is expressed from the member side in Azure's model.
Reference a pool through the load balancer's name-keyed map output,
e.g. valueFrom fieldPath "status.outputs.backend_pool_ids.web".
Each membership is realized as its own association resource, so
joining and leaving pools never touches the NIC itself.

- references: AzureLoadBalancer (`status.outputs.backend_pool_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLoadBalancer, name: <that resource's name>, fieldPath: status.outputs.backend_pool_ids}} -- a bare string does not parse

### spec.ipConfigurations[].loadBalancerInboundNatRuleIds

`[]string | valueFrom`

Single-target inbound NAT rules this configuration completes, by
rule ARM ID -- the load balancer declares the rule (frontend port ->
backend port) and the NIC-side association picks which instance
receives the forwarded traffic. Reference a rule through the load
balancer's name-keyed map output, e.g. valueFrom fieldPath
"status.outputs.nat_rule_ids.ssh-admin". Realized as association
resources.

- references: AzureLoadBalancer (`status.outputs.nat_rule_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLoadBalancer, name: <that resource's name>, fieldPath: status.outputs.nat_rule_ids}} -- a bare string does not parse

### spec.ipConfigurations[].applicationGatewayBackendAddressPoolIds

`[]string | valueFrom`

Application Gateway backend pools this configuration joins, by pool
ARM ID -- membership is expressed from the member side in Azure's
model. Reference a pool through the gateway's name-keyed map output,
e.g. valueFrom fieldPath
"status.outputs.backend_address_pool_ids.web". Realized as
association resources.

- references: AzureApplicationGateway (`status.outputs.backend_address_pool_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationGateway, name: <that resource's name>, fieldPath: status.outputs.backend_address_pool_ids}} -- a bare string does not parse

### spec.dnsServers

`[]string`

DNS servers workloads behind this NIC use, overriding the virtual
network's DNS settings for this NIC only. Rarely set -- prefer
configuring DNS on the virtual network so every workload inherits it;
this exists for appliances that need different resolution than their
network. Updatable in place.

### spec.internalDnsNameLabel

`string`

The (relative) DNS label other VMs in the same virtual network can
resolve this NIC's private IP by, e.g. "app-1" (the full name becomes
{label}.{vnet-internal-suffix}, surfaced after creation). Leave unset
for IP-only addressing. Updatable in place.

### spec.acceleratedNetworkingEnabled

`bool`

Whether accelerated networking (SR-IOV) is enabled: the NIC bypasses
the host's virtual switch for dramatically lower latency and higher
packets-per-second. Azure's default is false, but production NICs on
supported VM sizes (most current general-purpose sizes with 2+ vCPUs)
should enable it -- the constraint is the VM size, not the workload.
Updatable in place.

### spec.ipForwardingEnabled

`bool`

Whether the NIC forwards traffic not addressed to it. Azure's default
is false; enable ONLY on network virtual appliances (firewalls,
routers) that legitimately route other workloads' traffic -- a route
table pointing at an appliance's IP silently blackholes unless the
appliance NIC forwards. Updatable in place.

### spec.auxiliaryMode

`enum`

The auxiliary mode for network-virtual-appliance acceleration (a
preview Azure feature the subscription must be enrolled in): higher
connections-per-second and large concurrent connection counts for
NVAs. Unspecified sends nothing -- correct for every non-appliance
NIC. Must be set together with auxiliary_sku.

Allowed values (use exactly as shown):

- `azure_network_interface_auxiliary_mode_unspecified` -- Not specified: no auxiliary acceleration.
- `ACCELERATED_CONNECTIONS` -- Optimizes connections-per-second for appliance workloads.
- `FLOATING` -- Floating IP support for the auxiliary path.
- `MAX_CONNECTIONS` -- Optimizes for large numbers of simultaneous connections.

### spec.auxiliarySku

`enum`

The auxiliary SKU sizing the NVA acceleration (A1 smallest to A8
largest). Must be set together with auxiliary_mode.

Allowed values (use exactly as shown):

- `azure_network_interface_auxiliary_sku_unspecified` -- Not specified: no auxiliary SKU.
- `A1` -- The smallest acceleration tier.
- `A2` -- The second acceleration tier.
- `A4` -- The mid acceleration tier.
- `A8` -- The largest acceleration tier.

### spec.edgeZone

`string`

The Azure Edge Zone the NIC is deployed in, for edge-computing
workloads that pin resources to a metro-local extended zone. Leave
unset for regular regional deployment. Fixed at creation.

### spec.networkSecurityGroupId

`string | valueFrom`

The network security group that filters THIS NIC's traffic, by ARM
ID. NIC-level filtering is the per-workload complement to the
subnet-level NSG (AzureSubnet's network_security_group_id): when both
are attached, inbound traffic must pass the subnet NSG then the NIC
NSG, outbound the reverse. Omit to rely on subnet-level filtering
alone -- the common case; attach here when one workload needs rules
its subnet neighbors must not share.

- references: AzureNetworkSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureNetworkSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.applicationSecurityGroupIds

`[]string | valueFrom`

Application security groups this NIC joins. ASG membership lets NSG
rules target workload groups ("web-servers", "databases") instead of
IP ranges. Each entry is an application security group by ARM ID, or a
reference to an AzureApplicationSecurityGroup's output.

- references: AzureApplicationSecurityGroup (`status.outputs.application_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureApplicationSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.application_security_group_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the NIC, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them.
Updatable in place.

## Validation Rules

- `nic_auxiliary_mode_sku_paired`: auxiliary_mode and auxiliary_sku must be set together (both or neither)
- `nic_first_ip_configuration_primary_when_multiple`: when a NIC has multiple ip_configurations, the first must be marked primary (ARM's contract)
- `nic_at_most_one_primary_ip_configuration`: at most one ip_configuration may be marked primary

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureNetworkInterface, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_interface_id` | `string` | The Azure Resource Manager ID of the NIC. This is the primary output: AzureVirtualMachine's network_interface_ids references it to attach the NIC to a VM. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkInterfaces/{name} |
| `status.outputs.network_interface_name` | `string` | The name of the NIC. |
| `status.outputs.private_ip_address` | `string` | The primary configuration's private IP address -- what backends, firewall rules, and DNS records key on. |
| `status.outputs.private_ip_addresses` | `[]string` | The private IP addresses of ALL configurations, in configuration order (multi-IP and dual-stack NICs carry more than one). |
| `status.outputs.mac_address` | `string` | The NIC's MAC address, populated once the NIC is attached to a running virtual machine (empty until then) -- what license servers and appliance registrations key on. |
| `status.outputs.internal_domain_name_suffix` | `string` | The DNS suffix for the virtual network's internal name resolution (the {vnet-internal-suffix} completing internal_dns_name_label into a resolvable FQDN). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.ipConfigurations[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.ipConfigurations[].publicIpAddressId` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.ipConfigurations[].loadBalancerBackendAddressPoolIds` | AzureLoadBalancer | `status.outputs.backend_pool_ids` |
| `spec.ipConfigurations[].loadBalancerInboundNatRuleIds` | AzureLoadBalancer | `status.outputs.nat_rule_ids` |
| `spec.ipConfigurations[].applicationGatewayBackendAddressPoolIds` | AzureApplicationGateway | `status.outputs.backend_address_pool_ids` |
| `spec.networkSecurityGroupId` | AzureNetworkSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.applicationSecurityGroupIds` | AzureApplicationSecurityGroup | `status.outputs.application_security_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVirtualMachine | `spec.networkInterfaceIds` | `status.outputs.network_interface_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
