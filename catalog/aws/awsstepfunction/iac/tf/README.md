# AwsStepFunction — Terraform Module

## Overview

This Terraform module provisions an AWS Step Functions state machine with version publishing, optional CloudWatch logging, X-Ray tracing, and customer-managed KMS encryption.

## Module Structure

```
main.tf       — State machine resource with dynamic config blocks
locals.tf     — Identity tags, definition serialization, presence semantics
outputs.tf    — ARN, name, version ARN, revision, status, creation date
variables.tf  — Generator-owned typed contract (metadata, spec)
provider.tf   — AWS provider configuration (>= 6.0.0)
```

## Usage

```hcl
module "step_function" {
  source = "./path/to/module"

  metadata = {
    name = "my-workflow"
    org  = "my-org"
    env  = "dev"
    id   = "my-workflow-dev"
  }

  spec = {
    region   = "us-west-2"
    type     = "STANDARD"
    publish  = true
    role_arn = "arn:aws:iam::123456789012:role/StepFunctionsExecRole"
    definition = {
      StartAt = "Hello"
      States = {
        Hello = {
          Type   = "Pass"
          Result = "Hello, World!"
          End    = true
        }
      }
    }
  }
}
```

Note: `variables.tf` is generated from the proto spec by `planton tofu
generate-variables AwsStepFunction` and guarded against drift -- reference
fields (`role_arn`, `logging.log_destination`, `encryption.kms_key_id`)
arrive pre-resolved as plain strings.

## Outputs

| Output | Description |
|--------|-------------|
| `state_machine_arn` | ARN of the state machine |
| `state_machine_name` | Name of the state machine |
| `state_machine_version_arn` | ARN of the most recently published version (empty unless `publish` is true) |
| `revision_id` | Revision identifier of the current definition |
| `status` | Lifecycle status reported by AWS |
| `creation_date` | RFC3339 creation timestamp |

## Implementation Notes

- **Definition**: The `definition` arrives as a nested object and is serialized to JSON using `jsonencode()`; ASL key casing survives.
- **Versioning**: `publish = true` publishes an immutable version on create and on every configuration change; the latest version's ARN is exported.
- **Dynamic blocks**: Tracing, logging, and encryption configurations use `dynamic` blocks to conditionally include them only when specified.
- **Log destination suffix**: Auto-appends `:*` to CloudWatch Log Group ARNs (AWS requirement enforced at plan time by the provider).
- **Encryption type**: Automatically set to `CUSTOMER_MANAGED_KMS_KEY` when a KMS key is provided; an absent block leaves AWS-owned keys.
