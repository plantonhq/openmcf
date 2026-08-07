# AzureRedisCacheAccessPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureRedisCacheAccessPolicySpec** defines a CUSTOM data-plane access
policy on an Azure Cache for Redis -- a named permission set, written in
Redis's own ACL syntax, that Microsoft Entra identities are granted
through AzureRedisCacheAccessPolicyAssignment.

This is the Redis data-plane analog of a custom role definition: the
policy says WHAT is allowed (commands, command categories, key
patterns); the assignment says WHO gets it. Azure ships three built-in
policies -- "Data Owner", "Data Contributor", and "Data Reader" -- which
assignments reference by name without any policy resource existing;
create a custom policy only when the built-ins are too coarse (e.g.
read-write on one key prefix, no admin commands).

Entra authentication must be enabled on the cache
(redis_configuration.active_directory_authentication_enabled) for
policies to have any effect -- they gate token-authenticated clients,
not access-key clients.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRedisCacheAccessPolicy
metadata:
  name: test-redis-access-policy
spec:
  redisCacheId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cache/redis/planton-hack-redis
  policyName: app-read-only
  # Read-only on one key prefix -- the finer-than-built-in case a custom
  # policy exists for.
  permissions: "+@read +@connection ~app:*"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.redisCacheId` | `string \| valueFrom` | yes |  | AzureRedisCache (`status.outputs.redis_cache_id`) |
| `spec.policyName` | `string` | yes |  |  |
| `spec.permissions` | `string` | yes |  |  |

## Field Details

### spec.redisCacheId

`string | valueFrom` · required

The cache the policy is defined on, by ARM ID. References an
AzureRedisCache's redis_cache_id output. Fixed at creation: a policy
cannot move between caches.

- references: AzureRedisCache (`status.outputs.redis_cache_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCache, name: <that resource's name>, fieldPath: status.outputs.redis_cache_id}} -- a bare string does not parse

### spec.policyName

`string` · required

The policy's name -- what assignments reference, unique within the
cache. Must not collide with the built-in policy names ("Data Owner",
"Data Contributor", "Data Reader"). Changing the name replaces the
policy.

- rule: policy_name must not be one of the built-in policy names (Data Owner, Data Contributor, Data Reader) -- assignments reference built-ins directly without a policy resource
- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.permissions

`string` · required

The permission set, in Redis ACL syntax: command/category grants
followed by key patterns. Building blocks: "+@read" allows a command
category, "+get" a single command, "-@dangerous" carves one out, and
"~pattern" scopes which keys the grants apply to ("~*" for all keys,
"allkeys" is the same thing). Updatable in place.

Examples:
  "+@read +@connection ~*"          -- read-only on every key
  "+@all -@dangerous ~app1:*"       -- full data access on one prefix,
                                       no admin/dangerous commands
  "+get +set +del ~session:*"       -- three commands on one prefix

- rule: {"required":true}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureRedisCacheAccessPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.access_policy_id` | `string` | The Azure Resource Manager ID of the access policy. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cache/redis/{cache}/accessPolicies/{name} |
| `status.outputs.access_policy_name` | `string` | The policy's name -- the value AzureRedisCacheAccessPolicyAssignment.access_policy_name references to grant this policy to an identity. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.redisCacheId` | AzureRedisCache | `status.outputs.redis_cache_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureRedisCacheAccessPolicyAssignment | `spec.accessPolicyName` | `status.outputs.access_policy_name` |

## See Also

- [Overview](../README.md)
