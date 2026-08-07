# AzureLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureLoadBalancerSpec** defines the configuration for creating an Azure
Load Balancer: Azure's Layer 4 (TCP/UDP) traffic distributor, complete
with its frontends, backend address pools, health probes, load-balancing
rules, inbound NAT rules, and outbound (SNAT) rules.

The load balancer and its sub-resources are configured as ONE unit
because none of them has a life without the load balancer: a backend
pool, probe, or rule exists only inside its parent LB and is deployed,
versioned, and destroyed with it. What IS independent -- membership of a
pool -- is expressed from the member side (a network interface or a
virtual machine scale set referencing a pool by its exported ID), which
is Azure's own attachment model. The `backend_pool_ids` output keyed by
pool name is the seam members reference.

**Public vs internal is per-frontend, not per-LB.** Each frontend IP
configuration is either public (references a public IP or a public IP
prefix) or internal (references a subnet, with an optional static
private address). One load balancer can carry several frontends -- e.g.
a public frontend for ingress and an internal one for east-west traffic.

**SKU**: STANDARD (the default) is the production SKU -- zone-redundant,
SLA-backed, supports outbound rules and HA ports. GATEWAY is the niche
SKU for chaining network virtual appliances (requires the
Microsoft.Network/AllowGatewayLoadBalancer feature on the subscription).
Basic is not modeled: Azure retired it in September 2025.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureLoadBalancer
metadata:
  name: test-lb
