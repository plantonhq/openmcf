variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string
    id      = optional(string)
    org     = optional(string)
    env     = optional(string)
    labels  = optional(map(string))
    tags    = optional(list(string))
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for Cloudflare R2 Bucket"
  type = object({
    # DNS-compatible bucket name (3-63 characters, lowercase alphanumeric + hyphens)
    bucket_name = string

    # Cloudflare account ID (32 hex characters)
    account_id = string

    # Primary region for the bucket (location hint)
    # 1=WNAM, 2=ENAM, 3=WEUR, 4=EEUR, 5=APAC, 6=OC
    location = number

    # Expose bucket via public URL (r2.dev domain)
    public_access = optional(bool, false)

    # Enable object versioning (Note: R2 does not support versioning)
    versioning_enabled = optional(bool, false)

    # Custom domain configuration for the bucket
    custom_domain = optional(object({
      # Whether to enable custom domain access for the bucket
      enabled = bool

      # The Cloudflare Zone ID (literal value or reference)
      zone_id = object({
        value      = optional(string)
        value_from = optional(object({
          kind       = optional(string)
          env        = optional(string)
          name       = string
          field_path = optional(string)
        }))
      })

      # The full domain name to use for accessing the bucket
      domain = string
    }))
  })
}
