# AwsManagedPrometheusScraper

One AMP scraper: AWS's agentless Prometheus collector. It scrapes a SOURCE (an EKS cluster, or a bare VPC placement for non-EKS endpoints) and writes to a DESTINATION (an AMP workspace or a CloudWatch dataset) — its own kind because it needs no AMP workspace to exist.

## Highlights

- **Exactly one source and one destination arm** (CEL-mirrored from the provider's unions); the whole source replaces on change (AWS re-provisions the collector), while alias, destination, roles, and the scrape configuration update in place.
- **The scrape configuration is optional on the EKS arm** — AWS publishes a default, and the modules resolve it at deploy time (the data-source-at-plan idiom); the VPC arm requires your own (no default exists there).
- **Long lifecycles taught in place**: creates wait up to 30 minutes while AWS provisions collectors; deletes drain up to 20.
- **Folded logging satellite**: scraper component logs (service discovery / collector / exporter) to a CloudWatch log group — the modules append the `:*` suffix AWS requires, so charts wire AwsCloudwatchLogGroup's natural output.

## Both Engines

Both modules render the scraper and its logging configuration identically and export the same outputs: `scraper_id` (the import ID), `scraper_arn`, and `scraper_role_arn` (the AWS-managed writer role — grant it remote-write on cross-account destinations).

## Chart Wiring

`source_eks.cluster_arn` → AwsEksCluster; subnets/security groups → AwsSubnet/AwsSecurityGroup; `amp_workspace_arn` → AwsManagedPrometheus; `role_configuration` → AwsIamRole pairs for cross-account scraping.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
