# AwsManagedPrometheusScraper — Terraform/OpenTofu module

Manages one AMP scraper and its logging configuration (`aws_prometheus_scraper`, `aws_prometheus_scraper_logging_configuration`).

Module facts worth knowing before editing:

- **The whole source replaces on change** (AWS re-provisions the collector); alias, destination, role configuration, and the scrape configuration update in place.
- **Unset scrape configuration resolves AWS's published default at deploy** via the `aws_prometheus_default_scraper_configuration` data source — EKS sources only (spec-guaranteed; VPC sources bring their own).
- **The logging destination's log-group ARN gets `:*` appended here** — AWS requires the wildcard suffix while the log group resource exports the bare ARN.
- **Creates run long** — the provider waits up to 30 minutes for collector provisioning; deletes drain up to 20.
- **The scraper is taggable; its logging-configuration satellite is not.**

Outputs mirror the Pulumi module key-for-key: `scraper_id` (the import ID), `scraper_arn`, `scraper_role_arn` (the AWS-managed writer role — grant it remote-write on cross-account destinations).
