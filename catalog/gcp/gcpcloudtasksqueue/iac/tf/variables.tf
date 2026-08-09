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
  description = "Specification for the GCP Cloud Tasks queue"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Queue name (the GCP resource name). Immutable (ForceNew) — and a
    # deleted queue's ID is reserved by the Cloud Tasks API for up to
    # 7 days, so treat names as long-lived.
    queue_name = string

    # Location (region). Immutable (ForceNew).
    location = string

    # Queue-level HTTP task settings: method/header/URI overrides and the
    # OAuth-XOR-OIDC authorization pair applied to every HTTP task.
    http_target = optional(object({
      http_method = optional(string, "")

      header_overrides = optional(list(object({
        key   = string
        value = string
      })), [])

      oauth_token = optional(object({
        # Resolved from a GcpServiceAccount reference to the SA email.
        service_account_email = string
        scope                 = optional(string, "")
      }), null)

      oidc_token = optional(object({
        # Resolved from a GcpServiceAccount reference to the SA email.
        service_account_email = string
        audience              = optional(string, "")
      }), null)

      # Flattened in the spec: path/query_params map onto the provider's
      # nested path_override/query_override blocks (see main.tf).
      uri_override = optional(object({
        scheme       = optional(string, "")
        host         = optional(string, "")
        port         = optional(string, "")
        path         = optional(string, "")
        query_params = optional(string, "")
        enforce_mode = optional(string, "")
      }), null)
    }), null)

    # Routing override for App Engine tasks only; ignored for HTTP tasks.
    app_engine_routing_override = optional(object({
      service  = optional(string, "")
      version  = optional(string, "")
      instance = optional(string, "")
    }), null)

    rate_limits = optional(object({
      max_dispatches_per_second = optional(number, 0)
      max_concurrent_dispatches = optional(number, 0)
    }), null)

    retry_config = optional(object({
      max_attempts       = optional(number, 0)
      max_retry_duration = optional(string, "")
      min_backoff        = optional(string, "")
      max_backoff        = optional(string, "")
      max_doublings      = optional(number, 0)
    }), null)

    stackdriver_logging_config = optional(object({
      sampling_ratio = number
    }), null)

    # RUNNING (default) dispatches tasks; PAUSED holds them in the queue.
    # Reconciled on every apply.
    desired_state = optional(string, "")

    # DELETE (default) removes the queue and its backlog on destroy;
    # PREVENT fails the destroy; ABANDON leaves the queue running in GCP.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.queue_name != ""
    error_message = "queue_name is required."
  }

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }
}
