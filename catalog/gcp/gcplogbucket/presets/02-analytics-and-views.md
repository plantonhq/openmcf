# Analytics And Views

The queryable application-log bucket: Log Analytics on, a linked
BigQuery dataset for the data team, and an errors-only view for
on-call — one bucket serving three audiences with different eyes.

## What it configures

- `enableAnalytics: true` — SQL over the bucket from the Log Analytics
  UI. ONE-WAY: it never disables (which is why the preset arms it from
  day one instead of "later").
- `linkedBigqueryDataset` — the bucket appears in BigQuery as a
  read-only dataset named `app_logs`. Every field of the link is
  create-time-only.
- An `errors-only` view — grant `roles/logging.viewAccessor` on it
  instead of bucket-wide read.

## Adjust before deploying

- **linkId** — becomes the BigQuery DATASET ID (letters, numbers,
  underscores); pick the name your queries will read forever.
- **retentionDays** — analytics does not extend retention; entries
  leave BigQuery when they age out of the bucket.

## When to choose something else

For a long-term audit archive where nobody runs SQL, start from the
**Compliance Retention** preset — analytics is a one-way door you may
not want armed there.
