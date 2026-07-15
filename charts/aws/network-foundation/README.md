# AWS Network Foundation

The VPC baseline that production AWS workloads are built on: public and
private subnets striped across availability zones, an internet gateway for
the public tier, managed NAT egress for the private tier, and a free S3
gateway endpoint. Nearly everything you deploy afterwards — databases,
container services, Kubernetes clusters, load balancers — asks for the subnet
and VPC ids this chart produces.

It is opinionated where AWS is unopinionated: the public/private split, the
address plan, and the egress layout are decided here so every workload that
follows composes onto a network that was designed, not accreted.

## Architecture

```
                     ┌─────────────────────────────────────────────────┐
                     │  AwsVpc (10.0.0.0/16)                           │
                     │  DNS support + hostnames on                     │
  internet ◀───────▶ │  ┌──────────────────┐   ┌──────────────────┐   │
   AwsInternetGateway│  │ public-us-east-1a │   │ public-us-east-1b │  │
        ▲            │  │ 10.0.0.0/20       │   │ 10.0.16.0/20      │  │
        │            │  │ 0.0.0.0/0 → IGW   │   │ 0.0.0.0/0 → IGW   │  │
        │            │  │ ┌───────────────┐ │   └──────────────────┘   │
        └────────────┼──┼─│ AwsNatGateway │◀┼──────────────┐           │
                     │  │ │ + AwsElasticIp│ │              │           │
                     │  │ └───────────────┘ │              │           │
                     │  └──────────────────┘               │           │
                     │  ┌──────────────────┐   ┌───────────┼───────┐   │
                     │  │ private-us-east-1a│  │ private-us-east-1b│   │
                     │  │ 10.0.128.0/20     │  │ 10.0.144.0/20     │   │
                     │  │ 0.0.0.0/0 → NAT ──┼──┼── 0.0.0.0/0 → NAT │   │
                     │  └──────────────────┘   └───────────────────┘   │
                     │        every route table ──▶ AwsVpcEndpoint     │
                     │                              (S3, gateway, free)│
                     └─────────────────────────────────────────────────┘
```

Deployment order is derived automatically from the references: VPC first,
then the internet gateway and subnets, then the Elastic IP + NAT gateway
(which needs a public subnet), then the private-subnet routes and the S3
endpoint.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| VPC | `AwsVpc` | The address space and DNS foundation everything attaches to |
| Internet gateway | `AwsInternetGateway` | Internet door for the public subnets (and NAT egress path) |
| Public subnets (per AZ) | `AwsSubnet` | Load balancers, bastions, NAT gateways — routed to the internet gateway |
| Private subnets (per AZ) | `AwsSubnet` | Databases, containers, instances — outbound-only through NAT (conditional) |
| NAT Elastic IP(s) | `AwsElasticIp` | Stable egress address(es) that survive gateway replacement (conditional) |
| NAT gateway(s) | `AwsNatGateway` | Managed outbound internet for the private tier (conditional) |
| S3 gateway endpoint | `AwsVpcEndpoint` | Free private path to S3 that bypasses NAT data charges (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for the entire network | `us-east-1` | string |
| `network_name` | Name prefix for every resource (allows multiple networks per environment) | `core` | string |
| `vpc_cidr` | The VPC's primary CIDR — immutable, must never overlap peered/VPN ranges | `10.0.0.0/16` | string |
| `availability_zones` | AZs to stripe across; subnet CIDR lists are index-matched to this | 2 AZs | list |
| `public_subnet_cidrs` | One /20 per AZ for the public tier | `10.0.0.0/20`, `10.0.16.0/20` | list |
| `private_subnet_cidrs` | One /20 per AZ for the private tier | `10.0.128.0/20`, `10.0.144.0/20` | list |
| `nat_gateway_enabled` | Managed outbound internet for private subnets (~$32/mo + per-GB per gateway) | `true` | bool |
| `nat_gateway_per_az` | One NAT per AZ (AZ-outage isolation, no cross-AZ charge) vs one shared | `false` | bool |
| `s3_endpoint_enabled` | Free S3 gateway endpoint on every route table | `true` | bool |

## The cost model, honestly

- **NAT gateway is the only meaningful cost**: each gateway bills hourly
  (~$32/month) plus per-GB data processing. The default single-gateway layout
  minimizes the hourly component; the per-AZ layout trades more hourly cost
  for AZ-failure isolation and no cross-AZ data charge. High-throughput
  egress workloads should do the math — data processing usually dominates.
- **Everything else is free or near-free**: VPC, subnets, route tables, the
  internet gateway, and the S3 gateway endpoint carry no hourly charge. An
  idle Elastic IP bills a small hourly fee only while unattached — attached
  to a NAT gateway it is free.
- **The S3 endpoint pays for itself**: S3 traffic from private subnets would
  otherwise route through NAT and bill per GB. The endpoint removes that
  entirely.

## After deploying

Everything downstream composes by reference. A database, cluster, or service
manifest points at this network like so:

```yaml
subnetIds:
  - valueFrom:
      kind: AwsSubnet
      name: core-private-us-east-1a
      fieldPath: status.outputs.subnet_id
  - valueFrom:
      kind: AwsSubnet
      name: core-private-us-east-1b
      fieldPath: status.outputs.subnet_id
```

The useful join points:

- `AwsVpc` → `status.outputs.vpc_id` (security groups, endpoints, peering)
- `AwsSubnet` → `status.outputs.subnet_id` (workload placement) and
  `status.outputs.route_table_id` (additional gateway endpoints)
- `AwsNatGateway` → `status.outputs.public_ip` (the egress address to give
  partners for allowlisting)

## Day-2 guidance

- **Adding a third AZ**: append one entry to each of the three lists and
  redeploy — the new subnet pair (and, with `nat_gateway_per_az`, its NAT)
  is created without touching the existing ones. Removing AZs from a live
  network is a migration, not an edit: drain workloads off those subnets
  first.
- **Growing the address space**: the primary CIDR is immutable, but
  `AwsVpc.secondaryIpv4CidrBlocks` can add ranges later; new subnet tiers can
  then be carved from them (managed directly or by extending this chart).
- **More gateway endpoints**: DynamoDB has the same free gateway-endpoint
  shape as S3 — add an `AwsVpcEndpoint` with
  `serviceName: com.amazonaws.<region>.dynamodb` reusing the same
  `routeTableIds` pattern. Interface endpoints (ECR, Secrets Manager, STS,
  …) compose the same way but live in the private subnets and bill hourly.
- **Isolated-network mode**: with `nat_gateway_enabled: false` private
  subnets have no internet path at all. Reach AWS APIs through endpoints
  (the S3 gateway is included; add interface endpoints per service), and
  push container images through ECR endpoints rather than public registries.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
