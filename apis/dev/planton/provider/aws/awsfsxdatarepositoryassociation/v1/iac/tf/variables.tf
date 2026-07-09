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
  description = "AwsFsxDataRepositoryAssociation specification"
  type = object({
    region = string
    file_system_id = string
    file_system_path = string
    data_repository_path = string
    auto_import_events = optional(list(string), [])
    auto_export_events = optional(list(string), [])
    imported_file_chunk_size = optional(number)
    batch_import_meta_data_on_create = optional(bool, false)
    delete_data_in_filesystem = optional(bool, false)
  })
}
