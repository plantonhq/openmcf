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
  description = "Azure Front Door security policy (WAF association) specification"
  type = object({
    # The Front Door profile the security policy lives in, by ARM ID.
    # References are resolved to a literal ID by the platform before the
    # module runs. ForceNew.
    profile_id = string

    # The security policy's name -- unique within the profile; begins
    # and ends alphanumeric, letters/digits/hyphens only. ForceNew.
    security_policy_name = string

    # The Front Door WAF policy to enforce, by ARM ID. Its sku must
    # match the profile's sku (Azure rejects the pairing at deploy time
    # otherwise). ForceNew.
    firewall_policy_id = string

    # The hostnames the WAF protects: endpoint ARM ids (the generated
    # *.azurefd.net hostname) and/or custom-domain ARM ids -- Azure
    # accepts both types, and the list can mix them. 1-500 entries (a
    # STANDARD profile caps the list at 100; the cap rides the
    # profile's sku, checked at deploy time). Updatable in place.
    domain_ids = list(string)
  })
}
