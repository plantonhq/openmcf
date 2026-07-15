# AWS MemoryDB ACL

Creates a MemoryDB Access Control List — the set of users a cluster
authenticates against, and the single place an application's cluster access
is granted or revoked.

## What Gets Created

When you deploy an AwsMemorydbAcl resource, Planton provisions:

- **MemoryDB ACL** — an `aws_memorydb_acl` / `memorydb.Acl` whose AWS name
  is `metadata.name` (create-time immutable, max 40 characters), with
  membership referencing AwsMemorydbUser resources

Membership edits apply in place: AWS diffs the user set on update, so
granting one application access never disturbs the others.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless
  SSO/OIDC).
- **MemoryDB users** (AwsMemorydbUser) for authenticated access — or deploy
  an empty ACL first and add members later.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
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
    - valueFrom:
        kind: AwsMemorydbUser
        name: analytics-service
        fieldPath: status.outputs.user_name
```

Attach the ACL to a cluster:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsMemorydbCluster
metadata:
  name: payments-sessions
spec:
  region: us-west-2
  engine: valkey
  nodeType: db.r7g.large
  aclName:
    valueFrom:
      kind: AwsMemorydbAcl
      name: payments-env-acl
      fieldPath: status.outputs.acl_name
```

## Notes

- **Empty ACLs are valid** — MemoryDB has no mandatory "default" member. A
  cluster attached to an empty ACL accepts no authenticated connections, so
  production ACLs carry one user per application.
- **"open-access" is built in** — the allow-everything, no-authentication
  ACL always exists in the account; clusters reference it by literal value
  (`aclName: { value: open-access }`), never through this resource.
- Every member user must live in the ACL's region (users are regional).

## Stack Outputs

| Output | Description |
|--------|-------------|
| `acl_name` | What clusters attach via their `aclName` (same as `metadata.name`) |
| `acl_arn` | The ACL's ARN, for IAM policies |
| `minimum_engine_version` | The minimum engine version the ACL's combined user set requires |
