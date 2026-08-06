# OpenStackRoleAssignment

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackRoleAssignmentSpec defines the configuration for an OpenStack
Identity (Keystone) role assignment.

A role assignment binds a role to a principal (user or group) on a scope
(project or domain). This is the fundamental authorization mechanism in
OpenStack -- it determines what actions a user or group can perform on
a specific project or domain.

Constraints:
  - Exactly one of project_id or domain_id must be set (scope)
  - Exactly one of user_id or group_id must be set (principal)
  - All fields are ForceNew -- changing any field recreates the assignment

This is an admin-level operation. The OpenStack credentials must have
sufficient permissions to assign roles.

Role assignments do not have a natural "name" -- metadata.name provides
KRM identity only and is not sent to the OpenStack API.

Terraform resource: openstack_identity_role_assignment_v3
Pulumi resource: openstack.identity.RoleAssignment

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackRoleAssignment
metadata:
  name: test-role-assignment
spec:
  roleId: "test-role-uuid"
  projectId:
    value: test-project-uuid
  userId: "test-user-uuid"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.roleId` | `string` | yes |  |  |
| `spec.projectId` | `string \| valueFrom` |  |  | OpenStackProject (`status.outputs.project_id`) |
| `spec.domainId` | `string` |  |  |  |
| `spec.userId` | `string` |  |  |  |
| `spec.groupId` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.roleId

`string` · required

role_id is the UUID of the role to assign.
Required. ForceNew: changing the role recreates the assignment.
Roles are admin-managed Keystone objects (e.g., "admin", "member", "reader").
Use `openstack role list` to find available role UUIDs.

- rule: {"string":{"minLen":"1"}}

### spec.projectId

`string | valueFrom`

project_id is the project to assign the role on.
Mutually exclusive with domain_id (exactly one must be set).
ForceNew: changing the scope recreates the assignment.
Can reference an OpenStackProject resource's output or be a literal UUID.

- references: OpenStackProject (`status.outputs.project_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackProject, name: <that resource's name>, fieldPath: status.outputs.project_id}} -- a bare string does not parse

### spec.domainId

`string`

domain_id is the domain to assign the role on.
Mutually exclusive with project_id (exactly one must be set).
ForceNew: changing the scope recreates the assignment.
Domains are admin-managed Keystone objects -- this is a plain UUID string.

### spec.userId

`string`

user_id is the UUID of the user to assign the role to.
Mutually exclusive with group_id (exactly one must be set).
ForceNew: changing the principal recreates the assignment.
Users are admin-managed Keystone objects -- this is a plain UUID string.

### spec.groupId

`string`

group_id is the UUID of the group to assign the role to.
Mutually exclusive with user_id (exactly one must be set).
ForceNew: changing the principal recreates the assignment.
Groups are admin-managed Keystone objects -- this is a plain UUID string.

### spec.region

`string`

region overrides the region from the provider config.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `scope.exactly_one`: exactly one of project_id or domain_id must be set -- they define the scope of the role assignment
- `principal.exactly_one`: exactly one of user_id or group_id must be set -- they define who receives the role

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackRoleAssignment, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | id is the composite identifier of the role assignment. Terraform generates this as a combination of role_id/scope_id/principal_id. |
| `status.outputs.role_id` | `string` | role_id is the UUID of the assigned role. |
| `status.outputs.project_id` | `string` | project_id is the project scope (if project-scoped assignment). |
| `status.outputs.domain_id` | `string` | domain_id is the domain scope (if domain-scoped assignment). |
| `status.outputs.user_id` | `string` | user_id is the user principal (if user assignment). |
| `status.outputs.group_id` | `string` | group_id is the group principal (if group assignment). |
| `status.outputs.region` | `string` | region is the OpenStack region where the assignment was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.projectId` | OpenStackProject | `status.outputs.project_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
