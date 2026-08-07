# Input variables for KubernetesPriorityClass Terraform module.
# These mirror the KubernetesPriorityClassSpec protobuf schema; the
# preemption_policy enum arrives as the proto enum value name
# (e.g. "preempt_lower_priority", "never").

variable "metadata" {
  description = "Metadata for the priority class resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes PriorityClass"
  type = object({
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    # The priority integer (user classes <= 1,000,000,000). Immutable.
    value = number

    # Cluster-wide default for pods that name no class.
    global_default = optional(bool, false)

    description = optional(string, "")

    # "preempt_lower_priority" (default) or "never".
    preemption_policy = optional(string, "preempt_lower_priority")
  })
}
