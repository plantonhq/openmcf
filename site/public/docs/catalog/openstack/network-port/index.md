---
title: "Network Port"
description: "Network Port deployment documentation"
icon: "package"
order: 100
componentName: "openstacknetworkport"
---

# OpenStack Network Port

Deploys a Neutron port on OpenStack -- a virtual switch port that provides stable network identity (MAC address, fixed IPs, security groups) for instances, load balancers, or other network-consuming resources. The port supports explicit IP allocation, security group assignment, and ValueFromRef wiring for networks, subnets, and security groups in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Port** -- a virtual network port on the specified network with configurable fixed IPs, security groups, MAC address, port security, and admin state
- **Fixed IP Allocations** -- created only when `fixedIps` entries are specified; each entry assigns an IP address from a subnet on the port's network
- **Security Group Bindings** -- created only when `securityGroupIds` entries are specified; attaches the listed security groups to the port instead of the project default
- **OpenStack Tags** -- user-defined tags applied to the port for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Network** -- the port requires an existing Neutron network. Provide the network UUID directly or reference an OpenStackNetwork Cloud Resource via ValueFromRef.
- **Subnet** (optional) -- if specifying `fixedIps`, each entry can reference a subnet on the port's network. Provide the subnet UUID directly or reference an OpenStackSubnet Cloud Resource via ValueFromRef.
- **Security groups** (optional) -- if omitted and `noSecurityGroups` is `false`, OpenStack applies the project's default security group. Reference OpenStackSecurityGroup Cloud Resources via ValueFromRef for InfraChart wiring, or provide UUIDs directly.

## Deploy

### Console

Open the deployment store, find **OpenStack Network Port**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Fixed IP** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackNetworkPort
metadata:
  name: app-port
  org: acme-corp
  env: prod
spec:
  networkId:
    value: "<network-id>"
  fixedIps:
    - subnetId:
        value: "<subnet-id>"
```

```shell
planton apply -f network-port.yaml
```

This creates a port on the specified network with an IP auto-assigned from the given subnet. The project's default security group is applied. No explicit security groups, MAC address override, or port security settings are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the port to other resources deployed in the same InfraPipeline:

```yaml
spec:
  networkId:
    valueFrom:
      kind: OpenStackNetwork
      name: app-network
      fieldPath: status.outputs.network_id
  fixedIps:
    - subnetId:
        valueFrom:
          kind: OpenStackSubnet
          name: app-subnet
          fieldPath: status.outputs.subnet_id
  securityGroupIds:
    - valueFrom:
        kind: OpenStackSecurityGroup
        name: web-sg
        fieldPath: status.outputs.security_group_id
```

The InfraPipeline resolves the dependency graph, deploys the network, subnet, and security group first, then provisions the port with the resolved values.

## Key Configuration

These are the most important decisions when configuring a network port. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Security group strategy** -- Choose between default security groups (omit `securityGroupIds`), explicit security groups (`securityGroupIds` list), or no security groups (`noSecurityGroups: true`). The three options are mutually exclusive. Use `noSecurityGroups` for load balancer VIPs or network appliance ports that manage their own traffic filtering.

**Fixed IP allocation** -- `fixedIps` entries control IP assignment. Each entry can specify a `subnetId` (for subnet-specific allocation) and optionally an `ipAddress` (for a specific IP). If omitted entirely, OpenStack auto-assigns one IP from any subnet on the network.

**Port security** -- `portSecurityEnabled` controls whether security group enforcement applies to this port. If omitted, the port inherits the network's port security setting. Disable for ports that need promiscuous traffic (network appliances, VPN gateways, monitoring interfaces).

**MAC address** -- `macAddress` specifies a fixed MAC address. If omitted, OpenStack auto-assigns one. This is a ForceNew field -- use it for DPDK, network bonding, or license-tied MAC addresses.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackNetwork** | `networkId` | `status.outputs.network_id` |
| **OpenStackSecurityGroup** (optional) | `securityGroupIds` | `status.outputs.security_group_id` |
| **OpenStackSubnet** (optional) | `fixedIps[].subnetId` | `status.outputs.subnet_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `port_id` | UUID of the port in OpenStack | Instance network attachments, floating IP associations |
| `mac_address` | MAC address assigned to the port | Network configuration, DHCP reservations |
| `all_fixed_ips` | List of all IP addresses assigned to the port | DNS records, load balancer members, monitoring targets |
| `all_security_group_ids` | List of all security group UUIDs applied to the port | Audit, security compliance verification |
| `region` | OpenStack region where the port was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard port with fixed IP** -- A port on a network with an IP auto-assigned from a specific subnet and the project's default security group. The most common configuration for attaching to instances or creating floating IP targets. Start from the **Standard Fixed IP** preset.

**Unrestricted appliance port** -- A port with all security groups removed (`noSecurityGroups: true`). Traffic flows unrestricted through this port. Suitable for load balancer VIPs, network appliance ports, or trunk ports where security group filtering would interfere with traffic. Start from the **No Security Groups** preset.

## Works With

- [**OpenStack Network**](/cloud-catalog/openstack-network) -- provides the network UUID that the port is created on
- [**OpenStack Security Group**](/cloud-catalog/openstack-security-group) -- provides security group IDs for traffic filtering on the port
- [**OpenStack Subnet**](/cloud-catalog/openstack-subnet) -- provides subnet IDs for fixed IP allocation on the port