# OpenStack Role Assignment

Deploys a Keystone role assignment on OpenStack -- binding a role to a user or group on a project or domain. This is the fundamental authorization mechanism that determines what actions a principal can perform on a specific scope. The assignment supports ValueFromRef wiring for project references in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Role Assignment** -- a Keystone identity binding that grants a specified role to a user or group on a project or domain scope

## Before You Deploy

### OpenStack Account

- **Admin privileges** -- assigning roles is an admin-level operation. The OpenStack credentials must have sufficient permissions to create role assignments in Keystone.
- **Role UUID** -- the `roleId` must reference an existing Keystone role. Run `openstack role list` to find available role UUIDs (e.g., `admin`, `member`, `reader`).
- **Scope** -- exactly one of `projectId` or `domainId` must be set. For project-scoped assignments, the project must exist. For domain-scoped assignments, provide the domain UUID directly.
- **Principal** -- exactly one of `userId` or `groupId` must be set. The user or group must exist in Keystone. Run `openstack user list` or `openstack group list` to find UUIDs.
- **Immutability** -- all fields are ForceNew. Any change to role, scope, or principal destroys and recreates the assignment.

## Deploy

### Console

Open the deployment store, find **OpenStack Role Assignment**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Project User Member** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackRoleAssignment
metadata:
  name: platform-member
  org: acme-corp
  env: prod
spec:
  roleId: "<member-role-id>"
  projectId:
    value: "<project-id>"
  userId: "<user-id>"
```

```shell
planton apply -f role-assignment.yaml
```

This assigns the `member` role to a user on a project. The user can create and manage cloud resources within that project. No domain scope or group assignment is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the role assignment to a project deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: OpenStackProject
      name: team-platform
      fieldPath: status.outputs.project_id
```

The InfraPipeline resolves the dependency graph, deploys the project first, then provisions the role assignment with the resolved project ID.

## Key Configuration

These are the most important decisions when configuring a role assignment. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scope selection** -- Choose between project scope (`projectId`) and domain scope (`domainId`). Project-scoped assignments are the most common -- they grant access within a single tenant. Domain-scoped assignments grant access across all projects in a domain and are typically reserved for cloud administrators.

**Principal type** -- Choose between user (`userId`) and group (`groupId`). User assignments are direct and simple. Group assignments are preferred for teams -- all members of the Keystone group inherit the role, and adding/removing group members updates access without modifying the assignment.

**Role selection** -- The `roleId` determines the privilege level. Common roles are `admin` (full control), `member` (standard operations), and `reader` (read-only access). Custom roles may exist in your deployment for fine-grained access control.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Composite identifier of the role assignment | Audit trails, automation reference |
| `role_id` | UUID of the assigned role | Verifying role alignment |
| `project_id` | Project scope of the assignment (if project-scoped) | Downstream resources scoped to the same project |
| `domain_id` | Domain scope of the assignment (if domain-scoped) | Downstream resources scoped to the same domain |
| `user_id` | User principal (if user assignment) | Audit trails |
| `group_id` | Group principal (if group assignment) | Audit trails |
| `region` | OpenStack region where the assignment was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Project user member** -- Assigns the `member` role to a user on a project. The standard pattern for granting day-to-day cloud operations access to a team member. Start from the **Project User Member** preset.

## Works With

- [**OpenStack Project**](/cloud-catalog/openstack-project) -- provides the project ID as the scope for project-level role assignments