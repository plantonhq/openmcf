# AwsCloudwatchDashboard — Pulumi module (Go)

Manages one CloudWatch dashboard (`cloudwatch.Dashboard`).

Module facts worth knowing before editing:

- **The name is `spec.dashboard_name`, never `metadata.name`** — dashboard names carry uppercase metadata.name cannot. Renames replace the dashboard.
- **The body is a Struct, JSON-encoded here** — AWS normalizes it server-side and the provider diffs it semantically, so key order and whitespace never cause drift.
- **PutDashboard is a pure upsert** — create and update are the same AWS call; every body change applies in place.
- **No tags** — dashboards are untaggable at AWS, the one deliberate absence against the catalog's tag convention (mirrored in the Terraform module).

Outputs mirror the Terraform module key-for-key: `dashboard_name` (the import ID), `dashboard_arn`.
