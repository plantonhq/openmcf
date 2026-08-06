# AwsInternetGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsInternetGatewaySpec defines an internet gateway for an AWS Virtual Private
Cloud (VPC).

An internet gateway is the VPC's door to the public internet: a horizontally
scaled, redundant, AWS-managed component that allows bidirectional IPv4 (and,
for dual-stack VPCs, inbound/outbound IPv6) traffic between the VPC and the
internet. It performs network address translation for instances that have a
public IPv4 address.

Attaching a gateway to a VPC does not, by itself, expose anything: a subnet
only becomes "public" once its route table sends a default route
(0.0.0.0/0, or ::/0 for IPv6) to this gateway. The standard public-subnet
recipe is therefore an AwsInternetGateway attached here plus an AwsSubnet
whose route targets it (target_type = internet_gateway).

A VPC can have at most one internet gateway attached at a time. For
IPv6-only outbound access without inbound exposure, use an egress-only
internet gateway instead.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsInternetGateway
metadata:
  name: awsinternetgateway-demo
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

AWS region the internet gateway is created in. Must match the region of the
VPC it attaches to. Example: "us-west-2", "eu-west-1". This drives provider
construction, so it is required even though the gateway logically inherits
the VPC's region.

- rule: {"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

The VPC this internet gateway attaches to. Supply a literal vpc-id or
reference an AwsVpc and the platform resolves its vpc_id output. Unlike a
subnet's vpc_id, this is updatable: changing it detaches the gateway from
the old VPC and attaches it to the new one (the gateway itself is not
replaced). A VPC may have only one internet gateway attached at a time.

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsInternetGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.internet_gateway_id` | `string` | The internet gateway's id (e.g. "igw-0abc123"). This is the value a subnet route uses as its target_id when target_type is internet_gateway. |
| `status.outputs.internet_gateway_arn` | `string` | The internet gateway's ARN. |
| `status.outputs.vpc_id` | `string` | The id of the VPC this gateway is attached to. |
| `status.outputs.region` | `string` | The AWS region the internet gateway was created in. Echoed so downstream tooling and verifiers can target the correct region. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
