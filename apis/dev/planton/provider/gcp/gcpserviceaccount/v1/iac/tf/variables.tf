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
  description = "Specification for the GCP Service Account"
  type = object({
    # The short account ID that forms the email:
    # <service_account_id>@<project>.iam.gserviceaccount.com.
    # 6-30 chars; lowercase letters, digits, hyphens; starts with a letter,
    # cannot end with a hyphen. Immutable in GCP (ForceNew).
    service_account_id = string

    # The GCP project ID in which the service account is created. The CLI's
    # tfvars converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Human-readable display name shown in the GCP console.
    # Falls back to metadata.name when omitted.
    display_name = optional(string)

    # Human-readable description of what this identity is for (max 256 bytes).
    description = optional(string)

    # Whether the account is disabled. Disabled accounts keep their IAM
    # bindings but cannot authenticate until re-enabled.
    disabled = optional(bool)

    # Whether to create a user-managed JSON key. Defaults to false (keyless).
    create_key = optional(bool)

    # IAM roles granted at the project level (additive member grants).
    # Example: ["roles/logging.logWriter", "roles/storage.objectViewer"]
    project_iam_roles = optional(list(string))

    # Numeric organization ID; required when org_iam_roles is set.
    org_id = optional(string)

    # IAM roles granted at the organization level (additive member grants).
    # Example: ["roles/resourcemanager.organizationViewer"]
    org_iam_roles = optional(list(string))
  })

  validation {
    condition     = can(regex("^[a-z]([-a-z0-9]*[a-z0-9])?$", var.spec.service_account_id)) && length(var.spec.service_account_id) >= 6 && length(var.spec.service_account_id) <= 30
    error_message = "service_account_id must be 6-30 characters, start with a lowercase letter, and contain only lowercase letters, digits, and hyphens (cannot end with a hyphen)."
  }

  validation {
    condition     = var.spec.description == null || length(coalesce(var.spec.description, "")) <= 256
    error_message = "description must be at most 256 bytes."
  }

  validation {
    condition     = length(coalesce(var.spec.org_iam_roles, [])) == 0 || (var.spec.org_id != null && var.spec.org_id != "")
    error_message = "org_id must be specified when org_iam_roles is not empty."
  }
}
