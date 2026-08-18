locals {
  # Resource-identity tags match the Pulumi module key-for-key (the
  # scraper is taggable; its logging-configuration satellite is not).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsManagedPrometheusScraper"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Unset scrape configuration resolves to AWS's published default (EKS
  # sources only, spec-guaranteed).
  scrape_configuration = var.spec.scrape_configuration != "" ? var.spec.scrape_configuration : data.aws_prometheus_default_scraper_configuration.this[0].configuration

  # AWS requires log-group ARNs ending in ":*" on the scraper's logging
  # destination, while the log group resource (and the
  # AwsCloudwatchLogGroup kind's output) exports the bare ARN. The
  # module owns that quirk.
  logging_log_group_arn = var.spec.logging != null ? (
    endswith(var.spec.logging.log_group_arn, ":*")
    ? var.spec.logging.log_group_arn
    : "${var.spec.logging.log_group_arn}:*"
  ) : ""
}
