# OpenStack Project

Deploys a Keystone project on OpenStack -- the fundamental organizational unit that provides resource isolation, quota boundaries, and access control scoping. All cloud resources (instances, volumes, networks) belong to a project. The project supports domain assignment, hierarchical nesting, and enable/disable controls.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Keystone Project** -- an identity project with configurable domain, enable state, parent hierarchy, and description
- **OpenStack Tags** -- user-defined tags applied to the project for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Admin privileges** -- creating projects is an admin-level operation. The OpenStack credentials must have the `admin` role or equivalent project-creation permissions in Keystone.
- **Domain** (optional) -- if your deployment uses multiple Keystone domains, provide the `domainId` for the target domain. Single-domain deployments can omit this field to use the default domain.
- **Parent project** (optional) -- if using hierarchical multi-tenancy, provide the `parentId` of the parent project. Nested projects inherit quota limits from their parent. Most deployments use flat (top-level) projects.

## Deploy

### Console

Open the deployment store, find **OpenStack Project**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackProject
metadata:
  name: team-platform
  org: acme-corp
  env: prod
spec:
  description: Platform team workspace
```

```shell
planton apply -f project.yaml
```

This creates an enabled project in the default domain with no parent hierarchy. The project is ready for role assignments and resource provisioning immediately after creation.

## Key Configuration

These are the most important decisions when configuring a project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Enable state** -- `enabled` defaults to `true`. Set to `false` to create the project in a disabled state -- all users lose access to its resources, but existing resources are preserved. Useful for suspending a tenant without deleting infrastructure.

**Domain assignment** -- `domainId` assigns the project to a specific Keystone domain. This is a ForceNew field -- changing the domain recreates the project. Most single-domain deployments leave this empty to use the default domain.

**Hierarchical nesting** -- `parentId` places the project under a parent project in Keystone's hierarchy. Nested projects inherit quota limits and can be used for organizational structuring in large deployments. This is a ForceNew field.

**Tags** -- `tags` attaches simple string labels for filtering and organizational purposes. Tags are queryable through the OpenStack API and visible in the Horizon dashboard.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `project_id` | UUID of the project in OpenStack | Role assignments, resource provisioning scope |
| `name` | Name of the project | Monitoring labels, configuration references |
| `domain_id` | Keystone domain the project belongs to | Domain-aware downstream resources |
| `enabled` | Whether the project is currently active | Health checks, conditional provisioning |
| `region` | OpenStack region where the project was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard tenant project** -- An enabled project in the default domain with no hierarchy. The starting point for most workloads -- create the project, assign roles to users, then provision instances, networks, and volumes within it. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.