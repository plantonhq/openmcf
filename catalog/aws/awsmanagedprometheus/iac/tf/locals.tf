locals {
  # Resource-identity tags match the Pulumi module key-for-key (applied
  # to the workspace, rule group namespaces, and anomaly detectors; the
  # configuration/alert-manager/query-logging/resource-policy
  # satellites are untaggable at AWS).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsManagedPrometheus"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # AWS requires log-group ARNs ending in ":*" on AMP's logging fields,
  # while the log group resource (and the AwsCloudwatchLogGroup kind's
  # output) exports the bare ARN. The module owns that quirk: append the
  # suffix when absent so specs wire the natural output.
  workspace_log_group_arn = var.spec.logging != null ? (
    endswith(var.spec.logging.log_group_arn, ":*")
    ? var.spec.logging.log_group_arn
    : "${var.spec.logging.log_group_arn}:*"
  ) : ""

  rule_group_namespaces = { for n in var.spec.rule_group_namespaces : n.name => n }

  anomaly_detectors = { for d in var.spec.anomaly_detectors : d.alias => d }
}
