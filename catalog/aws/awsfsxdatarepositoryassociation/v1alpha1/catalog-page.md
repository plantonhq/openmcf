# AWS FSx Data Repository Association

Link a directory on an FSx for Lustre file system to an S3 bucket or prefix, with independent bidirectional sync policies — the modern S3 integration for Lustre and the only one available on PERSISTENT_2 file systems.

## Prerequisites

- **An FSx for Lustre file system** (reference its `file_system_id` output). The file system must not use the legacy in-spec `import_path` link.
- **An S3 bucket** readable (and writable, for auto-export) by the FSx service.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsFsxDataRepositoryAssociation
metadata:
  name: training-data-link
spec:
  region: us-west-2
  fileSystemId:
    valueFrom:
      kind: AwsFsxLustreFileSystem
      name: ml-training-fsx
      fieldPath: status.outputs.file_system_id
  fileSystemPath: /datasets/2026
  dataRepositoryPath: s3://training-data/2026/
  autoImportEvents: [NEW, CHANGED, DELETED]
  batchImportMetaDataOnCreate: true
```

Deploy:

```shell
planton apply -f data-repository-association.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The file system's AWS region. | Required; non-empty |
| `fileSystemId` | `StringValueOrRef` | The Lustre file system to attach to. Default kind: `AwsFsxLustreFileSystem`. ForceNew. | Required |
| `fileSystemPath` | `string` | Path on the file system (e.g., `/datasets`). Must not overlap other associations. ForceNew. | Begins with `/`, max 4096 chars |
| `dataRepositoryPath` | `string` | S3 URI with optional prefix. ForceNew. | `s3://` URI, 3–900 chars |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `autoImportEvents` | `string[]` | `[]` | S3 → Lustre sync: `NEW`, `CHANGED`, `DELETED`. Empty means creation-time state only. |
| `autoExportEvents` | `string[]` | `[]` | Lustre → S3 sync: `NEW`, `CHANGED`, `DELETED`. Empty means manual export tasks only. |
| `importedFileChunkSize` | `int32` | `1024` | Stripe size in MiB for imported files (1–512000). |
| `batchImportMetaDataOnCreate` | `bool` | `false` | Import the existing S3 metadata at creation. |
| `deleteDataInFilesystem` | `bool` | `false` | Delete the linked files from the file system when the association is deleted. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `association_id` | AWS-assigned association ID (`dra-...`), the join key for FSx data repository tasks. |
| `association_arn` | ARN for IAM resource-level permissions. |
| `file_system_id` | The attached file system's ID, echoed for composition. |

## Related Components

- [AWS FSx Lustre File System](/docs/catalog/aws/awsfsxlustrefilesystem) — the file system associations attach to.
- [AWS S3 Bucket](/docs/catalog/aws/awss3bucket) — the data repository.
