# AWS ElastiCache User Group

Collects Redis/Valkey RBAC users into the attachment unit caches reference —
a user group id (`metadata.name`), an engine, and a membership list wired
from `AwsElasticacheUser` outputs. Grant or revoke cache access by editing
group membership; the cache itself never changes.

## What Gets Created

When you deploy an AwsElasticacheUserGroup resource, Planton provisions:

- **ElastiCache user group** — an `aws_elasticache_user_group` /
  `elasticache.UserGroup` with the chosen engine and member user ids; the
  AWS user group id is `metadata.name` (create-time immutable)

The group never modifies users it merely references — membership updates
replace the group's user id list in place.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **Member users** (`AwsElasticacheUser`) in the same region and engine,
  including one whose `userName` is exactly `default`.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsElasticacheUserGroup
metadata:
  name: orders-rbac
spec:
  region: us-west-2
  engine: redis
  userIds:
    - valueFrom:
        kind: AwsElasticacheUser
        name: rbac-default-user
        fieldPath: status.outputs.user_id
    - valueFrom:
        kind: AwsElasticacheUser
        name: orders-service
        fieldPath: status.outputs.user_id
```

```shell
planton apply -f user-group.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; must match member users and target caches. | Required; non-empty |
| `engine` | `string` | `redis` or `valkey`. Create-only. | Required; CEL-enforced |
| `userIds` | `(string \| valueFrom)[]` | Member user ids; at least one entry. | Required; min 1; defaults to `AwsElasticacheUser` `user_id` output |

### Membership Rules

| Rule | Enforcement |
| --- | --- |
| A user whose `userName` is `default` must be included | AWS at group-create time |
| Every member's `engine` and `region` must match the group | AWS at group-create time |
| Membership updates in place | AWS API |

## Examples

### Valkey user group

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsElasticacheUserGroup
metadata:
  name: platform-valkey-rbac
spec:
  region: us-west-2
  engine: valkey
  userIds:
    - valueFrom:
        kind: AwsElasticacheUser
        name: valkey-default-user
        fieldPath: status.outputs.user_id
    - valueFrom:
        kind: AwsElasticacheUser
        name: analytics-service
        fieldPath: status.outputs.user_id
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `user_group_id` | The group's AWS identifier (same as `metadata.name`) |
| `arn` | The group's ARN — for IAM policies |

## Related Components

- [AwsElasticacheUser](/docs/catalog/aws/awselasticacheuser) — RBAC identities referenced in `userIds`
- [AwsRedisElasticache](/docs/catalog/aws/awsrediselasticache) — replication groups that attach user groups
