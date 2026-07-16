# Batch ETL — basic

A single-task Cloud Run job for nightly extract-transform-load work. One execution runs one task that must succeed for the run to succeed (`taskCount: 1`). Retries are disabled (`maxRetries: 0`) so a bad input fails fast instead of burning quota.

Trigger with Cloud Scheduler (`gcloud run jobs execute`) or Eventarc on a schedule. Pair with [GcpCloudSql](/docs/catalog/gcp/gcpcloudsql) via a Cloud SQL volume when the worker reads from a managed database.
