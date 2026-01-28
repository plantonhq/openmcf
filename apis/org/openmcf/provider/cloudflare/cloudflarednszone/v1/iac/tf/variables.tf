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
  })
}