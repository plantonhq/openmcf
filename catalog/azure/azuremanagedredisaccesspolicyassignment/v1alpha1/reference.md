# AzureManagedRedisAccessPolicyAssignment

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureManagedRedisAccessPolicyAssignmentSpec** grants Redis data-plane
access to a Microsoft Entra identity on an Azure Managed Redis instance
-- the Redis analog of a role assignment.

This is the grant half of Managed Redis's keyless-by-default story:
access keys are off unless explicitly enabled, so a granted identity is
how clients connect at all -- the identity presents its object ID as
the Redis username and an Entra token as the password, and no secret
exists to leak or rotate.

Managed Redis exposes one built-in access policy ("default", full data
access); every assignment grants it. Azure names the assignment after
the granted object ID, so an identity is granted at most once per
database -- there is nothing else to name or configure. (Classic Azure
Cache for Redis offered custom scoped policies; Managed Redis does not
yet -- when Azure adds them, a policy reference lands here.)

Every field is fixed at creation -- changing anything replaces the
assignment (safe: replacing a grant momentarily revokes and re-grants).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedRedisAccessPolicyAssignment
metadata:
  name: test-managed-redis-grant
spec:
  managedRedisId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cache/redisEnterprise/planton-hack-managed-redis
  # The PRINCIPAL id of the granted identity (never the client id).
  # Azure also names the assignment after this value.
  objectId:
    value: 11111111-2222-3333-4444-555555555555
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.managedRedisId` | `string \| valueFrom` | yes |  | AzureManagedRedis (`status.outputs.managed_redis_id`) |
| `spec.objectId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |

## Field Details

### spec.managedRedisId

`string | valueFrom` · required

The Managed Redis instance the grant applies to, by ARM ID (the
grant is created on its default database). References an
AzureManagedRedis's managed_redis_id output.

- references: AzureManagedRedis (`status.outputs.managed_redis_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureManagedRedis, name: <that resource's name>, fieldPath: status.outputs.managed_redis_id}} -- a bare string does not parse

### spec.objectId

`string | valueFrom` · required

The Entra object (principal) ID being granted -- a user, group,
service principal, or managed identity, as a GUID. References an
AzureUserAssignedIdentity's principal_id output by default (the
workload-identity case); pass a literal GUID for users and groups.
NOTE: for a managed identity this is the PRINCIPAL id, not the
client id -- granting the client id fails at connect time, not at
deploy time. Azure also uses this value as the assignment's name,
which is why an identity can be granted at most once per database.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureManagedRedisAccessPolicyAssignment, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.access_policy_assignment_id` | `string` | The Azure Resource Manager ID of the access policy assignment. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cache/redisEnterprise/{cluster}/databases/default/accessPolicyAssignments/{objectId} |
| `status.outputs.access_policy_assignment_name` | `string` | The assignment's name -- Azure names it after the granted object ID, so this equals the principal's GUID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.managedRedisId` | AzureManagedRedis | `status.outputs.managed_redis_id` |
| `spec.objectId` | AzureUserAssignedIdentity | `status.outputs.principal_id` |

## See Also

- [Overview](../README.md)
