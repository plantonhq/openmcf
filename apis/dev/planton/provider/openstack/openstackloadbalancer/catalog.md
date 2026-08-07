# OpenStack Load Balancer

Deploys an Octavia load balancer on OpenStack with a Virtual IP (VIP) allocated on a specified subnet. The load balancer serves as the entry point for all Octavia traffic distribution -- attach listeners, pools, members, and monitors to complete the traffic path. ValueFromRef wiring connects the VIP subnet from an OpenStackSubnet Cloud Resource in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Octavia Load Balancer** -- a load balancer instance with a VIP address allocated from the specified subnet, acting as the client-facing endpoint for traffic distribution
- **VIP Port** -- a Neutron port created automatically for the Virtual IP, which can be used to attach floating IPs or security groups
- **OpenStack Tags** -- user-defined tags applied to the load balancer for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Subnet** -- a subnet with available IP addresses where the VIP will be allocated. Provide the subnet ID directly or reference an OpenStackSubnet Cloud Resource via ValueFromRef.
- **Octavia service** -- the Octavia load balancing service must be enabled in your OpenStack deployment. Run `openstack loadbalancer list` to verify availability.
- **Flavor** (optional) -- if your deployment defines Octavia flavors for custom bandwidth or connection limits, note the flavor ID for the `flavorId` field.

## Deploy

### Console

Open the deployment store, find **OpenStack Load Balancer**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Load Balancer** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  vipSubnetId:
    value: "<subnet-id>"
```

```shell
planton apply -f loadbalancer.yaml
```

This creates a load balancer with an auto-assigned VIP on the specified subnet. No flavor, custom VIP address, or tags are configured. Attach listeners and pools to begin routing traffic.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the load balancer to a subnet deployed in the same InfraPipeline:

```yaml
spec:
  vipSubnetId:
    valueFrom:
      kind: OpenStackSubnet
      name: app-subnet
      fieldPath: status.outputs.subnet_id
```

The InfraPipeline resolves the dependency graph, deploys the subnet first, then provisions the load balancer with the resolved subnet ID.

## Key Configuration

These are the most important decisions when configuring a load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**VIP address allocation** -- By default, Octavia auto-assigns an available IP from the subnet. Set `vipAddress` to request a specific IP within the subnet's CIDR range. Changing the VIP address requires recreating the load balancer.

**Flavor selection** -- Leave `flavorId` empty to use Octavia's default resource limits. Set it when your deployment provides custom flavors with specific bandwidth caps or connection limits. Changing the flavor requires recreating the load balancer.

**Administrative state** -- The load balancer is active by default (`adminStateUp: true`). Set to `false` to provision the load balancer in a disabled state, useful for staged deployments where listeners and pools need configuration before accepting traffic.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackSubnet** | `vipSubnetId` | `status.outputs.subnet_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `loadbalancer_id` | UUID of the load balancer | Listener attachment via `loadbalancerId` |
| `name` | Name of the load balancer | Monitoring labels, resource identification |
| `vip_address` | Virtual IP address allocated to the load balancer | DNS records, client configuration |
| `vip_port_id` | Neutron port ID of the VIP | Floating IP associations, security group attachments |
| `vip_subnet_id` | Subnet where the VIP was allocated | Reference for downstream subnet-aware resources |
| `region` | OpenStack region where the load balancer was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard Load Balancer** -- Creates an Octavia load balancer with an auto-assigned VIP on the specified subnet. The entry point for building a complete load balancing stack (LB, Listener, Pool, Members, Monitor). Start from the **Standard Load Balancer** preset.

## Works With

- [**OpenStack Subnet**](/cloud-catalog/openstack-subnet) -- provides the subnet ID where the load balancer's VIP address is allocated