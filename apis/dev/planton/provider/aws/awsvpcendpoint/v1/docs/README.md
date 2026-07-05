# AWS VPC Endpoint: Private Service Access

## What This Component Is

A VPC endpoint is a private connection from a VPC to a service --
AWS's own (S3, DynamoDB, STS, ECR, CloudWatch, ...), a third party's
PrivateLink service, or a VPC Lattice resource. Traffic to the service
stays on the AWS network: no internet gateway, no NAT, no public IPs on
the path. `AwsVpcEndpoint` models one endpoint on one VPC -- a VPC
typically carries several (each service its own endpoint, each with its
own policy), which is exactly why the endpoint is a first-class node
rather than a list on the VPC.

Create-time immutable in AWS: the endpoint type, the service target,
and `serviceRegion`. Attachments (route tables, subnets, security
groups), the policy, DNS options, and IP address type update in place.

## The Two Types That Matter (and the Three That Are Modeled Anyway)

- **Gateway** -- S3 and DynamoDB only. Free. Works by route injection:
  AWS maintains a prefix-list route in every route table you attach, so
  traffic from subnets on those tables reaches the service privately.
  Gateway endpoints are the default cost move in any VPC that touches
  S3: S3 traffic through a NAT gateway is billed per GB; through a
  gateway endpoint it is free.
- **Interface** -- everything else, including S3/DynamoDB when
  on-premises or cross-VPC clients need the private path. An ENI per
  attached subnet (each AZ billed per hour + per GB), security-group
  controlled, and with private DNS able to override the service's
  public name inside the VPC so clients need zero configuration.
- **GatewayLoadBalancer** fronts a GWLB appliance fleet;
  **Resource** and **ServiceNetwork** attach VPC Lattice targets. All
  three are ENI-based like Interface but take neither security groups
  nor `privateDnsEnabled` -- the spec's CEL gating encodes exactly this.

## Type Gating, Enforced Before AWS Sees It

AWS rejects mismatched attachments at create time (route tables on an
interface endpoint, subnets on a gateway endpoint, private DNS on
anything but Interface). The spec enforces every one of those rules
with CEL at validation time -- misconfiguration fails in seconds with a
message naming the offending field, not minutes into a deploy.

The service target is an exactly-one-of trio (`serviceName` /
`resourceConfigurationArn` / `serviceNetworkArn`), and the Lattice ARNs
are coupled to their endpoint types.

## Route Tables: Where Gateway Endpoints Compose

A gateway endpoint attaches to route tables, and the resource graph
offers two honest sources:

- A subnet that OWNS its route table (inline `routes` on `AwsSubnet`)
  exports it as the `route_table_id` output -- reference that.
- Subnets riding the VPC main route table export an EMPTY
  `route_table_id`; reference the `AwsVpc`'s `main_route_table_id` or
  `default_route_table_id` outputs instead.

The endpoint never edits the referenced tables' own routes -- AWS
manages the prefix-list route as part of the endpoint, adding and
removing it with the endpoint's lifecycle.

## Private DNS, Including the S3 Dual-Stack Pattern

`privateDnsEnabled` associates an AWS-managed private hosted zone that
resolves the service's PUBLIC name (e.g. `sts.us-west-2.amazonaws.com`)
to the endpoint's private IPs -- clients keep their default SDK
endpoints and privately reach the service with zero code changes. It
requires the VPC to have both DNS support and DNS hostnames enabled.

`dnsOptions.privateDnsOnlyForInboundResolverEndpoint` exists for one
scenario AWS built explicitly: a service (S3) with BOTH a gateway and
an interface endpoint in the same VPC. In-VPC traffic keeps riding the
free gateway endpoint; only inbound resolver traffic (on-premises over
VPN/DX) resolves to the interface endpoint.

The `privateDnsPreference` / `privateDnsSpecifiedDomains` pair applies
to the Lattice types; AWS requires the domain list exactly when the
preference includes specified domains, and the spec couples them with
CEL.

## Endpoint Policies: The Overlooked Security Control

Gateway and most interface endpoints accept an IAM policy document
scoping which principals may use the endpoint to reach which resources.
An S3 gateway endpoint with a policy allowing only the organization's
buckets turns the endpoint into a data-exfiltration control: workloads
in the VPC physically cannot reach an attacker's bucket through it.
Empty keeps AWS's full-access default.

## Deliberately Skipped (with reasons)

- **The PrivateLink provider side** (`aws_vpc_endpoint_service` +
  allowed-principal / connection-accepter / connection-notification /
  private-DNS-verification satellites): publishing your own service
  behind an NLB/GWLB is a separate product surface from consuming
  services privately; deferred until real demand appears.
- **`aws_vpc_endpoint_policy`** as a standalone resource: a policy
  belongs to exactly one endpoint -- folded into the spec's `policy`
  field.
- **The standalone association resources** (subnet / route-table /
  security-group associations): glue with no independent lifecycle --
  folded into the spec's repeated reference fields, updated in place.

## Operational Notes

- Interface endpoints take a few minutes to provision (ENI creation per
  subnet); gateway endpoints are near-instant.
- A deleted endpoint lingers describable in a `deleted` state briefly
  before disappearing -- automation checking for absence should treat
  both signals as gone.
- `pendingAcceptance` in the `state` output means the PrivateLink
  service requires manual acceptance on the provider's side (or
  `autoAccept` for same-account services).
- For interface endpoints without explicit `securityGroupIds`, AWS
  attaches the VPC's DEFAULT security group -- fine for smoke tests,
  usually wrong for production; give endpoints a dedicated group that
  allows 443 from the VPC CIDR.
