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
  description = "Azure Storage Data Lake Gen2 Filesystem specification"
  type = object({
    # The storage account the filesystem lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs. The account should be HNS-enabled -- POSIX access control is
    # rejected on flat-namespace accounts.
    storage_account_id = string

    # The filesystem's name: 3-63 lowercase letters, digits, and
    # hyphens, not starting with a hyphen ($root is the one special
    # name); unique within the account.
    filesystem_name = string

    # The encryption scope applied to data that doesn't name its own.
    # References are resolved to a literal scope name by the platform
    # before the module runs. Fixed at creation.
    default_encryption_scope = optional(string)

    # The Entra object ID (or $superuser) owning the root path.
    owner = optional(string)

    # The Entra object ID (or $superuser) of the root path's owning
    # group.
    group = optional(string)

    # The root path's POSIX ACL. scope is the spec enum's name string
    # (ACCESS, DEFAULT -- unset means ACCESS); type likewise (USER,
    # GROUP, MASK, OTHER); object_id only on USER/GROUP entries;
    # permissions in the three-character rwx form.
    aces = optional(list(object({
      scope       = optional(string)
      type        = string
      object_id   = optional(string)
      permissions = string
    })), [])

    # Free-form properties stored on the filesystem; Azure requires the
    # VALUES to be base64-encoded.
    properties = optional(map(string), {})
  })
}
