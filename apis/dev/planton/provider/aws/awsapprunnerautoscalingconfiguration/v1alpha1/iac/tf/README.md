# AwsAppRunnerAutoScalingConfiguration — Terraform Module

## Overview

This Terraform module registers an App Runner auto scaling configuration version -- the reusable scaling policy (warm floor, scale-out ceiling, per-instance concurrency) that App Runner services adopt by ARN.

## Module Structure

```
main.tf       — aws_apprunner_auto_scaling_configuration_version
locals.tf     — identity tags
outputs.tf    — configuration_arn, configuration_revision, latest
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "auto_scaling_configuration" {
  source = "./path/to/module"

  metadata = {
    name = "prod-api-scaling"
    org  = "my-org"
    env  = "prod"
    id   = "prod-api-scaling-prod"
  }

  spec = {
    region          = "us-east-1"
    max_concurrency = 50
    max_size        = 15
    min_size        = 3
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsAppRunnerAutoScalingConfiguration` and guarded against
drift.

## Outputs

| Output | Description |
|--------|-------------|
| `configuration_arn` | The revision-carrying ARN services reference |
| `configuration_revision` | The revision this deployment registered |
| `latest` | Whether this revision is the newest under the name |

## Implementation Notes

- **Everything is ForceNew by design**: AWS versions this resource, so any value change destroys the state entry and registers the NEXT revision under the same name. The new revision-carrying ARN rolls referencing services through the resource graph.
- **Deletion never hard-removes a revision**: it flips to `inactive` and stays describable.
