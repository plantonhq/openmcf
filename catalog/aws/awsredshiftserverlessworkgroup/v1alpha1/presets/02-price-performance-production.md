# Price-Performance Production Workgroup

This preset creates a production workgroup where AWS owns the capacity baseline: the price-performance target sits at the balanced level (50), a 512-RPU ceiling bounds worst-case spend, all COPY/UNLOAD data movement is forced through the VPC, connections require TLS, and a four-hour query limit guards against runaway work. Ingress rules live on the referenced `AwsSecurityGroup` node.

## When to Use

- Production analytics where the query mix varies too much for a hand-picked RPU baseline
- Workloads with data-governance requirements (VPC-visible data movement, TLS-only connections)
- Teams that want AWS to balance cost against speed instead of tuning capacity by hand

## Key Configuration Choices

- **Price-performance dial** (`pricePerformanceTarget.level: 50`) -- AWS picks and adjusts the baseline; 1 leans cheapest, 100 leans fastest. `baseCapacity` must stay unset while the dial is enabled
- **512-RPU ceiling** (`maxCapacity: 512`) -- The spend guardrail still applies even when AWS owns the baseline
- **Enhanced VPC routing** (`enhancedVpcRouting: true`) -- COPY/UNLOAD traffic moves through the VPC where flow logs and endpoints can see and govern it
- **TLS required** (`require_ssl: "true"`) -- Plaintext connections are refused at the endpoint
- **Query time limit** (`max_query_execution_time: "14400"`) -- Queries exceeding four hours are cancelled; a query-monitoring guardrail, not a tuning knob
- **Security group by reference** (`securityGroupIds` → `AwsSecurityGroup`) -- Warehouse ingress (5439 from BI tooling) belongs on that node, never inside the workgroup

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<private-subnet-id-az1>` | Private subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az2>` | Private subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az3>` | Private subnet in the third Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `my-warehouse-sg` | Name of the `AwsSecurityGroup` carrying warehouse ingress rules | Your resource graph |

## Related Presets

- **01-capped-dev** -- Use for development, with a small fixed baseline and a tight spend cap
