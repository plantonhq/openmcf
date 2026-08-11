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
  description = "Specification for the GCP BigQuery table"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    project_id = optional(string, "")
    dataset_id = string
    table_id   = string

    friendly_name = optional(string, "")
    description   = optional(string, "")
    labels        = optional(map(string), {})
    resource_tags = optional(map(string), {})

    schema = optional(string, "")

    time_partitioning = optional(object({
      type          = string
      field         = optional(string, "")
      expiration_ms = optional(number, 0)
    }), null)

    range_partitioning = optional(object({
      field = string
      range = object({
        start    = optional(number, 0)
        end      = optional(number, 0)
        interval = optional(number, 0)
      })
    }), null)

    clustering               = optional(list(string), [])
    require_partition_filter = optional(bool, false)
    expiration_time          = optional(number, 0)
    max_staleness            = optional(string, "")

    kms_key_name = optional(string, "")

    view = optional(object({
      query          = string
      use_legacy_sql = optional(bool, false)
    }), null)

    materialized_view = optional(object({
      query                            = string
      enable_refresh                   = optional(bool)
      refresh_interval_ms              = optional(number, 0)
      allow_non_incremental_definition = optional(bool, false)
    }), null)

    external_data_configuration = optional(object({
      autodetect                = optional(bool, false)
      source_uris               = list(string)
      source_format             = optional(string, "")
      object_metadata           = optional(string, "")
      compression               = optional(string, "")
      schema                    = optional(string, "")
      ignore_unknown_values     = optional(bool, false)
      max_bad_records           = optional(number, 0)
      connection_id             = optional(string, "")
      reference_file_schema_uri = optional(string, "")
      metadata_cache_mode       = optional(string, "")
      file_set_spec_type        = optional(string, "")
      json_extension            = optional(string, "")
      decimal_target_types      = optional(list(string), [])

      csv_options = optional(object({
        quote                 = optional(string)
        allow_jagged_rows     = optional(bool, false)
        allow_quoted_newlines = optional(bool, false)
        encoding              = optional(string, "")
        field_delimiter       = optional(string, "")
        skip_leading_rows     = optional(number, 0)
        source_column_match   = optional(string, "")
      }), null)

      json_options = optional(object({
        encoding = optional(string, "")
      }), null)

      google_sheets_options = optional(object({
        range             = optional(string, "")
        skip_leading_rows = optional(number, 0)
      }), null)

      hive_partitioning_options = optional(object({
        mode                     = optional(string, "")
        require_partition_filter = optional(bool, false)
        source_uri_prefix        = optional(string, "")
      }), null)

      avro_options = optional(object({
        use_avro_logical_types = optional(bool, false)
      }), null)

      parquet_options = optional(object({
        enum_as_string        = optional(bool, false)
        enable_list_inference = optional(bool, false)
      }), null)

      bigtable_options = optional(object({
        ignore_unspecified_column_families = optional(bool, false)
        read_rowkey_as_string              = optional(bool, false)
        output_column_families_as_json     = optional(bool, false)
        column_families = optional(list(object({
          family_id        = optional(string, "")
          type             = optional(string, "")
          encoding         = optional(string, "")
          only_read_latest = optional(bool, false)
          columns = optional(list(object({
            qualifier_encoded = optional(string, "")
            qualifier_string  = optional(string, "")
            field_name        = optional(string, "")
            type              = optional(string, "")
            encoding          = optional(string, "")
            only_read_latest  = optional(bool, false)
          })), [])
        })), [])
      }), null)
    }), null)

    table_constraints = optional(object({
      primary_key = optional(object({
        columns = list(string)
      }), null)
      foreign_keys = optional(list(object({
        name = optional(string, "")
        referenced_table = object({
          project_id = string
          dataset_id = string
          table_id   = string
        })
        column_references = object({
          referencing_column = string
          referenced_column  = string
        })
      })), [])
    }), null)

    table_replication_info = optional(object({
      source_project_id       = string
      source_dataset_id       = string
      source_table_id         = string
      replication_interval_ms = optional(number, 0)
    }), null)

    biglake_configuration = optional(object({
      connection_id = string
      storage_uri   = string
      file_format   = string
      table_format  = string
    }), null)

    schema_foreign_type_info = optional(object({
      type_system = string
    }), null)

    external_catalog_table_options = optional(object({
      parameters    = optional(map(string), {})
      connection_id = optional(string, "")
      storage_descriptor = optional(object({
        location_uri  = optional(string, "")
        input_format  = optional(string, "")
        output_format = optional(string, "")
        serde_info = optional(object({
          name                  = optional(string, "")
          serialization_library = string
          parameters            = optional(map(string), {})
        }), null)
      }), null)
    }), null)

    deletion_protection = optional(bool, true)

    # DELETE (default) attempts table deletion on destroy (still gated by
    # deletion_protection); PREVENT fails the destroy outright; ABANDON
    # leaves the table serving, bypassing deletion_protection.
    deletion_policy = optional(string, "")

    # Hide diffs from columns BigQuery adds on its own.
    ignore_auto_generated_schema = optional(bool, false)

    # Schema sub-fields treated as non-authoritative per column; the
    # provider currently honors exactly "dataPolicies".
    ignore_schema_changes = optional(list(string), [])

    # How much table metadata the provider requests on reads:
    # BASIC, STORAGE_STATS (API default), or FULL.
    table_metadata_view = optional(string, "")
  })
}
