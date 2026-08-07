# OpenStack Network

Deploys a Neutron network on OpenStack -- the foundational Layer 2 broadcast domain that subnets, ports, routers, and instances attach to. The network supports configurable MTU, port security, DNS integration, and tenant/shared/external modes.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Network** -- an isolated Layer 2 network with configurable admin state, MTU, port security, DNS domain, and shared/external flags
- **OpenStack Tags** -- user-defined tags applied to the network for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Project access** -- the network is created in the project associated with the OpenStack credentials. Shared or external networks require admin privileges.
- **Network type** -- the underlying network type (VXLAN, VLAN, flat) is determined by the Neutron backend configuration, not by this spec. If your deployment uses a specific network type, confirm the MTU setting matches (e.g., 1450 for VXLAN overlays).
- **DNS integration** (optional) -- if you plan to use `dnsDomain` for automatic DNS name assignment on ports, confirm that the Neutron deployment has the dns-integration extension enabled.

## Deploy

### Console

Open the deployment store, find **OpenStack Network**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Tenant** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackNetwork
metadata:
  name: app-network
  org: acme-corp
  env: prod
spec: {}
```

```shell
planton apply -f network.yaml
```

This creates a tenant network with all OpenStack defaults: admin state up, standard MTU, port security enabled, not shared, and not external. Attach subnets to assign IP ranges before launching instances.

## Key Configuration

These are the most important decisions when configuring a network. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Admin state** -- `adminStateUp` defaults to `true`. Set to `false` to create the network in a disabled state (useful for pre-staging infrastructure before a maintenance window).

**Shared vs tenant** -- `shared` makes the network visible to all projects in the OpenStack deployment. Most workloads use tenant-scoped networks (the default). Creating shared networks requires admin privileges.

**External network** -- `external` marks the network as a provider network for floating IP allocation and router gateways. This is an admin-level operation. Tenant workloads consume external networks but do not create them.

**MTU** -- The `mtu` field overrides the default Maximum Transmission Unit. Set to `1450` for VXLAN overlays or `9000` for jumbo frames on flat/VLAN networks. Mismatched MTU values between the network and its physical infrastructure cause packet fragmentation or drops.

**Port security** -- `portSecurityEnabled` controls whether security group enforcement applies to ports on this network. Disable for networks that need promiscuous traffic (e.g., network appliances, VPN gateways, or monitoring interfaces).

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `network_id` | UUID of the network in OpenStack | Subnets, routers, network ports, floating IPs, instance network attachments |
| `name` | Name of the network | DNS records, monitoring labels |
| `region` | OpenStack region where the network was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard tenant network** -- A basic isolated network with all defaults. The starting point for most workloads -- attach subnets, connect a router for external access, and launch instances. Start from the **Standard Tenant** preset.

## Works With

This component operates independently and does not reference other deployment components.