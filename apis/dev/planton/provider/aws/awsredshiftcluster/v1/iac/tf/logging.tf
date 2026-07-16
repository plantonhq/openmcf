# Audit logging is a cluster setting keyed by the cluster identifier
# (AWS EnableLogging/DisableLogging), not a resource with its own
# identity -- which is why it is folded into this module rather than
# modeled as a standalone node. The provider models it as the separate
# aws_redshift_logging resource; attaching it after cluster creation is
# the supported ordering.
resource "aws_redshift_logging" "this" {
  count = var.spec.logging != null ? 1 : 0

  cluster_identifier   = aws_redshift_cluster.this.cluster_identifier
  log_destination_type = var.spec.logging.log_destination_type

  # S3 delivery needs the bucket (with a policy granting the Redshift
  # service write access); CloudWatch delivery needs the export list.
  # The spec's CEL rules enforce each destination's requirement.
  bucket_name   = var.spec.logging.s3_bucket_name != "" ? var.spec.logging.s3_bucket_name : null
  s3_key_prefix = var.spec.logging.s3_key_prefix != "" ? var.spec.logging.s3_key_prefix : null
  log_exports   = length(var.spec.logging.log_exports) > 0 ? var.spec.logging.log_exports : null
}
