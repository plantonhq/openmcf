# AwsHttpApiVpcLink — Terraform Module

## Overview

This Terraform module provisions an API Gateway v2 VPC link -- the managed ENI set that HTTP API private integrations route through to reach internal ALB/NLB/Cloud Map backends.

## Module Structure

```
main.tf       — aws_apigatewayv2_vpc_link
locals.tf     — identity tags
outputs.tf    — vpc_link_id, vpc_link_arn
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.29.0)
```

## Usage

```hcl
module "vpc_link" {
  source = "./path/to/module"

  metadata = {
    name = "private-services-link"
    org  = "my-org"
    env  = "prod"
    id   = "private-services-link-prod"
  }

  spec = {
    region             = "us-east-1"
    subnet_ids         = ["subnet-0abc123", "subnet-0def456"]
    security_group_ids = ["sg-0abc123"]
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsHttpApiVpcLink` and guarded against drift -- the
subnet/security-group references arrive pre-resolved as plain string lists.

## Outputs

| Output | Description |
|--------|-------------|
| `vpc_link_id` | The VPC link ID (what private integrations set as `connection_id`) |
| `vpc_link_arn` | ARN of the VPC link |

## Implementation Notes

- **Immutable attachment**: AWS has no update API for `subnet_ids`/`security_group_ids` -- changing either replaces the link. Only the name mutates in place.
- **Empty security groups**: a link without security groups applies no filtering on its side; reachability is governed solely by the target's security groups.
