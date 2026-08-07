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
  description = "Specification for the GKE workload-identity IAM grant"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.

    # The project hosting the GKE cluster — and therefore the implicit
    # workload-identity pool <project>.svc.id.goog the principal lives in.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Email of the Google Service Account the KSA impersonates. The grant
    # attaches to THIS service account's IAM policy.
    service_account_email = string

    # Kubernetes namespace and ServiceAccount name that form the principal.
    ksa_namespace = string
    ksa_name      = string

    # Optional IAM Condition restricting when the grant applies. The
    # condition is part of the grant's identity — the same grant with and
    # without a condition are two independent grants.
    condition = optional(object({
      title       = string
      expression  = string
      description = optional(string)
    }))
  })
}
