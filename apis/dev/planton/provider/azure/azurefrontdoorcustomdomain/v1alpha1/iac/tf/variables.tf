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
  description = "Azure Front Door custom domain specification"
  type = object({
    # The Front Door profile the domain lives in, by ARM ID. References
    # are resolved to a literal ID by the platform before the module
    # runs. ForceNew.
    profile_id = string

    # The ARM resource name (NOT the hostname) -- unique within the
    # profile; convention is the hostname with dots as hyphens.
    # ForceNew.
    domain_name = string

    # The hostname the domain serves (FQDN; wildcard first label needs
    # a customer certificate). ForceNew -- the hostname IS the domain's
    # identity.
    host_name = string

    # The Azure DNS zone hosting the hostname's records, by ARM ID.
    # Absent when DNS is hosted outside Azure DNS.
    dns_zone_id = optional(string)

    # TLS is always on for custom domains. Enum fields arrive as the
    # spec enum's FULL value names and are mapped in locals.
    tls = object({
      # MANAGED_CERTIFICATE (default when absent) / CUSTOMER_CERTIFICATE.
      certificate_type = optional(string)
      # The Front Door secret carrying the BYO certificate; required
      # with CUSTOMER_CERTIFICATE, absent with MANAGED (spec-enforced).
      secret_id = optional(string)
      # A hardened cipher policy; absent serves Azure's default suites.
      cipher_suite = optional(object({
        # TLS12_2022 / TLS12_2023 / CUSTOMIZED.
        type = string
        # Required with CUSTOMIZED (spec-enforced).
        custom_ciphers = optional(object({
          tls12 = list(string)
          # Empty means Azure's TLS 1.3 defaults; when set the spec
          # guarantees both mandatory suites.
          tls13 = optional(list(string))
        }))
      }))
    })
  })
}