spec:
  region: eastus
  resource_group:
    value: test-rg
  name: test-lb
  frontend_ip_configurations:
    - name: public
      public_ip_address_id:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/test-pip
    - name: internal
      subnet_id:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default
      private_ip_address: 10.0.1.100
      zones: ["1", "2", "3"]
  backend_pools:
    - name: web
    - name: appliances
      virtual_network_id:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet
      # Exercises the synchronous_mode enum seam (conditionally-mapped
      # enums must ride the hack manifest so the plan proves the mapping).
      synchronous_mode: MANUAL
      addresses:
        - name: appliance-1
          ip_address: 10.0.1.10
  health_probes:
    - name: http-health
      protocol: PROBE_HTTP
      port: 80
      request_path: /health
      probe_threshold: 2
    - name: tcp-health
      protocol: PROBE_TCP
      port: 443
  rules:
    - name: http
      frontend_ip_configuration_name: public
      protocol: TCP
      frontend_port: 80
      backend_port: 80
      backend_pool_names: [web]
      probe_name: http-health
      load_distribution: SOURCE_IP
      tcp_reset_enabled: true
      disable_outbound_snat: true
    - name: https
      frontend_ip_configuration_name: public
      protocol: TCP
      frontend_port: 443
      backend_port: 443
      backend_pool_names: [web]
      probe_name: tcp-health
  nat_rules:
    - name: ssh-admin
      frontend_ip_configuration_name: public
      protocol: TCP
      frontend_port: 2222
      backend_port: 22
    - name: per-instance-ssh
      frontend_ip_configuration_name: public
      protocol: TCP
      backend_port: 22
      backend_pool_name: web
      frontend_port_start: 50000
      frontend_port_end: 50099
  outbound_rules:
    - name: egress
      frontend_ip_configuration_names: [public]
      backend_pool_name: web
      protocol: ALL
      allocated_outbound_ports: 2048
      tcp_reset_enabled: true
  tags:
    cost-center: networking
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.skuTier` | `enum` |  |  |  |
| `spec.edgeZone` | `string` |  |  |  |
| `spec.frontendIpConfigurations` | `[]AzureLoadBalancerFrontendIpConfiguration` | yes |  |  |
| `spec.frontendIpConfigurations[].name` | `string` | yes |  |  |
| `spec.frontendIpConfigurations[].subnetId` | `string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.frontendIpConfigurations[].publicIpAddressId` | `string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.frontendIpConfigurations[].publicIpPrefixId` | `string \| valueFrom` |  |  | AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`) |
| `spec.frontendIpConfigurations[].privateIpAddress` | `string` |  |  |  |
| `spec.frontendIpConfigurations[].privateIpAddressVersion` | `enum` |  |  |  |
| `spec.frontendIpConfigurations[].zones` | `[]string` |  |  |  |
| `spec.frontendIpConfigurations[].gatewayLoadBalancerFrontendIpConfigurationId` | `string` |  |  |  |
| `spec.backendPools` | `[]AzureLoadBalancerBackendPool` |  |  |  |
| `spec.backendPools[].name` | `string` | yes |  |  |
| `spec.backendPools[].virtualNetworkId` | `string \| valueFrom` |  |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.backendPools[].synchronousMode` | `enum` |  |  |  |
| `spec.backendPools[].tunnelInterfaces` | `[]AzureLoadBalancerBackendPoolTunnelInterface` |  |  |  |
| `spec.backendPools[].tunnelInterfaces[].identifier` | `int32` |  |  |  |
| `spec.backendPools[].tunnelInterfaces[].port` | `int32` |  |  |  |
| `spec.backendPools[].tunnelInterfaces[].protocol` | `enum` |  |  |  |
| `spec.backendPools[].tunnelInterfaces[].type` | `enum` |  |  |  |
| `spec.backendPools[].addresses` | `[]AzureLoadBalancerBackendPoolAddress` |  |  |  |
| `spec.backendPools[].addresses[].name` | `string` | yes |  |  |
| `spec.backendPools[].addresses[].ipAddress` | `string` |  |  |  |
| `spec.backendPools[].addresses[].loadBalancerFrontendIpConfigurationId` | `string` |  |  |  |
| `spec.healthProbes` | `[]AzureLoadBalancerHealthProbe` |  |  |  |
| `spec.healthProbes[].name` | `string` | yes |  |  |
| `spec.healthProbes[].protocol` | `enum` |  |  |  |
| `spec.healthProbes[].port` | `int32` | yes |  |  |
| `spec.healthProbes[].requestPath` | `string` |  |  |  |
| `spec.healthProbes[].intervalInSeconds` | `int32` |  | `15` |  |
| `spec.healthProbes[].numberOfProbes` | `int32` |  | `2` |  |
| `spec.healthProbes[].probeThreshold` | `int32` |  | `1` |  |
| `spec.rules` | `[]AzureLoadBalancerRule` |  |  |  |
| `spec.rules[].name` | `string` | yes |  |  |
| `spec.rules[].frontendIpConfigurationName` | `string` |  |  |  |
| `spec.rules[].protocol` | `enum` |  |  |  |
| `spec.rules[].frontendPort` | `int32` |  |  |  |
| `spec.rules[].backendPort` | `int32` |  |  |  |
| `spec.rules[].backendPoolNames` | `[]string` | yes |  |  |
| `spec.rules[].probeName` | `string` |  |  |  |
| `spec.rules[].loadDistribution` | `enum` |  |  |  |
| `spec.rules[].idleTimeoutInMinutes` | `int32` |  | `4` |  |
| `spec.rules[].floatingIpEnabled` | `bool` |  | `false` |  |
| `spec.rules[].tcpResetEnabled` | `bool` |  | `false` |  |
| `spec.rules[].disableOutboundSnat` | `bool` |  | `false` |  |
| `spec.natRules` | `[]AzureLoadBalancerNatRule` |  |  |  |
| `spec.natRules[].name` | `string` | yes |  |  |
| `spec.natRules[].frontendIpConfigurationName` | `string` |  |  |  |
| `spec.natRules[].protocol` | `enum` |  |  |  |
| `spec.natRules[].frontendPort` | `int32` |  |  |  |
| `spec.natRules[].backendPort` | `int32` | yes |  |  |
| `spec.natRules[].backendPoolName` | `string` |  |  |  |
| `spec.natRules[].frontendPortStart` | `int32` |  |  |  |
| `spec.natRules[].frontendPortEnd` | `int32` |  |  |  |
| `spec.natRules[].floatingIpEnabled` | `bool` |  | `false` |  |
| `spec.natRules[].tcpResetEnabled` | `bool` |  | `false` |  |
| `spec.natRules[].idleTimeoutInMinutes` | `int32` |  | `4` |  |
| `spec.outboundRules` | `[]AzureLoadBalancerOutboundRule` |  |  |  |
| `spec.outboundRules[].name` | `string` | yes |  |  |
| `spec.outboundRules[].frontendIpConfigurationNames` | `[]string` | yes |  |  |
| `spec.outboundRules[].backendPoolName` | `string` | yes |  |  |
| `spec.outboundRules[].protocol` | `enum` |  |  |  |
| `spec.outboundRules[].allocatedOutboundPorts` | `int32` |  | `1024` |  |
| `spec.outboundRules[].tcpResetEnabled` | `bool` |  | `false` |  |
| `spec.outboundRules[].idleTimeoutInMinutes` | `int32` |  | `4` |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the load balancer will be created.
Must match the region of the backend resources (VMs, VMSS, AKS).
Examples: "eastus", "westus2", "westeurope", "southeastasia".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group where the load balancer will be created.
Can be a literal string or a reference to an AzureResourceGroup output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the load balancer.
Must be unique within the resource group.
Allowed characters: alphanumeric, underscores, hyphens, and periods.
Must start with alphanumeric. Length: 1 to 80 characters.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.sku

`enum`

The SKU of the load balancer. Unspecified applies STANDARD -- the
production SKU (zone redundancy, SLA, outbound rules, HA ports).
GATEWAY is for chaining network virtual appliances: every backend
pool must then declare tunnel interfaces, and the subscription needs
the Microsoft.Network/AllowGatewayLoadBalancer feature registered
(via an Azure support ticket). Fixed at creation.

Allowed values (use exactly as shown):

- `azure_load_balancer_sku_unspecified` -- Not specified: STANDARD.
- `STANDARD` -- The production SKU: zone redundancy, SLA, outbound rules, HA ports. (Basic is not modeled -- Azure retired it in September 2025.)
- `GATEWAY` -- The NVA-chaining SKU: backend pools carry tunnel interfaces and other load balancers chain their frontends through this one. Requires the Microsoft.Network/AllowGatewayLoadBalancer subscription feature (registered via an Azure support ticket).

### spec.skuTier

`enum`

The SKU tier. Unspecified applies REGIONAL -- a load balancer serving
one region, the shape virtually every deployment uses. GLOBAL creates
a cross-region load balancer whose backend pool members are the
frontends of REGIONAL load balancers (set each pool address's
load_balancer_frontend_ip_configuration_id). Requires STANDARD SKU.
Fixed at creation.

Allowed values (use exactly as shown):

- `azure_load_balancer_sku_tier_unspecified` -- Not specified: REGIONAL.
- `REGIONAL` -- A load balancer serving one region -- the shape virtually every deployment uses.
- `GLOBAL` -- A cross-region load balancer whose backend members are REGIONAL load balancers' frontends. STANDARD SKU only.

### spec.edgeZone

`string`

The Azure Edge Zone the load balancer is deployed in, for
edge-computing workloads that pin resources to a metro-local extended
zone. Leave unset for regular regional deployment. Fixed at creation.

### spec.frontendIpConfigurations

`[]AzureLoadBalancerFrontendIpConfiguration` · required

The frontend IP configurations -- the addresses that receive traffic.
At least one. Each frontend is public (public IP or public IP prefix)
or internal (subnet + optional static private address); one load
balancer can mix both. Rules, NAT rules, and outbound rules target a
frontend by its `name` (optional when exactly one frontend exists).

Azure does not allow removing ALL frontends from an existing load
balancer -- going from some frontends to none replaces the resource.

- rule: {"repeated":{"minItems":"1"}}
- rule: a frontend is either internal (subnet_id) or public (public_ip_address_id or public_ip_prefix_id) -- set at most one address source
- rule: private_ip_address (a pinned internal address) requires subnet_id
- rule: zones apply to internal (subnet) frontends; a public frontend's zone posture comes from its public IP resource

### spec.frontendIpConfigurations[].name

`string` · required

A label for this frontend, unique within the load balancer. Rules,
NAT rules, and outbound rules target the frontend by this name.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.frontendIpConfigurations[].subnetId

`string | valueFrom`

The subnet for an INTERNAL frontend, by ARM ID: the frontend gets a
private address in this subnet. Mutually exclusive with the public
address sources.

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.frontendIpConfigurations[].publicIpAddressId

`string | valueFrom`

The public IP for a PUBLIC frontend, by ARM ID. References a
first-class AzurePublicIp (Standard SKU) so the address is visible in
the resource graph, allowlistable, and reusable.

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.frontendIpConfigurations[].publicIpPrefixId

`string | valueFrom`

A public IP PREFIX for a public frontend, by ARM ID -- the frontend
draws from the reserved contiguous range, used with outbound rules to
scale SNAT ports across a known CIDR that downstream partners can
allowlist as one block.

- references: AzurePublicIpPrefix (`status.outputs.public_ip_prefix_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIpPrefix, name: <that resource's name>, fieldPath: status.outputs.public_ip_prefix_id}} -- a bare string does not parse

