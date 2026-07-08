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
  description = "AwsGlueCatalogDatabase specification"
  type = object({
    region = string
    description = optional(string, "")
    catalog_id = optional(string, "")
    location_uri = optional(string, "")
    parameters = optional(map(string), {})
    create_table_default_permissions = optional(list(object({
      permissions = optional(list(string), [])
      principal = optional(string, "")
    })), [])
    target_database = optional(object({
      catalog_id = string
      database_name = string
      region = optional(string, "")
    }))
    federated_database = optional(object({
      identifier = optional(string, "")
      connection_name = optional(string, "")
    }))
  })
}