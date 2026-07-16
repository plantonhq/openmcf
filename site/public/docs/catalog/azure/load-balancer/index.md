---
title: "Load Balancer"
description: "Load Balancer deployment documentation"
icon: "package"
order: 100
componentName: "azureloadbalancer"
---

# Azure Load Balancer

Deploys an Azure Load Balancer with its complete traffic surface: frontend IP configurations (public and/or internal), backend address pools, health probes, load-balancing rules, inbound NAT rules, and outbound (SNAT) rules -- configured as one unit because none of these sub-resources has a life outside its load balancer. What IS independent -- pool membership -- is expressed from the member side via the exported per-name pool IDs.

## What Gets Created

When you deploy an AzureLoadBalancer resource, Planton provisions:

- **Load Balancer** -- the `Microsoft.Network/loadBalancers` resource in the chosen SKU (STANDARD default, GATEWAY for NVA chaining) and tier (REGIONAL default, GLOBAL for cross-region), carrying every declared frontend IP configuration
- **Backend Address Pools** -- one per entry in `backendPools`, with optional virtual-network scoping, synchronous mode, Gateway tunnel interfaces, and inline IP-based member addresses
- **Health Probes** -- one per entry in `healthProbes`: TCP, HTTP, or HTTPS with configurable interval, failure count, and recovery threshold
- **Load Balancing Rules** -- one per entry in `rules`, mapping a frontend port/protocol to a backend pool and port, with session persistence, TCP reset, floating IP, and SNAT control
- **Inbound NAT Rules** -- one per entry in `natRules`: single-target port forwards (completed by a NIC-side association) or pool-style port ranges (one frontend port per pool member)
- **Outbound Rules** -- one per entry in `outboundRules`: explicit SNAT through public frontends with a deliberately sized per-instance port budget
- **Azure Tags** -- Planton-derived metadata tags merged with your `tags` (your keys win on collision)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the load balancer will be created (can reference an AzureResourceGroup resource)
- **Frontend addresses** -- a Standard SKU public IP or public IP prefix for public frontends; a VNet subnet for internal frontends (each referenceable from AzurePublicIp / AzurePublicIpPrefix / AzureSubnet resources)
- **Backend members** -- network interfaces or scale sets that join the pools after deployment by referencing the exported pool IDs

## Quick Start

Create a file `loadbalancer.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLoadBalancer
metadata:
  name: my-lb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureLoadBalancer.my-lb
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: my-lb
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/publicIPAddresses/my-pip
  backendPools:
    - name: web
  healthProbes:
    - name: http-health
      protocol: PROBE_HTTP
      port: 80
      requestPath: /healthz
  rules:
    - name: http
      protocol: TCP
      frontendPort: 80
      backendPort: 80
      backendPoolNames: [web]
      probeName: http-health
```

Deploy:

```shell
planton apply -f loadbalancer.yaml
```

