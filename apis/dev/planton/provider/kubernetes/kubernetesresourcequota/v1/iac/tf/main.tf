# Kubernetes namespace resource governance Terraform module: the
# ResourceQuota plus an optional companion LimitRange — two API objects, one
# governance story ("how much may this namespace consume, and what does a
# workload get when it doesn't say?").

resource "kubernetes_resource_quota_v1" "resource_quota" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    hard   = var.spec.hard
    scopes = length(local.scopes) > 0 ? local.scopes : null

    dynamic "scope_selector" {
      for_each = length(try(var.spec.scope_selector, [])) > 0 ? [1] : []
      content {
        dynamic "match_expression" {
          for_each = var.spec.scope_selector
          content {
            scope_name = lookup(local.scope_map, match_expression.value.scope_name, match_expression.value.scope_name)
            operator   = match_expression.value.operator
            values     = length(match_expression.value.values) > 0 ? match_expression.value.values : null
          }
        }
      }
    }
  }
}

# The companion LimitRange: per-object defaults and bounds. This is what
# keeps a compute quota livable — workloads that omit requests/limits inherit
# the defaults instead of being rejected by the quota's admission check.
resource "kubernetes_limit_range_v1" "limit_range" {
  count = local.create_limit_range ? 1 : 0

  metadata {
    name        = local.limit_range_name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    dynamic "limit" {
      for_each = var.spec.limit_defaults
      content {
        type                    = lookup(local.limit_type_map, limit.value.type, "Container")
        max                     = length(limit.value.max) > 0 ? limit.value.max : null
        min                     = length(limit.value.min) > 0 ? limit.value.min : null
        default                 = length(limit.value.default_limit) > 0 ? limit.value.default_limit : null
        default_request         = length(limit.value.default_request) > 0 ? limit.value.default_request : null
        max_limit_request_ratio = length(limit.value.max_limit_request_ratio) > 0 ? limit.value.max_limit_request_ratio : null
      }
    }
  }
}
