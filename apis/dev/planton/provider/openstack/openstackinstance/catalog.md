# OpenStack Instance

Deploys a Nova compute instance on OpenStack with configurable flavor, image, network attachments, security groups, block device mappings, and placement controls. The instance supports both ephemeral image-based and persistent volume-based boot modes, with ValueFromRef wiring for keypairs, networks, security groups, and server groups in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Instance** -- a Nova virtual machine with the specified flavor, image or block device mapping, network attachments, and optional SSH keypair
- **Block Device Mappings** -- created only when `blockDevice` entries are specified; maps Cinder volumes, Glance images, or snapshots as boot or data disks
- **Scheduler Hints** -- created only when `serverGroupId` is provided; places the instance according to the server group's affinity or anti-affinity policy
- **OpenStack Tags** -- user-defined tags applied to the instance for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Flavor** -- the instance requires a valid flavor (by name or ID). Run `openstack flavor list` to find available options in your deployment.
- **Image or boot volume** -- provide a Glance image name/ID for ephemeral boot, or configure `blockDevice` for persistent boot-from-volume. At least one boot source is required.
- **Network** -- at least one network attachment is required. Provide the network UUID directly or reference an OpenStackNetwork Cloud Resource via ValueFromRef.
- **SSH keypair** (optional) -- if you need key-based SSH access, create or import a keypair. Reference an OpenStackKeypair Cloud Resource via ValueFromRef, or provide the keypair name directly.
- **Security groups** (optional) -- if omitted, OpenStack applies the default security group. Reference OpenStackSecurityGroup Cloud Resources by name via ValueFromRef for InfraChart wiring.

## Deploy

### Console

Open the deployment store, find **OpenStack Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard VM** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackInstance
metadata:
  name: app-server
  org: acme-corp
  env: prod
spec:
  flavorName: m1.medium
  imageName: ubuntu-22.04
  networks:
    - uuid:
        value: "<network-id>"
```

```shell
planton apply -f instance.yaml
```

This creates an instance with an ephemeral root disk, a single network attachment, and the default security group. No SSH keypair, block device mappings, or placement controls are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to other resources deployed in the same InfraPipeline:

```yaml
spec:
  keyPair:
    valueFrom:
      kind: OpenStackKeypair
      name: admin-key
      fieldPath: status.outputs.name
  networks:
    - uuid:
        valueFrom:
          kind: OpenStackNetwork
          name: app-network
          fieldPath: status.outputs.network_id
  securityGroups:
    - valueFrom:
        kind: OpenStackSecurityGroup
        name: web-sg
        fieldPath: status.outputs.name
  serverGroupId:
    valueFrom:
      kind: OpenStackServerGroup
      name: app-spread
      fieldPath: status.outputs.server_group_id
```

The InfraPipeline resolves the dependency graph, deploys the keypair, network, security group, and server group first, then provisions the instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring an instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Boot mode** -- Choose between image-based (ephemeral) and volume-based (persistent) boot. For production, use `blockDevice` with `sourceType: image` and `destinationType: volume` for a root disk that survives instance deletion. For development, `imageName` with ephemeral boot is simpler and faster to provision.

**Flavor selection** -- Use `flavorName` for human-readable sizing (e.g., `m1.medium`) or `flavorId` for exact UUID references. Exactly one must be provided. Changing the flavor after creation triggers a resize operation.

**Network mode** -- Each `networks` entry connects the instance via a network UUID (auto-creates a port) or a pre-provisioned port UUID. Use ports for stable MAC/IP addresses or when you need port-level security group assignments. Mark one network as `accessNetwork: true` to control which IP appears in `accessIpV4`.

**User data** -- The `userData` field accepts cloud-init scripts or cloud-config YAML for first-boot automation. Changing user data recreates the instance.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackKeypair** (optional) | `keyPair` | `status.outputs.name` |
| **OpenStackNetwork** (optional) | `networks[].uuid` | `status.outputs.network_id` |
| **OpenStackNetworkPort** (optional) | `networks[].port` | `status.outputs.port_id` |
| **OpenStackSecurityGroup** (optional) | `securityGroups` | `status.outputs.name` |
| **OpenStackServerGroup** (optional) | `serverGroupId` | `status.outputs.server_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | UUID of the instance in OpenStack | Volume attachments, floating IP associations |
| `name` | Name of the instance | DNS records, monitoring labels |
| `access_ip_v4` | Primary IPv4 address for accessing the instance | Load balancer members, DNS A records |
| `access_ip_v6` | Primary IPv6 address for accessing the instance | DNS AAAA records |
| `region` | OpenStack region where the instance was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard VM** -- Boots from a Glance image with a specified flavor, keypair, and network attachment. Ephemeral root disk on the hypervisor. Suitable for development instances, stateless workers, and general-purpose VMs. Start from the **Standard VM** preset.

**Boot-from-volume** -- Boots from a Cinder volume created from a Glance image. The root disk persists independently of the instance lifecycle, enabling snapshots, migration, and data preservation. Start from the **Boot-from-Volume** preset.

## Works With

- [**OpenStack Keypair**](/cloud-catalog/openstack-keypair) -- provides the SSH keypair name for key-based instance access
- [**OpenStack Network**](/cloud-catalog/openstack-network) -- provides the network UUID for instance network attachments
- [**OpenStack Network Port**](/cloud-catalog/openstack-network-port) -- provides a pre-provisioned port for stable network identity
- [**OpenStack Security Group**](/cloud-catalog/openstack-security-group) -- provides the security group name for traffic filtering
- [**OpenStack Server Group**](/cloud-catalog/openstack-server-group) -- provides the server group ID for placement control (affinity/anti-affinity)