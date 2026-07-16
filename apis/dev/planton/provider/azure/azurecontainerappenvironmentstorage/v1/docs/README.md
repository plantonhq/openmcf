# AzureContainerAppEnvironmentStorage -- Design Research

## The Resource

An environment storage registration
(`Microsoft.App/managedEnvironments/storages`) makes an Azure Files share
mountable by Container Apps and Jobs: workloads declare share-backed
volumes that reference the registration by name. The component maps onto
`azurerm_container_app_environment_storage` (azurerm v4.x,
`internal/services/containerapps/container_app_environment_storage_resource.go`),
parity-verified against pulumi-azure v6 (`containerapp.EnvironmentStorage`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `container_app_environment_id` | `container_app_environment_id` | Single parent ARM-id FK to `AzureContainerAppEnvironment.environment_id`. ForceNew |
| `name` | `storage_name` | The registration handle app/job volumes reference; the provider's name contract (lowercase alphanumerics/hyphens, <=32) is mirrored as a CEL. ForceNew |
| `share_name` | `share_name` | Required; FK-defaults to `AzureStorageShare.share_name`. ForceNew |
| `access_mode` | `access_mode` | Required closed enum READ_ONLY/READ_WRITE mapping to ARM's ReadOnly/ReadWrite. ForceNew |
| `account_name` | `account_name` | The SMB path; FK-defaults to `AzureStorageAccount.storage_account_name`; provider contract `RequiredWith: access_key`, `ConflictsWith: nfs_server_url`. ForceNew |
| `access_key` | `access_key` | Sensitive; FK-defaults to `AzureStorageAccount.primary_access_key`. The ONE updatable field (key rotation) |
| `nfs_server_url` | `nfs_server_url` | The NFS path (`{account}.file.core.windows.net`); provider contract `ConflictsWith: account_name`. ForceNew |

No tags: ARM does not support tags on environment storage registrations.

## Decomposition Decision

First-class kind, not a fold into the environment: registrations are
many-per-environment with independent lifecycles (shares come and go as
workloads change), and they are FK-referenced -- app and job volumes bind
a registration through its `storage_name` output. Folding them into the
environment spec would force environment updates for every share change
and make the volume seam a dangling string.

## Front-Loaded Contracts

- **SMB XOR NFS** (message CEL): exactly one protocol per registration --
  `account_name` + `access_key` together, or `nfs_server_url` alone.
  Mirrors the provider's `RequiredWith`/`ConflictsWith` schema contracts.
- The NFS-requires-VNet-injection pairing is a cross-resource contract
  ARM enforces at apply time (a spec rule cannot see the referenced
  environment's network configuration).

## Recorded Skips (with reasons)

Nothing skipped: the azurerm surface is exactly the seven fields above,
and the spec models all of them.

## Operational Behavior Worth Knowing

- Recreating a registration (any change except `access_key`) briefly
  breaks the volume mounts that reference it -- plan rotations of
  anything other than the key around workload restarts.
- NFS registrations require the environment to be VNet-injected AND the
  storage account to allow the environment's subnet; both are apply-time
  cross-resource contracts.
