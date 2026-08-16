# AwsCloudwatchLogAnomalyDetector — Terraform/OpenTofu module

Manages one CloudWatch Logs anomaly detector (`aws_cloudwatch_log_anomaly_detector`).

Module facts worth knowing before editing:

- **`enabled` is required by the provider and always rendered** — false pauses analysis without losing the trained model.
- **`kms_key_id` replaces the detector** — AWS cannot re-encrypt a trained model in place; everything else updates in place.
- **`anomaly_visibility_time` is presence-typed** so the 7 and 90 boundary values are expressible; unset keeps AWS's default (21 days).
- **AWS currently accepts ONE log group ARN** in `log_group_arn_list` even though the API models a list — extra entries fail at apply.
- **AccessDenied reads look like deletion** — the provider drops the detector from state on AccessDeniedException; a permissions regression can masquerade as a deleted detector.

Outputs mirror the Pulumi module key-for-key: `anomaly_detector_arn` (identity and import ID).
