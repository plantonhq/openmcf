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
  description = "Azure Storage Queue specification"
  type = object({
    # The storage account the queue lives in. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    storage_account_id = string

    # The queue's name: 3-63 lowercase letters, digits, and hyphens;
    # unique within the account.
    queue_name = string

    # Free-form metadata key/value pairs stored on the queue.
    metadata = optional(map(string), {})
  })
}
