# Kubernetes HorizontalPodAutoscaler (autoscaling/v2) Terraform module.
#
# min_replicas is ALWAYS sent explicitly (Kubernetes default 1 applied in
# locals) so both engines submit identical objects. When the spec lists no
# metrics, the metric blocks are OMITTED — the API server then applies its
# own default (80% average CPU utilization).
#
# The five metric source families and the target's three value forms pass
# through 1:1; the spec's validations guarantee exactly the source matching
# each metric's type (and exactly the value form matching each target's
# type) is present, so the dynamic blocks below render at most one source
# per metric.

resource "kubernetes_horizontal_pod_autoscaler_v2" "horizontal_pod_autoscaler" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    scale_target_ref {
      api_version = local.target_api_version
      kind        = local.target_kind
      name        = var.spec.scale_target.name
    }

    min_replicas = local.min_replicas
    max_replicas = var.spec.max_replicas

    dynamic "metric" {
      for_each = try(var.spec.metrics, [])
      content {
        type = lookup(local.metric_type_map, metric.value.type, metric.value.type)

        dynamic "resource" {
          for_each = metric.value.resource != null ? [metric.value.resource] : []
          content {
            name = resource.value.name
            target {
              type                = lookup(local.target_type_map, resource.value.target.type, resource.value.target.type)
              average_utilization = try(resource.value.target.average_utilization, null)
              value               = resource.value.target.value != "" ? resource.value.target.value : null
              average_value       = resource.value.target.average_value != "" ? resource.value.target.average_value : null
            }
          }
        }

        dynamic "container_resource" {
          for_each = metric.value.container_resource != null ? [metric.value.container_resource] : []
          content {
            name      = container_resource.value.name
            container = container_resource.value.container
            target {
              type                = lookup(local.target_type_map, container_resource.value.target.type, container_resource.value.target.type)
              average_utilization = try(container_resource.value.target.average_utilization, null)
              value               = container_resource.value.target.value != "" ? container_resource.value.target.value : null
              average_value       = container_resource.value.target.average_value != "" ? container_resource.value.target.average_value : null
            }
          }
        }

        dynamic "pods" {
          for_each = metric.value.pods != null ? [metric.value.pods] : []
          content {
            metric {
              name = pods.value.metric.name
              dynamic "selector" {
                for_each = length(try(pods.value.metric.match_labels, {})) > 0 ? [1] : []
                content {
                  match_labels = pods.value.metric.match_labels
                }
              }
            }
            target {
              type                = lookup(local.target_type_map, pods.value.target.type, pods.value.target.type)
              average_utilization = try(pods.value.target.average_utilization, null)
              value               = pods.value.target.value != "" ? pods.value.target.value : null
              average_value       = pods.value.target.average_value != "" ? pods.value.target.average_value : null
            }
          }
        }

        dynamic "object" {
          for_each = metric.value.object != null ? [metric.value.object] : []
          content {
            described_object {
              api_version = object.value.described_object.api_version
              kind        = object.value.described_object.kind
              name        = object.value.described_object.name
            }
            metric {
              name = object.value.metric.name
              dynamic "selector" {
                for_each = length(try(object.value.metric.match_labels, {})) > 0 ? [1] : []
                content {
                  match_labels = object.value.metric.match_labels
                }
              }
            }
            target {
              type                = lookup(local.target_type_map, object.value.target.type, object.value.target.type)
              average_utilization = try(object.value.target.average_utilization, null)
              value               = object.value.target.value != "" ? object.value.target.value : null
              average_value       = object.value.target.average_value != "" ? object.value.target.average_value : null
            }
          }
        }

        dynamic "external" {
          for_each = metric.value.external != null ? [metric.value.external] : []
          content {
            metric {
              name = external.value.metric.name
              dynamic "selector" {
                for_each = length(try(external.value.metric.match_labels, {})) > 0 ? [1] : []
                content {
                  match_labels = external.value.metric.match_labels
                }
              }
            }
            target {
              type                = lookup(local.target_type_map, external.value.target.type, external.value.target.type)
              average_utilization = try(external.value.target.average_utilization, null)
              value               = external.value.target.value != "" ? external.value.target.value : null
              average_value       = external.value.target.average_value != "" ? external.value.target.average_value : null
            }
          }
        }
      }
    }

    dynamic "behavior" {
      for_each = try(var.spec.behavior, null) != null ? [var.spec.behavior] : []
      content {
        dynamic "scale_up" {
          for_each = behavior.value.scale_up != null ? [behavior.value.scale_up] : []
          content {
            stabilization_window_seconds = try(scale_up.value.stabilization_window_seconds, null)
            # Always sent with the API default (Max) applied, identical to
            # the Pulumi module.
            select_policy = lookup(local.select_policy_map, try(scale_up.value.select_policy, "max_change"), "Max")
            dynamic "policy" {
              for_each = try(scale_up.value.policies, [])
              content {
                type           = lookup(local.policy_type_map, policy.value.type, policy.value.type)
                value          = policy.value.value
                period_seconds = policy.value.period_seconds
              }
            }
          }
        }
        dynamic "scale_down" {
          for_each = behavior.value.scale_down != null ? [behavior.value.scale_down] : []
          content {
            stabilization_window_seconds = try(scale_down.value.stabilization_window_seconds, null)
            select_policy                = lookup(local.select_policy_map, try(scale_down.value.select_policy, "max_change"), "Max")
            dynamic "policy" {
              for_each = try(scale_down.value.policies, [])
              content {
                type           = lookup(local.policy_type_map, policy.value.type, policy.value.type)
                value          = policy.value.value
                period_seconds = policy.value.period_seconds
              }
            }
          }
        }
      }
    }
  }
}
