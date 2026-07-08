# Azure Managed Redis Access Policy Assignment

Grants Managed Redis data-plane access to a Microsoft Entra identity --
the Redis analog of a role assignment. Managed Redis is keyless by
default, so grants are how clients connect at all: the identity
presents its object ID as the Redis username and an Entra token as the
password, and no secret exists anywhere in the composition.

## What Gets Created

When you deploy an AzureManagedRedisAccessPolicyAssignment resource,
Planton provisions:

- **Access Policy Assignment** -- an
  `azurerm_managed_redis_access_policy_assignment` granting the
  built-in "default" policy to the principal on the referenced
  instance's default database

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **An AzureManagedRedis instance**
- **A principal to grant** -- an AzureUserAssignedIdentity (referenced),
  or a user/group object id (literal)

## Quick Start

Create a file `grant.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedRedisAccessPolicyAssignment
metadata:
  name: app-cache-grant
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureManagedRedisAccessPolicyAssignment.app-cache-grant
spec:
  managedRedisId:
    valueFrom:
      kind: AzureManagedRedis
      name: app-cache
      fieldPath: status.outputs.managed_redis_id
  objectId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: app-identity
      fieldPath: status.outputs.principal_id
```

Deploy:

```shell
planton pulumi up --manifest grant.yaml
```

## Spec Highlights

- `object_id` -- the PRINCIPAL id for managed identities (never the
  client id -- the wrong one fails at connect time, not deploy time);
  a literal GUID for users and groups (grant a group to cover a whole
  on-call rotation)
- **One grant per identity per instance** -- Azure names the assignment
  after the object ID
- **Revocation is deletion** -- removing the resource revokes access
  immediately

## Stack Outputs

| Output | Description |
| --- | --- |
| `access_policy_assignment_id` | The assignment's ARM ID |
| `access_policy_assignment_name` | The assignment's name (equals the granted object ID) |
