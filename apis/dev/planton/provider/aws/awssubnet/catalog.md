# AWS Subnet

Deploys a subnet inside an AWS VPC -- a contiguous range of IP addresses pinned to a single availability zone. A subnet is not inherently "public" or "private": that identity comes entirely from its route table. This component folds routing into the subnet, so you declare a subnet's intent (public, private-with-egress, or isolated) in one place, and references the VPC and any gateways as first-class building blocks.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Subnet** -- the IP range in the chosen VPC and availability zone, with its IPv4 CIDR and optional IPv6 `/64`
- **Launch behaviour** -- map-public-IP-on-launch, IPv6 auto-assignment, DNS64, and resource-name DNS A/AAAA records, applied to instances started in the subnet
- **Dedicated route table** -- created only when inline `routes` are supplied; the rules are written into a table owned and associated by this subnet
- **Route table association** -- created when `routeTableId` references an existing table; otherwise the subnet uses the VPC's main route table
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the subnet

To build a working network, compose the companion components listed under [Works With](#works-with).

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A VPC** -- the subnet must reference a VPC. Deploy an [AWS VPC](/cloud-catalog/aws-vpc) first, or reference an existing one by id.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **Region and zone** -- the subnet is created in the specified `region` and `availabilityZone`, and both are permanent. Spanning zones means one subnet per zone.
- **CIDR planning** -- the IPv4 CIDR must fall within the VPC's range and not overlap any sibling subnet. AWS reserves 5 addresses per subnet, so a `/24` yields 251 usable.
- **Gateways for egress** -- a public subnet needs an internet gateway attached to the VPC; a private subnet's outbound access needs a NAT gateway (IPv4) or an egress-only internet gateway (IPv6).

## Deploy

### Console

Open the deployment store, find **AWS Subnet**, and click **Deploy**. The creation wizard walks you through placement (region, VPC, zone), addressing (IPv4 and optional IPv6), DNS and launch options, and routing. Start from a preset in the [Presets](#presets) tab -- **Public**, **Private**, or **Isolated** -- to pre-populate a typical configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSubnet
metadata:
  name: public-usw2a
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: production-vpc
      fieldPath: status.outputs.vpc_id
  availabilityZone: us-west-2a
  cidrBlock: "10.0.0.0/24"
  mapPublicIpOnLaunch: true
  routes:
    - destinationCidrBlock: "0.0.0.0/0"
      targetType: internet_gateway
      targetId:
        valueFrom:
          kind: AwsInternetGateway
          name: production-igw
          fieldPath: status.outputs.internet_gateway_id
```

```shell
planton apply -f subnet.yaml
```

This creates a public subnet whose dedicated route table sends all IPv4 traffic to an internet gateway, both wired by reference to other Planton resources. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a subnet. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Placement** -- `region`, `vpcId`, and `availabilityZone` define the subnet and are immutable. Reference the VPC (`valueFrom`) rather than pasting a literal id so the dependency is an explicit edge in your InfraChart graph.

**IPv4 CIDR** -- the `cidrBlock` is immutable; size it for the tier it serves (a `/24` per application tier is a good default) and keep it within the VPC range without overlapping siblings.

**IPv6** -- set `ipv6CidrBlock` to a `/64` from the VPC's IPv6 range to make the subnet dual-stack, then optionally enable IPv6 auto-assignment and DNS64.

**Routing decides identity** -- choose one: inherit the VPC main route table (isolated/default), associate an existing `routeTableId`, or supply inline `routes`. A default route (`0.0.0.0/0`) to an `internet_gateway` makes the subnet public; one to a `nat_gateway` makes it private with egress; an `::/0` route to an `egress_only_internet_gateway` gives private IPv6 egress.

## Outputs and Dependencies

### What This Component Consumes

Via ValueFromRef, this component references:

| Input | Source Resource | Source Output |
|-------|-----------------|---------------|
| `vpcId` | [AWS VPC](/cloud-catalog/aws-vpc) | `status.outputs.vpc_id` |
| `routes[].targetId` (internet_gateway) | [AWS Internet Gateway](/cloud-catalog/aws-internet-gateway) | `status.outputs.internet_gateway_id` |
| `routes[].targetId` (nat_gateway) | [AWS NAT Gateway](/cloud-catalog/aws-nat-gateway) | `status.outputs.nat_gateway_id` |
| `routes[].targetId` (egress_only) | [AWS Egress-Only Internet Gateway](/cloud-catalog/aws-egress-only-internet-gateway) | `status.outputs.egress_only_internet_gateway_id` |
| `routeTableId` | An existing route table | route-table id |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `subnet_id` | ID of the created subnet | EKS node groups, RDS subnet groups, load balancers, NAT gateways, ENIs |
| `subnet_arn` | ARN of the subnet | IAM policies, resource sharing |
| `availability_zone` | The zone the subnet lives in | Zone-aware placement of dependent resources |
| `cidr_block` | The subnet's IPv4 CIDR | Security group rules, network ACLs |
| `route_table_id` | ID of the associated route table (inline, existing, or VPC main) | Adding routes, inspection |
| `region` | Region the subnet was created in | Downstream region wiring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public subnet** -- a default route to an internet gateway with public-IP-on-launch; home for load balancers, bastions, and NAT gateways. Start from the **Public** preset.

**Private subnet** -- a default route to a NAT gateway for outbound-only access; home for application and worker tiers. Start from the **Private** preset.

**Isolated subnet** -- no internet route at all; home for data tiers that should never reach or be reached from the internet. Start from the **Isolated** preset.

## Works With

A subnet sits between the VPC and the gateways that give it reachability:

- [**AWS VPC**](/cloud-catalog/aws-vpc) -- the network the subnet draws its range from, referenced by `status.outputs.vpc_id`
- [**AWS Internet Gateway**](/cloud-catalog/aws-internet-gateway) -- the default-route target that makes a subnet public
- [**AWS NAT Gateway**](/cloud-catalog/aws-nat-gateway) -- the default-route target that gives a private subnet outbound IPv4 access
- [**AWS Egress-Only Internet Gateway**](/cloud-catalog/aws-egress-only-internet-gateway) -- the `::/0` target for outbound-only IPv6 from a private dual-stack subnet
