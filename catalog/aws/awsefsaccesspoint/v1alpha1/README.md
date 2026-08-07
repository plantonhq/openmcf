# AwsEfsAccessPoint

An EFS access point — an application-specific entry point into an Elastic File System that enforces a POSIX user/group identity and optionally restricts the visible root directory.

## What It Is

An access point is the least-privilege front door to a shared EFS file system. When a client mounts through an access point, EFS overrides the client's claimed POSIX identity with the access point's enforced UID/GID and pins the visible tree to the access point's root directory — the application cannot wander the file system regardless of what it runs as.

One file system serves many applications, each behind its own access point (up to 1,000 per file system). The access point — not the file system — is what Lambda file-system configs and ECS task-definition EFS volumes reference.

## When to Use It

| Use Case | Description |
|----------|-------------|
| **Lambda file access** | Lambda's file system config requires an access point ARN. Enforce a fixed identity and expose only the needed subtree (e.g., `/ml-models`). |
| **ECS task volumes** | Authorize the task's EFS volume through the access point; the task sees the access point root as `/` and writes as the enforced UID/GID. |
| **Multi-tenant file systems** | Give each application its own directory + identity on one shared file system, without the applications managing POSIX permissions. |

## When NOT to Use It

| Need | Use Instead |
|------|-------------|
| **Whole-file-system admin access** (backups, migrations) | Mount the file system directly from EC2 with an appropriately privileged identity. |
| **Per-client identity preservation** | Mount the file system directly — access points override the client identity by design. |

## Key Facts

- **Create-time immutable.** Everything except tags (the file system, POSIX user, root directory) replaces the access point when changed. Plan the identity and root path upfront.
- **Directory auto-creation.** If `root_directory.path` does not exist, provide `creation_info` — EFS creates the directory with that ownership and permissions on first mount. Without it, mounting a missing path fails.
- **Name lives in tags.** An access point has no name argument; the `Name` tag (from `metadata.name`) is its console display name.

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region (must match the file system's region). |
| `file_system_id` | StringValueOrRef | **Yes** | The Elastic File System this access point enters. **ForceNew**. |
| `posix_user` | PosixUser | No | Enforced UID/GID (+ up to 16 secondary GIDs) for all operations. **ForceNew**. |
| `root_directory` | RootDirectory | No | Path exposed as `/` (absolute, ≤4 subdirectories, ≤100 chars) + optional `creation_info`. **ForceNew**. |

## Outputs

| Field | Type | Description |
|-------|------|-------------|
| `access_point_id` | string | Access point ID (e.g., `fsap-0123456789abcdef0`). ECS EFS volume authorization references this. |
| `access_point_arn` | string | Access point ARN. Lambda file system configs and IAM policy conditions reference this. |
| `file_system_id` | string | The file system this access point enters (for consumers composing on this one node). |
| `file_system_arn` | string | The file system's ARN for IAM policies. |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEfsAccessPoint
metadata:
  name: app-data
  org: my-org
spec:
  region: us-east-1
  fileSystemId:
    valueFrom:
      kind: AwsElasticFileSystem
      name: shared-efs
      fieldPath: status.outputs.file_system_id
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

A Lambda function then mounts it by reference:

```yaml
spec:
  fileSystemConfig:
    accessPointArn:
      valueFrom:
        kind: AwsEfsAccessPoint
        name: app-data
        fieldPath: status.outputs.access_point_arn
    localMountPath: /mnt/data
```

See [docs/README.md](docs/README.md) for the underlying concepts and [AwsElasticFileSystem](../../awselasticfilesystem/v1alpha1/README.md) for the file system itself.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
