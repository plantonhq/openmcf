---
title: "DocumentDB Serverless"
description: "This preset creates a DocumentDB Serverless cluster: a single `db.serverless` instance that scales between 0.5 and 16 DCUs as demand moves, an AWS-managed master password in Secrets Manager,..."
type: "preset"
rank: "02"
presetSlug: "02-serverless"
componentSlug: "documentdb"
componentTitle: "DocumentDB"
provider: "aws"
icon: "package"
order: 2
---

# DocumentDB Serverless

This preset creates a DocumentDB Serverless cluster: a single `db.serverless` instance that scales between 0.5 and 16 DCUs as demand moves, an AWS-managed master password in Secrets Manager, encrypted storage, and seven days of continuous backup. Capacity tracking is automatic -- there is no instance class to re-pick as the workload grows.

## When to Use

- Variable or spiky workloads where a fixed instance class is either oversized at night or undersized at peak
- New applications whose steady-state load is not yet known -- start serverless, switch to provisioned classes once the demand curve is clear
- Development and staging environments that should cost near the floor while idle

## Key Configuration Choices

- **Serverless capacity bounds** (`serverlessV2Scaling`) -- 0.5 DCU minimum keeps the idle floor as low as DocumentDB allows (it does not pause to zero); 16 DCUs caps the spend ceiling. Widen `maxCapacity` as the workload grows -- it applies in place.
- **A one-way door, marked** -- removing the serverless block from a live cluster replaces it; the choice to leave serverless is a migration, not an edit.
- **Managed master password** (`manageMasterUserPassword: true`) -- no secret in the manifest or IaC state; the ARN is exported as `master_user_secret_arn`.
- **Scale out later without touching the cluster** -- append reader entries to `instances` (also `db.serverless`) for read scale and failover.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 27017 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-production-managed-password** -- Use instead for steady production load on provisioned instance classes
