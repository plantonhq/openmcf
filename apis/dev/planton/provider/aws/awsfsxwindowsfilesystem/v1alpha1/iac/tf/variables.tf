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
  description = "AwsFsxWindowsFileSystem specification"
  type = object({
    region = string
    deployment_type = optional(string)
    storage_capacity_gib = optional(number)
    storage_type = optional(string)
    throughput_capacity = optional(number, 0)
    subnet_ids = list(string)
    preferred_subnet_id = optional(string, "")
    security_group_ids = optional(list(string), [])
    kms_key_id = optional(string, "")
    backup_id = optional(string, "")
    active_directory_id = optional(string, "")
    self_managed_active_directory = optional(object({
      domain_name = string
      dns_ips = list(string)
      username = optional(string, "")
      password = optional(string, "")
      domain_join_service_account_secret_arn = optional(string, "")
      file_system_administrators_group = optional(string)
      organizational_unit_distinguished_name = optional(string, "")
    }))
    aliases = optional(list(string), [])
    audit_log_configuration = optional(object({
      file_access_audit_log_level = optional(string)
      file_share_access_audit_log_level = optional(string)
      audit_log_destination = optional(string, "")
    }))
    disk_iops_configuration = optional(object({
      mode = optional(string)
      iops = optional(number)
    }))
    automatic_backup_retention_days = optional(number)
    daily_automatic_backup_start_time = optional(string, "")
    copy_tags_to_backups = optional(bool, false)
    skip_final_backup = optional(bool)
    final_backup_tags = optional(map(string), {})
    weekly_maintenance_start_time = optional(string, "")
  })
}
