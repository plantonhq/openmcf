# A Cost Explorer cost category: a named dimension YOU define that
# groups every line of spend into exactly one of its values. Rules
# evaluate IN ORDER and the first match wins.
#
# Lifecycle facts the render below depends on:
#   - the category's name is spec.category_name (an explicit field -
#     names legally carry spaces metadata.name cannot) and changing it
#     replaces the category;
#   - rule_version is module-pinned to "CostCategoryExpression.v1",
#     the only value the AWS API accepts - never a spec knob;
#   - rules are ORDERED (first match wins), so the module renders them
#     exactly in spec order;
#   - the rule expression's LEVELED spec shape (root -> node -> leaf)
#     is exactly the nesting AWS accepts, so the dynamic blocks below
#     unroll it 1:1 with no depth checks - neither side can express an
#     illegal tree (the Pulumi module walks the same levels).
resource "aws_ce_cost_category" "this" {
  name         = var.spec.category_name
  rule_version = "CostCategoryExpression.v1"

  default_value   = var.spec.default_value != "" ? var.spec.default_value : null
  effective_start = var.spec.effective_start != "" ? var.spec.effective_start : null

  dynamic "rule" {
    for_each = var.spec.rules
    content {
      # Always sent: AWS materializes type = "REGULAR" on expression
      # rules, so an omitted type is a perpetual post-apply diff (the
      # materialized-default class; server-verified 2026-08-25).
      type  = rule.value.type != "" ? rule.value.type : "REGULAR"
      value = rule.value.value != "" ? rule.value.value : null

      # The INHERITED_VALUE shape: the category value comes from the
      # dimension itself (the account name or a tag's value).
      dynamic "inherited_value" {
        for_each = rule.value.inherited_value != null ? [rule.value.inherited_value] : []
        content {
          dimension_name = inherited_value.value.dimension_name
          dimension_key  = inherited_value.value.dimension_key != "" ? inherited_value.value.dimension_key : null
        }
      }

      # The REGULAR shape: a Cost Explorer expression selects the spend
      # this rule matches. The explicit iterator name keeps the inner
      # block from shadowing the outer rule iterator.
      dynamic "rule" {
        for_each = rule.value.rule != null ? [rule.value.rule] : []
        iterator = expr
        content {
          dynamic "dimension" {
            for_each = expr.value.dimension != null ? [expr.value.dimension] : []
            content {
              key           = dimension.value.key
              match_options = dimension.value.match_options
              values        = dimension.value.values
            }
          }
          dynamic "tags" {
            for_each = expr.value.tag != null ? [expr.value.tag] : []
            content {
              key           = tags.value.key != "" ? tags.value.key : null
              match_options = tags.value.match_options
              values        = tags.value.values
            }
          }
          dynamic "cost_category" {
            for_each = expr.value.cost_category != null ? [expr.value.cost_category] : []
            content {
              key           = cost_category.value.key != "" ? cost_category.value.key : null
              match_options = cost_category.value.match_options
              values        = cost_category.value.values
            }
          }
          dynamic "and" {
            for_each = expr.value.and
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
              dynamic "and" {
                for_each = and.value.and
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
                for_each = and.value.or
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
                for_each = and.value.not != null ? [and.value.not] : []
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
          dynamic "or" {
            for_each = expr.value.or
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
              dynamic "and" {
                for_each = or.value.and
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
                for_each = or.value.or
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
                for_each = or.value.not != null ? [or.value.not] : []
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
          dynamic "not" {
            for_each = expr.value.not != null ? [expr.value.not] : []
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
              dynamic "and" {
                for_each = not.value.and
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
                for_each = not.value.or
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
                for_each = not.value.not != null ? [not.value.not] : []
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
        }
      }
    }
  }

  # Split-charge rules: re-allocate one value's costs across targets.
  dynamic "split_charge_rule" {
    for_each = var.spec.split_charge_rules
    content {
      method  = split_charge_rule.value.method
      source  = split_charge_rule.value.source
      targets = split_charge_rule.value.targets

      dynamic "parameter" {
        for_each = split_charge_rule.value.parameters
        content {
          type   = parameter.value.type
          values = parameter.value.values
        }
      }
    }
  }

  tags = local.aws_tags
}
