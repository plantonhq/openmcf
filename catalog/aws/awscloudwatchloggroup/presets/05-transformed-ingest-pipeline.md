# Transformed Ingest Pipeline

**Use case:** Normalize messy service logs at ingestion so queries, filters,
and alarms all see one clean shape.

This preset attaches an ingest-time transformer to the log group: every event
is parsed as JSON, its message key normalized, a `service` identity stamped
in, and numeric fields converted to real numbers. A metric filter then counts
server errors against the TRANSFORMED events — the filter pattern reads the
pipeline's output, not the raw line.

## What You Get

- A STANDARD class CloudWatch Log Group with 90-day retention
- A four-step transformer pipeline (parse → rename → enrich → type-convert)
  applied to every event at ingestion — Logs Insights, metric filters, and
  subscription filters all see the transformed shape
- A metric filter with `applyOnTransformedLogs: true` publishing
  `ServerErrorCount` from the parsed `statusCode` field
- Outputs: `log_group_arn`, `log_group_name`

## When to Use

- Services whose raw logs are JSON-in-a-string (the default for most
  container runtimes) and need real fields for querying
- Fleets where each service logs slightly different key names and you want
  one canonical shape downstream
- Alarming on numeric conditions (`statusCode >= 500`) that require typed
  fields rather than substring matches

## Key Configuration Choices

- **Parser first** — AWS requires the pipeline's first processor to be a
  parser (`parseJson`, `grok`, `csv`, `parseKeyValue`, or a vended-log
  parser); validation enforces it before the apply
- **STANDARD class only** — transformers are not supported on
  INFREQUENT_ACCESS or DELIVERY class groups (validation enforces this too)
- **`applyOnTransformedLogs`** — leave it unset on filters that should keep
  matching the raw line; only an explicit `false` switches an existing
  filter back once it has been true

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region for the log group | Your deployment region |
| `<service-name>` | Identity stamped into every event | Your service catalog |
| `<app-namespace>` | Custom metric namespace (e.g. `MyApp/Errors`) | Your naming convention |
