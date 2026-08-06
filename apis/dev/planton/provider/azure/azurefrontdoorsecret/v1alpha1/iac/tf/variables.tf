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
  description = "Azure Front Door secret specification"
  type = object({
    # The Front Door profile the secret lives in, by ARM ID. References
    # are resolved to a literal ID by the platform before the module
    # runs. ForceNew.
    profile_id = string

    # The secret's name -- unique within the profile. ForceNew.
    secret_name = string

    # The Key Vault CERTIFICATE the secret wraps (a vault data-plane
    # URL). A versionless id follows the certificate's latest version
    # (rotation propagates); a versioned id pins one version. ForceNew
    # -- the whole secret is immutable.
    key_vault_certificate_id = string
  })
}
