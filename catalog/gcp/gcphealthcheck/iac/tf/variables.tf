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
  description = "Specification for the GCP Compute Engine health check"
  type = object({
    # The GCP project that owns the health check. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the health check in GCP (RFC1035). Empty defaults to
    # metadata.name (see locals.tf). Immutable (ForceNew).
    health_check_name = optional(string, "")

    # Region for a REGIONAL health check; empty means GLOBAL. The scope
    # selects which provider resource is created (see main.tf). Immutable.
    region = optional(string, "")

    # What this health check probes and which backends rely on it.
    description = optional(string)

    # Probe cadence and verdict thresholds. Defaults (5s/5s/2/2) are applied
    # by Planton middleware from the spec's proto defaults; the coalesce in
    # locals.tf is only a safety net for direct tfvars invocations.
    check_interval_sec  = optional(number)
    timeout_sec         = optional(number)
    healthy_threshold   = optional(number)
    unhealthy_threshold = optional(number)

    # Export a log entry on every health status change.
    enable_logging = optional(bool, false)

    # GLOBAL checks only: probe from exactly 3 specific regions instead of
    # Google's default prober set.
    source_regions = optional(list(string), [])

    # What happens to the health check in GCP on destroy:
    # DELETE (provider default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")

    # The probe protocol — exactly one of these objects is set (enforced by
    # the proto oneof upstream; main.tf trusts that contract). Ports left
    # null fall through to the GCP API defaults (http/tcp 80, https/http2/ssl
    # 443) — the module never hardcodes them.
    http = optional(object({
      host               = optional(string)
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
      proxy_header       = optional(string)
      request_path       = optional(string)
      response           = optional(string)
    }))
    https = optional(object({
      host               = optional(string)
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
      proxy_header       = optional(string)
      request_path       = optional(string)
      response           = optional(string)
    }))
    http2 = optional(object({
      host               = optional(string)
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
      proxy_header       = optional(string)
      request_path       = optional(string)
      response           = optional(string)
    }))
    tcp = optional(object({
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
      proxy_header       = optional(string)
      request            = optional(string)
      response           = optional(string)
    }))
    ssl = optional(object({
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
      proxy_header       = optional(string)
      request            = optional(string)
      response           = optional(string)
    }))
    grpc = optional(object({
      grpc_service_name  = optional(string)
      port               = optional(number)
      port_name          = optional(string)
      port_specification = optional(string)
    }))
    grpc_tls = optional(object({
      grpc_service_name  = optional(string)
      port               = optional(number)
      port_specification = optional(string)
    }))
  })

  validation {
    condition = length([
      for p in [var.spec.http, var.spec.https, var.spec.http2, var.spec.tcp, var.spec.ssl, var.spec.grpc, var.spec.grpc_tls] : p if p != null
    ]) == 1
    error_message = "exactly one protocol block (http, https, http2, tcp, ssl, grpc, or grpc_tls) must be set."
  }

  validation {
    condition     = length(var.spec.source_regions) == 0 || length(var.spec.source_regions) == 3
    error_message = "source_regions must list exactly 3 GCP regions, or be omitted."
  }

  validation {
    condition     = length(var.spec.source_regions) == 0 || var.spec.region == ""
    error_message = "source_regions is only supported on global health checks — remove it or clear region."
  }

  # HCL's || does not short-circuit, so nullable optionals are guarded with
  # try() — a null timeout/interval must not crash validation.
  validation {
    condition     = try(coalesce(var.spec.timeout_sec, 5) <= coalesce(var.spec.check_interval_sec, 5), true)
    error_message = "timeout_sec must not exceed check_interval_sec — GCP rejects a probe timeout longer than the probe interval."
  }
}
