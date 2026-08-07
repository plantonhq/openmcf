# AwsAppRunnerVpcConnector — Terraform Module

## Overview

This Terraform module provisions an App Runner VPC connector -- the managed ENI set that App Runner services route their outbound traffic through to reach private VPC resources (databases, caches, internal APIs).

## Module Structure

```
main.tf       — aws_apprunner_vpc_connector
locals.tf     — identity tags
outputs.tf    — vpc_connector_arn, vpc_connector_revision, status
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "vpc_connector" {
  source = "./path/to/module"

  metadata = {
    name = "private-backend-access"
    org  = "my-org"
    env  = "prod"
    id   = "private-backend-access-prod"
  }

  spec = {
    region             = "us-east-1"
    subnet_ids         = ["subnet-0abc123", "subnet-0def456"]
    security_group_ids = ["sg-0abc123"]
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsAppRunnerVpcConnector` and guarded against drift -- the
subnet/security-group references arrive pre-resolved as plain string lists.

## Outputs

| Output | Description |
|--------|-------------|
| `vpc_connector_arn` | The ARN services set as their egress `vpc_connector_arn` |
| `vpc_connector_revision` | The revision of this connector under its name |
| `status` | Lifecycle status (`ACTIVE` when ready) |

## Implementation Notes

- **Immutable attachment**: AWS has no update API for connectors -- changing subnets or security groups replaces the connector (a new revision under the same name).
- **Security groups are required**: AWS mandates at least one group; it governs what the connected services can reach, and the targets' groups must admit ingress from it.
- **Egress only**: inbound access to private services is a separate resource (VPC Ingress Connection), composed against the service's ARN.
