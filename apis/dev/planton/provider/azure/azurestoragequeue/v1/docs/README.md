# AzureStorageQueue -- Design Research

## The Resource

A Storage queue (`Microsoft.Storage/storageAccounts/queueServices/
queues`) is the simple, massive-scale message buffer of Azure storage:
64 KB messages, poll-and-delete consumption, capacity bounded only by
the account. The component maps onto `azurerm_storage_queue` (azurerm
v4.x, `internal/services/storage/storage_queue_resource.go`),
parity-verified against pulumi-azure v6 (`storage.Queue`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | The v4/v5-forward ARM-id parent path (the deprecated `storage_account_name` data-plane form is not modeled); single authoritative parent FK, ForceNew |
| `name` | `queue_name` | Required, ForceNew, 3-63 lowercase/digits/hyphens, no leading/trailing hyphen, no consecutive hyphens (Azure's service contract; azurerm's validator omits the consecutive-hyphen check) |
| `metadata` | `metadata` | Queue metadata is NOT Azure tags -- ARM does not support tags on queues, so the platform's identity tags live on the account |

The queue's management surface is deliberately this small -- everything
else about a queue (messages, visibility timeouts, poison handling) is
data-plane behavior owned by producers and consumers at runtime.

## Decomposition Decisions

- **First-class kind, not a fold**: queues are many-per-account with
  independent lifecycles, and consumers (Functions triggers, worker
  deployments, role assignments) reference them individually.

## Recorded Skips (with reasons)

- **`resource_manager_id` attribute** -- a 4.x-only compat shim
  deprecated in favor of the resource id (removed in azurerm v5); the
  `queue_id` output IS the resource-manager id.
- **No `url` output** -- the queue's data-plane URL is the ACCOUNT's
  queue endpoint plus the queue name, and only the account knows its
  real endpoint (partitioned-DNS accounts differ). Compose client URLs
  from `AzureStorageAccount.primary_queue_endpoint` + `queue_name`.

## Operational Behavior Worth Knowing

- **Queues create in seconds** via the ARM control plane -- creation and
  verification are unaffected by the account's data-plane firewall.
- **The poison-queue convention is a NAMING convention**: Functions
  moves repeatedly-failing messages to `{queue}-poison` automatically --
  declare the companion queue explicitly so it is owned infrastructure,
  not an accident.
- **A queue cannot move between accounts** -- the parent reference is
  ForceNew.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `queue_id` output ← data-plane role-assignment scopes
  (AzureRoleAssignment with Storage Queue Data roles)
- `queue_name` + the account's `primary_queue_endpoint` ← SDK clients,
  Functions queue triggers
