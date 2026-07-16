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
  description = "Specification for the additive project-level IAM grant"
  type = object({
    # The GCP project whose IAM policy receives this grant. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string before
    # the module runs, so these arrive as plain strings.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # The role to grant: a predefined role ("roles/storage.objectViewer") or a
    # custom role's fully-qualified name ("projects/<project>/roles/<role_id>").
    role = string

    # The identity receiving the grant, in GCP IAM member format
    # (serviceAccount:<email>, user:<email>, group:<email>, domain:<domain>,
    # principal://..., principalSet://..., allUsers, allAuthenticatedUsers).
    member = string

    # Optional IAM Condition restricting when the grant applies. The condition
    # is part of the grant's identity — same role with and without a condition
    # are two independent grants.
    condition = optional(object({
      title       = string
      expression  = string
      description = optional(string)
    }))
  })

  validation {
    condition     = var.spec.role != ""
    error_message = "role is required."
  }

  validation {
    # Mirrors the IAM member format: either <type>:<value> or one of the bare
    # public/legacy literals. Deleted principals cannot be granted to.
    condition = (
      var.spec.member != "" &&
      !startswith(var.spec.member, "deleted:") &&
      (can(regex("^.+:.+$", var.spec.member)) || contains(["allUsers", "allAuthenticatedUsers", "projectOwners", "projectReaders", "projectWriters"], var.spec.member))
    )
    error_message = "member must be in IAM member format (e.g. serviceAccount:<email>, user:<email>, group:<email>) or one of allUsers / allAuthenticatedUsers; grants to deleted principals are not supported."
  }
}
