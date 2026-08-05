# OpenStackServerGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackServerGroupSpec defines the configuration for an OpenStack Compute
server group.

A server group controls the placement of compute instances on hypervisors.
When instances reference a server group via scheduler_hints, Nova's scheduler
enforces the group's policy:

  - affinity:          all instances land on the SAME hypervisor
  - anti-affinity:     instances are spread across DIFFERENT hypervisors
  - soft-affinity:     best-effort same hypervisor (no hard failure)
  - soft-anti-affinity: best-effort spread (no hard failure)

Server groups are immutable -- all fields are ForceNew. Changing the policy
recreates the group, which orphans existing member instances from the old
group (they are NOT migrated).

The server group name is derived from metadata.name.

Terraform resource: openstack_compute_servergroup_v2
Pulumi resource: openstack.compute.ServerGroup

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackServerGroup
metadata:
  name: test-server-group
spec:
  policy: anti-affinity
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.policy` | `string` | yes |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.policy

`string` · required

policy is the placement policy for instances in this server group.

Valid values:
  "affinity"           - all members on the same hypervisor (hard constraint)
  "anti-affinity"      - all members on different hypervisors (hard constraint)
  "soft-affinity"      - prefer same hypervisor, no failure if unavailable (requires API 2.15+)
  "soft-anti-affinity" - prefer different hypervisors, no failure if unavailable (requires API 2.15+)

ForceNew: changing the policy recreates the server group.

Note: the OpenStack API accepts this as a list, but only one policy is
allowed per server group. We model it as a singular string for clarity.

- rule: {"string":{"minLen":"1","in":["affinity","anti-affinity","soft-affinity","soft-anti-affinity"]}}

### spec.region

`string`

region overrides the region from the provider config for this server group.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackServerGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.server_group_id` | `string` | server_group_id is the unique identifier (UUID) of the server group in OpenStack. This is the primary output used as a foreign key by OpenStackInstance. |
| `status.outputs.name` | `string` | name is the name of the server group (derived from metadata.name). |
| `status.outputs.members` | `[]string` | members is the list of instance UUIDs that belong to this server group. This is computed by OpenStack -- instances join the group when they are launched with scheduler_hints referencing this server group's ID. The list is empty at creation time and grows as instances are added. |
| `status.outputs.region` | `string` | region is the OpenStack region where the server group was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackInstance | `spec.serverGroupId` | `status.outputs.server_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
