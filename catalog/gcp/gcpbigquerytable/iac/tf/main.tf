# Enable the BigQuery API so a fresh project can host tables.
resource "google_project_service" "bigquery_api" {
  project = local.project_id
  service = "bigquery.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The BigQuery table — a native table, a logical view, a materialized view,
# or an external/BigLake table, all arms of one resource. table_id,
# dataset_id, and project are immutable, as are the encryption key, BigLake
# configuration, replication info, a materialized view's query, and each
# partitioning field.
resource "google_bigquery_table" "this" {
  table_id   = var.spec.table_id
  dataset_id = var.spec.dataset_id
  project    = local.project_id

  friendly_name   = local.friendly_name
  description     = local.description
  schema          = local.schema
  expiration_time = local.expiration_time
  max_staleness   = local.max_staleness

  # Up to four clustering columns in precedence order.
  clustering               = local.clustering
  require_partition_filter = var.spec.require_partition_filter

  labels        = local.final_labels
  resource_tags = local.resource_tags

  # Client-side destroy guard, default TRUE on both engines: a destroy
  # fails until the spec explicitly sets it false.
  deletion_protection = var.spec.deletion_protection

  # Checked BEFORE deletion_protection: PREVENT fails the destroy
  # outright, ABANDON drops the table from state (bypassing the guard),
  # DELETE (provider default) proceeds to the guard above.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # Hide diffs from columns BigQuery adds on its own.
  ignore_auto_generated_schema = var.spec.ignore_auto_generated_schema

  # Schema sub-fields the provider stops reconciling ("dataPolicies").
  ignore_schema_changes = length(var.spec.ignore_schema_changes) > 0 ? var.spec.ignore_schema_changes : null

  # Read-tuning: how much table metadata the provider requests back.
  table_metadata_view = var.spec.table_metadata_view != "" ? var.spec.table_metadata_view : null

  dynamic "time_partitioning" {
    for_each = var.spec.time_partitioning != null ? [var.spec.time_partitioning] : []
    content {
      type          = time_partitioning.value.type
      field         = time_partitioning.value.field != "" ? time_partitioning.value.field : null
      expiration_ms = time_partitioning.value.expiration_ms > 0 ? time_partitioning.value.expiration_ms : null
    }
  }

  dynamic "range_partitioning" {
    for_each = var.spec.range_partitioning != null ? [var.spec.range_partitioning] : []
    content {
      field = range_partitioning.value.field
      range {
        start    = range_partitioning.value.range.start
        end      = range_partitioning.value.range.end
        interval = range_partitioning.value.range.interval
      }
    }
  }

  dynamic "view" {
    for_each = var.spec.view != null ? [var.spec.view] : []
    content {
      query = view.value.query
      # Always sent explicitly so the API's own legacy-SQL-by-default
      # behavior for views never silently applies (spec default: false).
      use_legacy_sql = view.value.use_legacy_sql
    }
  }

  dynamic "materialized_view" {
    for_each = var.spec.materialized_view != null ? [var.spec.materialized_view] : []
    content {
      query = materialized_view.value.query
      # null when unset in the spec — the API's default (enabled) applies.
      enable_refresh                   = materialized_view.value.enable_refresh
      refresh_interval_ms              = materialized_view.value.refresh_interval_ms > 0 ? materialized_view.value.refresh_interval_ms : null
      allow_non_incremental_definition = materialized_view.value.allow_non_incremental_definition
    }
  }

  # CMEK. Immutable — changing the key recreates the table. The BigQuery
  # service agent must hold cryptoKeyEncrypterDecrypter on the key.
  dynamic "encryption_configuration" {
    for_each = local.kms_key_name != null ? [local.kms_key_name] : []
    content {
      kms_key_name = encryption_configuration.value
    }
  }

  dynamic "external_data_configuration" {
    for_each = var.spec.external_data_configuration != null ? [var.spec.external_data_configuration] : []
    content {
      # Always sent explicitly (the API requires the field to be set).
      autodetect                = external_data_configuration.value.autodetect
      source_uris               = external_data_configuration.value.source_uris
      source_format             = external_data_configuration.value.source_format != "" ? external_data_configuration.value.source_format : null
      object_metadata           = external_data_configuration.value.object_metadata != "" ? external_data_configuration.value.object_metadata : null
      compression               = external_data_configuration.value.compression != "" ? external_data_configuration.value.compression : null
      schema                    = external_data_configuration.value.schema != "" ? external_data_configuration.value.schema : null
      ignore_unknown_values     = external_data_configuration.value.ignore_unknown_values
      max_bad_records           = external_data_configuration.value.max_bad_records > 0 ? external_data_configuration.value.max_bad_records : null
      connection_id             = external_data_configuration.value.connection_id != "" ? external_data_configuration.value.connection_id : null
      reference_file_schema_uri = external_data_configuration.value.reference_file_schema_uri != "" ? external_data_configuration.value.reference_file_schema_uri : null
      metadata_cache_mode       = external_data_configuration.value.metadata_cache_mode != "" ? external_data_configuration.value.metadata_cache_mode : null
      file_set_spec_type        = external_data_configuration.value.file_set_spec_type != "" ? external_data_configuration.value.file_set_spec_type : null
      json_extension            = external_data_configuration.value.json_extension != "" ? external_data_configuration.value.json_extension : null
      decimal_target_types      = length(external_data_configuration.value.decimal_target_types) > 0 ? external_data_configuration.value.decimal_target_types : null

      dynamic "csv_options" {
        for_each = external_data_configuration.value.csv_options != null ? [external_data_configuration.value.csv_options] : []
        content {
          # The provider requires quote; an unset spec value means the API
          # default double-quote, while an explicit "" means unquoted data
          # (why the spec field is presence-tracked).
          quote                 = csv_options.value.quote != null ? csv_options.value.quote : "\""
          allow_jagged_rows     = csv_options.value.allow_jagged_rows
          allow_quoted_newlines = csv_options.value.allow_quoted_newlines
          encoding              = csv_options.value.encoding != "" ? csv_options.value.encoding : null
          field_delimiter       = csv_options.value.field_delimiter != "" ? csv_options.value.field_delimiter : null
          skip_leading_rows     = csv_options.value.skip_leading_rows > 0 ? csv_options.value.skip_leading_rows : null
          source_column_match   = csv_options.value.source_column_match != "" ? csv_options.value.source_column_match : null
        }
      }

      dynamic "json_options" {
        for_each = external_data_configuration.value.json_options != null ? [external_data_configuration.value.json_options] : []
        content {
          encoding = json_options.value.encoding != "" ? json_options.value.encoding : null
        }
      }

      dynamic "google_sheets_options" {
        for_each = external_data_configuration.value.google_sheets_options != null ? [external_data_configuration.value.google_sheets_options] : []
        content {
          range             = google_sheets_options.value.range != "" ? google_sheets_options.value.range : null
          skip_leading_rows = google_sheets_options.value.skip_leading_rows > 0 ? google_sheets_options.value.skip_leading_rows : null
        }
      }

      dynamic "hive_partitioning_options" {
        for_each = external_data_configuration.value.hive_partitioning_options != null ? [external_data_configuration.value.hive_partitioning_options] : []
        content {
          mode                     = hive_partitioning_options.value.mode != "" ? hive_partitioning_options.value.mode : null
          require_partition_filter = hive_partitioning_options.value.require_partition_filter
          source_uri_prefix        = hive_partitioning_options.value.source_uri_prefix != "" ? hive_partitioning_options.value.source_uri_prefix : null
        }
      }

      dynamic "avro_options" {
        for_each = external_data_configuration.value.avro_options != null ? [external_data_configuration.value.avro_options] : []
        content {
          use_avro_logical_types = avro_options.value.use_avro_logical_types
        }
      }

      dynamic "parquet_options" {
        for_each = external_data_configuration.value.parquet_options != null ? [external_data_configuration.value.parquet_options] : []
        content {
          enum_as_string        = parquet_options.value.enum_as_string
          enable_list_inference = parquet_options.value.enable_list_inference
        }
      }

      dynamic "bigtable_options" {
        for_each = external_data_configuration.value.bigtable_options != null ? [external_data_configuration.value.bigtable_options] : []
        content {
          ignore_unspecified_column_families = bigtable_options.value.ignore_unspecified_column_families
          read_rowkey_as_string              = bigtable_options.value.read_rowkey_as_string
          output_column_families_as_json     = bigtable_options.value.output_column_families_as_json

          dynamic "column_family" {
            for_each = bigtable_options.value.column_families
            content {
              family_id        = column_family.value.family_id != "" ? column_family.value.family_id : null
              type             = column_family.value.type != "" ? column_family.value.type : null
              encoding         = column_family.value.encoding != "" ? column_family.value.encoding : null
              only_read_latest = column_family.value.only_read_latest

              dynamic "column" {
                for_each = column_family.value.columns
                content {
                  qualifier_encoded = column.value.qualifier_encoded != "" ? column.value.qualifier_encoded : null
                  qualifier_string  = column.value.qualifier_string != "" ? column.value.qualifier_string : null
                  field_name        = column.value.field_name != "" ? column.value.field_name : null
                  type              = column.value.type != "" ? column.value.type : null
                  encoding          = column.value.encoding != "" ? column.value.encoding : null
                  only_read_latest  = column.value.only_read_latest
                }
              }
            }
          }
        }
      }
    }
  }

  # Unenforced primary/foreign keys for the optimizer and lineage tools.
  dynamic "table_constraints" {
    for_each = var.spec.table_constraints != null ? [var.spec.table_constraints] : []
    content {
      dynamic "primary_key" {
        for_each = table_constraints.value.primary_key != null ? [table_constraints.value.primary_key] : []
        content {
          columns = primary_key.value.columns
        }
      }

      dynamic "foreign_keys" {
        for_each = table_constraints.value.foreign_keys
        content {
          name = foreign_keys.value.name != "" ? foreign_keys.value.name : null
          referenced_table {
            project_id = foreign_keys.value.referenced_table.project_id
            dataset_id = foreign_keys.value.referenced_table.dataset_id
            table_id   = foreign_keys.value.referenced_table.table_id
          }
          column_references {
            referencing_column = foreign_keys.value.column_references.referencing_column
            referenced_column  = foreign_keys.value.column_references.referenced_column
          }
        }
      }
    }
  }

  # Cross-dataset/region replica of a source materialized view. Immutable.
  dynamic "table_replication_info" {
    for_each = var.spec.table_replication_info != null ? [var.spec.table_replication_info] : []
    content {
      source_project_id       = table_replication_info.value.source_project_id
      source_dataset_id       = table_replication_info.value.source_dataset_id
      source_table_id         = table_replication_info.value.source_table_id
      replication_interval_ms = table_replication_info.value.replication_interval_ms > 0 ? table_replication_info.value.replication_interval_ms : null
    }
  }

  # BigLake managed table: open-format (Iceberg) files in your own bucket,
  # managed by BigQuery. Immutable.
  dynamic "biglake_configuration" {
    for_each = var.spec.biglake_configuration != null ? [var.spec.biglake_configuration] : []
    content {
      connection_id = biglake_configuration.value.connection_id
      storage_uri   = biglake_configuration.value.storage_uri
      file_format   = biglake_configuration.value.file_format
      table_format  = biglake_configuration.value.table_format
    }
  }

  dynamic "schema_foreign_type_info" {
    for_each = var.spec.schema_foreign_type_info != null ? [var.spec.schema_foreign_type_info] : []
    content {
      type_system = schema_foreign_type_info.value.type_system
    }
  }

  # Hive Metastore compatibility metadata for open-source engines.
  dynamic "external_catalog_table_options" {
    for_each = var.spec.external_catalog_table_options != null ? [var.spec.external_catalog_table_options] : []
    content {
      parameters    = length(external_catalog_table_options.value.parameters) > 0 ? external_catalog_table_options.value.parameters : null
      connection_id = external_catalog_table_options.value.connection_id != "" ? external_catalog_table_options.value.connection_id : null

      dynamic "storage_descriptor" {
        for_each = external_catalog_table_options.value.storage_descriptor != null ? [external_catalog_table_options.value.storage_descriptor] : []
        content {
          location_uri  = storage_descriptor.value.location_uri != "" ? storage_descriptor.value.location_uri : null
          input_format  = storage_descriptor.value.input_format != "" ? storage_descriptor.value.input_format : null
          output_format = storage_descriptor.value.output_format != "" ? storage_descriptor.value.output_format : null

          dynamic "serde_info" {
            for_each = storage_descriptor.value.serde_info != null ? [storage_descriptor.value.serde_info] : []
            content {
              name                  = serde_info.value.name != "" ? serde_info.value.name : null
              serialization_library = serde_info.value.serialization_library
              parameters            = length(serde_info.value.parameters) > 0 ? serde_info.value.parameters : null
            }
          }
        }
      }
    }
  }

  depends_on = [google_project_service.bigquery_api]
}
