# AwsCloudwatchLogAnomalyDetector — Pulumi module (Go)

Manages one CloudWatch Logs anomaly detector (`cloudwatch.LogAnomalyDetector`).

Module facts worth knowing before editing:

- **`Enabled` is required by the provider and always rendered** — false pauses analysis without losing the trained model.
- **`KmsKeyId` replaces the detector** — AWS cannot re-encrypt a trained model in place; everything else updates in place.
- **`AnomalyVisibilityTime` is presence-typed** so the 7 and 90 boundary values are expressible; unset keeps AWS's default (21 days).
- **AWS currently accepts ONE log group ARN** in `LogGroupArnLists` even though the API models a list — extra entries fail at apply.
- **AccessDenied reads look like deletion** — the provider drops the detector from state on AccessDeniedException; a permissions regression can masquerade as a deleted detector.

Outputs mirror the Terraform module key-for-key: `anomaly_detector_arn` (identity and import ID).
