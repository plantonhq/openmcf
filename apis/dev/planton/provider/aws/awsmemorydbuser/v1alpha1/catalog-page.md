# AWS MemoryDB User

Creates one MemoryDB ACL identity — a user name (`metadata.name`), a Redis
ACL access string, and an authentication mode — that composes into ACLs
(AwsMemorydbAcl) and reaches clusters without ever modifying the cluster
when access changes.

## What Gets Created

When you deploy an AwsMemorydbUser resource, Planton provisions:

- **MemoryDB user** — an `aws_memorydb_user` / `memorydb.User` with the
  chosen access string and authentication mode; the AWS user name is
  `metadata.name` (create-time immutable, max 40 characters)

The user never modifies resources it merely references: ACLs carry their
own membership lists, and clusters reference ACLs — credential rotation is
an in-place update on the user.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless
  SSO/OIDC).
- **For IAM auth** — TLS enabled on the target cluster, and the connecting
  principal granted `memorydb:Connect` on both the user ARN and the cluster
  ARN.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsMemorydbUser
metadata:
  name: orders-service
spec:
  region: us-west-2
  accessString: "on ~orders:* +@read +@write"
  authenticationMode:
    type: password
    passwords:
      - use-a-16-plus-character-secret
```

Grant the user cluster access by adding it to an ACL:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsMemorydbAcl
metadata:
  name: payments-env-acl
spec:
  region: us-west-2
  userNames:
    - valueFrom:
        kind: AwsMemorydbUser
        name: orders-service
        fieldPath: status.outputs.user_name
```

## Authentication Modes

| Mode | Credential | Best for |
|------|-----------|----------|
| `password` | 1–2 passwords, 16–128 chars each | Standard application access; two passwords rotate with zero downtime |
| `iam` | Short-lived IAM-signed tokens | Zero long-lived secrets; requires TLS on the cluster |

## Access Strings

The same syntax as Redis ACL SETUSER:

| Shape | Meaning |
|-------|---------|
| `on ~* +@all` | Full access (an admin user) |
| `on ~app1:* +@read +@write` | Read/write scoped to one key prefix |
| `on ~* +@read -@dangerous` | Read-only, no admin commands |
| `off ~* +@all` | A disabled user (kept for later re-enable) |

Access strings update in place — tightening permissions takes effect on new
connections without recreating the user.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `user_name` | The AUTH identity ACLs reference (same as `metadata.name`) |
| `user_arn` | The user's ARN, for IAM `memorydb:Connect` policies |
| `minimum_engine_version` | The minimum engine version the user's configuration requires |
