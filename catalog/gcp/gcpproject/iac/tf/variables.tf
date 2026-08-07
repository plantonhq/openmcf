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
  description = "Specification for the GCP project"
  type = object({
    project_id   = string
    display_name = optional(string, "")

    # "organization" or "folder"; empty creates a parentless project (only
    # possible outside an organization).
    parent_type = optional(string, "")
    parent_id   = optional(string, "")

    billing_account_id = optional(string, "")

    labels = optional(map(string), {})

    # Resource Manager tags (tagKeys/{id} -> tagValues/{id}), create-time
    # only.
    tags = optional(map(string), {})

    auto_create_network = optional(bool, false)

    enabled_apis = optional(list(string), [])

    # DELETE (default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })
}
