---
title: "MySQL (Production Multi-AZ)"
description: "This preset creates a production-shaped MySQL instance: Multi-AZ with a synchronous standby, gp3 storage with autoscaling headroom, an AWS-managed master password, encrypted storage, deletion..."
type: "preset"
rank: "02"
presetSlug: "02-mysql-production"
componentSlug: "rds-instance"
componentTitle: "RDS Instance"
provider: "aws"
icon: "package"
order: 2
---

# MySQL (Production Multi-AZ)

This preset creates a production-shaped MySQL instance: Multi-AZ with a synchronous standby, gp3 storage with autoscaling headroom, an AWS-managed master password, encrypted storage, deletion protection, seven days of backups, Performance Insights, and RDS Blue/Green Deployments for near-zero-downtime engine upgrades.

## When to Use

- Production MySQL that needs a single-node instance rather than an Aurora cluster
- Teams that upgrade engines regularly and want Blue/Green's under-a-minute switchover instead of an in-place upgrade outage
- Workloads standardizing on the community mysql engine across environments

## Key Configuration Choices

- **Blue/Green updates** (`blueGreenUpdateEnabled: true`) -- engine upgrades and parameter changes are applied to a synchronized green copy, then switched over in under a minute; a bad upgrade never takes the primary down.
- **Managed master password** (`manageMasterUserPassword: true`) -- AWS generates, stores, and rotates the credential; the secret's ARN is exported as `master_user_secret_arn`.
- **Multi-AZ + storage autoscaling + encryption + deletion safety** -- the same production posture as the PostgreSQL preset.
- **MySQL log export** -- `error` and `slowquery` logs stream to CloudWatch Logs.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 3306 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-postgresql-production** -- The same production shape for PostgreSQL
- **03-read-replica** -- Add read capacity to this instance without touching it
