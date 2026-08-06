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
  description = "Azure Front Door profile specification"
  type = object({
    # The resource group the profile is created in. References are
    # resolved to a literal name by the platform before the module runs.
    resource_group = string

    # The profile's name -- unique within the resource group; the ARM
    # namespace every Front Door child resource nests under. ForceNew.
    profile_name = string

    # The pricing tier as the spec enum's value name (STANDARD/PREMIUM).
    # Absent means STANDARD (tfvars drops zero-valued proto fields).
    # ForceNew, and Azure refuses PREMIUM -> STANDARD outright.
    sku = optional(string)

    # Origin response timeout in seconds (16-240). Azure defaults to 120
    # when omitted.
    response_timeout_seconds = optional(number)

    # Managed identity for keyless Key Vault certificate access
    # (bring-your-own TLS on custom domains).
    identity = optional(object({
      # SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED
      # (spec enum value names).
      type = string
      # AzureUserAssignedIdentity ARM ids -- resolved to literals by the
      # platform.
      user_assigned_identity_ids = optional(list(string))
    }))

    # Request parts to scrub from access logs, as spec enum value names
    # (QUERY_STRING_ARG_NAMES / REQUEST_IP_ADDRESS / REQUEST_URI).
    # Presence of at least one entry enables scrubbing.
    log_scrubbing_variables = optional(list(string))

    # Free-form user tags merged over the metadata-derived tags (user
    # tags win on collision).
    tags = optional(map(string), {})
  })
}
