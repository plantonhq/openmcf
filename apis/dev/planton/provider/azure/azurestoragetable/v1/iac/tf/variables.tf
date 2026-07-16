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
  description = "Azure Storage Table specification"
  type = object({
    # The storage account the table lives in. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    storage_account_id = string

    # The table's name: 3-63 alphanumeric characters starting with a
    # letter; unique within the account.
    table_name = string

    # Stored access policies (signed identifiers) anchoring SAS tokens.
    acls = optional(list(object({
      id = string
      access_policies = optional(list(object({
        # Permission letters in Azure's strict order: r, a, u, d.
        permissions = string
        # RFC 3339 UTC validity window; both ends are REQUIRED on table
        # policies (Azure's contract, enforced in the spec).
        start  = string
        expiry = string
      })), [])
    })), [])
  })
}
