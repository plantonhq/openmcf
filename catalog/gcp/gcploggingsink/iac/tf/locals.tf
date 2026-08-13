locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # The cloud-side sink name defaults to metadata.name when the spec leaves
  # sink_name empty — the same naming basis every kind uses.
  sink_name = (
    var.spec.sink_name != null && var.spec.sink_name != ""
    ? var.spec.sink_name
    : var.metadata.name
  )

  # Scope selection: exactly one of the four sink resources is created (the
  # count guards in main.tf key off these). An absent/empty scope means "a
  # project sink in the provider's default project".
  is_folder_sink  = var.spec.scope != null && var.spec.scope.folder_id != ""
  is_org_sink     = var.spec.scope != null && var.spec.scope.organization_id != ""
  is_billing_sink = var.spec.scope != null && var.spec.scope.billing_account != ""
  is_project_sink = !local.is_folder_sink && !local.is_org_sink && !local.is_billing_sink

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (null lets the google provider resolve its own project;
  # an empty string would be sent verbatim and rejected).
  project_id = (
    var.spec.scope != null && var.spec.scope.project_id != ""
    ? var.spec.scope.project_id
    : null
  )

  # The rendered destination URI — the exact string the Logging API expects,
  # assembled from whichever destination arm the spec set (the hand-assembly
  # this kind exists to remove):
  #   gcs_bucket       -> storage.googleapis.com/{bucket}
  #   bigquery_dataset -> bigquery.googleapis.com/projects/{p}/datasets/{d}
  #                       (accepts the dataset self_link URI or a bare
  #                       projects/... path — both normalize)
  #   pubsub_topic     -> pubsub.googleapis.com/projects/{p}/topics/{t}
  #   raw_uri          -> passed through verbatim
  destination = (
    var.spec.destination.gcs_bucket != ""
    ? "storage.googleapis.com/${var.spec.destination.gcs_bucket}"
    : var.spec.destination.bigquery_dataset != ""
    ? "bigquery.googleapis.com/${trimprefix(trimprefix(var.spec.destination.bigquery_dataset, "https://bigquery.googleapis.com/bigquery/v2/"), "bigquery.googleapis.com/")}"
    : var.spec.destination.pubsub_topic != ""
    ? "pubsub.googleapis.com/${var.spec.destination.pubsub_topic}"
    : var.spec.destination.raw_uri
  )

  # The bigquery_options block renders ONLY for BigQuery destinations — the
  # API rejects it elsewhere.
  render_bigquery_options = var.spec.destination.bigquery_dataset != "" && var.spec.destination.use_partitioned_tables
}
