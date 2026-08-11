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
  description = "Specification for the GCP Cloud Scheduler job"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}). If
    # project_id is empty, the provider's default project is used
    # (see locals.tf).
    project_id = optional(string, "")

    # Job name (the GCP resource name). Empty falls back to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    job_name = optional(string, "")

    # Region the job runs in. Immutable (ForceNew).
    location = string

    # Unix-cron schedule, interpreted in time_zone (default Etc/UTC).
    schedule = string

    time_zone        = optional(string, "")
    description      = optional(string, "")
    attempt_deadline = optional(string, "")
    paused           = optional(bool, false)

    # Exactly one of http_target / pubsub_target / app_engine_http_target
    # is set — enforced pre-deploy by the spec's CEL rule.
    http_target = optional(object({
      uri         = string
      http_method = optional(string, "")
      body        = optional(string, "")
      headers     = optional(map(string), {})

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
    }), null)

    pubsub_target = optional(object({
      # Resolved from a GcpPubSubTopic reference to the fully qualified
      # projects/{project}/topics/{name} path.
      topic_name = string
      data       = optional(string, "")
      attributes = optional(map(string), {})
    }), null)

    app_engine_http_target = optional(object({
      relative_uri = string
      http_method  = optional(string, "")
      body         = optional(string, "")
      headers      = optional(map(string), {})

      app_engine_routing = optional(object({
        service  = optional(string, "")
        version  = optional(string, "")
        instance = optional(string, "")
      }), null)
    }), null)

    retry_config = optional(object({
      retry_count          = optional(number, 0)
      max_retry_duration   = optional(string, "")
      min_backoff_duration = optional(string, "")
      max_backoff_duration = optional(string, "")
      max_doublings        = optional(number, 0)
    }), null)

    # DELETE (default) removes the job on destroy; PREVENT fails the
    # destroy; ABANDON leaves the job firing on schedule in GCP.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.schedule != ""
    error_message = "schedule is required."
  }
}
