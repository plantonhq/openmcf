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
  description = "Specification for the GCP Spanner instance"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")

    instance_name = optional(string, "")
    config        = string
    display_name  = string
    labels        = optional(map(string), {})

    num_nodes        = optional(number, 0)
    processing_units = optional(number, 0)

    autoscaling_config = optional(object({
      autoscaling_limits = object({
        min_nodes            = optional(number, 0)
        max_nodes            = optional(number, 0)
        min_processing_units = optional(number, 0)
        max_processing_units = optional(number, 0)
      })
      autoscaling_targets = optional(object({
        high_priority_cpu_utilization_percent = optional(number, 0)
        storage_utilization_percent           = optional(number, 0)
      }), null)
      asymmetric_autoscaling_options = optional(list(object({
        replica_location = string
        overrides = object({
          min_nodes = number
          max_nodes = number
        })
      })), [])
    }), null)

    instance_type                = optional(string, "")
    edition                      = optional(string, "")
    default_backup_schedule_type = optional(string, "")
    force_destroy                = optional(bool, false)
  })
}
