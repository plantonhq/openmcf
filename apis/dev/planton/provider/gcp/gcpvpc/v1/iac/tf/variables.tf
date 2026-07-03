variable "metadata" {
  description = "Metadata for the GCP VPC resource"
  type = object({
    name = string
    id   = string
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the GCP VPC"
  type = object({
    # The GCP project that owns the network. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used.
    project_id              = optional(string, "")
    auto_create_subnetworks = optional(bool, false)
    routing_mode            = optional(number, 0) # 0=REGIONAL (default), 1=GLOBAL
    network_name            = string
  })
  validation {
    condition     = can(regex("^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$", var.spec.network_name))
    error_message = "Network name must be 1-63 characters, lowercase letters, numbers, or hyphens, starting with a letter and ending with a letter or number."
  }
}
