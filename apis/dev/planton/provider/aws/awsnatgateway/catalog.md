# AWS NAT Gateway

Gives instances in a private subnet outbound network access while keeping them unreachable from inbound connections -- the standard way to let private workloads reach the internet (or other private networks) without exposing them. A NAT gateway is an AWS-managed, highly-available service that lives in a single subnet and is referenced by other subnets' route tables as the target of their default route. It is a first-class, independently composable building block rather than something bundled inside the VPC: create it as its own graph node, place it in exactly the subnet you intend, and reference its id from the subnets that should egress through it.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **NAT gateway** -- an EC2 NAT gateway created in the specified subnet, either `public` (fronted by an Elastic IP) or `private` (no Elastic IP)
- **Elastic IP association** -- for a public gateway, the referenced Elastic IP (`allocationId`) is bound as the gateway's stable outbound address; secondary allocations are attached when provided
- **AWS Tags** -- resource-identity tags (organization, environment, resource kind, resource ID) applied to the gateway

Creating a NAT gateway does not route anything on its own. To build a working egress topology, compose the companion components listed under [Works With](#works-with).

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A subnet** -- the gateway is created inside a subnet. Deploy an [AWS Subnet](/cloud-catalog/aws-subnet) first, or reference an existing one by id. A public gateway needs a *public* subnet.
- **An Elastic IP** (public gateways) -- a public gateway requires an [AWS Elastic IP](/cloud-catalog/aws-elastic-ip) for its outbound address. Deploy one first, or reference an existing `eipalloc-` id.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **Public gateways need an internet path** -- a public NAT gateway must live in a public subnet whose route table already reaches an [AWS Internet Gateway](/cloud-catalog/aws-internet-gateway), or it has no upstream internet path itself.
- **Zonal service** -- a NAT gateway lives in one Availability Zone. For high availability, run one per AZ and route each zone's private subnets to the gateway in their own zone.
- **Region** -- the gateway is created in the specified `region`, which must match the subnet's region.

## Deploy

### Console

Open the deployment store, find **AWS NAT Gateway**, and click **Deploy**. The creation wizard walks two steps: **Placement** (public vs private connectivity, region, and the subnet the gateway lives in) and **Addressing** (the Elastic IP for a public gateway, or private addressing for a private one). Start from a preset in the [Presets](#presets) tab -- **Public NAT Gateway** or **Private NAT Gateway**.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsNatGateway
metadata:
  name: prod-egress-nat
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  connectivityType: public
  subnetId:
    valueFrom:
      kind: AwsSubnet
      name: public-subnet-a
      fieldPath: status.outputs.subnet_id
  allocationId:
    valueFrom:
      kind: AwsElasticIp
      name: nat-eip-a
      fieldPath: status.outputs.allocation_id
```

```shell
planton apply -f nat-gateway.yaml
```

This creates a public NAT gateway in a Planton-managed public subnet, fronted by a referenced Elastic IP. A Stack Job tracks the provisioning in real time.

## Key Configuration

A NAT gateway's value is in how it composes; most deployments touch only a few fields. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Connectivity type** -- `connectivityType` is required and immutable: `public` (Elastic IP, internet egress) or `private` (no Elastic IP, egress to peered/transit/VPN networks only). Changing it replaces the gateway.

**Placement** -- `subnetId` is required and immutable; the gateway is built inside this subnet. A public gateway must sit in a public subnet. `region` must match the subnet's region.

**Elastic IP (public)** -- `allocationId` is required for a public gateway and references an `AwsElasticIp` (`status.outputs.allocation_id`) rather than embedding an address, so the IP keeps its own lifecycle. Add `secondaryAllocationIds` only for very high-throughput egress that exhausts a single EIP's source ports.

**Private addressing (private)** -- a private gateway can pin a fixed `privateIp` from the subnet's range (or let AWS choose), and add secondary private IPs either as an explicit list (`secondaryPrivateIpAddresses`) or an auto-assign count (`secondaryPrivateIpAddressCount`) -- the two are mutually exclusive.

**Placement is not routing** -- creating a gateway does nothing by itself. A private subnet gains outbound access only when its route table sends a default route (`0.0.0.0/0`) to this gateway. The egress recipe is this gateway plus an `AwsSubnet` whose route targets it (`targetType: nat_gateway`).

> **Cost note:** NAT gateways bill both an hourly charge per gateway and a per-GB data-processing charge on all traffic through them. High-volume egress -- and chatty cross-AZ traffic routed through a single-AZ gateway -- is a common surprise on an AWS bill. Traffic to S3 and DynamoDB can bypass the NAT gateway via VPC gateway endpoints.

## Outputs and Dependencies

### What This Component Consumes

Via ValueFromRef, this component references:

| Input | Source Resource | Source Output |
|-------|-----------------|---------------|
| `subnetId` | [AWS Subnet](/cloud-catalog/aws-subnet) | `status.outputs.subnet_id` |
| `allocationId` | [AWS Elastic IP](/cloud-catalog/aws-elastic-ip) | `status.outputs.allocation_id` |
| `secondaryAllocationIds` | [AWS Elastic IP](/cloud-catalog/aws-elastic-ip) | `status.outputs.allocation_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `nat_gateway_id` | ID of the created gateway | An `AwsSubnet` route's `targetId` (with `targetType: nat_gateway`) to give the subnet outbound egress |
| `public_ip` | The public IPv4 of a public gateway (the Elastic IP's address) | Allowlisting the egress source on third-party services |
| `private_ip` | The private IPv4 assigned within the subnet | Internal routing and diagnostics |
| `network_interface_id` | ID of the gateway's elastic network interface | Flow logs, troubleshooting |
| `subnet_id` | ID of the subnet the gateway lives in | Confirming placement |
| `region` | Region the gateway was created in | Downstream region wiring |

A NAT gateway has no ARN -- AWS exposes none, so none is published.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private-subnet internet egress** -- the canonical pattern: a public NAT gateway in a public subnet, with private subnets routing `0.0.0.0/0` to its `nat_gateway_id`. Start from the **Public NAT Gateway** preset.

**High-availability egress** -- one NAT gateway per Availability Zone, each in that zone's public subnet, so a zonal failure does not sever egress and cross-AZ data-processing charges are avoided.

**Private inter-network egress** -- a private NAT gateway for outbound communication to other VPCs or on-premises networks (via peering/Transit Gateway/VPN) without any internet exposure. Start from the **Private NAT Gateway** preset.

## Works With

A NAT gateway sits between the private subnets that need egress and the internet path a public subnet provides:

- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- the gateway lives in one subnet (public for a public gateway), and other subnets route a default route to this gateway's `nat_gateway_id` to gain egress
- [**AWS Elastic IP**](/cloud-catalog/aws-elastic-ip) -- a public gateway's stable outbound address, referenced by `status.outputs.allocation_id`
- [**AWS Internet Gateway**](/cloud-catalog/aws-internet-gateway) -- the internet path the public subnet (and the NAT gateway in it) routes through
- [**AWS VPC**](/cloud-catalog/aws-vpc) -- the network that contains the subnet, the gateway, and the routes that tie them together
- [**AWS Egress-Only Internet Gateway**](/cloud-catalog/aws-egress-only-internet-gateway) -- the IPv6 outbound-only counterpart for dual-stack VPCs
