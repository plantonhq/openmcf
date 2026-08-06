# Azure Storage Account

Creates an Azure Storage Account -- the multi-service storage primitive fronting Blob (objects), Files (SMB/NFS shares), Queues, Tables, and Data Lake Storage Gen2 behind one globally-unique DNS name. Blob containers are first-class AzureStorageContainer resources referencing the account, and app-hosting services (Function Apps, Web Apps) bind to its name and access-key outputs.

## What Gets Created

When you deploy an AzureStorageAccount resource, Planton provisions:

- **Storage Account** -- an `azurerm_storage_account` with your chosen kind/tier/replication SKU, security posture (TLS floor, key policy, anonymous-access lockdown, firewall), managed identity, customer-managed-key encryption, and the blob/file service settings (versioning, soft delete, point-in-time restore, SMB dials, CORS)
- **Lifecycle Management Policy** (when `lifecycleRules` are declared) -- the account's blob tiering/deletion schedule (Azure models this as one per-account policy document, so it lives here rather than as its own resource)
- **Static Website configuration** (when `staticWebsite` is declared) -- serves the auto-created `$web` container at the account's web endpoint

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the account in (referenced through `resourceGroup`)
- For customer-managed-key encryption: an **AzureKeyVaultKey** (in a purge-protected vault) and an **AzureUserAssignedIdentity** with wrap/unwrap access on the vault

## Quick Start

Create a file `storage.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageAccount
metadata:
  name: app-storage
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageAccount.app-storage
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      fieldPath: status.outputs.resource_group_name
  accountName: myappstorage001
  replicationType: ZRS
  blobProperties:
    versioningEnabled: true
    deleteRetentionPolicy:
      days: 14
    containerDeleteRetentionPolicy:
      days: 14
```

Deploy:

```shell
planton apply -f storage.yaml
```

> **Note:** `accountName` is globally unique across ALL of Azure and allows only 3-24 lowercase letters and digits -- no hyphens. It becomes the DNS prefix of every endpoint (`{name}.blob.core.windows.net`).

## Production Posture

Lock an account down with the firewall, Entra-first authorization, and customer-managed keys:

```yaml
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: my-rg
  accountName: myprodstorage001
  replicationType: GZRS
  allowNestedItemsToBePublic: false
  networkRules:
    defaultAction: DENY
    bypass:
      - AZURE_SERVICES
    virtualNetworkSubnetIds:
      - valueFrom:
          kind: AzureSubnet
          name: app-subnet
          fieldPath: status.outputs.subnet_id
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          kind: AzureUserAssignedIdentity
          name: storage-cmk-identity
          fieldPath: status.outputs.identity_id
  customerManagedKey:
    keyVaultKeyId:
      valueFrom:
        kind: AzureKeyVaultKey
        name: storage-cmk
        fieldPath: status.outputs.versionless_id
    userAssignedIdentityId:
      valueFrom:
        kind: AzureUserAssignedIdentity
        name: storage-cmk-identity
        fieldPath: status.outputs.identity_id
```

## Key Outputs

| Output | Purpose |
|--------|---------|
| `storage_account_id` | The ARM id containers and data-plane role assignments reference |
| `storage_account_name` | What Function App / Web App storage bindings consume |
| `primary_blob_endpoint` / `primary_blob_host` | Object access; CDN and custom-domain targets |
| `primary_web_endpoint` / `primary_web_host` | Static-website access; Front Door origin |
| `primary_access_key` | Static credential for consumers that require key auth (prefer Entra roles) |
| Secondary endpoints | The read-only paired-region mirror (RA_GRS / RA_GZRS only) |

## Related Resources

- [Azure Storage Container](/docs/catalog/azure/azurestoragecontainer) -- blob containers on this account
- [Azure Key Vault Key](/docs/catalog/azure/azurekeyvaultkey) -- customer-managed encryption keys
- [Azure Subnet](/docs/catalog/azure/azuresubnet) -- VNet rules for the firewall
- [Azure Function App](/docs/catalog/azure/azurefunctionapp) -- binds to the account for runtime state
