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
  description = "Specification for the GCP regional network endpoint group"
  type = object({
    # The GCP project that owns the NEG. The CLI's tfvars converter resolves
    # StringValueOrRef fields to their literal string before the module runs,
    # so this arrives as a plain string. Empty falls back to the provider's
    # default project (see locals.tf).
    project_id = optional(string, "")

    # Name of the NEG in GCP (RFC1035). Empty defaults to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    network_endpoint_group_name = optional(string, "")

    # Region the NEG lives in — required.
    region = string

    # Endpoint type. Middleware applies the proto default (SERVERLESS); empty
    # falls through to the same GCP API default.
    network_endpoint_type = optional(string, "")

    description = optional(string, "")

    # network/subnetwork arrive as resolved self-links (or literals). Only the
    # non-serverless endpoint types use them.
    network    = optional(string, "")
    subnetwork = optional(string, "")

    # PSC target service and settings (PRIVATE_SERVICE_CONNECT / INTERNET).
    psc_target_service = optional(string, "")
    psc_data = optional(object({
      producer_port = optional(string, "")
    }))

    # Serverless targets — exactly one for a SERVERLESS NEG (enforced by the
    # spec's CEL before deploy). service/function arrive as resolved strings.
    cloud_run = optional(object({
      service  = optional(string, "")
      tag      = optional(string, "")
      url_mask = optional(string, "")
    }))
    cloud_function = optional(object({
      function = optional(string, "")
      url_mask = optional(string, "")
    }))
    app_engine = optional(object({
      service  = optional(string, "")
      version  = optional(string, "")
      url_mask = optional(string, "")
    }))

    # DELETE (default) / PREVENT / ABANDON; empty falls through to the
    # provider default (DELETE).
    deletion_policy = optional(string, "")
  })

  # NOTE: never guard optional strings with coalesce() here — HCL's coalesce
  # skips empty strings as well as nulls, so coalesce("", "") errors and the
  # validation fails on a legitimately-empty value.
  validation {
    condition     = try(var.spec.network_endpoint_group_name, "") == "" || can(regex("^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$", var.spec.network_endpoint_group_name))
    error_message = "network_endpoint_group_name must be RFC1035-compliant: 1-63 lowercase letters, digits, or hyphens."
  }

  validation {
    condition     = var.spec.region != ""
    error_message = "region is required — a regional NEG is always scoped to one region."
  }

  validation {
    condition     = contains(["", "SERVERLESS", "PRIVATE_SERVICE_CONNECT", "INTERNET_IP_PORT", "INTERNET_FQDN_PORT", "GCE_VM_IP_PORTMAP"], var.spec.network_endpoint_type)
    error_message = "network_endpoint_type must be one of SERVERLESS, PRIVATE_SERVICE_CONNECT, INTERNET_IP_PORT, INTERNET_FQDN_PORT, or GCE_VM_IP_PORTMAP."
  }
}
