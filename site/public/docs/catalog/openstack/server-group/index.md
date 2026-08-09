---
title: "Server Group"
description: "Server Group deployment documentation"
icon: "package"
order: 100
componentName: "openstackservergroup"
---

# OpenStack Server Group

Deploys a Nova server group on OpenStack that controls compute instance placement across hypervisors. Server groups enforce affinity or anti-affinity scheduling policies -- placing instances on the same hypervisor for low-latency communication, or spreading them across different hypervisors for fault isolation. Instances reference the server group via scheduler hints when they are launched.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Server Group** -- a Nova server group with the specified placement policy (affinity, anti-affinity, soft-affinity, or soft-anti-affinity) that the scheduler enforces when instances reference this group
- **OpenStack Tags** -- resource metadata applied automatically for tracking

## Before You Deploy

### OpenStack Account

- **Nova Compute API 2.15+** (optional) -- required only when using `soft-affinity` or `soft-anti-affinity` policies. Standard `affinity` and `anti-affinity` policies work with all API versions.
- **Sufficient hypervisor capacity** -- anti-affinity groups require at least as many hypervisors as group members. If insufficient hypervisors are available, instance launches fail with a scheduling error.

## Deploy

### Console

Open the deployment store, find **OpenStack Server Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Anti-Affinity** preset in the [Presets](#presets) tab for production fault isolation.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackServerGroup
metadata:
  name: db-spread
  org: acme-corp
  env: prod
spec:
  policy: anti-affinity
```

```shell
planton apply -f server-group.yaml
```

This creates a server group with an anti-affinity policy. Instances that reference this group via scheduler hints are placed on different hypervisors. No region override is configured.

## Key Configuration

These are the most important decisions when configuring a server group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Placement policy** -- The `policy` field determines how Nova schedules member instances. Use `anti-affinity` for production workloads where a hypervisor failure should affect at most one instance (HA databases, load-balanced tiers). Use `affinity` for tightly coupled workloads that benefit from low-latency inter-instance communication (HPC, batch processing).

**Soft vs. hard constraints** -- `soft-affinity` and `soft-anti-affinity` are best-effort versions that do not fail instance launches when the constraint cannot be satisfied. Use soft policies when availability is more important than strict placement guarantees. Requires Nova API 2.15+.

**Immutability** -- All fields are immutable. Changing the policy recreates the server group, which orphans existing member instances from the old group. Instances are not automatically migrated to the new group.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `server_group_id` | UUID of the server group in OpenStack | OpenStackInstance `serverGroupId` scheduler hints |
| `name` | Server group name | Labels, audit trails |
| `members` | List of instance UUIDs in the group | Monitoring, capacity tracking (empty at creation, populated as instances join) |
| `region` | OpenStack region of the server group | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Anti-affinity for fault isolation** -- Spreads instances across different hypervisors so a single hardware failure affects at most one group member. Use for HA database replicas, load-balanced application tiers, and any workload requiring hardware-level fault domains. Start from the **Anti-Affinity** preset.

**Affinity for low latency** -- Places all instances on the same hypervisor to minimize network latency between them. Use for tightly coupled applications, HPC workloads, and batch processing pipelines where data locality reduces transfer time. Start from the **Affinity** preset.

## Works With

This component operates independently and does not reference other components.