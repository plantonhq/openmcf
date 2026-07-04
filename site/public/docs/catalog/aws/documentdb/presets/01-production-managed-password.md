---
title: "Production DocumentDB (Managed Password)"
description: "This preset creates a production-shaped DocumentDB cluster: one writer and one reader instance on shared cluster storage, an AWS-managed master password in Secrets Manager, encrypted storage,..."
type: "preset"
rank: "01"
presetSlug: "01-production-managed-password"
componentSlug: "documentdb"
componentTitle: "DocumentDB"
provider: "aws"
icon: "package"
order: 1
---

# Production DocumentDB (Managed Password)

This preset creates a production-shaped DocumentDB cluster: one writer and one reader instance on shared cluster storage, an AWS-managed master password in Secrets Manager, encrypted storage, deletion protection, seven days of continuous backup, and audit logs streaming to CloudWatch. The reader doubles as the failover target -- DocumentDB promotes it in seconds because it already reads from the same storage volume.

## When to Use

- Production MongoDB-compatible workloads with steady, predictable capacity needs
- Applications that need a reader endpoint for report/analytics traffic without touching the writer
- Anywhere failover time matters: a replica promotes in seconds

## Key Configuration Choices

- **Managed master password** (`manageMasterUserPassword: true`) -- AWS generates, stores, and rotates the credential in Secrets Manager; no secret in the manifest or IaC state. The secret's ARN is exported as `master_user_secret_arn`.
- **Writer + reader instances** (`instances`) -- each entry is its own managed resource; add readers by appending entries, never touching the cluster. `promotionTier: 1` on the reader makes it the failover target.
- **Encrypted storage** (`storageEncrypted: true`) -- a create-time one-way door; an unencrypted cluster cannot be encrypted later.
- **Deletion safety** -- `deletionProtection` plus a named final snapshot: a delete is a deliberate two-step, never an accident.
- **Audit logging, both halves** -- the `audit` log export only emits once the `audit_logs` cluster parameter is enabled, so the preset sets both.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 27017 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **02-serverless** -- Use instead when traffic is variable or spiky and capacity should track demand
