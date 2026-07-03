---
title: "PostgreSQL (Production Multi-AZ)"
description: "This preset creates a production-shaped PostgreSQL instance: Multi-AZ with a synchronous standby and automatic failover, gp3 storage with autoscaling headroom, an AWS-managed master password in..."
type: "preset"
rank: "01"
presetSlug: "01-postgresql-production"
componentSlug: "rds-instance"
componentTitle: "RDS Instance"
provider: "aws"
icon: "package"
order: 1
---

# PostgreSQL (Production Multi-AZ)

This preset creates a production-shaped PostgreSQL instance: Multi-AZ with a synchronous standby and automatic failover, gp3 storage with autoscaling headroom, an AWS-managed master password in Secrets Manager, encrypted storage, deletion protection, seven days of backups, and Performance Insights.

## When to Use

- Production PostgreSQL that needs a single-node instance rather than an Aurora cluster (cost, extension compatibility, or exact community-engine behavior)
- Workloads where a synchronous standby's zero-data-loss failover matters
- Teams standardizing on the community postgres engine across environments

## Key Configuration Choices

- **Managed master password** (`manageMasterUserPassword: true`) -- AWS generates, stores, and rotates the credential; no secret in the manifest or IaC state. The secret's ARN is exported as `master_user_secret_arn`.
- **Multi-AZ** (`multiAz: true`) -- a synchronous standby in a second AZ; failover is automatic and loses no committed transactions.
- **Storage autoscaling** (`maxAllocatedStorageGb: 200`) -- RDS grows the gp3 volume automatically as it fills; the cheap insurance against disk-full outages.
- **Encrypted storage** (`storageEncrypted: true`) -- a create-time one-way door; an unencrypted instance cannot be encrypted later.
- **Deletion safety** -- `deletionProtection` plus a named final snapshot: deleting is a deliberate two-step.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 5432 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **02-mysql-production** -- The same production shape for MySQL
- **03-read-replica** -- Add read capacity to this instance without touching it
