# An AWS Budgets budget - a spend/usage threshold AWS evaluates
# continuously - with its alert notifications and folded ACTIONS.
#
# Lifecycle facts the render below depends on:
#   - the budget's name is spec.budget_name (an explicit field - budget
#     names legally carry spaces metadata.name cannot) and changing it
#     replaces the budget;
#   - exactly ONE funding shape renders (spec-enforced): the fixed
#     limit, the planned per-period limits, or auto-adjustment;
#   - the two filter GENERATIONS are mutually exclusive (spec-enforced):
#     legacy cost_filter/cost_types vs the modern metrics +
#     filter_expression pair. The provider's `metrics` is a
#     single-element list - the spec's singular metric renders as
#     [metric];
#   - the filter expression's LEVELED spec shape (root -> node -> leaf)
#     is exactly the nesting AWS accepts, so the dynamic blocks below
#     unroll it 1:1 with no depth checks - neither side can express an
#     illegal tree (the Pulumi module walks the same levels);
#   - both resources are taggable; account_id (member-account budgets
#     managed from a payer account) is create-only on both.
resource "aws_budgets_budget" "this" {
  name        = var.spec.budget_name
  budget_type = var.spec.budget_type
  time_unit   = var.spec.time_unit

  account_id       = var.spec.account_id != "" ? var.spec.account_id : null
  billing_view_arn = var.spec.billing_view_arn != "" ? var.spec.billing_view_arn : null

  # The fixed-limit funding shape.
  limit_amount = var.spec.limit != null ? var.spec.limit.amount : null
  limit_unit   = var.spec.limit != null ? var.spec.limit.unit : null

  # The planned-limits funding shape (each period gets its own ceiling).
  dynamic "planned_limit" {
    for_each = var.spec.planned_limits
    content {
      start_time = planned_limit.value.start_time
      amount     = planned_limit.value.amount
      unit       = planned_limit.value.unit
    }
  }

  # The auto-adjusting funding shape: AWS recomputes the limit from
  # history or forecast.
  dynamic "auto_adjust_data" {
    for_each = var.spec.auto_adjust != null ? [var.spec.auto_adjust] : []
    content {
      auto_adjust_type = auto_adjust_data.value.auto_adjust_type
      dynamic "historical_options" {
        for_each = auto_adjust_data.value.budget_adjustment_period != 0 ? [1] : []
        content {
          budget_adjustment_period = auto_adjust_data.value.budget_adjustment_period
        }
      }
    }
  }

  time_period_start = var.spec.time_period_start != "" ? var.spec.time_period_start : null
  time_period_end   = var.spec.time_period_end != "" ? var.spec.time_period_end : null

  # Legacy-generation cost-component toggles. Presence-typed fields
  # render only when set - the provider then applies AWS's defaults
  # (include_* true; use_blended/use_amortized false) for the rest.
  dynamic "cost_types" {
    for_each = var.spec.cost_types != null ? [var.spec.cost_types] : []
    content {
      include_credit             = cost_types.value.include_credit
      include_discount           = cost_types.value.include_discount
      include_other_subscription = cost_types.value.include_other_subscription
      include_recurring          = cost_types.value.include_recurring
      include_refund             = cost_types.value.include_refund
      include_subscription       = cost_types.value.include_subscription
      include_support            = cost_types.value.include_support
      include_tax                = cost_types.value.include_tax
      include_upfront            = cost_types.value.include_upfront
      use_amortized              = cost_types.value.use_amortized
      use_blended                = cost_types.value.use_blended
    }
  }

  # Legacy-generation name/values filters (entries AND together).
  dynamic "cost_filter" {
    for_each = var.spec.cost_filters
    content {
      name   = cost_filter.value.name
      values = cost_filter.value.values
    }
  }

  # The modern generation's measure - the provider takes a
  # single-element list.
  metrics = var.spec.metric != "" ? [var.spec.metric] : null

  # The modern filter tree, unrolled level by level (root -> node ->
  # leaf, the exact nesting AWS accepts).
  dynamic "filter_expression" {
    for_each = var.spec.filter_expression != null ? [var.spec.filter_expression] : []
    content {
      dynamic "dimensions" {
        for_each = filter_expression.value.dimension != null ? [filter_expression.value.dimension] : []
        content {
          key           = dimensions.value.key
          match_options = dimensions.value.match_options
          values        = dimensions.value.values
        }
      }
      dynamic "tags" {
        for_each = filter_expression.value.tag != null ? [filter_expression.value.tag] : []
        content {
          key           = tags.value.key != "" ? tags.value.key : null
          match_options = tags.value.match_options
          values        = tags.value.values
        }
      }
      dynamic "cost_categories" {
        for_each = filter_expression.value.cost_category != null ? [filter_expression.value.cost_category] : []
        content {
          key           = cost_categories.value.key != "" ? cost_categories.value.key : null
          match_options = cost_categories.value.match_options
          values        = cost_categories.value.values
        }
      }
      dynamic "and" {
        for_each = filter_expression.value.and
        content {
          dynamic "dimensions" {
            for_each = and.value.dimension != null ? [and.value.dimension] : []
            content {
              key           = dimensions.value.key
              match_options = dimensions.value.match_options
              values        = dimensions.value.values
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
          dynamic "cost_categories" {
            for_each = and.value.cost_category != null ? [and.value.cost_category] : []
            content {
              key           = cost_categories.value.key != "" ? cost_categories.value.key : null
              match_options = cost_categories.value.match_options
              values        = cost_categories.value.values
            }
          }
          dynamic "and" {
            for_each = and.value.and
            content {
              dynamic "dimensions" {
                for_each = and.value.dimension != null ? [and.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = and.value.cost_category != null ? [and.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "or" {
            for_each = and.value.or
            content {
              dynamic "dimensions" {
                for_each = or.value.dimension != null ? [or.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = or.value.cost_category != null ? [or.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "not" {
            for_each = and.value.not != null ? [and.value.not] : []
            content {
              dynamic "dimensions" {
                for_each = not.value.dimension != null ? [not.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = not.value.cost_category != null ? [not.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
        }
      }
      dynamic "or" {
        for_each = filter_expression.value.or
        content {
          dynamic "dimensions" {
            for_each = or.value.dimension != null ? [or.value.dimension] : []
            content {
              key           = dimensions.value.key
              match_options = dimensions.value.match_options
              values        = dimensions.value.values
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
          dynamic "cost_categories" {
            for_each = or.value.cost_category != null ? [or.value.cost_category] : []
            content {
              key           = cost_categories.value.key != "" ? cost_categories.value.key : null
              match_options = cost_categories.value.match_options
              values        = cost_categories.value.values
            }
          }
          dynamic "and" {
            for_each = or.value.and
            content {
              dynamic "dimensions" {
                for_each = and.value.dimension != null ? [and.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = and.value.cost_category != null ? [and.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "or" {
            for_each = or.value.or
            content {
              dynamic "dimensions" {
                for_each = or.value.dimension != null ? [or.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = or.value.cost_category != null ? [or.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "not" {
            for_each = or.value.not != null ? [or.value.not] : []
            content {
              dynamic "dimensions" {
                for_each = not.value.dimension != null ? [not.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = not.value.cost_category != null ? [not.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
        }
      }
      dynamic "not" {
        for_each = var.spec.filter_expression.not != null ? [var.spec.filter_expression.not] : []
        content {
          dynamic "dimensions" {
            for_each = not.value.dimension != null ? [not.value.dimension] : []
            content {
              key           = dimensions.value.key
              match_options = dimensions.value.match_options
              values        = dimensions.value.values
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
          dynamic "cost_categories" {
            for_each = not.value.cost_category != null ? [not.value.cost_category] : []
            content {
              key           = cost_categories.value.key != "" ? cost_categories.value.key : null
              match_options = cost_categories.value.match_options
              values        = cost_categories.value.values
            }
          }
          dynamic "and" {
            for_each = not.value.and
            content {
              dynamic "dimensions" {
                for_each = and.value.dimension != null ? [and.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = and.value.cost_category != null ? [and.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "or" {
            for_each = not.value.or
            content {
              dynamic "dimensions" {
                for_each = or.value.dimension != null ? [or.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = or.value.cost_category != null ? [or.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
          dynamic "not" {
            for_each = not.value.not != null ? [not.value.not] : []
            content {
              dynamic "dimensions" {
                for_each = not.value.dimension != null ? [not.value.dimension] : []
                content {
                  key           = dimensions.value.key
                  match_options = dimensions.value.match_options
                  values        = dimensions.value.values
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
              dynamic "cost_categories" {
                for_each = not.value.cost_category != null ? [not.value.cost_category] : []
                content {
                  key           = cost_categories.value.key != "" ? cost_categories.value.key : null
                  match_options = cost_categories.value.match_options
                  values        = cost_categories.value.values
                }
              }
            }
          }
        }
      }
    }
  }

  # Threshold alerts. AWS requires at least one subscriber per
  # notification (spec-enforced).
  dynamic "notification" {
    for_each = var.spec.notifications
    content {
      comparison_operator        = notification.value.comparison_operator
      notification_type          = notification.value.notification_type
      threshold                  = notification.value.threshold
      threshold_type             = notification.value.threshold_type
      subscriber_email_addresses = length(notification.value.subscriber_email_addresses) > 0 ? notification.value.subscriber_email_addresses : null
      subscriber_sns_topic_arns  = length(notification.value.subscriber_sns_topic_arns) > 0 ? notification.value.subscriber_sns_topic_arns : null
    }
  }

  tags = local.aws_tags
}

# The folded budget actions - each entry is one aws_budgets_budget_action
# keyed by its spec name. The definition arm matches action_type
# (spec-enforced); references (execution role, policy ARNs, principals,
# instance IDs) arrive resolved.
resource "aws_budgets_budget_action" "this" {
  for_each = local.actions_by_name

  budget_name        = aws_budgets_budget.this.name
  action_type        = each.value.action_type
  approval_model     = each.value.approval_model
  notification_type  = each.value.notification_type
  execution_role_arn = each.value.execution_role_arn

  # The action inherits the budget's account scope: a member-account
  # budget's actions must live in the same account.
  account_id = var.spec.account_id != "" ? var.spec.account_id : null

  action_threshold {
    action_threshold_type  = each.value.action_threshold.action_threshold_type
    action_threshold_value = each.value.action_threshold.action_threshold_value
  }

  definition {
    dynamic "iam_action_definition" {
      for_each = each.value.iam_action_definition != null ? [each.value.iam_action_definition] : []
      content {
        policy_arn = iam_action_definition.value.policy_arn
        groups     = length(iam_action_definition.value.groups) > 0 ? iam_action_definition.value.groups : null
        roles      = length(iam_action_definition.value.roles) > 0 ? iam_action_definition.value.roles : null
        users      = length(iam_action_definition.value.users) > 0 ? iam_action_definition.value.users : null
      }
    }
    dynamic "scp_action_definition" {
      for_each = each.value.scp_action_definition != null ? [each.value.scp_action_definition] : []
      content {
        policy_id  = scp_action_definition.value.policy_id
        target_ids = scp_action_definition.value.target_ids
      }
    }
    dynamic "ssm_action_definition" {
      for_each = each.value.ssm_action_definition != null ? [each.value.ssm_action_definition] : []
      content {
        action_sub_type = ssm_action_definition.value.action_sub_type
        region          = ssm_action_definition.value.region
        instance_ids    = ssm_action_definition.value.instance_ids
      }
    }
  }

  dynamic "subscriber" {
    for_each = each.value.subscribers
    content {
      address           = subscriber.value.address
      subscription_type = subscriber.value.subscription_type
    }
  }

  tags = local.aws_tags
}
