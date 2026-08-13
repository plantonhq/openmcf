# NLB VPC Link

This preset fronts an internal Network Load Balancer with a REST API
VPC link — the only backend type API Gateway v1 links accept.

## When to Use

- REST APIs that must reach a private service behind an NLB
- The first VPC link in a VPC (one link per NLB is enough)

## What You Get

- A v1 VPC link pointing at the named AwsNlb
- `vpc_link_id` for REST API integrations' `connectionId`

## Customize

- Point `targetArn` at your internal NLB
- Share the link across every REST API that needs that backend

## Composing

```yaml
# In an AwsRestApiGateway route integration:
integration:
  type: HTTP_PROXY
  connectionType: VPC_LINK
  connectionId:
    valueFrom:
      kind: AwsRestApiVpcLink
      name: my-nlb-link
      fieldPath: status.outputs.vpc_link_id
```
