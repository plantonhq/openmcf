# MySQL Migration via S3 Import

This preset creates a new Multi-AZ RDS MySQL instance by restoring a Percona XtraBackup stored in S3 -- no logical dump, no replication bootstrap. The instance lands with an AWS-managed master password, storage autoscaling headroom, and a week of point-in-time recovery.

## When to Use

- Migrating an on-premises or EC2-hosted MySQL database onto RDS with minimal downtime preparation
- Large databases where a logical dump/restore would take too long -- XtraBackup restores at file-copy speed
- Teams that want the migration itself to be a declarative, reviewable manifest instead of a runbook

## Key Configuration Choices

- **S3 import** (`s3Import`) -- points at the XtraBackup files and the IAM role RDS assumes to read them; `sourceEngine`/`sourceEngineVersion` describe the backup's origin (only `mysql` sources are supported). Create-time only, and mutually exclusive with read-replica, snapshot, and point-in-time create sources -- and with the create-time-only character-set choice.
- **Credentials are required** -- an S3 import creates a brand-new instance, so `username` and the managed password apply (unlike snapshot restores, which inherit credentials).
- **Storage autoscaling** (`maxAllocatedStorageGb: 500`) -- headroom above the restored size, the cheap insurance against disk-full during post-migration growth.
- **Multi-AZ + named final snapshot** -- the availability and safety posture for freshly migrated data.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 3306 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |
| `replace-with-backup-bucket` | S3 bucket holding the Percona XtraBackup files | The bucket used by your XtraBackup upload |
| `replace-with-s3-ingestion-role` | IAM role RDS assumes to read the backup (S3 read access + RDS trust) | `AwsIamRole` status outputs or the AWS IAM console |

## Related Presets

- **02-mysql-production** -- The steady-state production instance this migration typically settles into
- **03-read-replica** -- Scale reads once the migrated instance carries traffic
