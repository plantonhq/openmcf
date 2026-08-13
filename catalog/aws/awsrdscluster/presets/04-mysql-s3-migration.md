# MySQL to Aurora via S3 Import

This preset migrates a self-managed MySQL database into a new Aurora MySQL cluster by restoring a Percona XtraBackup stored in S3 -- no logical dump, no replication bootstrap. The new cluster lands with an AWS-managed master password, a 72-hour backtrack window for cutover safety, and a week of point-in-time recovery.

## When to Use

- Migrating an on-premises or EC2-hosted MySQL database into Aurora with minimal downtime preparation
- Large databases where a logical dump/restore would take too long -- XtraBackup restores at file-copy speed
- Teams that want the migration itself to be a declarative, reviewable manifest instead of a runbook

## Key Configuration Choices

- **S3 import** (`s3Import`) -- points at the XtraBackup files and the IAM role RDS assumes to read them; `sourceEngine`/`sourceEngineVersion` describe the backup's origin (only `mysql` sources are supported, and the target must be `aurora-mysql`). Create-time only, and mutually exclusive with snapshot or point-in-time restores.
- **Credentials are required** -- an S3 import creates a brand-new cluster, so `masterUsername` and the managed password apply (unlike snapshot restores, which inherit credentials).
- **Backtrack** (`backtrackWindowSeconds: 259200`) -- 72 hours of in-place rewind, the fastest "undo" while cutover traffic settles; enabling it later is not supported by AWS, so it lands at create.
- **Deletion safety** -- `skipFinalSnapshot: false` with a named final snapshot, the posture for anything holding migrated data.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing the database port from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |
| `replace-with-backup-bucket` | S3 bucket holding the Percona XtraBackup files | The bucket used by your XtraBackup upload |
| `replace-with-s3-ingestion-role` | IAM role RDS assumes to read the backup (S3 read access + RDS trust) | `AwsIamRole` status outputs or the AWS IAM console |

## Related Presets

- **02-aurora-mysql** -- The steady-state production cluster this migration typically settles into
- **03-aurora-serverless-v2** -- The scale-to-zero variant for variable workloads
