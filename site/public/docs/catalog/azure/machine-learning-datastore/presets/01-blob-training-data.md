---
title: "Blob Training Data"
description: "This preset registers a blob-container datastore under workspace-identity auth -- the credential-free posture: no key or SAS in the manifest, one role assignment on the container instead."
type: "preset"
rank: "01"
presetSlug: "01-blob-training-data"
componentSlug: "machine-learning-datastore"
componentTitle: "Machine Learning Datastore"
provider: "azure"
icon: "package"
order: 1
---

# Blob Training Data

This preset registers a blob-container datastore under workspace-identity auth -- the credential-free posture: no key or SAS in the manifest, one role assignment on the container instead.

## When to Use

- Registering curated training datasets held in blob containers
- Estates standardizing on identity-based data access (no secret rotation)
- The everyday datastore shape for most ML pipelines

## Key Configuration Choices

- **`WORKSPACE_SYSTEM_ASSIGNED_IDENTITY`** -- the workspace's own identity reads the data; grant it Storage Blob Data Reader (or Contributor for writes) on the container BEFORE jobs run
- **No credentials** -- valid precisely because the identity mode covers access; with `serviceDataIdentity` unset or NONE, an account key or SAS becomes required (validated at manifest time)
- **`isDefault` left false** -- claiming the default re-points where job outputs land workspace-wide; flip it deliberately

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-machine-learning-workspace-id>` | ARM ID of the parent workspace | `AzureMachineLearningWorkspace` status outputs (`machine_learning_workspace_id`), or reference it with valueFrom |
| `<your-storage-container-id>` | ARM ID of the blob container (`.../blobServices/default/containers/{name}`) | `AzureStorageContainer` status outputs (`container_id`), or reference it with valueFrom |

The `name` field carries a realistic example (`training_data`) because the datastore name is pattern-validated -- replace it with your own.