### spec.frontendIpConfigurations[].privateIpAddress

`string`

For an internal frontend: pin a specific private address inside the
subnet's range (it must be unassigned). Leave unset for dynamic
allocation -- the address is stable for the frontend's lifetime
either way; pin only when DNS, firewall rules, or service discovery
are configured with the address elsewhere.

### spec.frontendIpConfigurations[].privateIpAddressVersion

`enum`

The address family of an internal frontend's private address.
Unspecified applies Azure's default (IPV4). An IPv6 internal frontend
requires a dual-stack subnet.

Allowed values (use exactly as shown):

- `azure_load_balancer_private_ip_version_unspecified` -- Not specified: IPv4.
- `IPV4` -- An IPv4 private address.
- `IPV6` -- An IPv6 private address (requires a dual-stack subnet).

### spec.frontendIpConfigurations[].zones

`[]string`

Availability zones for an INTERNAL frontend's private address:
["1","2","3"] for zone redundancy (the production default posture),
a single zone to pin, or empty for no zone guarantee. A PUBLIC
frontend's zone posture comes from the referenced public IP resource
instead. Changing zones replaces the frontend.

### spec.frontendIpConfigurations[].gatewayLoadBalancerFrontendIpConfigurationId

`string`

The frontend IP configuration of a GATEWAY-SKU load balancer to chain
this frontend behind, by ARM ID -- traffic arriving here is first
steered through the gateway's network virtual appliances. Plain ARM
ID because gateway frontends are addressed as sub-resources
(referenceable via the gateway LB's frontend_ip_configuration_ids
output). Only meaningful on a STANDARD SKU frontend.

