# Neptune Serverless

This preset creates a Neptune Serverless cluster: a single `db.serverless` instance that scales between 1 and 32 NCUs as traversal load moves, IAM database authentication, encrypted storage, and seven days of continuous backup. Capacity tracking is automatic -- there is no instance class to re-pick as the graph grows.

## When to Use

- Variable or spiky graph workloads where a fixed instance class is either oversized at night or undersized at peak
- New graph applications whose steady-state load is not yet known -- start serverless, switch to provisioned classes once the demand curve is clear
- Development and staging environments that should cost near the floor while idle

## Key Configuration Choices

- **Serverless capacity bounds** (`serverlessV2Scaling`) -- 1 NCU minimum is Neptune's floor (it does not pause to zero); 32 NCUs caps the spend ceiling. Widen `maxCapacity` (up to 128) as the workload grows -- it applies in place.
- **IAM database authentication** (`iamDatabaseAuthenticationEnabled: true`) -- Neptune's only credential mechanism; network-only access is too loose for anything shared.
- **Encrypted storage** (`storageEncrypted: true`) -- a create-time one-way door.
- **Scale out later without touching the cluster** -- append reader entries to `instances` (also `db.serverless`) for read scale and failover.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 8182 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-production-graph** -- Use instead for steady production load on provisioned instance classes
