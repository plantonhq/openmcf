# Aurora MySQL (Provisioned)

This preset creates a production-shaped Aurora MySQL cluster: one writer and one reader on shared cluster storage, an AWS-managed master password, encrypted storage, deletion protection, seven days of continuous backup, a 24-hour backtrack window, and Performance Insights.

## When to Use

- Production MySQL workloads with steady, predictable capacity needs
- Applications that want Aurora's in-place rewind (backtrack) as an "undo" for bad writes
- Read-heavy MySQL workloads that benefit from a reader endpoint

## Key Configuration Choices

- **Managed master password** (`manageMasterUserPassword: true`) -- AWS generates, stores, and rotates the credential in Secrets Manager; the secret's ARN is exported as `master_user_secret_arn`.
- **Backtrack window** (`backtrackWindowSeconds: 86400`) -- rewinds the cluster in place, no restore and no new endpoint; Aurora MySQL only, and only enableable at create time -- which is why the preset carries it from day one.
- **Writer + reader instances** (`instances`) -- each entry is its own managed resource; the reader is the failover target (`promotionTier: 1`).
- **Encrypted storage + deletion safety** -- the same create-time encryption and two-step deletion posture as every production preset.
- **MySQL log export** -- `error` and `slowquery` logs stream to CloudWatch Logs.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing the database port from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-aurora-postgresql** -- The same shape for PostgreSQL-compatible workloads
- **03-aurora-serverless-v2** -- Use instead when traffic is variable or spiky and capacity should track demand
