# AzureStorageContainer

A blob container inside an AzureStorageAccount: the namespace unit of
Azure blob storage. Applications organize objects into containers per
data domain (uploads, logs, backups, artifacts), and Azure scopes
anonymous access, data-plane RBAC roles, encryption scopes, and
lifecycle prefixes at the container level.

## When to Use

Use AzureStorageContainer when you need:

- **A namespace for application objects** -- one container per data
  domain on a shared account
- **A public website/CDN origin** -- `container_access_type: BLOB`
  serves objects anonymously by direct URL
- **A role-assignment scope** -- grant Storage Blob Data
  Reader/Contributor on `container_id` for container-level access
- **Sub-account key isolation** -- pin a `default_encryption_scope` so
  a tenant's blobs encrypt under their own key

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation)
- `container_name` -- 3-63 lowercase letters/digits/hyphens, unique
  within the account; becomes the URL path segment
- `container_access_type` -- PRIVATE (default), BLOB (anonymous object
  reads -- also requires the account to permit nested public access), or
  CONTAINER (adds anonymous listing; rarely appropriate)
- `metadata` -- free-form key/value pairs on the container (not
  secrets; not Azure tags)

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Blob URLs compose from the ACCOUNT's endpoint plus this container's
name: `{primary_blob_endpoint}{container_name}/{blob}`.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