### spec.backendPools

`[]AzureLoadBalancerBackendPool`

Backend address pools that receive load-balanced traffic. Optional:
a frontend-only load balancer carrying just inbound NAT rules is
legal. Pool MEMBERSHIP is expressed from the member side -- a network
interface or scale set references `status.outputs.backend_pool_ids.<name>`
-- except for IP-based members declared inline via `addresses`.

- rule: synchronous_mode requires virtual_network_id (it governs IP-based membership of a vnet-scoped pool)
- rule: IP-based addresses require the pool's virtual_network_id (unless they reference regional load balancer frontends for a GLOBAL-tier pool)

### spec.backendPools[].name

`string` · required

The name of the backend pool, unique within the load balancer. Rules
target the pool by this name, and the pool's ARM ID is exported as
`status.outputs.backend_pool_ids.<name>` for member-side references.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.backendPools[].virtualNetworkId

`string | valueFrom`

The virtual network of this pool's IP-based members, by ARM ID.
Required when `addresses` are declared or `synchronous_mode` is set;
leave unset for a plain NIC-membership pool.

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.backendPools[].synchronousMode

`enum`

How IP-based backend members synchronize with the pool (STANDARD SKU
only, and only for vnet-scoped pools). Unspecified sends nothing --
right for NIC-membership pools. AUTOMATIC lets Azure manage member
lifecycle; MANUAL leaves membership entirely to the addresses list.
Fixed at creation.

Allowed values (use exactly as shown):

- `azure_load_balancer_backend_pool_sync_mode_unspecified` -- Not specified: no synchronous mode (NIC-membership pool).
- `AUTOMATIC` -- Azure manages IP-based member lifecycle automatically.
- `MANUAL` -- Membership is entirely governed by the declared addresses list.

### spec.backendPools[].tunnelInterfaces

`[]AzureLoadBalancerBackendPoolTunnelInterface`

GATEWAY SKU only: the tunnel interfaces (VXLAN identifiers/ports)
through which chained traffic reaches the network virtual appliances
in this pool. Required on every pool of a GATEWAY load balancer,
forbidden otherwise (spec-level validation enforces both).

### spec.backendPools[].tunnelInterfaces[].identifier

`int32`

The VXLAN network identifier for this tunnel. The conventional pair
is 800 (internal) and 801 (external).

- rule: {"int32":{"gte":0}}

### spec.backendPools[].tunnelInterfaces[].port

`int32`

The port the tunnel listens on. The conventional pair is 10800
(internal) and 10801 (external).

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.backendPools[].tunnelInterfaces[].protocol

`enum`

The traffic encapsulation protocol. VXLAN is the shape gateway
chaining uses in practice.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_load_balancer_tunnel_protocol_unspecified` -- Not specified -- invalid; every tunnel declares its protocol.
- `TUNNEL_PROTOCOL_NONE` -- No encapsulation.
- `NATIVE` -- Native (unencapsulated) forwarding.
- `VXLAN` -- VXLAN encapsulation -- the shape gateway chaining uses in practice.

### spec.backendPools[].tunnelInterfaces[].type

`enum`

Whether this tunnel carries traffic toward the appliances (INTERNAL)
or back out (EXTERNAL). A chaining pool conventionally declares one
of each.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_load_balancer_tunnel_type_unspecified` -- Not specified -- invalid; every tunnel declares its type.
- `TUNNEL_TYPE_NONE` -- No direction (rarely used).
- `INTERNAL` -- Carries traffic toward the appliances.
- `EXTERNAL` -- Carries traffic back out of the appliances.

### spec.backendPools[].addresses

