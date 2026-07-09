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
  description = "AwsFsxLustreFileSystem specification"
  type = object({
    region = string
    deployment_type = optional(string)
    storage_capacity_gib = optional(number)
    storage_type = optional(string)
    per_unit_storage_throughput = optional(number)
    throughput_capacity = optional(number)
    data_compression_type = optional(string)
    file_system_type_version = optional(string, "")
    efa_enabled = optional(bool, false)
    drive_cache_type = optional(string, "")
    data_read_cache_configuration = optional(object({
      sizing_mode = string
      size_gib = optional(number)
    }))
    subnet_id = string
    security_group_ids = optional(list(string), [])
    kms_key_id = optional(string, "")
    backup_id = optional(string, "")
    import_path = optional(string, "")
    export_path = optional(string, "")
    auto_import_policy = optional(string, "")
    imported_file_chunk_size = optional(number)
    root_squash_configuration = optional(object({
      root_squash = optional(string, "")
      no_squash_nids = optional(list(string), [])
    }))
    log_configuration = optional(object({
      destination = optional(string, "")
      level = optional(string)
    }))
    metadata_configuration = optional(object({
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
