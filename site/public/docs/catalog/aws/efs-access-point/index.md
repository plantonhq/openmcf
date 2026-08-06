---
title: "EFS Access Point"
description: "EFS Access Point deployment documentation"
icon: "package"
order: 100
componentName: "awsefsaccesspoint"
---

# AWS EFS Access Point

Deploys an EFS access point — the application-specific, least-privilege entry point into an Elastic File System. The access point enforces a POSIX user/group identity for every file operation and exposes a chosen directory as the root, so applications (Lambda functions, ECS tasks) get exactly the subtree and identity they need on a shared file system and nothing more.

## What Gets Created

When you deploy an AwsEfsAccessPoint resource, Planton provisions:

- **EFS Access Point** — an `efs.AccessPoint` on the referenced file system, with the configured POSIX identity enforcement and root directory restriction

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An EFS file system** — an [AwsElasticFileSystem](/docs/catalog/aws/elastic-file-system) resource (or a literal `fs-` ID) in the same region

## Quick Start

Create a file `access-point.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEfsAccessPoint
metadata:
  name: app-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsEfsAccessPoint.app-data
spec:
  region: us-east-1
  fileSystemId:
    value: fs-0123456789abcdef0
  posixUser:
    uid: 1000
    gid: 1000
  rootDirectory:
    path: /app-data
    creationInfo:
      ownerUid: 1000
      ownerGid: 1000
      permissions: "0755"
```

Deploy:

```shell
planton apply -f access-point.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region (must match the file system's region). | Required; non-empty |
| `fileSystemId` | `StringValueOrRef` | The Elastic File System this access point enters. Can reference an AwsElasticFileSystem resource via `valueFrom`. ForceNew. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `posixUser.uid` | `int64` | — | POSIX user ID enforced for all file operations (0–4294967295). ForceNew. |
| `posixUser.gid` | `int64` | — | POSIX primary group ID (0–4294967295). ForceNew. |
| `posixUser.secondaryGids` | `int64[]` | `[]` | Up to 16 secondary group IDs for group permission checks. ForceNew. |
| `rootDirectory.path` | `string` | `/` | Absolute path exposed as root (≤4 subdirectories, ≤100 chars). ForceNew. |
| `rootDirectory.creationInfo.ownerUid` | `int64` | — | POSIX UID for the auto-created directory. |
| `rootDirectory.creationInfo.ownerGid` | `int64` | — | POSIX GID for the auto-created directory. |
| `rootDirectory.creationInfo.permissions` | `string` | — | Octal permissions, 3–4 digits (e.g., `755`, `0755`). Required within `creationInfo`. |

The ENTIRE access point (everything except tags) is create-time immutable — changing any field replaces it.

## Examples

### Lambda ML Model Store

Read-mostly model directory for a Lambda function:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEfsAccessPoint
metadata:
  name: lambda-models
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsEfsAccessPoint.lambda-models
spec:
  region: us-east-1
  fileSystemId:
    valueFrom:
      kind: AwsElasticFileSystem
      name: shared-efs
      fieldPath: status.outputs.file_system_id
  posixUser:
    uid: 1001
    gid: 1001
  rootDirectory:
    path: /ml-models
    creationInfo:
      ownerUid: 1001
      ownerGid: 1001
      permissions: "0750"
```

### ECS Task Volume Authorization

The task definition's EFS volume references both nodes:

```yaml
spec:
  volumes:
    - name: app-data
      efs:
        fileSystemId:
          valueFrom:
            kind: AwsElasticFileSystem
            name: shared-efs
            fieldPath: status.outputs.file_system_id
        accessPointId:
          valueFrom:
            kind: AwsEfsAccessPoint
            name: app-data
            fieldPath: status.outputs.access_point_id
        iamAuthorization: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `access_point_id` | `string` | Access point ID (e.g., `fsap-0123456789abcdef0`). ECS EFS volume authorization references this. |
| `access_point_arn` | `string` | Access point ARN. Lambda file system configs and IAM policy conditions reference this. |
| `file_system_id` | `string` | The file system this access point enters. |
| `file_system_arn` | `string` | The file system's ARN for IAM policies. |

## Related Components

- [AwsElasticFileSystem](/docs/catalog/aws/elastic-file-system) — the file system this access point enters
- [AwsLambda](/docs/catalog/aws/lambda) — mounts EFS through the access point ARN
- [AwsEcsTaskDefinition](/docs/catalog/aws/ecs-task-definition) — authorizes EFS volumes through the access point ID