`[]AzureLoadBalancerBackendPoolAddress`

IP-based backend members, declared inline: appliances or servers
addressed by IP rather than by NIC association (requires
virtual_network_id). For a GLOBAL-tier load balancer the members are
regional load balancers instead -- set each address's
load_balancer_frontend_ip_configuration_id.

- rule: each address is either an IP member (ip_address) or a regional load balancer frontend (load_balancer_frontend_ip_configuration_id) -- set exactly one

### spec.backendPools[].addresses[].name

`string` · required

A label for this member, unique within the pool.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.backendPools[].addresses[].ipAddress

`string`

The member's private IP address inside the pool's virtual network
(REGIONAL tier). Mutually exclusive with
load_balancer_frontend_ip_configuration_id.

### spec.backendPools[].addresses[].loadBalancerFrontendIpConfigurationId

`string`

For a GLOBAL-tier pool: the frontend IP configuration of the REGIONAL
load balancer this member represents, by ARM ID (referenceable via
that load balancer's frontend_ip_configuration_ids output). Mutually
exclusive with ip_address.

### spec.healthProbes

`[]AzureLoadBalancerHealthProbe`

Health probes that check backend instance availability. Probes run at
`interval_in_seconds` intervals; after `number_of_probes` consecutive
failures the backend is removed from rotation and re-added when it
recovers. Rules reference a probe by name; a rule without a probe
load-balances blindly (legal, but production rules should probe).

- rule: PROBE_HTTP/PROBE_HTTPS probes require request_path; PROBE_TCP probes must not set it

### spec.healthProbes[].name

`string` · required

The name of the health probe, unique within the load balancer.
Referenced by load-balancing rules; the probe's ARM ID is exported as
`status.outputs.probe_ids.<name>` (a virtual machine scale set's
rolling-upgrade health probe references it).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.healthProbes[].protocol

`enum`

The probe protocol. Unspecified applies Azure's default (PROBE_TCP:
port-open means healthy). PROBE_HTTP / PROBE_HTTPS send a GET to
`request_path` and require HTTP 200 -- prefer them when the workload
has a health endpoint, because a process can listen while unhealthy.

Allowed values (use exactly as shown):

- `azure_load_balancer_probe_protocol_unspecified` -- Not specified: PROBE_TCP.
- `PROBE_TCP` -- A TCP connection attempt: port open means healthy.
- `PROBE_HTTP` -- An HTTP GET to request_path: 200 OK means healthy.
- `PROBE_HTTPS` -- An HTTPS GET to request_path: 200 OK means healthy (certificate validity is not checked).

### spec.healthProbes[].port

`int32` · required

The port number to probe on the backend instances.
Examples: 80 (HTTP), 443 (HTTPS), 8080 (app), 3306 (MySQL).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.healthProbes[].requestPath

`string`

The URI path for PROBE_HTTP / PROBE_HTTPS probes -- required for
them, forbidden for PROBE_TCP (spec-level validation enforces both).
The probe sends a GET request to this path and expects HTTP 200.
Examples: "/health", "/api/healthz", "/ready".

### spec.healthProbes[].intervalInSeconds

`int32` · optional (explicit presence)

The interval between probe attempts, in seconds. Lower values detect
failures faster but generate more probe traffic.
Default: 15 seconds. Minimum: 5 seconds.

- default: `15`
- rule: {"int32":{"gte":5}}

### spec.healthProbes[].numberOfProbes

`int32` · optional (explicit presence)

The number of consecutive probe failures before marking a backend
unhealthy and removing it from rotation. Default: 2 (Azure default).

- default: `2`
- rule: {"int32":{"gte":1}}

### spec.healthProbes[].probeThreshold

`int32` · optional (explicit presence)

The number of consecutive SUCCESSFUL probes required before an
instance is considered healthy (the flap dampener). Default: 1
(Azure default -- one success re-admits). Range: 1 to 100.

- default: `1`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.rules

`[]AzureLoadBalancerRule`

Load-balancing rules mapping a frontend port/protocol to a backend
pool and port. Each rule targets one frontend by name (optional when
exactly one frontend exists) and one backend pool by name (two pools
only on GATEWAY SKU).

- rule: protocol ALL (HA ports) requires frontend_port and backend_port to be 0, and non-ALL rules must use non-zero ports

### spec.rules[].name

`string` · required

The name of the rule, unique within the load balancer.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.rules[].frontendIpConfigurationName

`string`

The frontend this rule listens on, by the frontend's `name`.
Optional when the load balancer has exactly one frontend; required
with several (spec-level validation enforces both).

### spec.rules[].protocol

`enum`

