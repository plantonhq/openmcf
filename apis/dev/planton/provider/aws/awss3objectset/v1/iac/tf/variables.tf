variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsS3ObjectSet specification"
  type = object({
    region = string
    bucket = string
    objects = list(object({
      key = string
      content = optional(string, "")
      content_base64 = optional(string, "")
      content_type = optional(string)
      cache_control = optional(string, "")
      content_encoding = optional(string, "")
      content_disposition = optional(string, "")
      content_language = optional(string, "")
      metadata = optional(map(string), {})
      website_redirect = optional(string, "")
      storage_class = optional(string, "")
      server_side_encryption = optional(string, "")
      kms_key = optional(string, "")
      bucket_key_enabled = optional(bool)
      checksum_algorithm = optional(string, "")
      object_lock_mode = optional(string, "")
      object_lock_retain_until_date = optional(string, "")
      object_lock_legal_hold_status = optional(string, "")
      acl = optional(string, "")
      force_destroy = optional(bool, false)
      tags = optional(map(string), {})
    }))
    tags = optional(map(string), {})
  })
}
