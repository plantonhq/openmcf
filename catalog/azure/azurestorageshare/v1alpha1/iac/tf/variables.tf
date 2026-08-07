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
  description = "Azure Storage Share specification"
  type = object({
    # The storage account the share lives in. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    storage_account_id = string

    # The share's name: 3-63 lowercase letters, digits, and hyphens;
    # unique within the account.
    share_name = string

    # The share's maximum size in gigabytes (the provisioned quota).
    # Standard accounts allow 1-5120 (102400 with large-file-share);
    # premium FileStorage accounts require at least 100.
    quota_gb = number

    # The file-sharing protocol, as the spec enum's name string (SMB,
    # NFS). Unset means SMB. NFS requires a premium FileStorage account.
    enabled_protocol = optional(string)

    # The performance/billing tier, as the spec enum's name string
    # (TRANSACTION_OPTIMIZED, HOT, COOL, PREMIUM). Unset lets Azure pick
    # its default for the account kind.
    access_tier = optional(string)

    # Stored access policies (signed identifiers) anchoring SAS tokens.
    acls = optional(list(object({
      id = string
      access_policies = optional(list(object({
        # Permission letters in Azure's strict order: r, w, d, l.
        permissions = string
        # RFC 3339 UTC validity window; both ends optional on shares
        # (the SAS token may carry them instead).
        start  = optional(string)
        expiry = optional(string)
      })), [])
    })), [])

    # Free-form metadata key/value pairs stored on the share.
    metadata = optional(map(string), {})
  })
}
