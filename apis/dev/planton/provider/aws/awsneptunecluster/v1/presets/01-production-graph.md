# Production Neptune (Provisioned)

This preset creates a production-shaped Neptune cluster: one writer and one reader instance on shared cluster storage, IAM database authentication, encrypted storage, deletion protection, seven days of continuous backup, and audit logs streaming to CloudWatch. The reader doubles as the failover target -- Neptune promotes it in seconds because it already reads from the same storage volume.

## When to Use

- Production graph workloads (Gremlin, openCypher, SPARQL) with steady, predictable capacity needs
- Applications that need a reader endpoint for analytics traversals without touching the writer
- Anywhere failover time matters: a replica promotes in seconds

## Key Configuration Choices

- **IAM database authentication** (`iamDatabaseAuthenticationEnabled: true`) -- Neptune has no master username or password; SigV4-signed requests from IAM identities are its only credential mechanism, and network-only access is too loose for production.
- **Writer + reader instances** (`instances`) -- each entry is its own managed resource; add readers by appending entries, never touching the cluster. `promotionTier: 1` on the reader makes it the failover target.
- **Encrypted storage** (`storageEncrypted: true`) -- a create-time one-way door; an unencrypted cluster cannot be encrypted later.
- **Deletion safety** -- `deletionProtection` plus a named final snapshot: a delete is a deliberate two-step, never an accident.
- **Audit logging, both halves** -- the `audit` log export only emits once the `neptune_enable_audit_log` cluster parameter is enabled, so the preset sets both.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 8182 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **02-serverless** -- Use instead when traversal load is variable or spiky and capacity should track demand
