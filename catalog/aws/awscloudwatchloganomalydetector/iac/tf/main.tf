# One CloudWatch Logs anomaly detector training over a list of log
# groups.
#
# Lifecycle facts the render below depends on:
#   - enabled is REQUIRED by the provider and always rendered - false
#     pauses analysis without losing the trained model;
#   - kms_key_id replaces the detector on change (AWS cannot re-encrypt
#     a trained model in place); everything else updates in place;
#   - anomaly_visibility_time is presence-typed in the spec so the 7 and
#     90 boundary values are expressible; unset keeps AWS's default (21);
#   - the provider treats an AccessDeniedException on read as "detector
#     gone" and drops it from state - a permissions regression can look
#     like a deleted detector.

resource "aws_cloudwatch_log_anomaly_detector" "this" {
  log_group_arn_list = var.spec.log_group_arns
  enabled            = var.spec.enabled

  detector_name           = var.spec.detector_name != "" ? var.spec.detector_name : null
  evaluation_frequency    = var.spec.evaluation_frequency != "" ? var.spec.evaluation_frequency : null
  filter_pattern          = var.spec.filter_pattern != "" ? var.spec.filter_pattern : null
  kms_key_id              = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  anomaly_visibility_time = var.spec.anomaly_visibility_time != null ? var.spec.anomaly_visibility_time : null

  tags = local.aws_tags
}
