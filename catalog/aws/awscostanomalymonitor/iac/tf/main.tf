# A Cost Explorer anomaly monitor - the ML-driven watcher that flags
# unusual spend - with its folded alert subscriptions.
#
# Lifecycle facts the render below depends on:
#   - the monitor's SHAPE arms are create-only: monitor_type,
#     monitor_dimension, and monitor_specification all force
#     replacement (only the display name updates in place);
#   - the spec's CUSTOM-arm Struct renders as the provider's raw
#     Expression JSON string (the provider takes the AWS document
#     verbatim, not typed blocks);
#   - each subscription is one aws_ce_anomaly_subscription bound to
#     THIS monitor's ARN, keyed by its spec name;
#   - the threshold expression's LEVELED spec shape (root -> leaf) is
#     exactly the nesting AWS accepts on subscriptions, so the dynamic
#     blocks below unroll it 1:1 with no depth checks (the Pulumi
#     module walks the same levels).
resource "aws_ce_anomaly_monitor" "this" {
  name         = var.spec.monitor_name
  monitor_type = var.spec.monitor_type

  monitor_dimension     = var.spec.monitor_dimension != "" ? var.spec.monitor_dimension : null
  monitor_specification = var.spec.monitor_specification != null ? jsonencode(var.spec.monitor_specification) : null

  tags = local.aws_tags
}

resource "aws_ce_anomaly_subscription" "this" {
  for_each = local.subscriptions_by_name

  name             = each.value.name
  frequency        = each.value.frequency
  monitor_arn_list = [aws_ce_anomaly_monitor.this.arn]

  dynamic "subscriber" {
    for_each = each.value.subscribers
    content {
      address = subscriber.value.address
      type    = subscriber.value.type
    }
  }

  dynamic "threshold_expression" {
    for_each = each.value.threshold_expression != null ? [each.value.threshold_expression] : []
    content {
      dynamic "dimension" {
        for_each = threshold_expression.value.dimension != null ? [threshold_expression.value.dimension] : []
        content {
          key           = dimension.value.key
          match_options = dimension.value.match_options
          values        = dimension.value.values
        }
      }
      dynamic "tags" {
        for_each = threshold_expression.value.tag != null ? [threshold_expression.value.tag] : []
        content {
          key           = tags.value.key != "" ? tags.value.key : null
          match_options = tags.value.match_options
          values        = tags.value.values
        }
      }
      dynamic "cost_category" {
        for_each = threshold_expression.value.cost_category != null ? [threshold_expression.value.cost_category] : []
        content {
          key           = cost_category.value.key != "" ? cost_category.value.key : null
          match_options = cost_category.value.match_options
          values        = cost_category.value.values
        }
      }
      dynamic "and" {
        for_each = threshold_expression.value.and
        content {
          dynamic "dimension" {
            for_each = and.value.dimension != null ? [and.value.dimension] : []
            content {
              key           = dimension.value.key
              match_options = dimension.value.match_options
              values        = dimension.value.values
            }
          }
          dynamic "tags" {
            for_each = and.value.tag != null ? [and.value.tag] : []
            content {
              key           = tags.value.key != "" ? tags.value.key : null
              match_options = tags.value.match_options
              values        = tags.value.values
            }
          }
          dynamic "cost_category" {
            for_each = and.value.cost_category != null ? [and.value.cost_category] : []
            content {
              key           = cost_category.value.key != "" ? cost_category.value.key : null
              match_options = cost_category.value.match_options
              values        = cost_category.value.values
            }
          }
        }
      }
      dynamic "or" {
        for_each = threshold_expression.value.or
        content {
          dynamic "dimension" {
            for_each = or.value.dimension != null ? [or.value.dimension] : []
            content {
              key           = dimension.value.key
              match_options = dimension.value.match_options
              values        = dimension.value.values
            }
          }
          dynamic "tags" {
            for_each = or.value.tag != null ? [or.value.tag] : []
            content {
              key           = tags.value.key != "" ? tags.value.key : null
              match_options = tags.value.match_options
              values        = tags.value.values
            }
          }
          dynamic "cost_category" {
            for_each = or.value.cost_category != null ? [or.value.cost_category] : []
            content {
              key           = cost_category.value.key != "" ? cost_category.value.key : null
              match_options = cost_category.value.match_options
              values        = cost_category.value.values
            }
          }
        }
      }
      dynamic "not" {
        for_each = threshold_expression.value.not != null ? [threshold_expression.value.not] : []
        content {
          dynamic "dimension" {
            for_each = not.value.dimension != null ? [not.value.dimension] : []
            content {
              key           = dimension.value.key
              match_options = dimension.value.match_options
              values        = dimension.value.values
            }
          }
          dynamic "tags" {
            for_each = not.value.tag != null ? [not.value.tag] : []
            content {
              key           = tags.value.key != "" ? tags.value.key : null
              match_options = tags.value.match_options
              values        = tags.value.values
            }
          }
          dynamic "cost_category" {
            for_each = not.value.cost_category != null ? [not.value.cost_category] : []
            content {
              key           = cost_category.value.key != "" ? cost_category.value.key : null
              match_options = cost_category.value.match_options
              values        = cost_category.value.values
            }
          }
        }
      }
    }
  }

  tags = local.aws_tags
}
