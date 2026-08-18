# AwsManagedPrometheus — Terraform/OpenTofu module

Manages one AMP workspace and its folded satellites (`aws_prometheus_workspace`, `_workspace_configuration`, `_alert_manager_definition`, `_rule_group_namespace`, `_query_logging_configuration`, `_resource_policy`, `_anomaly_detector`).

Module facts worth knowing before editing:

- **The alias can never be unset once set** — AWS offers no un-alias; clearing `spec.alias` replaces the workspace (the provider's ForceNewIfChange contract). `kms_key_arn` also replaces.
- **The workspace configuration persists after destroy** — AWS has no delete API for it (created via update, no-op delete); removing the block leaves the last-applied retention/limits in place.
- **Log-group ARNs get `:*` appended here** — AWS requires the wildcard suffix on the workspace-logging and query-logging fields while the log group resource exports the bare ARN; the module owns that quirk so specs wire the natural AwsCloudwatchLogGroup output.
- **The alert manager definition is strictly one per workspace** (its provider ID is the workspace ID).
- **The resource policy's revision_id is deliberately not modeled** — a state-managed concurrency token, not declarative config.
- **The anomaly detector's missing-data action** is the provider's exactly-one pair of must-be-true bools, rendered from the spec's honest enum.
- **Tags land on the taggable three** (workspace, rule group namespaces, anomaly detectors).

Outputs mirror the Pulumi module key-for-key: `workspace_id`, `workspace_arn`, `prometheus_endpoint`, `rule_group_namespace_arns` (keyed by name), `anomaly_detector_ids`/`anomaly_detector_arns` (keyed by alias).
