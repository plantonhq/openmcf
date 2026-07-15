# AWS EFS Access Point: Concepts

An access point is EFS's least-privilege mechanism: an application-specific entry point that enforces a POSIX identity and pins the visible tree to a chosen root directory. This reference covers what an access point actually enforces and how the pieces compose.

## Why Access Points Exist

A shared NFS file system has a classic multi-tenancy problem: every client claims its own UID/GID, and correctness depends on every application managing POSIX permissions perfectly. Access points invert that:

- **Identity enforcement** — every file operation through the access point uses the access point's UID/GID, regardless of what the client claims. A container running as root still reads and writes as the enforced identity.
- **Root directory restriction** — the client sees the access point's `path` as `/`. Files outside that subtree do not exist from the client's perspective.

The result: one file system, many applications, each confined to its own directory and identity — with zero permission management inside the applications.

## Design Notes

- **Create-time immutability.** The file system, POSIX user, and root directory are all fixed at creation; only tags mutate. Changing anything else replaces the access point (its `fsap-` ID changes). Consumers that reference the access point by Planton reference re-resolve automatically on replacement.
- **Directory auto-creation.** `creation_info` (owner UID/GID + octal permissions) lets EFS create a missing root path on first mount. Without it, mounting an access point whose path does not exist fails at the client. Provide it whenever the path is not guaranteed to pre-exist.
- **Secondary GIDs.** Up to 16 supplementary groups (the NFS AUTH_SYS limit) for applications that need access to files owned by multiple groups.
- **Limits.** Up to 1,000 access points per file system; root path up to 4 subdirectories deep and 100 characters.

## Composition

| Consumer | What it references | Why |
|----------|--------------------|-----|
| `AwsLambda.file_system_config.access_point_arn` | `status.outputs.access_point_arn` | Lambda's API requires the access point ARN (not the ID). |
| `AwsEcsTaskDefinition` EFS volume `access_point_id` | `status.outputs.access_point_id` | ECS volume authorization takes the ID; transit encryption is required (the modules always enable it). |
| IAM policies (`elasticfilesystem:AccessPointArn` condition) | `status.outputs.access_point_arn` | Scope `ClientMount`/`ClientWrite` grants to one access point. |

A common pattern pairs the access point with a file system resource policy that requires all access to go through access points:

```json
{
  "Effect": "Deny",
  "Principal": "*",
  "Action": "*",
  "Resource": "*",
  "Condition": {
    "Null": { "elasticfilesystem:AccessPointArn": "true" }
  }
}
```

See the [AwsElasticFileSystem architecture reference](../../awselasticfilesystem/v1/docs/README.md) for the file system side: mount targets, storage classes, throughput modes, encryption, and replication.
