# OpenStack Floating IP

Allocates a floating IP from an external (provider) network on OpenStack, providing public connectivity to instances or ports on tenant networks. The floating IP can be allocated without association (for DNS pre-configuration or IP reservation) or immediately bound to a port. ValueFromRef wiring connects the external network and optional port from their respective Cloud Resources in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Floating IP** -- a public IP address allocated from the specified external network, with optional immediate association to a port
- **Port Association** -- created only when `portId` is provided; binds the floating IP to a port for immediate external connectivity
- **OpenStack Tags** -- user-defined tags applied to the floating IP for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **External network** -- a provider (external) network with available floating IP addresses. Obtain the external network ID from your OpenStack administrator or reference an OpenStackNetwork Cloud Resource via ValueFromRef.
- **Target port** (optional) -- if associating the floating IP at creation time, have the port ID ready or reference an OpenStackNetworkPort Cloud Resource via ValueFromRef. For DAG-visible association in InfraCharts, allocate without `portId` and use the separate OpenStackFloatingIpAssociate component.
- **Router with external gateway** -- the tenant subnet where the target port resides must be connected (via a router interface) to a router that has the same external network as its gateway. Without this routing path, the floating IP association succeeds but traffic does not flow.

## Deploy

### Console

Open the deployment store, find **OpenStack Floating IP**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Allocation Only** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackFloatingIp
metadata:
  name: web-fip
  org: acme-corp
  env: prod
spec:
  floatingNetworkId:
    value: "<external-network-id>"
```

```shell
planton apply -f floating-ip.yaml
```

This allocates a floating IP from the external network with an auto-assigned address. The IP is reserved but not associated with any port. No specific address, tags, or port association is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the floating IP to an external network and a port deployed in the same InfraPipeline:

```yaml
spec:
  floatingNetworkId:
    valueFrom:
      kind: OpenStackNetwork
      name: provider-net
      fieldPath: status.outputs.network_id
  portId:
    valueFrom:
      kind: OpenStackNetworkPort
      name: web-port
      fieldPath: status.outputs.port_id
```

The InfraPipeline resolves the dependency graph, deploys the external network and port first, then provisions the floating IP with the resolved IDs.

## Key Configuration

These are the most important decisions when configuring a floating IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Allocation vs association** -- Omit `portId` for allocation-only mode -- the IP is reserved but unbound. Set `portId` for immediate association. For InfraChart deployments where allocation and association should be separate DAG nodes, use allocation-only mode and the OpenStackFloatingIpAssociate component.

**Specific address** -- By default, OpenStack auto-assigns any available address from the external network. Set `address` to request a specific IP when you need a predictable address for DNS pre-configuration, firewall whitelisting, or IP reservation. This is a create-time setting.

**Fixed IP selection** -- `fixedIp` is only relevant when `portId` is set and the target port has multiple IP addresses. It selects which fixed IP on the port receives the floating IP mapping. If the port has a single IP, this can be omitted.

**Subnet selection** -- `subnetId` restricts allocation to a specific subnet on the external network. If omitted, OpenStack allocates from any available subnet. Use this when the external network has multiple subnets with different IP ranges.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackNetwork** | `floatingNetworkId` | `status.outputs.network_id` |
| **OpenStackNetworkPort** (optional) | `portId` | `status.outputs.port_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `floating_ip_id` | UUID of the floating IP resource | Import, audit, and debugging |
| `address` | Allocated floating IP address (e.g., `203.0.113.42`) | DNS records, firewall rules, client configuration |
| `floating_network_id` | UUID of the external network the IP was allocated from | Network topology reference |
| `port_id` | UUID of the associated port (empty if allocation-only) | Association verification |
| `fixed_ip` | Fixed IP address mapped to the floating IP (empty if unassociated) | Debugging, network tracing |
| `region` | OpenStack region where the floating IP was allocated | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Allocation only** -- Reserves a floating IP without binding it to a port. Use for DNS pre-configuration, firewall whitelisting, or InfraChart deployments where association is a separate DAG node via OpenStackFloatingIpAssociate. Start from the **Allocation Only** preset.

**With port association** -- Allocates a floating IP and immediately binds it to an existing port. The simplest approach for standalone deployments where a single instance needs a public IP. Start from the **With Port Association** preset.

## Works With

- [**OpenStack Network**](/cloud-catalog/openstack-network) -- provides the external (provider) network from which the floating IP is allocated
- [**OpenStack Network Port**](/cloud-catalog/openstack-network-port) -- provides the port ID for immediate floating IP association