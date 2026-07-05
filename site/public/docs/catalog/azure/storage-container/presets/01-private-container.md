---
title: "Private Container"
description: "This preset creates a private blob container -- the default posture: every read requires authorization (Entra data-plane roles or the account's keys). The right starting point for everything that is..."
type: "preset"
rank: "01"
presetSlug: "01-private-container"
componentSlug: "storage-container"
componentTitle: "Storage Container"
provider: "azure"
icon: "package"
order: 1
---

# Private Container

This preset creates a private blob container -- the default posture:
every read requires authorization (Entra data-plane roles or the
account's keys). The right starting point for everything that is not a
public website or CDN origin.

## When to Use

- Application data domains: uploads, artifacts, exports, backups
- Any container whose objects should never be anonymously readable

## Key Configuration Choices

- **`containerAccessType: PRIVATE`** -- no anonymous access; grant
  Storage Blob Data Reader/Contributor on the container's `container_id`
  output for scoped access
- **`metadata.purpose`** -- a self-documenting marker for operators
  browsing the account

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<container-name>` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<data-domain>` | What lives in this container | Your data taxonomy |

## Downstream Wiring

Scope a data-plane grant to just this container:

```yaml
# On an AzureRoleAssignment
scope:
  valueFrom:
    kind: AzureStorageContainer
    name: my-app-uploads
    fieldPath: status.outputs.container_id
roleDefinitionName: Storage Blob Data Contributor
```
