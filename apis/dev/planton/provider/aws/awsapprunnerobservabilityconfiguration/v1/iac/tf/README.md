# AwsAppRunnerObservabilityConfiguration — Terraform Module

## Overview

This Terraform module registers an App Runner observability configuration version -- the reusable tracing policy (AWS X-Ray via a managed OpenTelemetry collector) that App Runner services adopt by ARN.

## Module Structure

```
main.tf       — aws_apprunner_observability_configuration
locals.tf     — identity tags
outputs.tf    — configuration_arn, configuration_revision, latest
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "observability_configuration" {
  source = "./path/to/module"

  metadata = {
    name = "xray-tracing"
    org  = "my-org"
    env  = "prod"
    id   = "xray-tracing-prod"
  }

  spec = {
    region = "us-east-1"
    trace_configuration = {
      vendor = "AWSXRAY"
    }
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsAppRunnerObservabilityConfiguration` and guarded against
drift.

## Outputs

| Output | Description |
|--------|-------------|
| `configuration_arn` | The revision-carrying ARN services reference |
| `configuration_revision` | The revision this deployment registered |
| `latest` | Whether this revision is the newest under the name |

## Implementation Notes

- **Trace settings are ForceNew by design**: AWS versions this resource; changes register the next revision under the same name, and the new ARN rolls referencing services through the resource graph.
- **The trace block is optional**: a configuration without it is valid but inert -- emitted only when the spec configures tracing.
