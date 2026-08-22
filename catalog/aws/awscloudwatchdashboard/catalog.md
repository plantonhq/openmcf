# AWS CloudWatch Dashboard

A monitoring dashboard your team opens first thing in an incident — metric graphs, log queries, alarm status, and runbook notes on one named canvas.

## What Gets Managed

- The dashboard's name and its full widget layout (the dashboard body document), applied as an in-place upsert on every change.
- Any widget type CloudWatch supports: metric graphs, logs insights queries, alarm status tiles, text/markdown panels — laid out on the 24-column grid.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with CloudWatch dashboard permissions (`cloudwatch:PutDashboard`, `cloudwatch:GetDashboard`, `cloudwatch:DeleteDashboards`).

### Naming

The dashboard's name is the spec's `dashboard_name` (letters, digits, hyphens, underscores — uppercase legal). Renaming replaces the dashboard.

## After You Deploy

- `dashboard_arn` and `dashboard_name` land in outputs; the console URL is `https://console.aws.amazon.com/cloudwatch/home?region=<region>#dashboards/dashboard/<name>`.
- Dashboards above the account's free tier (3) bill ~$3/month each, prorated.

## Common Changes

- Editing widgets: change `dashboard_body` and re-deploy — in place, always.
- Prototyping in the console: build visually, copy Actions → View/edit source into the manifest, and the dashboard is declarative from then on.
