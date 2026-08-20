# One S3 directory bucket (S3 Express One Zone).
#
# Lifecycle facts the render below depends on:
#   - AWS mandates the full bucket name "{base}--{zone_id}--x-s3"; it
#     is DERIVED here from metadata.name + spec.zone_id so the name
#     and the location can never disagree (and exported as the
#     bucket_name output);
#   - everything except force_destroy replaces the bucket - a
#     directory bucket is replaced, not edited;
#   - data_redundancy is sent only when set; the provider derives the
#     only-valid value from the location type otherwise;
#   - force_destroy is config-only at AWS (never read back), so
#     imports do not round-trip it.

locals {
  bucket_name = "${var.metadata.name}--${var.spec.zone_id}--x-s3"
}

resource "aws_s3_directory_bucket" "this" {
  bucket = local.bucket_name

  location {
    name = var.spec.zone_id
    type = var.spec.zone_type != "" ? var.spec.zone_type : "AvailabilityZone"
  }

  data_redundancy = var.spec.data_redundancy != "" ? var.spec.data_redundancy : null
  force_destroy   = var.spec.force_destroy

  tags = local.aws_tags
}
