locals {
  # ---------------------------------------------------------------------------
  # Name & tags
  # ---------------------------------------------------------------------------

  # metadata.name is the delivery stream's cloud name (ForceNew).
  delivery_stream_name = var.metadata.name

  # Resource-identity tags follow the catalog convention.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsKinesisFirehose"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # ---------------------------------------------------------------------------
  # Destination type (derived from which oneof field is populated)
  # ---------------------------------------------------------------------------
  # The proto uses a oneof for destination_config: exactly one of the four
  # destination blocks is non-null (enforced by validation). The derived
  # string selects the provider's destination argument.

  destination_type = (
    var.spec.extended_s3 != null ? "extended_s3" :
    var.spec.opensearch != null ? "opensearch" :
    var.spec.http_endpoint != null ? "http_endpoint" :
    "redshift"
  )

  # ---------------------------------------------------------------------------
  # Kinesis stream source (optional — Direct PUT when absent)
  # ---------------------------------------------------------------------------

  has_kinesis_source = var.spec.kinesis_stream_source != null

  # ---------------------------------------------------------------------------
  # Server-side encryption (Direct PUT only)
  # ---------------------------------------------------------------------------
  # When a Kinesis source is configured, SSE must NOT be set — the source
  # stream handles its own encryption. Proto-level CEL validates this.

  sse_enabled     = var.spec.sse_enabled
  sse_kms_key_arn = var.spec.sse_kms_key_arn != "" ? var.spec.sse_kms_key_arn : null
  sse_key_type    = local.sse_kms_key_arn != null ? "CUSTOMER_MANAGED_CMK" : "AWS_OWNED_CMK"
}
