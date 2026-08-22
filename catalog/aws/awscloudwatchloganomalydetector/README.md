# AwsCloudwatchLogAnomalyDetector

One CloudWatch Logs anomaly detector: a machine-learning model that trains over a list of log groups and surfaces unusual patterns — new error classes, volume spikes, format drift — on a chosen evaluation cadence.

## Highlights

- **Multi-parent by design**: the detector references log groups by ARN (chart-wired to AwsCloudwatchLogGroup); AWS's API models a list, though it currently accepts exactly ONE entry — taught on the spec field.
- **Presence-typed visibility window** so the 7- and 90-day boundary values are expressible; unset keeps AWS's 21-day default.
- **Contracts taught in place**: the KMS key replaces the detector (a trained model cannot re-encrypt); `enabled: false` pauses analysis without losing the model; the provider treats AccessDenied reads as "detector gone" — a permissions regression can masquerade as deletion.

## Both Engines

Both modules render the single resource identically and export the same output: `anomaly_detector_arn` (identity and import ID).

## Chart Wiring

`log_group_arns` → AwsCloudwatchLogGroup `log_group_arn` outputs; `kms_key_id` → AwsKmsKey `key_arn`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
