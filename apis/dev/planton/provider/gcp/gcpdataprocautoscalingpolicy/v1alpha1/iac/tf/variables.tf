variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP Dataproc autoscaling policy"
  type = object({
    # The GCP project for the policy. The CLI's tfvars converter resolves
    # StringValueOrRef fields to their literal string before the module
    # runs, so this arrives as a plain string. If empty, the provider's
    # default project is used (see locals.tf).
    project_id = optional(string, "")

    # Policy ID (the GCP resource name). Immutable (ForceNew).
    policy_id = string

    # Region the policy lives in. A cluster can only attach policies in
    # its own region. Immutable (ForceNew).
    location = string

    # Primary worker bounds. max_instances is required; min_instances 0
    # accepts the API default (2). weight steers the primary/secondary
    # capacity split.
    worker_config = object({
      max_instances = number
      min_instances = optional(number, 0)
      weight        = optional(number, 0)
    })

    # Secondary (preemptible/spot) worker bounds. All default to 0 —
    # the group can scale to zero when idle.
    secondary_worker_config = optional(object({
      max_instances = optional(number, 0)
      min_instances = optional(number, 0)
      weight        = optional(number, 0)
    }), null)

    # The YARN-based autoscaling algorithm and its evaluation cadence.
    basic_algorithm = object({
      cooldown_period = optional(string, "")
      yarn_config = object({
        graceful_decommission_timeout  = string
        scale_up_factor                = number
        scale_down_factor              = number
        scale_up_min_worker_fraction   = optional(number, 0)
        scale_down_min_worker_fraction = optional(number, 0)
      })
    })
  })

  validation {
    condition     = var.spec.policy_id != ""
    error_message = "policy_id is required."
  }

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.worker_config.max_instances >= 1
    error_message = "worker_config.max_instances must be at least 1."
  }
}
