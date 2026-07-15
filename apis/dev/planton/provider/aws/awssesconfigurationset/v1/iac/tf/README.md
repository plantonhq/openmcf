# AwsSesConfigurationSet Terraform Module

Provisions an Amazon SES (SESv2) configuration set with optional delivery
controls, suppression overrides, tracking, VDM options, and per-name event
destinations.

## Resources Created

- `aws_sesv2_configuration_set.this` — The configuration set (name from `metadata.name`)
- `aws_sesv2_configuration_set_event_destination.this` — One per named event destination

## Inputs

| Variable | Description |
|----------|-------------|
| `metadata` | Resource metadata (name, org, env, id, labels) |
| `spec` | AwsSesConfigurationSetSpec — desired configuration |

## Outputs

| Output | Description |
|--------|-------------|
| `configuration_set_arn` | ARN of the configuration set |
| `configuration_set_name` | Name of the configuration set (from `metadata.name`) |

## Provider

Requires `hashicorp/aws` provider version `>= 6.0.0`.
