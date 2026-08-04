# AwsEgressOnlyInternetGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsEgressOnlyInternetGatewaySpec defines an egress-only internet gateway for
an AWS Virtual Private Cloud (VPC).

An egress-only internet gateway is the IPv6 counterpart of a NAT gateway: it
lets instances in a dual-stack VPC initiate OUTBOUND IPv6 traffic to the
internet while AWS statefully blocks any unsolicited INBOUND IPv6 connections.
It is horizontally scaled, redundant, AWS-managed, and free of charge (no
per-hour or per-GB fee, unlike a NAT gateway). Use it when private instances
need IPv6 internet access (package mirrors, API calls, telemetry) but must not
be reachable from the internet.

IPv6 only: because every IPv6 address AWS assigns is globally routable, there
is no IPv6 network address translation -- the egress-only gateway provides the
"outbound-but-not-inbound" guarantee that a private IPv4 subnet gets from a NAT
gateway. For IPv4 outbound-without-inbound use an AwsNatGateway; for full
bidirectional internet access use an AwsInternetGateway.

Attaching a gateway does not by itself route anything: a subnet only uses it
once its route table sends the IPv6 default route (::/0) to this gateway. The
standard recipe is therefore an AwsEgressOnlyInternetGateway here plus an
AwsSubnet whose route targets it (target_type = egress_only_internet_gateway).

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsEgressOnlyInternetGateway
metadata:
  name: awsegressonlyinternetgateway-demo
spec:
  region: us-west-2
  vpcId:
    value: "vpc-0abc1234def567890"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AwsVpc (`status.outputs.vpc_id`) |

## Field Details

### spec.region

`string` · required

AWS region the egress-only internet gateway is created in. Must match the
region of the VPC it attaches to. Example: "us-west-2", "eu-west-1". This
drives provider construction, so it is required even though the gateway
logically inherits the VPC's region.

- rule: {"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

The VPC this egress-only internet gateway attaches to. Supply a literal
vpc-id or reference an AwsVpc and the platform resolves its vpc_id output.
The VPC should be dual-stack (have an IPv6 CIDR) for the gateway to be
useful. Unlike an internet gateway's vpc_id, this is immutable: AWS attaches
the gateway to the VPC at creation and provides no detach/re-attach API, so
changing it replaces the gateway (ForceNew).

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEgressOnlyInternetGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.egressOnlyInternetGatewayId` | `string` | The egress-only internet gateway's id (e.g. "eigw-0abc123"). This is the value a subnet route uses as its target_id when target_type is egress_only_internet_gateway. AWS exposes no ARN for this resource. |
| `status.outputs.vpcId` | `string` | The id of the VPC this gateway is attached to. |
| `status.outputs.region` | `string` | The AWS region the egress-only internet gateway was created in. Echoed so downstream tooling and verifiers can target the correct region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
