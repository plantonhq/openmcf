# AwsSesConfigurationSet Pulumi Module

Provisions an Amazon SES (SESv2) configuration set with delivery controls,
suppression overrides, tracking, VDM options, and per-name event destinations.

## Resources Created

- `aws:sesv2:ConfigurationSet` — The configuration set
- `aws:sesv2:ConfigurationSetEventDestination` — One per named event destination

## Outputs

| Key | Description |
|-----|-------------|
| `configuration_set_arn` | ARN of the configuration set |
| `configuration_set_name` | Name of the configuration set |

## Local Development

```bash
./debug.sh preview
./debug.sh up
./debug.sh destroy
```
