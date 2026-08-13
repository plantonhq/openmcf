# AwsRestApiVpcLink

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsRestApiVpcLink example (hack/dev manifest and refgen
# Example source): a REST API VPC link fronting an internal Network
# Load Balancer. Literal ARN stands in for the composed NLB reference
# so the offline `tofu plan` renders the resource.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRestApiVpcLink
metadata:
  name: orders-nlb-link
  id: orders-nlb-link
  org: test-org
  env: dev
spec:
  region: us-west-2
  description: Orders service NLB
  targetArn:
    value: arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/net/orders/abcdef
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.targetArn` | `string \| valueFrom` | yes |  | AwsNlb (`status.outputs.load_balancer_arn`) |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.targetArn

`string | valueFrom` · required

- references: AwsNlb (`status.outputs.load_balancer_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsNlb, name: <that resource's name>, fieldPath: status.outputs.load_balancer_arn}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRestApiVpcLink, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_link_id` | `string` |  |
| `status.outputs.vpc_link_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.targetArn` | AwsNlb | `status.outputs.load_balancer_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsRestApiGateway | `spec.routes[].integration.vpcLinkId` | `status.outputs.vpc_link_id` |

## See Also

- [Overview](../README.md)
