# AzureDiskSnapshot Terraform Module

## Overview

Creates a managed disk snapshot -- a point-in-time copy of a disk used for backup, cloning, and as the source of gallery image versions.

## Resources Created

- `azurerm_snapshot` -- the snapshot (source pairing per create_option, incremental mode, network posture, optional legacy ADE encryption settings, tags)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDiskSnapshotSpec fields; the resource-group, source-disk, storage-account, disk-access, and Key Vault references arrive as resolved literals

## Outputs

- `snapshot_id` -- the snapshot's ARM resource ID (what disks restore from and gallery image versions build from)
- `snapshot_name` -- the snapshot's name

## Behavior Notes

- **The source pairing is Azure-validated, not schema-validated**: the provider's own schema does not tie `create_option` to its source fields, so neither does the spec. The working pairs: "Copy" reads `source_resource_id` (a disk or another snapshot); "Import" reads `source_uri` + `storage_account_id`. The module sends each source field only when set.
- **Incremental snapshots store only the delta** since the previous snapshot of the same disk, on standard storage -- the right default for backup chains. `incremental_enabled` is fixed at creation.
- **Encryption settings are one-way**: removing a previously set `encryption_settings` block forces replacement (Azure cannot disable encryption in place -- the provider's CustomizeDiff, mirrored here as behavior).
- **Unset optionals ride the provider defaults**: `network_access_policy` "AllowAll", `public_network_access_enabled` true, `disk_size_gb` inherited from the source.
- **A snapshot bills only for the storage it holds** -- an incremental snapshot of a small, mostly-empty disk costs effectively nothing.
