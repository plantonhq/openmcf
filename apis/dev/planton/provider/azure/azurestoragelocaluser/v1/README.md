# AzureStorageLocalUser

A local user on an AzureStorageAccount: the credential identity the
account's SFTP endpoint authenticates. Partners, legacy pipelines, and
managed-file-transfer tools that only speak SFTP reach blob storage
through local users -- each with its own SSH credentials, home
directory, and per-container permission scopes, so one account serves
many isolated exchange partners.

## When to Use

Use AzureStorageLocalUser when you need:

- **Partner file exchange** -- one user per partner, each scoped to its
  own container, onboarded and offboarded independently
- **Legacy SFTP pipelines** -- systems that cannot speak the blob API
  land files straight into blob storage
- **Key-based automation** -- CI/MFT systems authenticate with SSH
  public keys; no shared secret to distribute

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation); the account needs
  `sftp_enabled: true` (which requires `is_hns_enabled`) for logins to
  work
- `user_name` -- 3-64 lowercase letters/digits; the SFTP login is
  `{account}.{user}` (the `sftp_username` output)
- `ssh_key_enabled` + `ssh_authorized_keys` -- public-key auth (the
  posture to prefer); `ssh_password_enabled` -- an Azure-GENERATED
  password, returned exactly once in the `password` output
- `permission_scopes` -- read/write/delete/list/create per container
  (or file share); a user with no scopes can log in but touch nothing
- `home_directory` -- where a session lands after login

## Composition

```yaml
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

The `sid` and `password` outputs are secret-bearing -- hand the
password to the partner over a secure channel.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
