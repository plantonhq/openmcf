# autoscaling/v2 HorizontalPodAutoscaler, created when autoscaling is enabled.
# CPU targets are average Utilization (percentage of requests); memory targets
# are an absolute AverageValue per pod — matching how each metric is
# meaningfully compared across replicas. The availability replica count is the
# floor; max_replicas the ceiling.
resource "kubernetes_horizontal_pod_autoscaler_v2" "this" {
  count = local.hpa_enabled ? 1 : 0

  metadata {
    name      = var.metadata.name
    namespace = local.namespace
    labels    = local.final_labels
  }

  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = var.metadata.name
    }

    min_replicas = local.replicas
    max_replicas = var.spec.availability.horizontal_pod_autoscaling.max_replicas

    dynamic "metric" {
      for_each = try(var.spec.availability.horizontal_pod_autoscaling.target_cpu_utilization_percent, null) != null ? [var.spec.availability.horizontal_pod_autoscaling.target_cpu_utilization_percent] : []
      content {
        type = "Resource"
        resource {
          name = "cpu"
          target {
            type                = "Utilization"
            average_utilization = metric.value
          }
        }
      }
    }

    dynamic "metric" {
      for_each = try(var.spec.availability.horizontal_pod_autoscaling.target_memory_utilization, "") != "" ? [var.spec.availability.horizontal_pod_autoscaling.target_memory_utilization] : []
      content {
        type = "Resource"
        resource {
          name = "memory"
          target {
            type          = "AverageValue"
            average_value = metric.value
          }
        }
      }
    }
  }

  depends_on = [kubernetes_deployment_v1.this]
}
