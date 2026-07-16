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
  description = "Azure Storage Local User specification"
  type = object({
    # The storage account the user lives on. References are resolved to
    # a literal ARM ID by the platform before the module runs. The
    # account needs sftp_enabled (and is_hns_enabled) for the user to
    # actually connect.
    storage_account_id = string

    # The user's name: 3-64 lowercase letters and digits; unique within
    # the account. The SFTP login is {account}.{user}.
    user_name = string

    # Whether the user authenticates with SSH public keys (requires
    # ssh_authorized_keys) and/or an Azure-generated password. At least
    # one must be on (enforced in the spec).
    ssh_key_enabled      = optional(bool)
    ssh_password_enabled = optional(bool)

    # The SSH public keys the user may authenticate with -- paired with
    # ssh_key_enabled (enforced in the spec).
    ssh_authorized_keys = optional(list(object({
      key         = string
      description = optional(string)
    })), [])

    # The directory an SFTP session lands in, as "{container}" or
    # "{container}/{path}".
    home_directory = optional(string)

    # Per-resource grants: service is the spec enum's name string (BLOB,
    # FILE); resource_name references are resolved to a literal
    # container/share name by the platform before the module runs; the
    # five booleans pick the granted operations.
    permission_scopes = optional(list(object({
      service       = string
      resource_name = string
      read          = optional(bool)
      write         = optional(bool)
      delete        = optional(bool)
      list          = optional(bool)
      create        = optional(bool)
    })), [])
  })
}
