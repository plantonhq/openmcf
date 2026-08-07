# AzureRedisCacheAccessPolicyAssignment

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureRedisCacheAccessPolicyAssignmentSpec** grants a Redis data-plane
access policy to a Microsoft Entra identity on an Azure Cache for Redis
-- the Redis analog of a role assignment. The policy (built-in or an
AzureRedisCacheAccessPolicy) says what is allowed; this assignment says
WHO gets it.

This is the grant half of the keyless Redis story: with Entra
authentication enabled on the cache
(redis_configuration.active_directory_authentication_enabled), a granted
identity connects with its object ID (or the alias set here) as the
Redis username and an Entra token as the password -- no access key
involved. Pair with access_keys_authentication_enabled: false on the
cache and no secret exists at all.

Every field is fixed at creation -- changing anything replaces the
assignment (safe: replacing a grant momentarily revokes and re-grants).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRedisCacheAccessPolicyAssignment
metadata:
  name: test-redis-grant
spec:
  redisCacheId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cache/redis/planton-hack-redis
  assignmentName: app-identity-data-reader
  # A built-in policy referenced by literal name (custom policies
  # reference an AzureRedisCacheAccessPolicy's name output instead).
  accessPolicyName:
    value: Data Reader
  # The PRINCIPAL id of the granted identity (never the client id).
  objectId:
    value: 11111111-2222-3333-4444-555555555555
  objectIdAlias: app-identity
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.redisCacheId` | `string \| valueFrom` | yes |  | AzureRedisCache (`status.outputs.redis_cache_id`) |
| `spec.assignmentName` | `string` | yes |  |  |
| `spec.accessPolicyName` | `string \| valueFrom` | yes |  | AzureRedisCacheAccessPolicy (`status.outputs.access_policy_name`) |
| `spec.objectId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.objectIdAlias` | `string` | yes |  |  |

## Field Details

### spec.redisCacheId

`string | valueFrom` · required

The cache the grant applies to, by ARM ID. References an
AzureRedisCache's redis_cache_id output.

- references: AzureRedisCache (`status.outputs.redis_cache_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCache, name: <that resource's name>, fieldPath: status.outputs.redis_cache_id}} -- a bare string does not parse

### spec.assignmentName

`string` · required

The assignment's name, unique within the cache. A label for the
grant itself (it does not affect authentication) -- name it after the
principal and policy, e.g. "app-identity-data-reader".

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.accessPolicyName

`string | valueFrom` · required

The policy being granted. The three built-ins are referenced by
literal name -- "Data Owner" (full access including admin commands),
"Data Contributor" (read-write), "Data Reader" (read-only) -- and a
custom policy by referencing an AzureRedisCacheAccessPolicy's
access_policy_name output.

- references: AzureRedisCacheAccessPolicy (`status.outputs.access_policy_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCacheAccessPolicy, name: <that resource's name>, fieldPath: status.outputs.access_policy_name}} -- a bare string does not parse

### spec.objectId

`string | valueFrom` · required

The Entra object (principal) ID being granted -- a user, group,
service principal, or managed identity. References an
AzureUserAssignedIdentity's principal_id output by default (the
workload-identity case); pass a literal GUID for users and groups.
NOTE: for a managed identity this is the PRINCIPAL id, not the
client id -- granting the client id fails at connect time, not at
deploy time.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.objectIdAlias

`string` · required

A human-readable alias for the object ID. Doubles as an alternative
Redis USERNAME at connect time (clients may authenticate as either
the raw object ID or this alias, with the Entra token as password).
Convention: the identity's display name, or "ServicePrincipal" /
"UserGroup" style labels for non-identity principals.

- rule: {"required":true,"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureRedisCacheAccessPolicyAssignment, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.access_policy_assignment_id` | `string` | The Azure Resource Manager ID of the access policy assignment. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cache/redis/{cache}/accessPolicyAssignments/{name} |
| `status.outputs.access_policy_assignment_name` | `string` | The assignment's name within the cache. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.redisCacheId` | AzureRedisCache | `status.outputs.redis_cache_id` |
| `spec.accessPolicyName` | AzureRedisCacheAccessPolicy | `status.outputs.access_policy_name` |
| `spec.objectId` | AzureUserAssignedIdentity | `status.outputs.principal_id` |

## See Also

- [Overview](../README.md)
