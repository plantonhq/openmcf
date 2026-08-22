# AzureMonitorAutoscaleSetting Pulumi Module

## Overview

Creates an Azure Monitor autoscale setting -- the rule book that automatically adds and removes instances of ONE scalable target (a VM Scale Set, an App Service plan, ...) based on metric rules and schedules. Up to 20 profiles carry capacity envelopes, up to 10 metric rules each, and optionally a fixed-date or weekly-recurrence schedule; exactly one profile is in effect at any moment (fixed-date beats recurrence beats the default profile).

## Resources Created

- `monitoring.AutoscaleSetting` -- the setting with its profiles, rules, schedules, predictive configuration, and notifications

## Outputs

- `autoscale_setting_id` -- the setting's ARM resource ID
- `autoscale_setting_name` -- the setting's resource name

## Behavior Notes

- **One setting per target resource** -- Azure rejects a second setting on the same target at apply time; the setting must live in the target's region.
- **`enabled` is always sent explicitly** (platform default true) so previews stay deterministic.
- **Omitting `predictive` IS the disabled state** -- the API exposes no Disabled mode (provider contract). Predictive applies to VM Scale Set targets only.
- **A profile carries at most one schedule** (fixed_date XOR recurrence -- spec-validated, mirroring the provider's documented contract); a profile with neither is the default profile.
- **Recurrence `hour`/`minute` are single values in the spec** -- the classic SDK flattens the provider's one-item lists to scalar ints; both engines write the same single-element ARM arrays.
- **Schedule timezones default to UTC, always sent explicitly**; both timezone fields validate against Azure's fixed 107-value vocabulary at manifest time.
- **`metric_namespace` and `look_ahead_time` are sent only when set** -- the provider rejects empty strings on both.
- **Everything except name, resource group, region, and the target updates in place.**
- **Billing**: the autoscale setting object itself is free; you pay for the instances it creates.
