---
title: "Read Replica"
description: "This preset creates a read replica of an existing RDS instance. The replica inherits engine, version, storage, and credentials from its source -- the manifest carries only what is genuinely the..."
type: "preset"
rank: "03"
presetSlug: "03-read-replica"
componentSlug: "rds-instance"
componentTitle: "RDS Instance"
provider: "aws"
icon: "package"
order: 3
---

# Read Replica

This preset creates a read replica of an existing RDS instance. The replica inherits engine, version, storage, and credentials from its source -- the manifest carries only what is genuinely the replica's own: placement, size, and networking. Promote it to a standalone instance at any time by removing `replicateSourceDb`.

## When to Use

- Offloading report, analytics, or search-indexing traffic from the primary
- Scaling read throughput horizontally without touching the write path
- Staging a promotion target ahead of a migration or a region evacuation

## Key Configuration Choices

- **Source by identifier** (`replicateSourceDb`) -- an instance identifier for a same-region replica, or an ARN for cross-region. The source needs `backupRetentionPeriod` above 0 (MySQL-family requirement).
- **No credentials, engine, or storage here** -- all inherited from the source; the spec's validation rejects them on a replica so drift is impossible.
- **No final snapshot** (`skipFinalSnapshot: true`) -- a replica holds no unique data; backups live with the source.
- **Independent sizing** -- the replica's `instanceClass` can differ from the source (bigger for heavy analytics, smaller for light read traffic).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-postgresql-instance` | The source instance identifier (or ARN for cross-region) | The source `AwsRdsInstance`'s `instance_identifier` output |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing the database port from readers | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-postgresql-production** -- The primary this replica typically hangs off
- **02-mysql-production** -- The MySQL-flavored primary
