# Terraform Module to Deploy AwsCloudwatchCompositeAlarm

This module provisions a CloudWatch composite alarm — the boolean combination
of other alarms' states used to page once for shared-cause outages, express
alerting dependencies, and gate actions behind maintenance suppression.

## Features

- **Boolean rule expression** over other alarms (`ALARM`/`OK`/
  `INSUFFICIENT_DATA` with `AND`/`OR`/`NOT`), addressed by alarm name
- **Three action lists** (alarm / ok / insufficient-data), pre-resolved from
  references to SNS topics
- **Actions suppressor** with wait/extension windows — the maintenance-window
  mechanism that silences paging without pausing evaluation
- **Presence-aware `actions_enabled`** (ForceNew on this resource — a change
  replaces the alarm)

Generated `variables.tf` reflects the proto schema for
`AwsCloudwatchCompositeAlarm`.

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest hack/manifest.yaml
planton tofu plan --manifest hack/manifest.yaml
planton tofu apply --manifest hack/manifest.yaml --auto-approve
planton tofu destroy --manifest hack/manifest.yaml --auto-approve
```

## Outputs

| Output | Description |
| --- | --- |
| `alarm_arn` | ARN of the composite alarm |
| `alarm_name` | Name of the composite alarm (the join key for parent composites) |
