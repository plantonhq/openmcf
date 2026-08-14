---
title: "Data Lake Filesystem"
description: "This preset registers a Data Lake Gen2 filesystem datastore with service-principal auth -- the lakehouse access shape: a Gen2 filesystem on a hierarchical-namespace account, read through an app..."
type: "preset"
rank: "02"
presetSlug: "02-datalake-filesystem"
componentSlug: "machine-learning-datastore"
componentTitle: "Machine Learning Datastore"
provider: "azure"
icon: "package"
order: 2
---

# Data Lake Filesystem

This preset registers a Data Lake Gen2 filesystem datastore with service-principal auth -- the lakehouse access shape: a Gen2 filesystem on a hierarchical-namespace account, read through an app registration.

## When to Use

- Lakehouse estates where curated data lives in Data Lake Gen2 filesystems
- Cross-tenant or tightly-scoped access through a dedicated app registration
- Pipelines that need POSIX-style ACL semantics on the data side

## Key Configuration Choices

- **The service-principal triad** -- `tenantId`, `clientId`, and `clientSecret` come together or not at all (validated at manifest time, mirroring the provider); drop all three to fall back to workspace-identity or credential-free registration
- **The filesystem IS a container** -- Gen2 filesystems are blob containers on HNS-enabled accounts; the reference target is the filesystem kind's `filesystem_id` output
- **Grant on the data side** -- the service principal needs Storage Blob Data Reader (or ACLs) on the filesystem; the datastore stores the credential, it does not grant anything

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-machine-learning-workspace-id>` | ARM ID of the parent workspace | `AzureMachineLearningWorkspace` status outputs (`machine_learning_workspace_id`), or reference it with valueFrom |
| `<your-datalake-filesystem-id>` | ARM ID of the Gen2 filesystem (`.../blobServices/default/containers/{name}`) | `AzureStorageDataLakeGen2Filesystem` status outputs (`filesystem_id`), or reference it with valueFrom |
| `<your-client-secret>` | The app registration's client secret | Your Entra ID app registration -- reference a secret rather than embedding the literal |

The `tenantId` and `clientId` fields carry realistic example UUIDs because both are UUID-validated -- replace them with your app registration's actual IDs. The `name` field carries a realistic example (`lake_features`) for the same reason.
