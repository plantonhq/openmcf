# One Amazon Managed Prometheus workspace with its folded satellites.
#
# Lifecycle facts the render below depends on:
#   - the workspace ALIAS can never be unset once set - AWS offers no
#     un-alias, so clearing spec.alias replaces the workspace (the
#     provider's ForceNewIfChange contract);
#   - kms_key_arn replaces the workspace on change;
#   - the workspace CONFIGURATION is created via update and has NO
#     delete API - removing the block is a no-op at AWS and the
#     last-applied retention/limits persist (the settings-retention
#     class);
#   - the alert manager definition is strictly one per workspace (its
#     provider ID is the workspace ID);
#   - the resource policy is revision-guarded server-side; the
#     provider's optional revision_id input is a state-managed
#     concurrency token, deliberately not modeled;
#   - every satellite waits on the workspace through its workspace_id
#     reference - no explicit depends_on needed.

resource "aws_prometheus_workspace" "this" {
  alias       = var.spec.alias != "" ? var.spec.alias : null
  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  dynamic "logging_configuration" {
    for_each = var.spec.logging != null ? [1] : []
    content {
      log_group_arn = local.workspace_log_group_arn
    }
  }

  tags = local.aws_tags
}

resource "aws_prometheus_workspace_configuration" "this" {
  count = var.spec.configuration != null ? 1 : 0

  workspace_id = aws_prometheus_workspace.this.id

  retention_period_in_days            = var.spec.configuration.retention_period_in_days != null ? var.spec.configuration.retention_period_in_days : null
  out_of_order_time_window_in_seconds = var.spec.configuration.out_of_order_time_window_in_seconds != null ? var.spec.configuration.out_of_order_time_window_in_seconds : null
  rule_query_offset_in_seconds        = var.spec.configuration.rule_query_offset_in_seconds != null ? var.spec.configuration.rule_query_offset_in_seconds : null

  dynamic "limits_per_label_set" {
    for_each = var.spec.configuration.limits_per_label_set
    content {
      label_set = limits_per_label_set.value.label_set
      limits {
        max_series = limits_per_label_set.value.max_series
      }
    }
  }
}

resource "aws_prometheus_alert_manager_definition" "this" {
  count = var.spec.alert_manager_definition != "" ? 1 : 0

  workspace_id = aws_prometheus_workspace.this.id
  definition   = var.spec.alert_manager_definition
}

resource "aws_prometheus_rule_group_namespace" "this" {
  for_each = local.rule_group_namespaces

  workspace_id = aws_prometheus_workspace.this.id
  name         = each.value.name
  data         = each.value.data

  tags = local.aws_tags
}

resource "aws_prometheus_query_logging_configuration" "this" {
  count = var.spec.query_logging != null ? 1 : 0

  workspace_id = aws_prometheus_workspace.this.id

  dynamic "destination" {
    for_each = var.spec.query_logging.destinations
    content {
      cloudwatch_logs {
        # The module appends the ":*" suffix AWS requires (see locals).
        log_group_arn = endswith(destination.value.log_group_arn, ":*") ? destination.value.log_group_arn : "${destination.value.log_group_arn}:*"
      }
      filters {
        qsp_threshold = destination.value.qsp_threshold
      }
    }
  }
}

resource "aws_prometheus_resource_policy" "this" {
  count = var.spec.resource_policy != null ? 1 : 0

  workspace_id    = aws_prometheus_workspace.this.id
  policy_document = jsonencode(var.spec.resource_policy)
}

resource "aws_prometheus_anomaly_detector" "this" {
  for_each = local.anomaly_detectors

  workspace_id = aws_prometheus_workspace.this.id
  alias        = each.value.alias

  evaluation_interval_in_seconds = each.value.evaluation_interval_in_seconds != null ? each.value.evaluation_interval_in_seconds : null
  labels                         = length(each.value.labels) > 0 ? each.value.labels : null

  configuration {
    # Random Cut Forest is AWS's only detection algorithm today.
    random_cut_forest {
      query        = each.value.query
      sample_size  = each.value.sample_size != null ? each.value.sample_size : null
      shingle_size = each.value.shingle_size != null ? each.value.shingle_size : null

      dynamic "ignore_near_expected_from_above" {
        for_each = each.value.ignore_near_expected_from_above != null ? [each.value.ignore_near_expected_from_above] : []
        content {
          amount = ignore_near_expected_from_above.value.amount != null ? ignore_near_expected_from_above.value.amount : null
          ratio  = ignore_near_expected_from_above.value.ratio != null ? ignore_near_expected_from_above.value.ratio : null
        }
      }

      dynamic "ignore_near_expected_from_below" {
        for_each = each.value.ignore_near_expected_from_below != null ? [each.value.ignore_near_expected_from_below] : []
        content {
          amount = ignore_near_expected_from_below.value.amount != null ? ignore_near_expected_from_below.value.amount : null
          ratio  = ignore_near_expected_from_below.value.ratio != null ? ignore_near_expected_from_below.value.ratio : null
        }
      }
    }
  }

  # The provider models the action as an exactly-one pair of
  # must-be-true bools; the spec models it honestly as an enum.
  missing_data_action {
    mark_as_anomaly = each.value.missing_data_action == "MARK_AS_ANOMALY" ? true : null
    skip            = each.value.missing_data_action == "SKIP" ? true : null
  }

  tags = local.aws_tags
}
