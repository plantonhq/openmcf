# AwsCloudwatchDashboard

One CloudWatch dashboard: a named canvas of widgets (metric graphs, logs insights queries, alarm status, text) defined by the dashboard body document — the same JSON the CloudWatch console shows under Actions → View/edit source.

## Highlights

- **The body is structured configuration**, not a string blob: the widget layout lives in the manifest as YAML, renders as the provider's JSON document, and diffs semantically — formatting never causes drift.
- **Everything updates in place**: AWS's PutDashboard is a pure upsert, so layout iterations never replace the dashboard.
- **The name is an explicit spec field** (`dashboard_name`) because dashboard names carry uppercase letters; changing it replaces the dashboard.
- **Untaggable at AWS** — the one deliberate absence from the catalog's tag convention (the resource has no tags argument).

## Both Engines

The Terraform/OpenTofu and Pulumi modules render the same single resource and export the same outputs: `dashboard_name` (the import ID) and `dashboard_arn`.

## Chart Wiring

Dashboards chart metrics by namespace/name, not by resource reference — widget properties carry the metric coordinates of whatever the chart deploys (Lambda names, ALB dimensions, custom namespaces).
