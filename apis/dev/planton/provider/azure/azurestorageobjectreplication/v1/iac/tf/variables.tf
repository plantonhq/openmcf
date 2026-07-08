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
  description = "Azure Storage Object Replication specification"
  type = object({
    # The account pair. References are resolved to literal ARM IDs by
    # the platform before the module runs. The source needs blob
    # versioning AND change feed; the destination needs versioning --
    # Azure rejects the policy at apply time otherwise.
    source_storage_account_id      = string
    destination_storage_account_id = string

    # Container-to-container mappings. Container references are resolved
    # to literal names by the platform before the module runs.
    rules = list(object({
      source_container_name      = string
      destination_container_name = string

      # OnlyNewObjects (the platform materializes this default),
      # Everything, or an RFC 3339 instant -- which existing blobs join
      # the copy.
      copy_blobs_created_after = optional(string)

      # INCLUDE-prefix filters: replicate only blobs whose names start
      # with one of these.
      prefix_match = optional(list(string), [])
    }))
  })
}
