# OpenStackProject

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackProjectSpec defines the configuration for an OpenStack Identity
(Keystone) project.

A project (historically called "tenant") is the fundamental organizational
unit in OpenStack. All cloud resources (instances, volumes, networks) belong
to a project. Projects provide resource isolation, quota boundaries, and
access control scoping.

Creating projects is typically an admin-level operation. Tenant users cannot
create projects -- only cloud administrators or users with the "admin" role
can provision new projects.

The project name is derived from metadata.name.

Terraform resource: openstack_identity_project_v3
Pulumi resource: openstack.identity.Project

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackProject
metadata:
  name: test-project
spec:
  description: Test project for local development
  enabled: true
  tags:
    - test
    - hack
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.domainId` | `string` |  |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.parentId` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.description

`string`

description is a human-readable description of the project.
Visible in the OpenStack API, Horizon dashboard, and CLI output.

### spec.domainId

`string`

domain_id is the Keystone domain to which this project belongs.
ForceNew: changing the domain requires recreating the project.
If omitted, OpenStack assigns the project to the default domain.
Most single-domain deployments can leave this empty.
Example: "default", or a domain UUID like "abcdef12-3456-7890-abcd-ef1234567890"

### spec.enabled

`bool` · optional (explicit presence)

enabled controls whether the project is active.
When disabled (false), all users in the project lose access to its
resources, but the resources are NOT deleted.
Default: true

- default: `true`

### spec.parentId

`string`

parent_id is the UUID of the parent project in the project hierarchy.
ForceNew: changing the parent requires recreating the project.
If omitted, the project is created as a top-level project in its domain.
Project hierarchies are used for nested quota management and organizational
structuring in large deployments.

### spec.tags

`[]string`

tags is a set of tags to associate with the project.
Tags are simple strings used for filtering and organizational purposes.
They are visible in the OpenStack API and can be queried via tag-based filters.

### spec.region

`string`

region overrides the region from the provider config for this project.
If omitted, the region from the OpenStack provider config is used.
Note: Keystone is typically a global service (not per-region), so this
field is rarely needed. It controls which Keystone endpoint is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackProject, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.project_id` | `string` | project_id is the unique identifier (UUID) of the project in OpenStack. This is the primary output used as a foreign key by OpenStackRoleAssignment and other downstream resources that need to reference the project. |
| `status.outputs.name` | `string` | name is the name of the project (derived from metadata.name). |
| `status.outputs.domain_id` | `string` | domain_id is the Keystone domain to which this project belongs. Computed by OpenStack if not specified in the spec. |
| `status.outputs.enabled` | `bool` | enabled indicates whether the project is currently active. |
| `status.outputs.region` | `string` | region is the OpenStack region where the project was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackRoleAssignment | `spec.projectId` | `status.outputs.project_id` |

## See Also

- [Overview](../README.md)
