# Azure Storage Container

Creates a blob container inside an AzureStorageAccount -- the namespace unit of Azure blob storage. Applications organize objects into containers per data domain (uploads, logs, backups, artifacts), and Azure scopes anonymous access, data-plane RBAC, encryption scopes, and lifecycle prefixes at the container level.

## What Gets Created

When you deploy an AzureStorageContainer resource, Planton provisions:

- **Blob Container** -- an `azurerm_storage_container` on the referenced account (via its ARM id -- the control-plane path), with your chosen access posture, optional default encryption scope, and metadata

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the container in (referenced through `storageAccountId`)

## Quick Start

Create a file `container.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageContainer
metadata:
  name: app-uploads
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageContainer.app-uploads
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-app-storage
      fieldPath: status.outputs.storage_account_id
  containerName: uploads
```

Deploy:

```shell
planton apply -f container.yaml
```

The container is private by default -- every read requires authorization. For a public website or CDN origin, set `containerAccessType: BLOB` (anonymous object reads by direct URL, no listing) -- which also requires the account's `allowNestedItemsToBePublic` to be true.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `container_id` | The ARM id data-plane role assignments (Storage Blob Data Reader/Contributor) scope to |
| `container_name` | What SDK clients and function bindings reference within the account |
| `storage_account_name` | The account/container pair, without a second reference |

Blob URLs compose from the ACCOUNT's endpoint output plus this container's name: `{primary_blob_endpoint}{container_name}/{blob-path}`.

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/azurestorageaccount) -- the parent account
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- container-scoped data-plane grants
