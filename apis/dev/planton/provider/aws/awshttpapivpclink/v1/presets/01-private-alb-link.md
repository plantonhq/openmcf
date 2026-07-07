# Preset: Private ALB Link

## When to Use

Use this preset to front an internal Application Load Balancer (or NLB / Cloud Map service) with HTTP APIs. One link per VPC is typically enough -- every API that needs private backends in that VPC shares it.

## Key Configuration Choices

- **Two availability zones** -- the link can only reach targets in AZs it has an ENI in, so it mirrors the load balancer's AZ spread.
- **Explicit security group** -- gives the target load balancer's security group a stable source to admit; allow egress from this group to the listener ports.
- **Composed networking** -- `valueFrom` references to Planton-managed subnets and security group keep the graph connected.

## What to Customize

1. **`<vpc-link-name>`** — Link name (e.g., `private-services-link`); one per VPC is the common shape
2. **`<private-subnet-a>` / `<private-subnet-b>`** — Names of the AwsSubnet resources in the target VPC
3. **`<vpc-link-security-group>`** — Name of the AwsSecurityGroup resource allowing egress to the backend listener ports

## Composing

Reference the link from an HTTP API's private integration:

```yaml
integration:
  integrationType: HTTP_PROXY
  integrationUri:
    value: <internal-alb-listener-arn>
  connectionType: VPC_LINK
  connectionId:
    valueFrom:
      kind: AwsHttpApiVpcLink
      name: <vpc-link-name>
      fieldPath: status.outputs.vpc_link_id
```
