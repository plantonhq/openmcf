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
  description = "Azure Storage Container specification"
  type = object({
    # The storage account the container lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs.
    storage_account_id = string

    # The container's name: 3-63 lowercase letters, digits, and
    # hyphens; unique within the account.
    container_name = string

    # Anonymous read access, as the spec enum's name string (PRIVATE,
    # BLOB, CONTAINER). Unset means PRIVATE.
    container_access_type = optional(string)

    # The encryption scope applied to blobs that don't name their own.
    # References an AzureStorageEncryptionScope's name output, resolved
    # to a literal by the platform before the module runs; the scope
    # must live on the same account. Fixed at creation.
    default_encryption_scope = optional(string)

    # Whether blob writes may override the default scope (Azure's
    # default is true). Only meaningful with default_encryption_scope.
    encryption_scope_override_enabled = optional(bool)

    # Free-form metadata key/value pairs stored on the container.
    metadata = optional(map(string), {})
  })
}
