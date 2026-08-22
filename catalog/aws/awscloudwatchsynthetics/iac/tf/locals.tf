locals {
  # Resource-identity tags match the Pulumi module key-for-key (applied
  # to the canary and every owned group; the association join is
  # untaggable at AWS).
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCloudwatchSynthetics"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The provider wants one "s3://bucket/prefix" artifact location string;
  # the spec models bucket and prefix separately for chart wiring.
  artifact_s3_location = var.spec.canary != null ? (
    var.spec.canary.artifact_prefix != ""
    ? "s3://${var.spec.canary.artifact_bucket}/${trim(var.spec.canary.artifact_prefix, "/")}"
    : "s3://${var.spec.canary.artifact_bucket}"
  ) : ""

  # Owned groups keyed by name for the for_each entries.
  groups_by_name = { for g in var.spec.groups : g.name => g }

  # Group joins render only alongside the canary (spec-guaranteed).
  group_names = var.spec.canary != null ? toset(var.spec.group_names) : toset([])
}
