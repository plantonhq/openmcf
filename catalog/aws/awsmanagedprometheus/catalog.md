# AWS Managed Prometheus

Prometheus without running Prometheus: an AWS-managed, PromQL-compatible metrics backend your clusters remote-write into — with recording/alerting rules, a managed Alertmanager, retention controls, and ML anomaly detection on your queries.

## What Gets Managed

- The workspace (alias, optional KMS encryption, event logging) and its configuration: metric retention, out-of-order ingest window, rule query offset, and per-label-set active-series caps.
- Rules as code: recording/alerting rule files (standard Prometheus rules YAML) as name-keyed namespaces, plus the alertmanager.yml document (SNS is the supported receiver).
- Query logging (which expensive queries get logged, and where), the workspace resource policy (cross-account remote-write/query grants), and PromQL anomaly detectors.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AMP permissions (plus CloudWatch Logs when logging arms are set, KMS for a customer-managed key).

### AWS Prerequisites

- Log groups for the logging arms (reference AwsCloudwatchLogGroup outputs — the modules add AWS's required `:*` suffix).

## After You Deploy

- `prometheus_endpoint` is the remote-write/query URL; `workspace_arn` is what scrapers and IAM policies reference.
- Idle workspaces bill nothing — AMP prices ingestion, storage, and queries.

## Common Changes

- Rules iterate in place (each namespace's `data` is one rules file); the alert manager document likewise.
- Retention/limits changes update in place; REMEMBER removing the whole configuration block is a no-op at AWS (the last-applied values persist).
