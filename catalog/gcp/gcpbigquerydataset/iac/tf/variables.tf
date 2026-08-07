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
  description = "Specification for the GCP BigQuery dataset"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")

    dataset_id = string
    location   = string

    friendly_name = optional(string, "")
    description   = optional(string, "")
    labels        = optional(map(string), {})
    resource_tags = optional(map(string), {})

    default_table_expiration_ms     = optional(number, 0)
    default_partition_expiration_ms = optional(number, 0)
    max_time_travel_hours           = optional(number, 0)
    is_case_insensitive             = optional(bool, false)
    default_collation               = optional(string, "")
    storage_billing_model           = optional(string, "")
    delete_contents_on_destroy      = optional(bool, false)

    kms_key_name = optional(string, "")

    access = optional(list(object({
      role           = optional(string, "")
      user_by_email  = optional(string, "")
      group_by_email = optional(string, "")
      domain         = optional(string, "")
      special_group  = optional(string, "")
      iam_member     = optional(string, "")
      view = optional(object({
        project_id = string
        dataset_id = string
        table_id   = string
      }), null)
      routine = optional(object({
        project_id = string
        dataset_id = string
        routine_id = string
      }), null)
      dataset = optional(object({
        project_id   = string
        dataset_id   = string
        target_types = list(string)
      }), null)
      condition = optional(object({
        expression  = string
        title       = optional(string, "")
        description = optional(string, "")
        location    = optional(string, "")
      }), null)
    })), [])

    external_dataset_reference = optional(object({
      external_source = string
      connection      = string
    }), null)

    external_catalog_options = optional(object({
      default_storage_location_uri = optional(string, "")
      parameters                   = optional(map(string), {})
    }), null)
  })
}
