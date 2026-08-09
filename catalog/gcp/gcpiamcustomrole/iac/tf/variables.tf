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
  description = "Specification for the GCP IAM custom role"
  type = object({
    # The role identifier, unique within the project. Forms the full role
    # name projects/<project>/roles/<role_id>. 3-64 chars; letters, digits,
    # underscores, periods; NO hyphens. Immutable in GCP (ForceNew).
    role_id = string

    # The GCP project that owns this custom role. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Human-readable title shown in the GCP console (max 100 chars). Mutable.
    title = string

    # What this role is for and who should hold it (max 256 bytes). Mutable.
    description = optional(string)

    # The IAM permissions the role grants; at least one required. Mutable —
    # edits propagate immediately to every grant of the role.
    permissions = list(string)

    # Launch stage label: ALPHA, BETA, GA (default), DEPRECATED, DISABLED, EAP.
    stage = optional(string)

    # DELETE (default) soft-deletes the role on destroy; PREVENT fails the
    # destroy; ABANDON leaves the role active with every binding working.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = can(regex("^[a-zA-Z0-9_.]{3,64}$", var.spec.role_id))
    error_message = "role_id must be 3-64 characters of letters, digits, underscores, or periods (hyphens are not allowed)."
  }

  validation {
    condition     = length(var.spec.title) > 0 && length(var.spec.title) <= 100
    error_message = "title is required and must be at most 100 characters."
  }

  validation {
    condition     = length(var.spec.permissions) > 0
    error_message = "at least one permission is required."
  }

  validation {
    condition     = var.spec.stage == null || contains(["ALPHA", "BETA", "GA", "DEPRECATED", "DISABLED", "EAP"], coalesce(var.spec.stage, "GA"))
    error_message = "stage must be one of ALPHA, BETA, GA, DEPRECATED, DISABLED, or EAP."
  }
}
