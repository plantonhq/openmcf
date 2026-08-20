# One S3 vector bucket with its vector indexes and resource policy -
# three provider resources under one declarative owner.
#
# Lifecycle facts the render below depends on:
#   - EVERY index property is create-only (RequiresReplace) - an index
#     is replaced, not edited, so the spec teaches sizing dimension to
#     the embedding model before the first vector lands;
#   - bucket encryption is likewise fixed for life
#     (RequiresReplaceIfConfigured);
#   - indexes key by name (stable across list reorders);
#   - the index data_type argument is module-pinned to float32 - the
#     provider's enum holds exactly that one value;
#   - the policy is JSON-normalized by AWS (importIgnore at the
#     provider); force_destroy is config-only and never round-trips.

locals {
  bucket_encryption = var.spec.encryption != null ? [{
    sse_type    = var.spec.encryption.sse_type
    kms_key_arn = var.spec.encryption.kms_key_arn != "" ? var.spec.encryption.kms_key_arn : null
  }] : null
}

resource "aws_s3vectors_vector_bucket" "this" {
  vector_bucket_name = var.metadata.name

  encryption_configuration = local.bucket_encryption
  force_destroy            = var.spec.force_destroy

  tags = local.aws_tags
}

resource "aws_s3vectors_vector_bucket_policy" "this" {
  count = var.spec.policy != "" ? 1 : 0

  vector_bucket_arn = aws_s3vectors_vector_bucket.this.vector_bucket_arn
  policy            = var.spec.policy
}

resource "aws_s3vectors_index" "this" {
  for_each = { for index in var.spec.indexes : index.name => index }

  index_name         = each.value.name
  vector_bucket_name = aws_s3vectors_vector_bucket.this.vector_bucket_name
  # The provider's enum holds exactly this one value.
  data_type       = "float32"
  dimension       = each.value.dimension
  distance_metric = each.value.distance_metric

  encryption_configuration = each.value.encryption != null ? [{
    sse_type    = each.value.encryption.sse_type
    kms_key_arn = each.value.encryption.kms_key_arn != "" ? each.value.encryption.kms_key_arn : null
  }] : null

  dynamic "metadata_configuration" {
    for_each = length(each.value.non_filterable_metadata_keys) > 0 ? [1] : []
    content {
      non_filterable_metadata_keys = each.value.non_filterable_metadata_keys
    }
  }

  tags = local.aws_tags
}
