variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name = string,
    id = optional(string),
    org = optional(string),
    env = optional(string),
    labels = optional(map(string)),
    tags = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "CloudflareDnsZoneSpec defines the configuration for creating a Cloudflare DNS Zone"
  type = object({
    # (Required) The fully qualified domain name of the DNS zone (e.g., "example.com")
    zone_name = string

    # (Required) The Cloudflare account identifier under which to create the zone
    account_id = string

    # (Optional) The subscription plan for the zone
    # Valid values: "free", "pro", "business", "enterprise"
    # Defaults to "free" if not specified
    plan = optional(string, "free")

    # (Optional) Indicates if the zone is created in a paused state
    # If true, the zone will be DNS-only with no proxy/CDN/WAF services
    # Defaults to false (active)
    paused = optional(bool, false)

    # (Optional) If true, new DNS records in this zone will default to being proxied (orange-cloud)
    # Defaults to false (grey-cloud)
    default_proxied = optional(bool, false)

    # (Optional) DNS records to create in this zone
    # This allows managing DNS records as part of the zone configuration
    records = optional(list(object({
      # (Required) The name of the DNS record (e.g., "www", "api", "@" for root)
      name = string

      # (Required) The type of DNS record
      # Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS", "CAA"
      type = string

      # (Required) The value/target of the DNS record
      value = string

      # (Optional) Whether the record is proxied through Cloudflare (orange cloud)
      # Only applicable to A, AAAA, and CNAME records
      # Defaults to false
      proxied = optional(bool, false)

      # (Optional) TTL for the DNS record in seconds
      # Set to 1 for automatic TTL (recommended for proxied records)
      # Valid values: 1 (auto), or 60-86400 seconds
      # Defaults to 1 (automatic)
      ttl = optional(number, 1)

      # (Optional) Priority for MX and SRV records
      # Lower values indicate higher priority
      # Required for MX records
      # Range: 0-65535
      priority = optional(number, 0)

      # (Optional) Comment/note for the DNS record
      # Maximum 100 characters
      comment = optional(string, "")
    })), [])
  })
}