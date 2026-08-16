# AwsCloudTrailEventDataStore — Pulumi module

Manages a CloudTrail Lake event data store
(`aws:cloudtrail/eventDataStore:EventDataStore`).

Module facts worth knowing before editing:

- **Provider-default toggles render only on an explicit choice.**
  `billing_mode`, `multi_region_enabled`, `retention_period`, and
  `termination_protection_enabled` all carry provider/AWS defaults;
  unset spec values are never sent, so the module never fights them.
- **`suspend` is the provider's nullable string.** The spec's
  tri-state bool renders as `"true"`/`"false"` only when explicitly
  set; AWS never reports it back (write-only), so applies re-assert
  it.
- **Deletion honors termination protection.** The provider does NOT
  auto-disable it — a destroy with protection on fails by AWS design;
  the two-step teardown is taught in the GUIDE, never worked around
  here.
- **Selector rendering matches the trail kind's idiom** (empty
  operator lists are omitted so both engines send identical payloads).

Outputs mirror the Terraform module key-for-key:
`event_data_store_arn`.
