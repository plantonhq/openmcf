# OpenStack Load Balancer Member

Deploys an Octavia pool member on OpenStack that registers a backend server with a load balancer pool. Each member has an IP address and port that receives traffic according to the pool's load-balancing algorithm. The member supports weighted distribution, administrative draining, and optional subnet specification for cross-subnet L3 routing. ValueFromRef wiring connects the member to an OpenStackLoadBalancerPool and optionally an OpenStackSubnet in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Octavia Pool Member** -- a backend member resource registered with the specified pool, receiving traffic at the configured IP address and port according to the pool's load-balancing algorithm
- **Subnet Routing Hint** -- created only when `subnetId` is specified; tells Octavia which subnet the member resides on for L3 routing when it differs from the VIP subnet
- **OpenStack Tags** -- user-defined tags applied to the member for filtering and organization

## Before You Deploy

### OpenStack Account

- **Pool** -- an Octavia pool must exist to add the member to. Provide the pool ID directly or reference an OpenStackLoadBalancerPool Cloud Resource via ValueFromRef.
- **Backend server** -- the server at the specified IP address and port must be reachable from the Octavia amphora network. Verify network connectivity between the load balancer's VIP subnet and the member's subnet.
- **Subnet** (optional) -- when the member is on a different subnet than the VIP, provide the member's subnet ID so Octavia can route traffic correctly. Reference an OpenStackSubnet Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **OpenStack Load Balancer Member**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Pool Member** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancerMember
metadata:
  name: web-server-1
  org: acme-corp
  env: prod
spec:
  poolId:
    value: "<pool-id>"
  address: "192.168.1.10"
  protocolPort: 8080
```

```shell
planton apply -f member.yaml
```

This registers a backend server at `192.168.1.10:8080` with the specified pool. No subnet hint, weight override, or tags are configured. The member receives traffic immediately with the default weight of 1.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the member to a pool and subnet deployed in the same InfraPipeline:

```yaml
spec:
  poolId:
    valueFrom:
      kind: OpenStackLoadBalancerPool
      name: web-pool
      fieldPath: status.outputs.pool_id
  subnetId:
    valueFrom:
      kind: OpenStackSubnet
      name: app-subnet
      fieldPath: status.outputs.subnet_id
```

The InfraPipeline resolves the dependency graph, deploys the pool and subnet first, then provisions the member with the resolved values.

## Key Configuration

These are the most important decisions when configuring a pool member. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Backend address and port** -- The `address` and `protocolPort` fields identify the backend server. Both are immutable after creation -- changing either requires recreating the member. Create one member resource per backend server.

**Weight** -- Leave unset for equal distribution (Octavia defaults to 1). Set `weight` to a value between 1 and 256 for weighted load balancing, where higher values receive proportionally more traffic. Set to `0` to drain a member without removing it from the pool.

**Subnet hint** -- When the member is on the same subnet as the VIP, the `subnetId` field is optional. When the member is on a different subnet, provide `subnetId` so Octavia can perform L3 routing. Changing the subnet requires recreating the member.

**Administrative state** -- The member is active by default (`adminStateUp: true`). Set to `false` to remove the member from the pool's rotation without deleting it, useful for maintenance windows.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackLoadBalancerPool** | `poolId` | `status.outputs.pool_id` |
| **OpenStackSubnet** (optional) | `subnetId` | `status.outputs.subnet_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `member_id` | UUID of the pool member | Resource identification, API operations |
| `name` | Name of the member | Monitoring labels |
| `address` | IP address of the backend server | Diagnostics, inventory tracking |
| `protocol_port` | Port on the backend server | Diagnostics, inventory tracking |
| `weight` | Weight in the load-balancing algorithm | Monitoring, capacity planning |
| `region` | OpenStack region where the member was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard Pool Member** -- Registers a backend server with explicit subnet routing and default weight. Suitable for any backend instance that should receive traffic from a pool. Start from the **Standard Pool Member** preset.

## Works With

- [**OpenStack Load Balancer Pool**](/cloud-catalog/openstack-load-balancer-pool) -- provides the pool ID that this member is registered with
- [**OpenStack Subnet**](/cloud-catalog/openstack-subnet) -- provides the subnet ID for L3 routing when the member is on a different subnet than the VIP