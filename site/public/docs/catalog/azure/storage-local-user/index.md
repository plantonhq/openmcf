---
title: "Storage Local User"
description: "Storage Local User deployment documentation"
icon: "package"
order: 100
componentName: "azurestoragelocaluser"
---

# Azure Storage Local User

Creates a local user on an AzureStorageAccount -- the credential identity the account's SFTP endpoint authenticates. Each user carries its own SSH credentials, home directory, and per-container permission scopes: the isolation unit that lets one account serve many file-exchange partners.

## What Gets Created

When you deploy an AzureStorageLocalUser resource, Planton provisions:

- **Storage Local User** -- an `azurerm_storage_account_local_user` on the referenced account, with SSH key and/or Azure-generated password authentication, an optional home directory, and per-resource permission scopes

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the user on (referenced through `storageAccountId`) -- for the user to actually CONNECT, the account needs `sftpEnabled: true` (which requires `isHnsEnabled: true`)
- **A container** (or file share) for the permission scope to grant access to

## Quick Start

Create a file `local-user.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageLocalUser
metadata:
  name: partner01
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageLocalUser.partner01
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-sftp-storage
      fieldPath: status.outputs.storage_account_id
  userName: partner01
  sshKeyEnabled: true
  sshAuthorizedKeys:
    - key: ssh-ed25519 AAAA... partner-pipeline
      description: partner pipeline key
  homeDirectory: partner-inbound
  permissionScopes:
    - service: BLOB
      resourceName:
        valueFrom:
          kind: AzureStorageContainer
          name: partner-inbound
          fieldPath: status.outputs.container_name
      read: true
      write: true
      list: true
      create: true
```

Deploy:

```shell
planton apply -f local-user.yaml
```

The partner connects as `{account-name}.partner01` to `{account-name}.blob.core.windows.net` on port 22.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `local_user_id` | The user's ARM id |
| `sftp_username` | The full login: `{account-name}.{user-name}` |
| `password` | The Azure-generated password (SECRET; returned exactly once, only when password auth is on) |
| `sid` | The user's Security Identifier (secret-bearing; Azure Files ACLs reference it) |
| `storage_account_name` | The account/user pair, without a second reference |

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/storage-account) -- the parent account (enable `sftpEnabled` + `isHnsEnabled`)
- [Azure Storage Container](/docs/catalog/azure/storage-container) -- what permission scopes grant access to
- [Azure Storage Share](/docs/catalog/azure/storage-share) -- the FILE-service scope target
