# AzureManagedRedisAccessPolicyAssignment -- Design Research

## The Resource

An access policy assignment
(`Microsoft.Cache/redisEnterprise/databases/accessPolicyAssignments`)
grants Managed Redis data-plane access to a Microsoft Entra identity:
the Redis analog of a role assignment, and the grant half of the
keyless-by-default story (Managed Redis ships with access keys OFF).
The component maps onto
`azurerm_managed_redis_access_policy_assignment` (azurerm v4.x,
`internal/services/managedredis/managed_redis_access_policy_assignment_resource.go`),
parity-verified against pulumi-azure v6
(`managedredis.AccessPolicyAssignment`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `managed_redis_id` | same | `StringValueOrRef` → `AzureManagedRedis.managed_redis_id`; ForceNew. The provider derives the default database's path from the cluster id |
| `object_id` | same | `StringValueOrRef` → `AzureUserAssignedIdentity.principal_id` (the workload-identity default); literal GUID for users/groups. The PRINCIPAL-vs-client-id trap is called out in the field comment |

## Decomposition Decision

**SPLIT.** This is the grant class: a module must never own grants (the
principle that extracted role assignments from the identity kind).
Grants are many-per-instance, span two other resources (instance +
principal), and their lifecycle -- revoke, re-grant, audit -- is
independent of both.

## Recorded Constants and Skips (with reasons)

- **The access policy is the built-in "default"** -- Managed Redis
  exposes exactly one policy today and the provider hardcodes it; a
  one-value knob is a constant, not configuration. When Azure ships
  custom policies for Managed Redis (classic Redis already has them), a
  policy reference lands here.
- **No assignment name** -- Azure names the assignment after the
  granted object ID (which is also why an identity is granted at most
  once per database); a field would be contradictable redundant state.
- **No alias** -- the Managed Redis assignment API carries no
  `object_id_alias` (a classic-Redis-only convenience); clients connect
  with the object ID as the username.

## Operational Behavior Worth Knowing

- **Connecting under a grant**: username = the object id, password = an
  Entra token (`az account get-access-token --scope
  https://redis.azure.com/.default` for humans; the identity SDK for
  workloads), against `{hostname}:10000`.
- **Everything is ForceNew** -- replacing an assignment momentarily
  revokes and re-grants; safe for the grant class.
- **With access keys off (the default), grants are the ONLY way in** --
  deleting the last grant locks everyone out of the data plane until a
  new one propagates.
