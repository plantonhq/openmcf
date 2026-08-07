---
title: "ElastiCache User"
description: "ElastiCache User deployment documentation"
icon: "package"
order: 100
componentName: "awselasticacheuser"
---

# AWS ElastiCache User

Deploys an ElastiCache RBAC user — one identity in the Redis/Valkey Role-Based Access Control system. RBAC is AWS's recommended authentication model for ElastiCache: instead of one shared AUTH token for every client, each application gets its own user with an access string scoping exactly which commands and keys it may touch. Rotating one application's credentials or revoking its access never disturbs the others. Users join user groups ([AwsElasticacheUserGroup](/cloud-catalog/aws-elasticache-user-group)), and a group attaches to a replication group or serverless cache. The user integrates with Planton's Provider Connections for AWS credential management and keeps password material in managed secrets — never in the manifest.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ElastiCache User** -- one RBAC identity whose AWS user id is the resource name (create-time immutable)
- **ACL Access String** -- the Redis `ACL SETUSER` rule list scoping keys and command categories; tightening it later applies in place
- **Authentication Mode** -- exactly one mechanism: password (1–2 secrets, enabling zero-downtime rotation), IAM-signed tokens, or no credential (for the disabled default user)
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Managed secrets for passwords** -- when using password authentication, create the password as an org secret first; the spec carries a `$secret/<slug>` reference and the runner resolves it just-in-time at deploy. Each password's value must be 16–128 printable characters.

### AWS Account

- **A Redis or Valkey posture** -- RBAC applies to the Redis and Valkey engines only (Memcached has no users). The user's engine must match the user group and cache it will serve.
- **For IAM authentication** -- the user name must equal the user id (the resource name), the attached cache must have transit encryption (TLS) enabled, and the client's IAM principal needs `elasticache:Connect` on both the user ARN and the cache ARN.

## Deploy

### Console

Open the deployment store, find **AWS ElastiCache User**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Password Auth** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticacheUser
metadata:
  name: orders-service
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  engine: redis
  userName: orders-service
  accessString: "on ~orders:* +@read +@write"
  authenticationMode:
    type: password
    passwords:
      - $secret/orders-cache-password
```

```shell
planton apply -f elasticache-user.yaml
```

This creates a password-authenticated user scoped to the `orders:` key prefix with read/write command categories. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, users deploy first, then the group that collects them, then the cache that attaches the group. The user itself has no upstream references — downstream groups reference its `user_id` output:

```yaml
# In the AwsElasticacheUserGroup manifest:
spec:
  userIds:
    - valueFrom:
        kind: AwsElasticacheUser
        name: orders-service
        fieldPath: status.outputs.user_id
```

## Key Configuration

These are the most important decisions when configuring an ElastiCache user. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Two names, two jobs** -- the resource name is the AWS user id (what groups reference; create-time immutable), while `userName` is what clients present in the AUTH command (also create-time immutable). User names need not be unique — AWS unions the credentials of same-named users at authentication time, which is how zero-downtime credential migrations work.

**The access string** -- the same syntax as the Redis `ACL SETUSER` rule list: the on/off switch enables or disables the user, key patterns (`~orders:*`) grant matching keys, and command categories (`+@read`, `-@dangerous`) grant or subtract command groups. Updates apply in place on new connections. Scope production applications to their key prefix; reserve `~* +@all` for admin identities.

**Authentication** -- exactly one mechanism. `password` works everywhere (1–2 secret references; two enable rotation: add the new password, roll clients, remove the old). `iam` replaces long-lived secrets with short-lived tokens signed by the workload's AWS identity — it requires `userName` to equal the resource name and a TLS-enabled cache. `no-password-required` exists for the disabled default user and migration windows — never for a production user that is "on". The whole authentication mode edits in place.

**The mandatory "default" user** -- every user group must contain a user whose user NAME is exactly `default`; it defines what unauthenticated clients may do. The production pattern is a locked-down default user (`accessString: "off ~* +@all"`, `type: no-password-required`) so anonymous connections are rejected outright — the **Disabled Default** preset ships exactly this shape.

## Outputs and Dependencies

### What This Component Consumes

This component has no upstream Cloud Resource dependencies — it is a leaf the RBAC graph builds on.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `user_id` | The user's AWS identifier (the resource name) | AwsElasticacheUserGroup membership (`userIds`) |
| `arn` | Amazon Resource Name of the user | IAM policies — `elasticache:Connect` for IAM-authenticated clients |
| `user_name` | The name clients present in AUTH | Application configuration wired from the resource graph |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Password-authenticated application user** -- one user per service, scoped to its own key prefix, authenticating with a managed-secret password. Start from the **Password Auth** preset.

**IAM-authenticated user** -- workloads on ECS/EKS/Lambda with an IAM role and a TLS-enabled cache skip password rotation entirely: short-lived signed tokens, no secret material anywhere. Start from the **IAM Auth** preset.

**Disabled default user** -- the mandatory `default`-named member of every user group, switched off with no credential so anonymous connections are rejected. Start from the **Disabled Default** preset.

## Works With

- [**AWS ElastiCache User Group**](/cloud-catalog/aws-elasticache-user-group) -- collects users into the RBAC attachment unit (references `user_id`)
- [**AWS ElastiCache Serverless**](/cloud-catalog/aws-serverless-elasticache) -- the cache the user ultimately authenticates against, via its attached user group