The transport protocol. ALL creates an HA-ports rule (every port,
both TCP and UDP) -- the appliance/NVA pattern, only valid on an
internal STANDARD frontend with frontend_port and backend_port 0.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_load_balancer_transport_protocol_unspecified` -- Not specified -- invalid; every rule declares its protocol.
- `TCP` -- TCP traffic.
- `UDP` -- UDP traffic.
- `ALL` -- Both TCP and UDP. On a load-balancing rule this creates an HA-ports rule (internal STANDARD frontends only, ports 0); on an outbound rule it SNATs all traffic.

### spec.rules[].frontendPort

`int32`

The port on the frontend that receives traffic.
Range: 0 to 65534. Use 0 only for HA ports (protocol ALL).

- rule: {"int32":{"lte":65534,"gte":0}}

### spec.rules[].backendPort

`int32`

The port on the backend instances that receives forwarded traffic.
Range: 0 to 65535. Use 0 only for HA ports (protocol ALL).

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.rules[].backendPoolNames

`[]string` · required

The backend pool(s) this rule routes to, by pool `name`. One pool for
STANDARD SKU; a GATEWAY rule may target two pools (the dual-tunnel
NVA pattern). At least one.

- rule: {"repeated":{"minItems":"1","maxItems":"2"}}

### spec.rules[].probeName

`string`

The health probe gating this rule's backends, by probe `name`.
Optional -- a rule without a probe load-balances blindly; production
rules should reference one.

### spec.rules[].loadDistribution

`enum`

Session persistence. Unspecified applies Azure's default (DEFAULT:
5-tuple hash -- effectively no persistence, best distribution).
SOURCE_IP (2-tuple) and SOURCE_IP_PROTOCOL (3-tuple) pin a client to
one backend -- for legacy stateful workloads only; prefer stateless
backends over persistence.

Allowed values (use exactly as shown):

- `azure_load_balancer_load_distribution_unspecified` -- Not specified: DEFAULT.
- `DEFAULT` -- 5-tuple hash (source IP+port, destination IP+port, protocol): effectively no persistence, best distribution.
- `SOURCE_IP` -- 2-tuple hash (source and destination IP): a client IP sticks to one backend across connections.
- `SOURCE_IP_PROTOCOL` -- 3-tuple hash (source IP, destination IP, protocol): a client IP sticks per protocol.

### spec.rules[].idleTimeoutInMinutes

`int32` · optional (explicit presence)

The TCP idle timeout in minutes; idle connections beyond it are
dropped. Default: 4. Range: 4 to 100. Raise for long-lived
connections (WebSocket, database pools) -- or enable tcp_reset so
clients at least learn of the drop.

- default: `4`
- rule: {"int32":{"lte":100,"gte":4}}

### spec.rules[].floatingIpEnabled

`bool` · optional (explicit presence)

Floating IP (Direct Server Return): the backend sees the FRONTEND's
IP as the destination instead of its own. Required for SQL Server
AlwaysOn listeners and some clustering schemes; the backend OS must
be configured with a loopback matching the frontend IP.
Default: false.

- default: `false`

### spec.rules[].tcpResetEnabled

`bool` · optional (explicit presence)

Send a TCP reset to both ends when a connection is dropped for idle
timeout, so clients fail fast instead of discovering a dead socket on
next write. Default: false (Azure's default), but production TCP
rules generally want it on.

- default: `false`

### spec.rules[].disableOutboundSnat

`bool` · optional (explicit presence)

Disable the default outbound SNAT this rule would otherwise provide
to its backend pool. Set it when the pool's egress is handled by an
explicit outbound rule (required to combine the two on one pool) or
by a NAT gateway -- implicit SNAT has a small, exhaustion-prone port
budget. Default: false.

- default: `false`

### spec.natRules

`[]AzureLoadBalancerNatRule`

Inbound NAT rules forwarding frontend ports to individual backend
instances. Two modes per rule: a single-target rule (`frontend_port`)
whose attachment is completed from the member side (a NIC's NAT-rule
association referencing `status.outputs.nat_rule_ids.<name>`), or a
pool-style rule (`frontend_port_start`/`_end` + `backend_pool_name`)
that gives every pool member its own frontend port -- the modern
replacement for the legacy NAT pool mechanism (which is deliberately
not modeled).

- rule: a NAT rule is single-target (frontend_port) or pool-style (backend_pool_name + frontend_port_start/end) -- set exactly one mode
- rule: pool-style NAT rules require both frontend_port_start and frontend_port_end (and single-target rules must not set them)

### spec.natRules[].name

`string` · required

The name of the NAT rule, unique within the load balancer. A
single-target rule's ARM ID is exported as
`status.outputs.nat_rule_ids.<name>` for the member-side NIC
association that completes the attachment.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.natRules[].frontendIpConfigurationName

`string`

The frontend this NAT rule listens on, by the frontend's `name`.
Optional when the load balancer has exactly one frontend; required
with several (spec-level validation enforces both).

### spec.natRules[].protocol

`enum`

The transport protocol for the forwarded traffic.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_load_balancer_transport_protocol_unspecified` -- Not specified -- invalid; every rule declares its protocol.
- `TCP` -- TCP traffic.
- `UDP` -- UDP traffic.
- `ALL` -- Both TCP and UDP. On a load-balancing rule this creates an HA-ports rule (internal STANDARD frontends only, ports 0); on an outbound rule it SNATs all traffic.

