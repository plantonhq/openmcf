# BigQuery Delivery

The zero-ETL analytics sink: Pub/Sub streams every message straight into
a BigQuery table — no Dataflow job, no custom consumer.

## What this preset creates

A BigQuery-delivery subscription composed entirely by reference: the
topic from a `GcpPubSubTopic`, the destination from a `GcpBigQueryTable`
(via its dotted `qualified_name` output). With `useTopicSchema`, the
topic's attached schema drives the column mapping; message metadata
lands in dedicated columns for lineage and debugging.

## Prerequisites

- A `GcpPubSubTopic` named `click-events` (ideally schema-validated) and
  a `GcpBigQueryTable` named `click-events-raw` whose columns match the
  schema.
- The Pub/Sub service agent needs `roles/bigquery.dataEditor` on the
  destination dataset (or set a dedicated writer via
  `serviceAccountEmail`).

## Remix ideas

- Swap `useTopicSchema` for `useTableSchema` to let the TABLE define the
  contract instead.
- Add a `deadLetterPolicy` so rows that repeatedly fail to write divert
  to a DLQ topic rather than backing up.
