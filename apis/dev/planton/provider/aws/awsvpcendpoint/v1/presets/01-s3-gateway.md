# S3 Gateway Endpoint

This preset gives a VPC's private subnets a free, private path to S3 by
injecting the S3 prefix-list route into their route tables. It is the
default cost-and-security move for any VPC that touches S3: traffic
stops flowing through the NAT gateway (billed per GB) and never crosses
the internet.

## When to Use

- Any VPC whose workloads read or write S3 -- almost every VPC
- Cutting NAT gateway data-processing charges (S3 is usually the
  biggest contributor)
- Private subnets with no internet path at all (`nat_mode: none`
  topologies) that still need S3

## Key Configuration Choices

- **`endpointType: Gateway`** -- free and route-based; only S3 and
  DynamoDB support it. Interface endpoints for S3 exist but cost per
  AZ-hour and per GB -- only needed for on-premises or cross-VPC
  access.
- **Route tables by reference** -- each subnet that owns its route
  table (inline routes) exports `route_table_id`; add one entry per
  subnet. For subnets riding the VPC main table, reference the
  `AwsVpc`'s `main_route_table_id` output instead.
- **No `policy`** -- full access by default. Add a policy scoping the
  endpoint to your organization's buckets to turn it into a
  data-exfiltration control.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<endpoint-resource-name>` | Name for this endpoint resource (e.g. `platform-s3-gateway`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) -- also update the region inside `serviceName` | Your deployment region |
| `<vpc-resource-name>` | Name of the AwsVpc resource | Your VPC manifest's `metadata.name` |
| `<private-subnet-resource-name>` | Name of an AwsSubnet that owns its route table | Your subnet manifest's `metadata.name` |

## Common Additions

- A second `routeTableIds` entry per additional private subnet
- A `policy` document restricting the endpoint to specific buckets
- A twin endpoint for DynamoDB (`com.amazonaws.<region>.dynamodb`) --
  the only other gateway-capable service

## Related Presets

- **02-interface-endpoint** -- ENI-based private access for every other
  AWS service
- **03-privatelink-service** -- consuming a third-party PrivateLink
  service
