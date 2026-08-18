# One AMP scraper: AWS's agentless Prometheus collector.
#
# Lifecycle facts the render below depends on:
#   - the whole SOURCE block replaces on change (AWS re-provisions the
#     collector); alias, destination, role configuration, and the
#     scrape configuration update in place;
#   - the spec guarantees exactly one source arm and exactly one
#     destination arm; VPC sources carry their own scrape configuration
#     (AWS's published default exists only for EKS);
#   - creates run long (AWS provisions collector infrastructure - the
#     provider waits up to 30 minutes) and deletes drain before removal
#     (up to 20);
#   - the scraper's logging configuration is created via update against
#     the scraper ID and replaces with it.

resource "aws_prometheus_scraper" "this" {
  alias = var.spec.alias != "" ? var.spec.alias : null

  scrape_configuration = local.scrape_configuration

  source {
    dynamic "eks" {
      for_each = var.spec.source_eks != null ? [var.spec.source_eks] : []
      content {
        cluster_arn        = eks.value.cluster_arn
        subnet_ids         = eks.value.subnet_ids
        security_group_ids = length(eks.value.security_group_ids) > 0 ? eks.value.security_group_ids : null
      }
    }

    dynamic "vpc" {
      for_each = var.spec.source_vpc != null ? [var.spec.source_vpc] : []
      content {
        subnet_ids         = vpc.value.subnet_ids
        security_group_ids = vpc.value.security_group_ids
      }
    }
  }

  destination {
    dynamic "amp" {
      for_each = var.spec.amp_workspace_arn != "" ? [var.spec.amp_workspace_arn] : []
      content {
        workspace_arn = amp.value
      }
    }

    dynamic "cloudwatch" {
      for_each = var.spec.cloudwatch_dataset_arn != "" ? [var.spec.cloudwatch_dataset_arn] : []
      content {
        dataset_arn = cloudwatch.value
      }
    }
  }

  dynamic "role_configuration" {
    for_each = var.spec.role_configuration != null ? [var.spec.role_configuration] : []
    content {
      source_role_arn = role_configuration.value.source_role_arn != "" ? role_configuration.value.source_role_arn : null
      target_role_arn = role_configuration.value.target_role_arn != "" ? role_configuration.value.target_role_arn : null
    }
  }

  tags = local.aws_tags
}

resource "aws_prometheus_scraper_logging_configuration" "this" {
  count = var.spec.logging != null ? 1 : 0

  scraper_id = aws_prometheus_scraper.this.id

  logging_destination {
    cloudwatch_logs {
      log_group_arn = local.logging_log_group_arn
    }
  }

  scraper_components = length(var.spec.logging.components) > 0 ? var.spec.logging.components : null
}
