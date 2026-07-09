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
  description = "AwsFsxOntapVolume specification"
  type = object({
    region = string
    storage_virtual_machine_id = string
    name = string
    size_in_megabytes = optional(number)
    size_in_bytes = optional(number)
    junction_path = optional(string, "")
    ontap_volume_type = optional(string)
    volume_style = optional(string)
    security_style = optional(string, "")
    snapshot_policy = optional(string, "")
    storage_efficiency_enabled = optional(bool)
    copy_tags_to_backups = optional(bool)
    skip_final_backup = optional(bool)
    final_backup_tags = optional(map(string), {})
    bypass_snaplock_enterprise_retention = optional(bool)
    tiering_policy = optional(object({
      name = optional(string, "")
      cooling_period = optional(number, 0)
    }))
    snaplock_configuration = optional(object({
      snaplock_type = string
      audit_log_volume = optional(bool)
      privileged_delete = optional(string)
      volume_append_mode_enabled = optional(bool)
      autocommit_period = optional(object({
        type = optional(string, "")
        value = optional(number, 0)
      }))
      retention_period = optional(object({
        default_retention = optional(object({
          type = optional(string, "")
          value = optional(number, 0)
        }))
        minimum_retention = optional(object({
          type = optional(string, "")
          value = optional(number, 0)
        }))
        maximum_retention = optional(object({
          type = optional(string, "")
          value = optional(number, 0)
        }))
      }))
    }))
    aggregate_configuration = optional(object({
      aggregates = optional(list(string), [])
      constituents_per_aggregate = optional(number, 0)
    }))
  })
}
