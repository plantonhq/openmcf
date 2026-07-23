# Input variables for KubernetesHorizontalPodAutoscaler Terraform module.
# These mirror the KubernetesHorizontalPodAutoscalerSpec protobuf schema; the
# namespace and scale_target.name StringValueOrRef fields arrive flattened to
# plain strings, and enum fields arrive as the proto enum value names
# (e.g. "resource", "utilization", "max_change", "pods").

variable "metadata" {
  description = "Metadata for the horizontal pod autoscaler resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes HorizontalPodAutoscaler"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    scale_target = object({
      api_version = optional(string, "apps/v1")
      kind        = optional(string, "Deployment")
      name        = string
    })

    min_replicas = optional(number, 1)
    max_replicas = number

    # Metric source families: "resource", "container_resource", "pods",
    # "object", "external" — exactly the matching source object is set.
    metrics = optional(list(object({
      type = string
      resource = optional(object({
        name   = string
        target = object({
          # "utilization", "raw_value", or "average_value".
          type                = string
          average_utilization = optional(number)
          value               = optional(string, "")
          average_value       = optional(string, "")
        })
      }))
      container_resource = optional(object({
        name      = string
        container = string
        target = object({
          type                = string
          average_utilization = optional(number)
          value               = optional(string, "")
          average_value       = optional(string, "")
        })
      }))
      pods = optional(object({
        metric = object({
          name         = string
          match_labels = optional(map(string), {})
        })
        target = object({
          type                = string
          average_utilization = optional(number)
          value               = optional(string, "")
          average_value       = optional(string, "")
        })
      }))
      object = optional(object({
        described_object = object({
          api_version = optional(string, "")
          kind        = string
          name        = string
        })
        metric = object({
          name         = string
          match_labels = optional(map(string), {})
        })
        target = object({
          type                = string
          average_utilization = optional(number)
          value               = optional(string, "")
          average_value       = optional(string, "")
        })
      }))
      external = optional(object({
        metric = object({
          name         = string
          match_labels = optional(map(string), {})
        })
        target = object({
          type                = string
          average_utilization = optional(number)
          value               = optional(string, "")
          average_value       = optional(string, "")
        })
      }))
    })), [])

    # Per-direction scaling tuning.
    behavior = optional(object({
      scale_up = optional(object({
        stabilization_window_seconds = optional(number)
        # "max_change" (default), "min_change", or "disabled".
        select_policy = optional(string, "max_change")
        policies = optional(list(object({
          # "pods" or "percent".
          type           = string
          value          = number
          period_seconds = number
        })), [])
      }))
      scale_down = optional(object({
        stabilization_window_seconds = optional(number)
        select_policy                = optional(string, "max_change")
        policies = optional(list(object({
          type           = string
          value          = number
          period_seconds = number
        })), [])
      }))
    }))
  })
}
