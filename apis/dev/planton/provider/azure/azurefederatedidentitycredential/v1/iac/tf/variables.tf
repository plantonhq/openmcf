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
  description = "Azure Federated Identity Credential specification"
  type = object({
    # The credential's ARM resource name under the parent identity, unique
    # within it. Name it after the external workload it trusts.
    name = string

    # The full ARM ID of the user-assigned managed identity the credential
    # is written on. References are resolved to a literal ID by the platform
    # before the module runs.
    user_assigned_identity = string

    # The OIDC issuer URL the incoming token's `iss` claim must equal.
    issuer = string

    # The workload identifier the incoming token's `sub` claim must equal.
    subject = string

    # The audience the incoming token's `aud` claim must equal. Defaults to
    # the Azure AD token-exchange audience every standard client requests.
    audience = optional(string, "api://AzureADTokenExchange")
  })
}
