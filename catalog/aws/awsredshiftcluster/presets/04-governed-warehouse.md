# Governed Warehouse with Cost Controls and Cross-VPC Access

This preset creates a production RA3 cluster whose spend and access are governed declaratively: usage limits cap what Spectrum and concurrency scaling may consume each day, scheduled actions pause the warehouse overnight and resume it before business hours (compute stops billing while paused; storage persists), a managed VPC endpoint exposes the cluster to a BI-tooling VPC without peering, and a second AWS account is authorized to create its own endpoints to the cluster.

## When to Use

- Warehouses with predictable working hours where overnight compute is pure waste
- Teams that need hard caps on Spectrum scans and concurrency-scaling burst spend
- Hub-and-spoke networks where BI tools live in their own VPC (or their own account) and peering is off the table

## Key Configuration Choices

- **Daily Spectrum cap** (`usageLimits: spectrum / data-scanned / 10 TB / emit-metric`) -- Publishes a CloudWatch metric when external S3 scans exceed 10 TB in a day, without blocking queries
- **Concurrency-scaling cutoff** (`usageLimits: concurrency-scaling / time / 60 min / disable`) -- Hard-disables burst clusters after an hour per day; the feature re-enables when the period resets
- **Nightly pause / morning resume** (`scheduledActions`) -- Cron schedules in UTC (05:00/13:00 UTC = 22:00/06:00 PT); the IAM role's trust policy MUST allow `scheduler.redshift.amazonaws.com` to assume it, or AWS rejects the action at create. Action names are unique per AWS ACCOUNT -- keep the cluster name in them
- **Cross-VPC endpoint** (`endpointAccesses`) -- A Redshift-managed VPC endpoint in the consuming VPC's subnet group (RA3 only); the endpoint's private address surfaces in the `endpoint_access_addresses` output keyed by endpoint name
- **Cross-account grant** (`endpointAuthorizations`) -- The grantor side: account 210987654321 may create endpoints to this cluster from its own VPCs; the grantee creates its endpoint in its own account

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<private-subnet-id-az1>` | Private subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az2>` | Private subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<redshift-scheduler-role-arn>` | IAM role whose trust policy allows scheduler.redshift.amazonaws.com | AWS IAM console or `AwsIamRole` status outputs |
| `<consumer-vpc-subnet-group-name>` | Redshift subnet group in the CONSUMING VPC for the managed endpoint | AWS Redshift console (subnet groups) |

Also replace the example grantee account `210987654321` with the real AWS account ID you are authorizing.

## Related Presets

- **02-multi-node-production** -- The compliance-focused production baseline (encryption CMK, audit logging, DR snapshot copy) without governance
- **03-analytics-workload** -- Large-scale analytics with Multi-AZ and Spectrum integration
