# Table-scoped settings AWS models as standalone resources but that are
# honestly part of the table's own configuration -- each is keyed by
# the table (replace-on-change), owned by exactly one table, and
# referenced by nothing else.

# A resource-based IAM policy on the table -- cross-account access
# grants without role assumption.
resource "aws_dynamodb_resource_policy" "this" {
  count = var.spec.resource_policy != "" ? 1 : 0

  resource_arn = aws_dynamodb_table.this.arn
  policy       = var.spec.resource_policy
}

# Item-level change data into a Kinesis Data Stream (independent of
# DynamoDB Streams). AWS allows exactly one destination per table.
resource "aws_dynamodb_kinesis_streaming_destination" "this" {
  count = var.spec.kinesis_streaming_destination != null ? 1 : 0

  table_name = aws_dynamodb_table.this.name
  stream_arn = var.spec.kinesis_streaming_destination.stream_arn

  approximate_creation_date_time_precision = var.spec.kinesis_streaming_destination.approximate_creation_date_time_precision != "" ? var.spec.kinesis_streaming_destination.approximate_creation_date_time_precision : null
}

# CloudWatch contributor insights: one provider resource for the table,
# plus one per opted-in GSI -- materialized per-name so an index list
# edit updates in place.
resource "aws_dynamodb_contributor_insights" "table" {
  count = var.spec.contributor_insights != null ? (var.spec.contributor_insights.enabled ? 1 : 0) : 0

  table_name = aws_dynamodb_table.this.name
  mode       = var.spec.contributor_insights.mode != "" ? var.spec.contributor_insights.mode : null
}

resource "aws_dynamodb_contributor_insights" "index" {
  for_each = toset(local.insights_gsi_names)

  table_name = aws_dynamodb_table.this.name
  index_name = each.value
  mode       = var.spec.contributor_insights.mode != "" ? var.spec.contributor_insights.mode : null
}
