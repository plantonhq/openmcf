# OpenStack Floating IP Associate

Associates an existing OpenStack Neutron floating IP with a port, providing external connectivity to the resource attached to that port. This is a join resource -- it binds a pre-allocated floating IP (from OpenStackFloatingIp) to a port (from OpenStackNetworkPort) as a separate DAG node. Supports ValueFromRef wiring for both the floating IP address and port ID in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Floating IP Association** -- a Neutron floating IP association that maps the specified floating IP address to a fixed IP on the target port, enabling inbound and outbound external connectivity
- **OpenStack Tags** -- resource metadata applied automatically for tracking

## Before You Deploy

### OpenStack Account

- **A floating IP** already allocated in Neutron. Provide the IP address directly (e.g., `203.0.113.42`) or reference an OpenStackFloatingIp Cloud Resource via ValueFromRef. The association targets the address string, not the floating IP UUID.
- **A port** with at least one fixed IP address. Provide the port UUID directly or reference an OpenStackNetworkPort Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **OpenStack Floating IP Associate**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate the required fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackFloatingIpAssociate
metadata:
  name: web-server-fip
  org: acme-corp
  env: prod
spec:
  floatingIp:
    value: "203.0.113.42"
  portId:
    value: "<port-uuid>"
```

```shell
planton apply -f floating-ip-associate.yaml
```

This binds the floating IP `203.0.113.42` to the specified port. If the port has a single fixed IP, the floating IP maps to it automatically. No `fixedIp` override is needed for single-IP ports.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the association to a floating IP and port deployed in the same InfraPipeline:

```yaml
spec:
  floatingIp:
    valueFrom:
      kind: OpenStackFloatingIp
      name: web-fip
      fieldPath: status.outputs.address
  portId:
    valueFrom:
      kind: OpenStackNetworkPort
      name: web-port
      fieldPath: status.outputs.port_id
```

The InfraPipeline resolves the dependency graph, deploys the floating IP and port first, then creates the association with the resolved values.

## Key Configuration

These are the most important decisions when configuring a floating IP association. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Floating IP vs. port lifecycle** -- This component decouples floating IP allocation from association, making them separate DAG nodes. Use this when the floating IP and port are managed by different teams or when you need to re-point a floating IP between ports. For simple cases where allocation and association happen together, use OpenStackFloatingIp with its built-in `portId` field instead.

**Fixed IP selection** -- When the target port has multiple fixed IP addresses, set `fixedIp` to specify which one the floating IP maps to. If omitted, OpenStack uses the first fixed IP on the port. For single-IP ports, this field is unnecessary.

**Immutability** -- Changing `floatingIp` or `region` recreates the association. Only `portId` can be updated in place, allowing you to re-point a floating IP from one port to another without re-allocating the IP address.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackFloatingIp** | `floatingIp` | `status.outputs.address` |
| **OpenStackNetworkPort** | `portId` | `status.outputs.port_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values for operational visibility:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Resource identifier for the association | Terraform state reference |
| `floating_ip` | The floating IP address that was associated | DNS records, monitoring dashboards |
| `port_id` | UUID of the port the floating IP is bound to | Debugging, audit logs |
| `fixed_ip` | Fixed IP on the port mapped to the floating IP | Network debugging, firewall rules |
| `region` | OpenStack region of the association | Region-aware operations |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard association** -- Binds a pre-allocated floating IP to a port with automatic fixed IP selection. Covers the common case of giving external access to an instance or load balancer. Start from the **Standard** preset.

## Works With

- [**OpenStack Floating IP**](/cloud-catalog/openstack-floating-ip) -- provides the floating IP address to associate
- [**OpenStack Network Port**](/cloud-catalog/openstack-network-port) -- provides the port ID to bind the floating IP to