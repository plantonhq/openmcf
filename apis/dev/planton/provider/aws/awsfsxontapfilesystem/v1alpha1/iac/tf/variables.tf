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
  description = "AwsFsxOntapFileSystem specification"
  type = object({
    region = string
    deployment_type = optional(string)
    storage_capacity_gib = optional(number, 0)
    storage_type = optional(string)
    throughput_capacity = optional(number)
    throughput_capacity_per_ha_pair = optional(number)
    ha_pairs = optional(number)
    subnet_ids = list(string)
    preferred_subnet_id = optional(string, "")
    security_group_ids = optional(list(string), [])
    endpoint_ip_address_range = optional(string, "")
    route_table_ids = optional(list(string), [])
    kms_key_id = optional(string, "")
    fsx_admin_password = optional(string, "")
    disk_iops_configuration = optional(object({
      mode = optional(string)
      iops = optional(number, 0)
    }))
    automatic_backup_retention_days = optional(number)
    daily_automatic_backup_start_time = optional(string, "")
    weekly_maintenance_start_time = optional(string, "")
  })
}
