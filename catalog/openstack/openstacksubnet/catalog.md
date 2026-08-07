# OpenStack Subnet

Deploys a Neutron subnet on OpenStack, defining an IP address range within a network with DHCP, DNS, gateway, and allocation pool configuration. The subnet references a parent OpenStackNetwork via a foreign key and supports ValueFromRef wiring for InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Subnet** -- an IPv4 or IPv6 subnet attached to the specified network, with configurable CIDR, gateway, DHCP settings, DNS nameservers, and allocation pools
- **OpenStack Tags** -- user-defined tags applied to the subnet for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Parent network** -- the subnet must belong to exactly one network. Provide the network ID directly or reference an OpenStackNetwork Cloud Resource via ValueFromRef.
- **CIDR planning** -- choose a CIDR block that does not overlap with other subnets on the same network. Common choices are `192.168.x.0/24` for workload subnets and `10.0.x.0/24` for infrastructure subnets.
- **DNS servers** -- decide which DNS nameservers to push to instances via DHCP. Public resolvers (`8.8.8.8`, `8.8.4.4`) or internal resolvers are typical choices.

## Deploy

### Console

Open the deployment store, find **OpenStack Subnet**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard DHCP** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackSubnet
metadata:
  name: workload-subnet
  org: acme-corp
  env: prod
spec:
  networkId:
    value: "<network-id>"
  cidr: "192.168.1.0/24"
```

```shell
planton apply -f subnet.yaml
```

This creates an IPv4 subnet with DHCP enabled and an auto-assigned gateway (first usable IP). No custom DNS nameservers or allocation pools are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the subnet to a network deployed in the same InfraPipeline:

```yaml
spec:
  networkId:
    valueFrom:
      kind: OpenStackNetwork
      name: app-network
      fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, deploys the network first, then provisions the subnet with the resolved network ID.

## Key Configuration

These are the most important decisions when configuring a subnet. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CIDR and IP version** -- `cidr` defines the address range (e.g., `192.168.1.0/24` for IPv4, `2001:db8::/64` for IPv6). Set `ipVersion` to `4` (default) or `6`. The CIDR cannot be changed after creation without replacing the subnet.

**Gateway** -- By default, OpenStack assigns the first usable IP as the gateway. Set `gatewayIp` to override, or `noGateway: true` for isolated subnets that do not need routing (e.g., storage replication networks). Gateway and no-gateway are mutually exclusive.

**DHCP** -- `enableDhcp` defaults to `true`. Disable for subnets where IPs are assigned statically (e.g., via port fixed_ips or instance configuration). When enabled, the Neutron DHCP agent assigns IPs from the allocation pool.

**Allocation pools** -- `allocationPools` restricts which IPs within the CIDR are assignable via DHCP. Use this to reserve ranges for static assignments, VIPs, or infrastructure addresses (e.g., allocate only `192.168.1.100` to `192.168.1.200`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackNetwork** | `networkId` | `status.outputs.network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `subnet_id` | UUID of the subnet in OpenStack | Router interfaces, load balancer VIP placement, container cluster templates |
| `name` | Name of the subnet | DNS records, monitoring labels |
| `cidr` | CIDR block of the subnet | Security group rules, network ACLs |
| `gateway_ip` | Gateway IP address (empty if no_gateway) | Instance routing configuration |
| `network_id` | ID of the parent network | Convenience reference for downstream resources |
| `region` | OpenStack region where the subnet was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard DHCP subnet** -- IPv4 `/24` subnet with DHCP enabled and public DNS servers. The most common configuration for workload connectivity -- instances receive IPs automatically via DHCP. Start from the **Standard DHCP** preset.

**Isolated no-gateway subnet** -- A subnet with no gateway and DHCP disabled. Designed for backend networks (storage replication, database clusters, heartbeat links) where traffic stays within the Layer 2 domain and IPs are statically assigned. Start from the **Isolated No Gateway** preset.

## Works With

- [**OpenStack Network**](/cloud-catalog/openstack-network) -- provides the parent network where the subnet is created