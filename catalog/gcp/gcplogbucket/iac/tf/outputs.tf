# The bucket's full resource name — THE composition handle: a
# GcpLoggingSink routes into this bucket with destination raw_uri
# "logging.googleapis.com/{bucket_name}", and a bucket-scoped GcpLogMetric
# references it directly.
output "bucket_name" {
  description = "Full resource name of the bucket (projects/{p}/locations/{l}/buckets/{b} or the scope's form)"
  value       = local.bucket_name
}

# The BigQuery dataset ID of the linked dataset, when armed. The provider
# reports it in the linked dataset's bigquery_dataset block.
output "linked_dataset_id" {
  description = "BigQuery dataset ID of the linked dataset (empty when not configured)"
  value = (
    var.spec.linked_bigquery_dataset != null
    ? google_logging_linked_dataset.this[0].bigquery_dataset[0].dataset_id
    : ""
  )
}
