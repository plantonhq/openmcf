# Input variables for KubernetesPodDisruptionBudget Terraform module.
# These mirror the KubernetesPodDisruptionBudgetSpec protobuf schema; the
# namespace StringValueOrRef arrives flattened to a plain string, and the
# unhealthy_pod_eviction_policy enum arrives as the proto enum value name
# (e.g. "if_healthy_budget", "always_allow").

variable "metadata" {
  description = "Metadata for the pod disruption budget resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes PodDisruptionBudget"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # Selects the protected pods; empty (but present) selects ALL pods in
    # the namespace.
    selector = object({
      match_labels = optional(map(string), {})
      match_expressions = optional(list(object({
        key      = string
        operator = string
        values   = optional(list(string), [])
      })), [])
    })

    # Exactly one of min_available / max_unavailable — absolute ("2") or
    # percentage ("50%").
    min_available   = optional(string, "")
    max_unavailable = optional(string, "")

    # "if_healthy_budget" (default) or "always_allow".
    unhealthy_pod_eviction_policy = optional(string, "if_healthy_budget")
  })
}
