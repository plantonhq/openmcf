# Database Activity Stream: AWS creates and owns a Kinesis stream that
# receives every audited database event, encrypted with the given KMS
# key. Create/delete-only lifecycle -- every argument forces
# replacement, and AWS walks the stream through starting/stopping
# states with a waiter on each side, so a slow apply or destroy here is
# the service's state machine, not this module. The stream needs an
# available instance to start against, hence the explicit dependency on
# the folded instances.
resource "aws_rds_cluster_activity_stream" "this" {
  count = var.spec.activity_stream != null ? 1 : 0

  resource_arn                        = aws_rds_cluster.this.arn
  mode                                = var.spec.activity_stream.mode
  kms_key_id                          = var.spec.activity_stream.kms_key_id
  engine_native_audit_fields_included = var.spec.activity_stream.engine_native_audit_fields_included

  depends_on = [aws_rds_cluster_instance.this]
}
