# AwsManagedPrometheus

One Amazon Managed Prometheus (AMP) workspace with its folded satellites: workspace configuration (retention and label-set series limits), workspace logging, the alert manager definition, name-keyed rule group namespaces, query logging, the workspace resource policy, and alias-keyed anomaly detectors.

## Highlights

- **Scrapers are deliberately NOT here** — an AMP scraper can target CloudWatch with zero AMP workspaces, so it is its own kind (AwsManagedPrometheusScraper) referencing this workspace's ARN.
- **Two AWS contracts taught in place**: the alias can never be unset once set (clearing it replaces the workspace), and the workspace configuration persists after destroy (AWS has no delete API for it — the settings-retention class).
- **Log-group wiring just works**: AWS demands `:*`-suffixed log-group ARNs on the logging fields while AwsCloudwatchLogGroup exports the bare ARN — both modules append the suffix, so charts wire the natural output.
- **Anomaly detectors flattened honestly**: Random Cut Forest is AWS's only algorithm today, so its knobs sit on the detector entry; the missing-data action is an enum, rendered as the provider's exactly-one-true bool pair.

## Both Engines

Both modules render the seven-resource surface identically and export the same outputs: `workspace_id`, `workspace_arn`, `prometheus_endpoint`, `rule_group_namespace_arns` (by name), `anomaly_detector_ids`/`anomaly_detector_arns` (by alias).

## Chart Wiring

`kms_key_arn` → AwsKmsKey; logging and query-logging destinations → AwsCloudwatchLogGroup; downstream, scrapers and remote-write clients consume `workspace_arn`/`prometheus_endpoint`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
