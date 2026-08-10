# AWS VPC

Deploys a Virtual Private Cloud on AWS -- the isolated network foundation that other AWS resources launch into. The VPC is a clean primitive: an IPv4 address space (with optional secondary ranges and IPv6), a tenancy mode, and DNS behaviour. Subnets, internet gateways, NAT gateways, and route tables are separate components that reference this VPC, so the whole topology is composed from first-class building blocks.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPC** -- the virtual network with your primary IPv4 CIDR (an explicit block or an IPAM allocation), DNS resolution and hostname settings, and the instance tenancy mode
- **Secondary IPv4 CIDRs** -- created only when `secondaryIpv4Cidrs` is set; each entry (an explicit block, an IPAM-sized allocation, or a pool-pinned block) is associated as its own resource and can be added or removed without recreating the VPC
- **Secondary IPv6 CIDRs** -- created only when `secondaryIpv6Cidrs` is set; each entry names exactly one source (an Amazon-provided block, a BYOIP public pool, or an IPAM pool)
- **Encryption Control** -- created only when `encryptionControl` is set; monitors or enforces encryption in transit VPC-wide, with per-service exclusions in enforce mode
- **IPv6 association** -- created only when IPv6 is enabled; either an Amazon-provided `/56` or an IPAM-allocated block, making the VPC dual-stack
- **Default route table, security group, and network ACL** -- AWS creates these automatically with every VPC; their IDs are surfaced as outputs
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the VPC

To build a working network on top of this VPC, compose the companion components listed under [Works With](#works-with).

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **Target region** -- the VPC is created in the specified `region` and is permanently bound to it. A VPC cannot be moved between regions.
- **CIDR planning** -- choose a primary CIDR (`/16`-`/28`) that does not overlap with VPCs you intend to peer or connect via Transit Gateway. Common choices are `10.0.0.0/16` for production and `10.1.0.0/16` for staging.
- **Service Quotas** -- the default AWS limit is 5 VPCs per region (and 5 CIDR associations per VPC). Request a quota increase if needed.

## Deploy

### Console

Open the deployment store, find **AWS VPC**, and click **Deploy**. The creation wizard walks you through region, IPv4 addressing, optional IPv6, and DNS/options. Start from the **Production Dual-Stack** preset in the [Presets](#presets) tab to pre-populate a production-ready configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsVpc
metadata:
  name: production-vpc
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  cidrBlock: "10.0.0.0/16"
  assignGeneratedIpv6CidrBlock: true
  enableDnsHostnames: true
  enableDnsSupport: true
```

```shell
planton apply -f vpc.yaml
```

This creates a dual-stack VPC with a `/16` IPv4 range and an Amazon-provided IPv6 `/56`, DNS resolution and hostnames on. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Primary IPv4 source** -- Provide either an explicit `cidrBlock` (the common path) or allocate it from an IPAM pool (`ipv4IpamPoolId` + `ipv4NetmaskLength`). The two are mutually exclusive. A `/16` block gives 65,536 addresses; the primary range cannot change after creation without replacing the VPC.

**Secondary CIDRs** -- Add `secondaryIpv4Cidrs` when a single range runs low or you need a distinct space (such as `100.64.0.0/10`). Each entry takes an explicit `cidrBlock`, an IPAM allocation (`ipamPoolId` + `netmaskLength`), or a pool-pinned block (both) -- the same three modes as the primary range. IPv6 grows the same way through `secondaryIpv6Cidrs`, whose entries choose an Amazon-provided block (`assignGenerated`), a BYOIP `ipv6Pool`, or an IPAM pool. Unlike the primary, secondaries can be added and removed in place.

**Encryption in transit** -- `encryptionControl` turns on AWS's VPC-wide encryption-in-transit control: `monitor` observes and reports unencrypted traffic paths; `enforce` blocks them, honoring per-service exclusions (internet gateway, NAT gateway, Lambda, EFS, ...) for paths that cannot encrypt. Roll out monitor-first, review the findings, then enforce with the minimum exclusion list.

**IPv6** -- Opt in for globally routable addresses. **Amazon-provided** is the simplest (one setting yields a `/56`); **IPAM** suits orgs that govern IPv6 space centrally. Pair a dual-stack VPC with an egress-only internet gateway for outbound-only IPv6.

**DNS settings** -- Leave `enableDnsSupport` on its default (resolution on) unless you run your own resolver. Enable `enableDnsHostnames` for services that depend on DNS names within the VPC (EKS, RDS, private hosted zones) -- it requires resolution to be on.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpc_id` | ID of the created VPC | Subnets, gateways, security groups, EKS/RDS, load balancers |
| `vpc_arn` | ARN of the VPC | IAM policies, resource sharing |
| `cidr_block` | The primary IPv4 CIDR | Security group rules, network ACLs, peering configuration |
| `ipv6_cidr_block` | The associated IPv6 CIDR (empty when IPv4-only) | IPv6 subnet ranges, egress-only routing |
| `main_route_table_id` | ID of the VPC's main route table | Subnet route-table associations |
| `default_route_table_id` | ID of the default route table | Default routing inspection |
| `default_security_group_id` | ID of the default security group | Baseline rule management |
| `default_network_acl_id` | ID of the default network ACL | Subnet-level network ACLs |
| `owner_id` | AWS account ID that owns the VPC | Cross-account references |
| `region` | Region the VPC was created in | Downstream region wiring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production dual-stack** -- A `/16` IPv4 range plus an Amazon-provided IPv6 `/56`, with DNS resolution and hostnames on. The standard foundation for a production network; compose subnets and gateways onto it. Start from the **Production Dual-Stack** preset.

**Development** -- A minimal single-CIDR, IPv4-only VPC for development environments. Start from the **Development** preset.

## Works With

A VPC is the root of an AWS network topology. Compose these components, each referencing this VPC by `status.outputs.vpc_id`:

- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- carves subnets (public, private, or isolated) from the VPC's address space, with their own routing
- [**AWS Internet Gateway**](/cloud-catalog/aws-internet-gateway) -- attaches to the VPC to give public subnets inbound and outbound internet access
- [**AWS NAT Gateway**](/cloud-catalog/aws-nat-gateway) -- gives private subnets outbound internet access, composing an Elastic IP by reference
- [**AWS Egress-Only Internet Gateway**](/cloud-catalog/aws-egress-only-internet-gateway) -- the IPv6 outbound-only counterpart of a NAT gateway for dual-stack VPCs
