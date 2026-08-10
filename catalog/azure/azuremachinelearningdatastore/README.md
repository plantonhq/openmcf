# Overview

The **AzureMachineLearningDatastore** component registers a datastore on an Azure Machine Learning workspace -- the saved connection that tells the workspace where data lives and how to reach it. ONE kind covers the three storage flavors as variants: a blob container, a Data Lake Gen2 filesystem, or an Azure Files share. Exactly one variant block is set; the block IS the datastore type.

## Purpose

- **Data connections as declarative infrastructure**: where training data lives, which credentials reach it -- reviewed and versioned like everything else.
- **One kind, three flavors**: blob / data-lake / file-share variants in one contract instead of three near-identical kinds.
- **Typed references end-to-end**: the workspace, container, filesystem, and share all wire by reference -- chart-ready.
- **Credential honesty**: keys, SAS tokens, and client secrets are sensitive fields referencing secrets, never returned by ARM, and recorded as write-normalized for imports.

## Key Features

- Full azurerm v5 surface across all three variant resources (`_blobstorage`, `_datalake_gen2`, `_fileshare`), each variant's auth grammar validated at manifest time: blob needs a key or SAS unless a workspace-identity mode covers it; the file share requires exactly one credential; the data lake takes the service-principal triad all-or-none.
- The provider's per-variant identity-argument rename (`service_data_auth_identity` vs `service_data_identity`) folded into ONE spec field.
- `is_default` modeled where the provider allows it (blob only) and read back elsewhere.

## Use Cases

- **Training data registration**: point the workspace at the container holding curated datasets.
- **Lakehouse access**: register a Data Lake Gen2 filesystem with service-principal auth.
- **Shared file workloads**: an Azure Files share for jobs that need SMB-style file semantics.

## Future Enhancements

- Workspace-connection-based datastores as the service's connection surface settles.
