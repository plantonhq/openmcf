# AWS ElastiCache User

Creates one Redis/Valkey RBAC identity — a user id (`metadata.name`), an
AUTH user name, an ACL access string, and an authentication mode — that
composes into user groups and attaches to caches without ever modifying
the cache itself when access changes.

## What Gets Created

When you deploy an AwsElasticacheUser resource, Planton provisions:

- **ElastiCache user** — an `aws_elasticache_user` / `elasticache.User`
  with the chosen engine, access string, and authentication mode; the AWS
  user id is `metadata.name` (create-time immutable)

The user never modifies resources it merely references: user groups carry
their own membership lists, and caches reference groups — credential
rotation is an in-place update on the user.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **A matching engine** — `redis` or `valkey`; Memcached has no RBAC users.
- **For IAM auth** — transit encryption enabled on the target cache, and
  `userName` equal to `metadata.name`.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticacheUser
metadata:
  name: orders-service
spec:
  region: us-west-2
  engine: redis
  userName: orders-service
  accessString: "on ~orders:* +@read +@write"
  authenticationMode:
    type: password
    passwords:
      - "<strong-password-from-secret-manager>"
```

```shell
planton apply -f user.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; users are regional. | Required; non-empty |
| `engine` | `string` | `redis` or `valkey`. Create-only. | Required; CEL-enforced |
| `userName` | `string` | Name clients present in AUTH. Create-only. | Required |
| `accessString` | `string` | Redis ACL rule list. Updates in place. | Required |
| `authenticationMode.type` | `string` | `password`, `iam`, or `no-password-required`. | Required; CEL-enforced |
| `authenticationMode.passwords` | `string[]` | 1–2 passwords when type is `password`; empty otherwise. | CEL-enforced; sensitive |

### Identity Notes

| Concept | Source | Notes |
| --- | --- | --- |
| AWS user id | `metadata.name` | Create-time immutable; what user groups reference |
| AUTH user name | `spec.userName` | Need not be unique; several users may share a name |
| Mandatory "default" user | `userName: default` | Every user group must include one |

## Examples

### IAM-authenticated user

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticacheUser
metadata:
  name: platform-cache-reader
spec:
  region: us-west-2
  engine: redis
  userName: platform-cache-reader
  accessString: "on ~* +@read -@dangerous"
  authenticationMode:
    type: iam
```

### Locked-down default user (mandatory group member)

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticacheUser
metadata:
  name: rbac-default-user
spec:
  region: us-west-2
  engine: redis
  userName: default
  accessString: "off ~* +@all"
  authenticationMode:
    type: no-password-required
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `user_id` | The user's AWS identifier (same as `metadata.name`) |
| `arn` | The user's ARN — for IAM `elasticache:Connect` policies |
| `user_name` | The name clients present in the AUTH command |

## Related Components

- [AwsElasticacheUserGroup](/docs/catalog/aws/awselasticacheusergroup) — collects users and attaches to caches
- [AwsRedisElasticache](/docs/catalog/aws/awsrediselasticache) — replication groups that reference user groups