This creates a public Standard load balancer with one backend pool, an HTTP health probe, and a TCP rule forwarding port 80. Network interfaces join the `web` pool by referencing `status.outputs.backend_pool_ids.web`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region (must match backend resources). | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name; can reference an AzureResourceGroup. | Required |
| `name` | `string` | Load balancer name, unique within the resource group. | Required, 1-80 characters |
| `frontendIpConfigurations` | `list` | The frontends that receive traffic. Each is public (`publicIpAddressId` or `publicIpPrefixId`) or internal (`subnetId`) -- at most one address source per frontend. | Required, minimum 1 |
| `frontendIpConfigurations[].name` | `string` | Frontend label; rules target frontends by this name. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sku` | `enum` | `STANDARD` | `STANDARD` (production) or `GATEWAY` (NVA chaining; every pool must declare `tunnelInterfaces`, and the subscription needs the `Microsoft.Network/AllowGatewayLoadBalancer` feature). Basic is not modeled (retired by Azure September 2025). |
| `skuTier` | `enum` | `REGIONAL` | `GLOBAL` creates a cross-region load balancer whose pool members are regional load balancers' frontends. Requires `STANDARD`. |
| `edgeZone` | `string` | | Azure Edge Zone pinning. Fixed at creation. |
| `frontendIpConfigurations[].privateIpAddress` | `string` | _(dynamic)_ | Pin a static private address (internal frontends only). |
| `frontendIpConfigurations[].privateIpAddressVersion` | `enum` | `IPV4` | `IPV6` for an IPv6 internal frontend (dual-stack subnet required). |
| `frontendIpConfigurations[].zones` | `string[]` | | Availability zones for an internal frontend's address (`["1","2","3"]` for zone redundancy). A public frontend's zone posture comes from its public IP resource. |
| `backendPools[].virtualNetworkId` | `StringValueOrRef` | | The pool's virtual network; required for IP-based `addresses` or `synchronousMode`. |
| `backendPools[].synchronousMode` | `enum` | | `AUTOMATIC` or `MANUAL` IP-based membership sync (STANDARD SKU, vnet-scoped pools). |
| `backendPools[].tunnelInterfaces` | `list` | | GATEWAY SKU: VXLAN tunnels (identifier/port/protocol/type) to the pool's NVAs. |
| `backendPools[].addresses` | `list` | | Inline IP-based members (`ipAddress`), or regional LB frontends for GLOBAL-tier pools (`loadBalancerFrontendIpConfigurationId`). |
| `healthProbes[].protocol` | `enum` | `PROBE_TCP` | `PROBE_HTTP`/`PROBE_HTTPS` GET `requestPath` and require HTTP 200 (path required for them, forbidden for TCP). |
| `healthProbes[].intervalInSeconds` | `int` | `15` | Seconds between probes (min 5). |
| `healthProbes[].numberOfProbes` | `int` | `2` | Consecutive failures before unhealthy. |
| `healthProbes[].probeThreshold` | `int` | `1` | Consecutive successes before a recovered instance is re-admitted (1-100). |
| `rules[].frontendIpConfigurationName` | `string` | _(sole frontend)_ | Which frontend the rule listens on; required when several frontends exist. |
| `rules[].backendPoolNames` | `string[]` | | Target pool(s) by name (two only on GATEWAY SKU). Required, 1-2 items. |
| `rules[].probeName` | `string` | | Gating probe by name. Optional, but production rules should probe. |
| `rules[].loadDistribution` | `enum` | `DEFAULT` | Session persistence: `SOURCE_IP` (2-tuple) or `SOURCE_IP_PROTOCOL` (3-tuple). |
| `rules[].idleTimeoutInMinutes` | `int` | `4` | TCP idle timeout (4-100). |
| `rules[].floatingIpEnabled` | `bool` | `false` | Direct Server Return (SQL AlwaysOn, HA clustering). |
| `rules[].tcpResetEnabled` | `bool` | `false` | Send TCP reset on idle drop so clients fail fast. |
| `rules[].disableOutboundSnat` | `bool` | `false` | Disable the rule's implicit SNAT (required when the pool uses an explicit outbound rule). |
| `natRules[]` | `list` | | Single-target (`frontendPort`) or pool-style (`backendPoolName` + `frontendPortStart`/`End`) port forwarding; `idleTimeoutInMinutes` range 4-30. |
| `outboundRules[]` | `list` | | Explicit SNAT: `frontendIpConfigurationNames` (public frontends), `backendPoolName`, `protocol`, `allocatedOutboundPorts` (default 1024). |
| `tags` | `map<string,string>` | | Free-form tags merged over the Planton-derived tags (yours win). |

## Examples

### Internal Load Balancer with a Zone-Redundant Static Frontend

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLoadBalancer
metadata:
  name: internal-lb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLoadBalancer.internal-lb
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  name: internal-lb
  frontendIpConfigurations:
    - name: internal
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app
      privateIpAddress: "10.0.1.100"
      zones: ["1", "2", "3"]
  backendPools:
    - name: api
  healthProbes:
    - name: tcp-health
      protocol: PROBE_TCP
      port: 8080
  rules:
    - name: api
      protocol: TCP
      frontendPort: 8080
      backendPort: 8080
      backendPoolNames: [api]
      probeName: tcp-health
      idleTimeoutInMinutes: 30
      tcpResetEnabled: true
```

