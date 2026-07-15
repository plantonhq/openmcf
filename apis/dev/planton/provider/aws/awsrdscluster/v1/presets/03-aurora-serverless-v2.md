# Aurora Serverless v2 (Scale-to-Zero)

This preset creates an Aurora PostgreSQL Serverless v2 cluster: one `db.serverless` instance that scales between 0 and 16 ACUs with demand, automatic pause after five idle minutes (compute cost drops to zero; storage still billed), an AWS-managed master password, and the Data API for connection-less SQL access.

## When to Use

- Development and staging environments where cost should track actual usage -- an idle cluster costs storage only
- Applications with variable, unpredictable, or spiky traffic
- Serverless architectures calling the database from Lambda or Step Functions through the Data API
- Workloads that want Aurora's reliability without capacity planning

## Key Configuration Choices

- **Scale-to-zero** (`minCapacity: 0`) -- the instance pauses after `secondsUntilAutoPause` (AWS default: 300s) of idleness and resumes on the next connection in roughly fifteen seconds. Latency-sensitive production should raise `minCapacity` to 0.5 or more so the cluster never pauses.
- **Serverless v2 is provisioned mode** -- the scaling block plus a `db.serverless` instance, NOT `engineMode: serverless` (that selects the legacy Serverless v1 offering).
- **Data API** (`enableHttpEndpoint: true`) -- SQL over HTTPS with IAM auth; no connection pools to manage from Lambda.
- **Managed master password** -- AWS generates, stores, and rotates the credential in Secrets Manager.
- **Capacity ceiling** (`maxCapacity: 16`) -- the hard spend/performance bound per instance; raise it for heavier peaks.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing the database port from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-aurora-postgresql** -- Use instead for steady production capacity with a writer/reader pair
- **02-aurora-mysql** -- The provisioned MySQL-compatible variant
