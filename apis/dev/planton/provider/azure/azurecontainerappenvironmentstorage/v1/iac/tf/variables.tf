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
  description = "Azure Container App Environment storage registration specification"
  type = object({
    # The Container App Environment ARM ID. ForceNew.
    container_app_environment_id = string

    # The registration name app/job volumes reference (max 32 lowercase
    # alphanumerics/hyphens, starts with a letter). ForceNew.
    storage_name = string

    # The Azure Files share name. References are resolved to a literal
    # name by the platform before the module runs. ForceNew.
    share_name = string

    # How workloads may use the share, as the spec enum's name string
    # (READ_ONLY / READ_WRITE). ForceNew.
    access_mode = string

    # SMB path: the storage account holding the share. Pairs with
    # access_key; mutually exclusive with nfs_server_url (spec-enforced).
    # ForceNew.
    account_name = optional(string)

    # SMB path: the account access key. The one field that updates in
    # place (key rotation).
    access_key = optional(string)

    # NFS path: the account's file endpoint
    # ({account}.file.core.windows.net). Requires a VNet-injected
    # environment. ForceNew.
    nfs_server_url = optional(string)
  })
}
