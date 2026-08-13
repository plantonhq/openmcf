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
  # The proto uses a oneof for destination_config: exactly one of the eight
  # destination blocks is non-null (enforced by validation). The derived
  # string selects the provider's destination argument.

  destination_type = (
    var.spec.extended_s3 != null ? "extended_s3" :
    var.spec.opensearch != null ? "opensearch" :
    var.spec.opensearch_serverless != null ? "opensearchserverless" :
    var.spec.http_endpoint != null ? "http_endpoint" :
    var.spec.redshift != null ? "redshift" :
    var.spec.splunk != null ? "splunk" :
    var.spec.snowflake != null ? "snowflake" :
    "iceberg"
  )

  # ---------------------------------------------------------------------------
  # Stream sources (optional — Direct PUT when both absent)
  # ---------------------------------------------------------------------------

  has_kinesis_source = var.spec.kinesis_stream_source != null
  has_msk_source     = var.spec.msk_source != null

  # ---------------------------------------------------------------------------
  # Server-side encryption (Direct PUT only)
  # ---------------------------------------------------------------------------
  # When a Kinesis or MSK source is configured, SSE must NOT be set — the
  # source handles its own encryption. Proto-level CEL validates this.

  sse_enabled     = var.spec.sse_enabled
  sse_kms_key_arn = var.spec.sse_kms_key_arn != "" ? var.spec.sse_kms_key_arn : null
  sse_key_type    = local.sse_kms_key_arn != null ? "CUSTOMER_MANAGED_CMK" : "AWS_OWNED_CMK"

  # ---------------------------------------------------------------------------
  # Processing pipeline normalization
  # ---------------------------------------------------------------------------
  # The spec models the transformation pipeline as an ordered list of TYPED
  # processor blocks (lambda, metadata_extraction, decompression,
  # cloudwatch_log_processing, append_delimiter, record_deaggregation) so a
  # manifest reads as intent instead of AWS's internal parameter vocabulary.
  # The provider takes the raw {type, parameters[]} shape, so the typed arms
  # are translated here ONCE — every destination block then renders the same
  # normalized list. Exactly one destination is active (proto oneof), so
  # picking that destination's processing block yields a single pipeline.

  active_processing = (
    var.spec.extended_s3 != null ? var.spec.extended_s3.processing :
    var.spec.opensearch != null ? var.spec.opensearch.processing :
    var.spec.opensearch_serverless != null ? var.spec.opensearch_serverless.processing :
    var.spec.http_endpoint != null ? var.spec.http_endpoint.processing :
    var.spec.redshift != null ? var.spec.redshift.processing :
    var.spec.splunk != null ? var.spec.splunk.processing :
    var.spec.snowflake != null ? var.spec.snowflake.processing :
    var.spec.iceberg != null ? var.spec.iceberg.processing :
    null
  )

  processing_enabled = try(local.active_processing.enabled, false)

  # Each typed arm expands to the provider's processor type plus its
  # parameter name/value pairs. Optional numeric knobs are only sent when
  # set (> 0) so AWS applies its own defaults otherwise. Parameter names are
  # AWS's ProcessorParameterName constants.
  processors_normalized = local.processing_enabled ? [
    for p in local.active_processing.processors :
    p.lambda != null ? {
      type = "Lambda"
      parameters = concat(
        [{ name = "LambdaArn", value = p.lambda.lambda_arn }],
        p.lambda.buffer_size_in_mbs > 0 ? [{ name = "BufferSizeInMBs", value = tostring(p.lambda.buffer_size_in_mbs) }] : [],
        p.lambda.buffer_interval_in_seconds > 0 ? [{ name = "BufferIntervalInSeconds", value = tostring(p.lambda.buffer_interval_in_seconds) }] : [],
        p.lambda.number_of_retries > 0 ? [{ name = "NumberOfRetries", value = tostring(p.lambda.number_of_retries) }] : [],
        # Only sent when set: AWS defaults RoleArn to the delivery role and
        # the provider drops default-valued parameters from state, so
        # sending the delivery role itself would cause perpetual diffs.
        p.lambda.role_arn != "" ? [{ name = "RoleArn", value = p.lambda.role_arn }] : [],
      )
      } : p.metadata_extraction != null ? {
      type = "MetadataExtraction"
      # JsonParsingEngine is mandatory for MetadataExtraction; JQ-1.6 is the
      # only engine AWS supports today.
      parameters = [
        { name = "MetadataExtractionQuery", value = p.metadata_extraction.query },
        { name = "JsonParsingEngine", value = p.metadata_extraction.json_parsing_engine != "" ? p.metadata_extraction.json_parsing_engine : "JQ-1.6" },
      ]
      } : p.decompression != null ? {
      type = "Decompression"
      parameters = [
        { name = "CompressionFormat", value = p.decompression.compression_format },
      ]
      } : p.cloudwatch_log_processing != null ? {
      type = "CloudWatchLogProcessing"
      parameters = [
        { name = "DataMessageExtraction", value = p.cloudwatch_log_processing.data_message_extraction ? "true" : "false" },
      ]
      } : p.append_delimiter != null ? {
      type = "AppendDelimiterToRecord"
      parameters = [
        { name = "Delimiter", value = p.append_delimiter.delimiter },
      ]
      } : {
      type = "RecordDeAggregation"
      parameters = concat(
        [{ name = "SubRecordType", value = p.record_deaggregation.sub_record_type }],
        p.record_deaggregation.delimiter != "" ? [{ name = "Delimiter", value = p.record_deaggregation.delimiter }] : [],
      )
    }
  ] : []
}
