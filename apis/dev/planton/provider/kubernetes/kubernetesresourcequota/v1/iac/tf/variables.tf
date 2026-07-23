# Input variables for KubernetesResourceQuota Terraform module.
# These mirror the KubernetesResourceQuotaSpec protobuf schema; the namespace
# StringValueOrRef arrives flattened to a plain string, and enum fields
# arrive as the proto enum value names (e.g. "best_effort", "container").

variable "metadata" {
  description = "Metadata for the resource quota resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes namespace resource governance pair"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # Aggregate caps: resource name -> quantity.
    hard = map(string)

    # Coarse tracking filters: terminating / not_terminating / best_effort /
    # not_best_effort / priority_class / cross_namespace_pod_affinity /
    # volume_attributes_class.
    scopes = optional(list(string), [])

    # Fine-grained scope filters (most usefully priority_class with In/NotIn).
    scope_selector = optional(list(object({
      scope_name = string
      operator   = string
      values     = optional(list(string), [])
    })), [])

    # Per-object defaults and bounds, managed as a companion LimitRange.
    limit_defaults = optional(list(object({
      # "container", "pod", or "persistent_volume_claim".
      type                    = string
      max                     = optional(map(string), {})
      min                     = optional(map(string), {})
      default_limit           = optional(map(string), {})
      default_request         = optional(map(string), {})
      max_limit_request_ratio = optional(map(string), {})
    })), [])
  })
}
