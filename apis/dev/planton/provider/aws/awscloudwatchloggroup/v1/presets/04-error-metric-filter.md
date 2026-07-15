# Preset: Error-Count Metric Filter

**Use case:** Alarm on errors that only exist in application logs.

This preset pairs a standard 30-day log group with a metric filter that counts
every log event containing `ERROR` and publishes the count as a custom
CloudWatch metric. An `AwsCloudwatchAlarm` on `ErrorCount` then closes the
loop: log line in, page out.

## What You Get

- A STANDARD class CloudWatch Log Group with 30-day retention
- A metric filter publishing `ErrorCount` in your custom namespace
- `defaultValue: 0` so quiet periods report zero instead of missing data —
  alarms behave predictably without special missing-data handling
- Outputs: `log_group_arn`, `log_group_name`

## When to Use

- Alerting on application errors that never surface as AWS service metrics
- Counting business events (signups, payments) logged by the application
- Building error-rate dashboards from structured or unstructured logs

## Key Configuration Choices

- **Plain-term pattern** (`pattern: "ERROR"`) — matches any event containing
  the term; switch to a JSON pattern like `{ $.level = "error" }` for
  structured logs
- **Custom namespace** — keeps log-derived metrics separate from AWS service
  namespaces; every unique namespace+dimensions combination is a billed
  custom metric

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region for the log group | Your deployment region |
| `<app-namespace>` | Custom metric namespace (e.g. `MyApp/Errors`) | Your naming convention |
