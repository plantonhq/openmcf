# AwsFsxDataRepositoryAssociation

Amazon FSx data repository association — the composable link between a directory on an FSx for Lustre file system and an S3 bucket or prefix, with independent bidirectional sync policies.

## What It Is

A data repository association maps one file-system path (e.g., `/datasets/2026`) to one S3 data repository (e.g., `s3://training-data/2026/`). Object metadata appears in the Lustre namespace as files; file data is lazy-loaded from S3 on first access and served at Lustre speed afterwards. Auto-import keeps the namespace in sync as objects change in S3; auto-export writes file-system changes back to the bucket.

Associations are the modern S3 integration for Lustre — and the **only** one available on PERSISTENT_2 file systems (the legacy in-spec `import_path` arm is limited to SCRATCH and PERSISTENT_1 generations). Each file system supports up to 8 associations (25 per account), each with its own lifecycle: create, update, and delete links without ever touching the file system.

## When to Use It

| Use Case | Description |
|----------|-------------|
| **ML training data** | Link S3 training corpora into a PERSISTENT_2 file system; jobs read at Lustre speed, new objects appear automatically. |
| **Results export** | Link an output directory with auto-export so processed results land back in S3 as they are written. |
| **Multi-dataset layouts** | One association per dataset/prefix — different sync policies per link, added and removed independently. |
| **S3-backed archives with POSIX access** | Cold data stays in S3; the namespace stays browsable and data hydrates on access. |

## When NOT to Use It

| Need | Use Instead |
|------|-------------|
| **A one-time S3 import on a scratch file system** | The Lustre spec's legacy `import_path`/`export_path` arm. |
| **Object storage without POSIX** | Amazon S3 directly. |

## Prerequisites

- An **AwsFsxLustreFileSystem** that does NOT use the legacy `import_path` arm (AWS forbids mixing the two S3-link generations).
- An **S3 bucket** the FSx service can read (and write, for auto-export).

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | The file system's AWS region. |
| `file_system_id` | StringValueOrRef | **Yes** | The Lustre file system to attach to. **ForceNew**. |
| `file_system_path` | string | **Yes** | Path on the file system (begins with `/`, max 4096 chars). Must not overlap other associations. **ForceNew**. |
| `data_repository_path` | string | **Yes** | S3 URI with optional prefix (`s3://bucket/prefix/`, 3–900 chars). **ForceNew**. |
| `auto_import_events` | []string | No | S3 → Lustre sync events: `NEW`, `CHANGED`, `DELETED`. Empty = creation-time state only. In place. |
| `auto_export_events` | []string | No | Lustre → S3 sync events: `NEW`, `CHANGED`, `DELETED`. Empty = manual export tasks only. In place. |
| `imported_file_chunk_size` | int32 | No | Stripe size in MiB for imported files (1–512000; AWS default 1024). In place. |
| `batch_import_meta_data_on_create` | bool | No | Import the existing S3 metadata at creation (otherwise only later changes appear). |
| `delete_data_in_filesystem` | bool | No | Delete the linked files from the file system when the association is deleted (default keeps them). |

## Outputs

| Field | Type | Description |
|-------|------|-------------|
| `association_id` | string | AWS-assigned association ID (`dra-...`). |
| `association_arn` | string | ARN for IAM resource-level permissions. |
| `file_system_id` | string | The attached file system's ID, echoed for composition. |

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsFsxDataRepositoryAssociation
metadata:
  name: training-data-link
  org: my-org
spec:
  region: us-west-2
  file_system_id:
    valueFrom:
      kind: AwsFsxLustreFileSystem
      name: ml-training-fsx
      fieldPath: status.outputs.file_system_id
  file_system_path: /datasets/2026
  data_repository_path: s3://training-data/2026/
  auto_import_events: [NEW, CHANGED, DELETED]
  batch_import_meta_data_on_create: true
```

## ForceNew Warnings

The association's identity is the (file system, path, bucket) triple — changing `file_system_id`, `file_system_path`, or `data_repository_path` replaces the association. The sync policies and chunk size update in place.

See [docs/README.md](docs/README.md) for architecture details.
