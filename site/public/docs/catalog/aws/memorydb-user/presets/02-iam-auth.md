---
title: "IAM-Authenticated Application User"
description: "This preset creates a MemoryDB user that authenticates with short-lived IAM-signed tokens instead of a password — no long-lived secret exists anywhere. The connecting workload signs an auth token..."
type: "preset"
rank: "02"
presetSlug: "02-iam-auth"
componentSlug: "memorydb-user"
componentTitle: "MemoryDB User"
provider: "aws"
icon: "package"
order: 2
---

# IAM-Authenticated Application User

This preset creates a MemoryDB user that authenticates with short-lived
IAM-signed tokens instead of a password — no long-lived secret exists
anywhere. The connecting workload signs an auth token with its own AWS
identity (an EC2 instance profile, an ECS task role, a Lambda execution
role) and presents it in the AUTH command.

## When to Use

- Workloads already running with an AWS IAM identity — zero secret
  distribution and zero rotation burden
- Security postures that forbid long-lived database credentials
- Clusters with TLS enabled (IAM authentication requires it)

## Key Configuration Choices

- **`type: iam`** — no password material; CEL rejects passwords on this
  mode
- **`accessString`** — least-privilege still applies: scope the key
  patterns and command categories exactly as for password users

## IAM Wiring (outside this resource)

The connecting principal needs `memorydb:Connect` on BOTH the user ARN
(`status.outputs.user_arn`) and the cluster ARN
(`status.outputs.cluster_arn` on AwsMemorydbCluster).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-name>` | The AUTH identity (e.g. `orders-service`); becomes `metadata.name`, max 40 chars | Your naming convention |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<key-prefix>` | Key namespace this user may touch (e.g. `orders`) | Your data model |

## Related Presets

- **01-password-auth** — password authentication for workloads without an
  AWS identity