### spec.natRules[].frontendPort

`int32`

SINGLE-TARGET mode: the one frontend port to forward. The rule then
attaches to exactly one backend NIC from the member side (the NIC's
NAT-rule association). Mutually exclusive with the pool-style fields.

- rule: {"int32":{"lte":65534,"gte":0}}

### spec.natRules[].backendPort

`int32` · required

The port on the backend instance that receives the forwarded traffic.

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.natRules[].backendPoolName

`string`

POOL-STYLE mode: the backend pool whose members each get a dedicated
frontend port from [frontend_port_start, frontend_port_end] -- e.g.
per-instance SSH across a scale set. Mutually exclusive with
frontend_port; requires both range bounds.

### spec.natRules[].frontendPortStart

`int32`

The first frontend port of the pool-style range (inclusive).

- rule: {"int32":{"lte":65534,"gte":0}}

### spec.natRules[].frontendPortEnd

`int32`

The last frontend port of the pool-style range (inclusive). The range
must be at least as large as the pool's member count.

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.natRules[].floatingIpEnabled

`bool` · optional (explicit presence)

Floating IP (Direct Server Return) for the forwarded traffic --
the backend sees the frontend's IP as destination. Default: false.

- default: `false`

### spec.natRules[].tcpResetEnabled

`bool` · optional (explicit presence)

Send a TCP reset on idle-timeout drop so both ends fail fast.
Default: false.

- default: `false`

### spec.natRules[].idleTimeoutInMinutes

`int32` · optional (explicit presence)

The TCP idle timeout in minutes for forwarded connections.
Default: 4. Range: 4 to 30 (tighter than load-balancing rules).

- default: `4`
- rule: {"int32":{"lte":30,"gte":4}}

### spec.outboundRules

`[]AzureLoadBalancerOutboundRule`

Outbound rules configuring explicit SNAT: which frontend public IPs
egress traffic from a backend pool uses, and how many SNAT ports each
instance gets. Only legal with PUBLIC frontends on STANDARD SKU.
Production pools that egress should either use an explicit outbound
rule (set disable_outbound_snat on the load-balancing rules for the
same pool) or a NAT gateway on the subnet.

### spec.outboundRules[].name

`string` · required

The name of the outbound rule, unique within the load balancer.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.outboundRules[].frontendIpConfigurationNames

`[]string` · required

The PUBLIC frontends whose addresses egress traffic is SNATed to, by
frontend `name`. Multiple frontends multiply the available SNAT port
budget. At least one.

- rule: {"repeated":{"minItems":"1"}}

### spec.outboundRules[].backendPoolName

`string` · required

The backend pool whose members' outbound traffic this rule governs,
by pool `name`.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.outboundRules[].protocol

`enum`