### Explicit Outbound SNAT with Admin NAT Rules

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLoadBalancer
metadata:
  name: egress-lb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLoadBalancer.egress-lb
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: egress-lb
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/publicIPAddresses/egress-pip
  backendPools:
    - name: web
  healthProbes:
    - name: http-health
      protocol: PROBE_HTTP
      port: 80
      requestPath: /healthz
  rules:
    - name: http
      protocol: TCP
      frontendPort: 80
      backendPort: 80
      backendPoolNames: [web]
      probeName: http-health
      disableOutboundSnat: true
  natRules:
    - name: ssh-admin
      protocol: TCP
      frontendPort: 2222
      backendPort: 22
  outboundRules:
    - name: egress
      frontendIpConfigurationNames: [public]
      backendPoolName: web
      protocol: ALL
      allocatedOutboundPorts: 2048
```

### SQL AlwaysOn with Floating IP

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLoadBalancer
metadata:
  name: sql-lb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLoadBalancer.sql-lb
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: sql-lb
  frontendIpConfigurations:
    - name: listener
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/data
      privateIpAddress: "10.0.2.50"
  backendPools:
    - name: sql-nodes
  healthProbes:
    - name: sql-health
      protocol: PROBE_TCP
      port: 59999
      intervalInSeconds: 5
  rules:
    - name: sql
      protocol: TCP
      frontendPort: 1433
      backendPort: 1433
      backendPoolNames: [sql-nodes]
      probeName: sql-health
      floatingIpEnabled: true
      disableOutboundSnat: true
```

### Using Foreign Key References

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLoadBalancer
metadata:
  name: ref-lb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureLoadBalancer.ref-lb
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  name: ref-lb
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        valueFrom:
          kind: AzurePublicIp
          name: my-pip
          fieldPath: status.outputs.public_ip_id
  backendPools:
    - name: web
  healthProbes:
    - name: http-health
      protocol: PROBE_HTTP
      port: 80
      requestPath: /healthz
  rules:
    - name: http
      protocol: TCP
      frontendPort: 80
      backendPort: 80
      backendPoolNames: [web]
      probeName: http-health
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `load_balancer_id` | `string` | Azure Resource Manager ID of the load balancer |
| `load_balancer_name` | `string` | Name of the load balancer |
| `private_ip_address` | `string` | First internal frontend's private address (empty when every frontend is public) |
| `private_ip_addresses` | `string[]` | All internal frontends' private addresses |
| `frontend_ip_configuration_ids` | `map<string,string>` | Frontend IDs keyed by frontend name (gateway chaining, GLOBAL-tier pool members) |
| `backend_pool_ids` | `map<string,string>` | Pool IDs keyed by pool name -- reference `status.outputs.backend_pool_ids.<pool>` from a NIC or scale set to join |
| `probe_ids` | `map<string,string>` | Probe IDs keyed by probe name -- a scale set's rolling-upgrade health probe references one |
| `nat_rule_ids` | `map<string,string>` | Inbound NAT rule IDs keyed by rule name -- a NIC's NAT-rule association completes a single-target rule |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/resource-group) -- the resource group for load balancer placement
- [AzurePublicIp](/docs/catalog/azure/public-ip) -- Standard SKU public IPs for public frontends
- [AzurePublicIpPrefix](/docs/catalog/azure/public-ip-prefix) -- reserved contiguous ranges for SNAT-scaling frontends
- [AzureSubnet](/docs/catalog/azure/subnet) -- subnets for internal frontends
- [AzureNetworkInterface](/docs/catalog/azure/network-interface) -- joins backend pools and completes NAT rules from the member side
- [AzureDnsRecord](/docs/catalog/azure/dns-record) -- DNS records pointing at frontend addresses
