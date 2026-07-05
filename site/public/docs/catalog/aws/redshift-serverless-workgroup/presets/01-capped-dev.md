---
title: "Capped Development Workgroup"
description: "This preset creates a cost-bounded serverless workgroup: the smallest practical RPU baseline (8) with a hard scaling ceiling (32 RPU), attached to a namespace from the resource graph. Billing follows..."
type: "preset"
rank: "01"
presetSlug: "01-capped-dev"
componentSlug: "redshift-serverless-workgroup"
componentTitle: "Redshift Serverless Workgroup"
provider: "aws"
icon: "package"
order: 1
---

# Capped Development Workgroup

This preset creates a cost-bounded serverless workgroup: the smallest practical RPU baseline (8) with a hard scaling ceiling (32 RPU), attached to a namespace from the resource graph. Billing follows the compute -- RPU-hours accrue only while queries execute, so an idle dev workgroup costs nothing -- and the cap bounds the worst case when someone runs a monster query.

## When to Use

- Development and testing against a live serverless warehouse without production spend risk
- Ad-hoc analytics teams that want a self-serve endpoint with a built-in budget guardrail
- A second, cheaper workgroup over the SAME namespace as production -- the two compute planes share one copy of the data

## Key Configuration Choices

- **Namespace by reference** (`namespaceName` → `AwsRedshiftServerlessNamespace`) -- The data plane comes from the resource graph; the workgroup can be destroyed and recreated without touching data
- **8-RPU baseline** (`baseCapacity: 8`) -- The floor each query starts from; small baselines suit intermittent dev workloads
- **32-RPU ceiling** (`maxCapacity: 32`) -- The spend guardrail; AWS will not scale past it no matter the query mix
- **Three subnets, three AZs** (`subnetIds`) -- The Redshift Serverless minimum; use private subnets (the workgroup stays private by default)
- **Default port and SG** -- Port 5439 and the VPC's default security group; put ingress rules on a referenced `AwsSecurityGroup` when you need them

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<private-subnet-id-az1>` | Private subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az2>` | Private subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az3>` | Private subnet in the third Availability Zone | AWS VPC console or `AwsSubnet` status outputs |

## Related Presets

- **02-price-performance-production** -- Use for production, where AWS picks the capacity baseline against a price-performance dial and query limits guard runaway work