The transport protocol the SNAT applies to. ALL covers both TCP and
UDP -- the usual choice.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_load_balancer_transport_protocol_unspecified` -- Not specified -- invalid; every rule declares its protocol.
- `TCP` -- TCP traffic.
- `UDP` -- UDP traffic.
- `ALL` -- Both TCP and UDP. On a load-balancing rule this creates an HA-ports rule (internal STANDARD frontends only, ports 0); on an outbound rule it SNATs all traffic.

### spec.outboundRules[].allocatedOutboundPorts

`int32` · optional (explicit presence)

SNAT ports allocated to EACH backend instance. Default: 1024. Set 0
to let Azure divide the frontend's port budget (64,000 per frontend
IP) evenly across the pool's current size -- convenient but
reallocation churns connections as the pool scales; production pools
should size it explicitly (budget / max instances, in multiples of 8).

- default: `1024`
- rule: {"int32":{"gte":0}}

### spec.outboundRules[].tcpResetEnabled

`bool` · optional (explicit presence)

Send a TCP reset on idle-timeout drop for outbound connections.
Default: false.

- default: `false`

### spec.outboundRules[].idleTimeoutInMinutes

`int32` · optional (explicit presence)

The idle timeout in minutes for outbound connections. Default: 4.

- default: `4`
- rule: {"int32":{"gte":4}}

### spec.tags

`map<string, string>`

Free-form tags applied to the load balancer, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Validation Rules

- `lb_global_tier_requires_standard_sku`: GLOBAL sku_tier (cross-region load balancing) is only supported on the STANDARD SKU
- `lb_gateway_pools_require_tunnel_interfaces`: GATEWAY SKU requires every backend pool to declare tunnel_interfaces (the NVA chaining contract)
- `lb_tunnel_interfaces_require_gateway`: tunnel_interfaces are only supported on the GATEWAY SKU
- `lb_rule_multi_pool_requires_gateway`: a load-balancing rule may target two backend pools only on the GATEWAY SKU
- `lb_rule_pool_names_declared`: every rule's backend_pool_names must match pools declared in backend_pools
- `lb_rule_probe_declared`: a rule's probe_name must match a probe declared in health_probes
- `lb_rule_frontend_declared`: a rule's frontend_ip_configuration_name must match a declared frontend (optional only when exactly one frontend exists)
- `lb_nat_rule_frontend_declared`: a NAT rule's frontend_ip_configuration_name must match a declared frontend (optional only when exactly one frontend exists)
- `lb_nat_rule_pool_declared`: a pool-style NAT rule's backend_pool_name must match a pool declared in backend_pools
- `lb_outbound_rule_pool_declared`: an outbound rule's backend_pool_name must match a pool declared in backend_pools
- `lb_outbound_rule_frontends_declared`: every outbound rule frontend name must match a declared frontend

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | The Azure Resource Manager ID of the load balancer. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/loadBalancers/{name} |
| `status.outputs.load_balancer_name` | `string` | The name of the load balancer. |
| `status.outputs.private_ip_address` | `string` | The first internal frontend's private IP address -- the address DNS records for an internal load balancer point at. Empty when every frontend is public (a public frontend's address lives on its referenced AzurePublicIp resource). |
| `status.outputs.private_ip_addresses` | `[]string` | The private IP addresses of ALL internal frontends, in declaration order. |
| `status.outputs.frontend_ip_configuration_ids` | `map<string, string>` | The ARM ID of each frontend IP configuration, keyed by the frontend's name. Referenced when chaining a frontend behind a Gateway load balancer or registering a regional frontend in a GLOBAL-tier pool. Example valueFrom fieldPath: status.outputs.frontend_ip_configuration_ids.public |
| `status.outputs.backend_pool_ids` | `map<string, string>` | The ARM ID of each backend address pool, keyed by the pool's name. THE membership seam: NIC ip_configurations and scale-set network profiles reference a pool ID here to join the pool. Example valueFrom fieldPath: status.outputs.backend_pool_ids.web |
| `status.outputs.probe_ids` | `map<string, string>` | The ARM ID of each health probe, keyed by the probe's name. A virtual machine scale set's rolling-upgrade health_probe_id references one. Example valueFrom fieldPath: status.outputs.probe_ids.http-health |
| `status.outputs.nat_rule_ids` | `map<string, string>` | The ARM ID of each inbound NAT rule, keyed by the rule's name. A NIC's NAT-rule association references a single-target rule here to complete the attachment. Example valueFrom fieldPath: status.outputs.nat_rule_ids.ssh-admin |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.frontendIpConfigurations[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.frontendIpConfigurations[].publicIpAddressId` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.frontendIpConfigurations[].publicIpPrefixId` | AzurePublicIpPrefix | `status.outputs.public_ip_prefix_id` |
| `spec.backendPools[].virtualNetworkId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureNetworkInterface | `spec.ipConfigurations[].loadBalancerBackendAddressPoolIds` | `status.outputs.backend_pool_ids` |
| AzureNetworkInterface | `spec.ipConfigurations[].loadBalancerInboundNatRuleIds` | `status.outputs.nat_rule_ids` |
| AzureVirtualMachineScaleSet | `spec.networkInterfaces[].ipConfigurations[].loadBalancerBackendAddressPoolIds` | `status.outputs.backend_pool_ids` |
| AzureVirtualMachineScaleSet | `spec.networkInterfaces[].ipConfigurations[].loadBalancerInboundNatRuleIds` | `status.outputs.nat_rule_ids` |
| AzureVirtualMachineScaleSet | `spec.upgradePolicy.healthProbeId` | `status.outputs.probe_ids` |

## See Also

- [Overview](../README.md)
