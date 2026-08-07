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
  description = "Specification for the additive service-account-scoped IAM grant"
  type = object({
    # The service account whose IAM policy receives this grant, as its
    # fully-qualified resource name (projects/<project>/serviceAccounts/<email>).
    # The CLI's tfvars converter resolves StringValueOrRef fields to their
    # literal string before the module runs, so these arrive as plain strings.
    # There is no separate project input: the account's project is embedded
    # in the resource name.
    service_account_id = string

    # The role to grant on the service account: a predefined role (typically
    # roles/iam.workloadIdentityUser, roles/iam.serviceAccountTokenCreator, or
    # roles/iam.serviceAccountUser) or a custom role's fully-qualified name
    # ("projects/<project>/roles/<role_id>").
    role = string

    # The identity receiving the grant, in GCP IAM member format
    # (serviceAccount:<email>, principal://..., principalSet://...,
    # user:<email>, group:<email>, domain:<domain>).
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
    # The fully-qualified resource-name shape the service account IAM API
    # addresses; the account segment is the email or unique numeric ID.
    condition     = can(regex("^projects/[^/]+/serviceAccounts/[^/]+$", var.spec.service_account_id))
    error_message = "service_account_id must be a fully-qualified service account resource name (projects/<project>/serviceAccounts/<email>)."
  }

  validation {
    condition     = var.spec.role != ""
    error_message = "role is required."
  }

  validation {
    # Mirrors the IAM member format: either <type>:<value> or one of the bare
    # public literals. Deleted principals cannot be granted to.
    condition = (
      var.spec.member != "" &&
      !startswith(var.spec.member, "deleted:") &&
      (can(regex("^.+:.+$", var.spec.member)) || contains(["allUsers", "allAuthenticatedUsers"], var.spec.member))
    )
    error_message = "member must be in IAM member format (e.g. serviceAccount:<email>, principalSet://..., user:<email>, group:<email>) or one of allUsers / allAuthenticatedUsers; grants to deleted principals are not supported."
  }
}
