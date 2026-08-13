# One kind, 13 provider resources: Azure stores every dataset shape
# in the SAME factory-scoped dataset namespace
# ({factory_id}/datasets/{name}), so the spec's variant block selects
# which resource is created. Shared fields (name, factory, the linked
# service reference, description, annotations, parameters,
# additional_properties, folder) travel identically on every shape;
# each resource below adds only its variant's own arguments.
#
# Optional string arguments are sent only when non-empty: the
# provider's own schema defaults then apply (e.g. delimited text's
# "," column delimiter), and several arguments reject an explicit
# empty string outright.

# Azure Blob Storage files addressed by a flat path + filename pair.
resource "azurerm_data_factory_dataset_azure_blob" "main" {
  count = var.spec.azure_blob != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  path                     = var.spec.azure_blob.path != "" ? var.spec.azure_blob.path : null
  filename                 = var.spec.azure_blob.filename != "" ? var.spec.azure_blob.filename : null
  dynamic_path_enabled     = var.spec.azure_blob.dynamic_path_enabled
  dynamic_filename_enabled = var.spec.azure_blob.dynamic_filename_enabled

  dynamic "schema_column" {
    for_each = var.spec.azure_blob.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# An Azure SQL Database table -- the one variant that references its
# linked service by ARM ID (same factory enforced by Azure).
resource "azurerm_data_factory_dataset_azure_sql_table" "main" {
  count = var.spec.azure_sql_table != null ? 1 : 0

  name              = var.spec.name
  data_factory_id   = var.spec.data_factory_id
  linked_service_id = var.spec.azure_sql_table.linked_service_id

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  schema = var.spec.azure_sql_table.schema != "" ? var.spec.azure_sql_table.schema : null
  table  = var.spec.azure_sql_table.table != "" ? var.spec.azure_sql_table.table : null

  dynamic "schema_column" {
    for_each = var.spec.azure_sql_table.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# Opaque binary files -- exactly one location (enforced by the spec);
# the one format that can also live on an SFTP server.
resource "azurerm_data_factory_dataset_binary" "main" {
  count = var.spec.binary != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  # The binary format's HTTP location requires path and filename (the
  # spec enforces both non-empty).
  dynamic "http_server_location" {
    for_each = var.spec.binary.http_server_location != null ? [var.spec.binary.http_server_location] : []
    content {
      relative_url             = http_server_location.value.relative_url
      path                     = http_server_location.value.path
      dynamic_path_enabled     = http_server_location.value.dynamic_path_enabled
      filename                 = http_server_location.value.filename
      dynamic_filename_enabled = http_server_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_storage_location" {
    for_each = var.spec.binary.azure_blob_storage_location != null ? [var.spec.binary.azure_blob_storage_location] : []
    content {
      container                 = azure_blob_storage_location.value.container
      dynamic_container_enabled = azure_blob_storage_location.value.dynamic_container_enabled
      path                      = azure_blob_storage_location.value.path != "" ? azure_blob_storage_location.value.path : null
      dynamic_path_enabled      = azure_blob_storage_location.value.dynamic_path_enabled
      filename                  = azure_blob_storage_location.value.filename != "" ? azure_blob_storage_location.value.filename : null
      dynamic_filename_enabled  = azure_blob_storage_location.value.dynamic_filename_enabled
    }
  }

  dynamic "sftp_server_location" {
    for_each = var.spec.binary.sftp_server_location != null ? [var.spec.binary.sftp_server_location] : []
    content {
      path                     = sftp_server_location.value.path
      dynamic_path_enabled     = sftp_server_location.value.dynamic_path_enabled
      filename                 = sftp_server_location.value.filename
      dynamic_filename_enabled = sftp_server_location.value.dynamic_filename_enabled
    }
  }

  dynamic "compression" {
    for_each = var.spec.binary.compression != null ? [var.spec.binary.compression] : []
    content {
      type  = compression.value.type
      level = compression.value.level != "" ? compression.value.level : null
    }
  }
}

# An Azure Cosmos DB (SQL API) collection.
resource "azurerm_data_factory_dataset_cosmosdb_sqlapi" "main" {
  count = var.spec.cosmosdb_sqlapi != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  collection_name = var.spec.cosmosdb_sqlapi.collection_name != "" ? var.spec.cosmosdb_sqlapi.collection_name : null

  dynamic "schema_column" {
    for_each = var.spec.cosmosdb_sqlapi.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# Any other dataset type, as raw type-properties JSON -- the escape
# hatch for dataset types azurerm has no first-class resource for. The
# only form whose linked service reference can carry parameter values.
resource "azurerm_data_factory_custom_dataset" "main" {
  count = var.spec.custom != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  type                 = var.spec.custom.type
  type_properties_json = var.spec.custom.type_properties_json
  schema_json          = var.spec.custom.schema_json != "" ? var.spec.custom.schema_json : null

  linked_service {
    name       = var.spec.custom.linked_service.name
    parameters = length(var.spec.custom.linked_service.parameters) > 0 ? var.spec.custom.linked_service.parameters : null
  }
}

# Delimited text (CSV) files -- exactly one location (enforced by the
# spec). Omitted parse settings fall back to the provider's own
# defaults ("," column delimiter, '"' quote, "\" escape, first row
# not a header).
resource "azurerm_data_factory_dataset_delimited_text" "main" {
  count = var.spec.delimited_text != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  # The delimited text format's HTTP location requires path and
  # filename (the spec enforces both non-empty).
  dynamic "http_server_location" {
    for_each = var.spec.delimited_text.http_server_location != null ? [var.spec.delimited_text.http_server_location] : []
    content {
      relative_url             = http_server_location.value.relative_url
      path                     = http_server_location.value.path
      dynamic_path_enabled     = http_server_location.value.dynamic_path_enabled
      filename                 = http_server_location.value.filename
      dynamic_filename_enabled = http_server_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_storage_location" {
    for_each = var.spec.delimited_text.azure_blob_storage_location != null ? [var.spec.delimited_text.azure_blob_storage_location] : []
    content {
      container                 = azure_blob_storage_location.value.container
      dynamic_container_enabled = azure_blob_storage_location.value.dynamic_container_enabled
      path                      = azure_blob_storage_location.value.path != "" ? azure_blob_storage_location.value.path : null
      dynamic_path_enabled      = azure_blob_storage_location.value.dynamic_path_enabled
      filename                  = azure_blob_storage_location.value.filename != "" ? azure_blob_storage_location.value.filename : null
      dynamic_filename_enabled  = azure_blob_storage_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_fs_location" {
    for_each = var.spec.delimited_text.azure_blob_fs_location != null ? [var.spec.delimited_text.azure_blob_fs_location] : []
    content {
      file_system                 = azure_blob_fs_location.value.file_system != "" ? azure_blob_fs_location.value.file_system : null
      dynamic_file_system_enabled = azure_blob_fs_location.value.dynamic_file_system_enabled
      path                        = azure_blob_fs_location.value.path != "" ? azure_blob_fs_location.value.path : null
      dynamic_path_enabled        = azure_blob_fs_location.value.dynamic_path_enabled
      filename                    = azure_blob_fs_location.value.filename != "" ? azure_blob_fs_location.value.filename : null
      dynamic_filename_enabled    = azure_blob_fs_location.value.dynamic_filename_enabled
    }
  }

  column_delimiter    = var.spec.delimited_text.column_delimiter != "" ? var.spec.delimited_text.column_delimiter : null
  row_delimiter       = var.spec.delimited_text.row_delimiter != "" ? var.spec.delimited_text.row_delimiter : null
  quote_character     = var.spec.delimited_text.quote_character != "" ? var.spec.delimited_text.quote_character : null
  escape_character    = var.spec.delimited_text.escape_character != "" ? var.spec.delimited_text.escape_character : null
  encoding            = var.spec.delimited_text.encoding != "" ? var.spec.delimited_text.encoding : null
  first_row_as_header = coalesce(var.spec.delimited_text.first_row_as_header, false)
  null_value          = var.spec.delimited_text.null_value != "" ? var.spec.delimited_text.null_value : null
  compression_codec   = var.spec.delimited_text.compression_codec != "" ? var.spec.delimited_text.compression_codec : null
  compression_level   = var.spec.delimited_text.compression_level != "" ? var.spec.delimited_text.compression_level : null

  dynamic "schema_column" {
    for_each = var.spec.delimited_text.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# A file served by an HTTP endpoint (through a web linked service).
resource "azurerm_data_factory_dataset_http" "main" {
  count = var.spec.http != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  relative_url   = var.spec.http.relative_url != "" ? var.spec.http.relative_url : null
  request_body   = var.spec.http.request_body != "" ? var.spec.http.request_body : null
  request_method = var.spec.http.request_method != "" ? var.spec.http.request_method : null

  dynamic "schema_column" {
    for_each = var.spec.http.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# JSON files -- exactly one location; the JSON format requires path
# and filename in BOTH location shapes (enforced by the spec).
resource "azurerm_data_factory_dataset_json" "main" {
  count = var.spec.json != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  dynamic "http_server_location" {
    for_each = var.spec.json.http_server_location != null ? [var.spec.json.http_server_location] : []
    content {
      relative_url             = http_server_location.value.relative_url
      path                     = http_server_location.value.path
      dynamic_path_enabled     = http_server_location.value.dynamic_path_enabled
      filename                 = http_server_location.value.filename
      dynamic_filename_enabled = http_server_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_storage_location" {
    for_each = var.spec.json.azure_blob_storage_location != null ? [var.spec.json.azure_blob_storage_location] : []
    content {
      container                 = azure_blob_storage_location.value.container
      dynamic_container_enabled = azure_blob_storage_location.value.dynamic_container_enabled
      path                      = azure_blob_storage_location.value.path
      dynamic_path_enabled      = azure_blob_storage_location.value.dynamic_path_enabled
      filename                  = azure_blob_storage_location.value.filename
      dynamic_filename_enabled  = azure_blob_storage_location.value.dynamic_filename_enabled
    }
  }

  encoding = var.spec.json.encoding != "" ? var.spec.json.encoding : null

  dynamic "schema_column" {
    for_each = var.spec.json.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# A MySQL table.
resource "azurerm_data_factory_dataset_mysql" "main" {
  count = var.spec.mysql != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  table_name = var.spec.mysql.table_name != "" ? var.spec.mysql.table_name : null

  dynamic "schema_column" {
    for_each = var.spec.mysql.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# Parquet files -- exactly one location; Parquet's HTTP location is
# the one shape whose folder path may be omitted.
resource "azurerm_data_factory_dataset_parquet" "main" {
  count = var.spec.parquet != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  dynamic "http_server_location" {
    for_each = var.spec.parquet.http_server_location != null ? [var.spec.parquet.http_server_location] : []
    content {
      relative_url             = http_server_location.value.relative_url
      path                     = http_server_location.value.path != "" ? http_server_location.value.path : null
      dynamic_path_enabled     = http_server_location.value.dynamic_path_enabled
      filename                 = http_server_location.value.filename
      dynamic_filename_enabled = http_server_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_storage_location" {
    for_each = var.spec.parquet.azure_blob_storage_location != null ? [var.spec.parquet.azure_blob_storage_location] : []
    content {
      container                 = azure_blob_storage_location.value.container
      dynamic_container_enabled = azure_blob_storage_location.value.dynamic_container_enabled
      path                      = azure_blob_storage_location.value.path != "" ? azure_blob_storage_location.value.path : null
      dynamic_path_enabled      = azure_blob_storage_location.value.dynamic_path_enabled
      filename                  = azure_blob_storage_location.value.filename != "" ? azure_blob_storage_location.value.filename : null
      dynamic_filename_enabled  = azure_blob_storage_location.value.dynamic_filename_enabled
    }
  }

  dynamic "azure_blob_fs_location" {
    for_each = var.spec.parquet.azure_blob_fs_location != null ? [var.spec.parquet.azure_blob_fs_location] : []
    content {
      file_system                 = azure_blob_fs_location.value.file_system != "" ? azure_blob_fs_location.value.file_system : null
      dynamic_file_system_enabled = azure_blob_fs_location.value.dynamic_file_system_enabled
      path                        = azure_blob_fs_location.value.path != "" ? azure_blob_fs_location.value.path : null
      dynamic_path_enabled        = azure_blob_fs_location.value.dynamic_path_enabled
      filename                    = azure_blob_fs_location.value.filename != "" ? azure_blob_fs_location.value.filename : null
      dynamic_filename_enabled    = azure_blob_fs_location.value.dynamic_filename_enabled
    }
  }

  # compression_level exists in the provider's schema for this
  # resource but its create/update code never reads it (dead
  # argument) -- deliberately not modeled; recorded in
  # iac/provider-parity.yaml.
  compression_codec = var.spec.parquet.compression_codec != "" ? var.spec.parquet.compression_codec : null

  dynamic "schema_column" {
    for_each = var.spec.parquet.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# A PostgreSQL table.
resource "azurerm_data_factory_dataset_postgresql" "main" {
  count = var.spec.postgresql != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  table_name = var.spec.postgresql.table_name != "" ? var.spec.postgresql.table_name : null

  dynamic "schema_column" {
    for_each = var.spec.postgresql.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}

# A Snowflake table -- its columns use Snowflake's own type
# vocabulary with precision/scale instead of the interim-type column
# form.
resource "azurerm_data_factory_dataset_snowflake" "main" {
  count = var.spec.snowflake != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  table_name  = var.spec.snowflake.table_name != "" ? var.spec.snowflake.table_name : null
  schema_name = var.spec.snowflake.schema_name != "" ? var.spec.snowflake.schema_name : null

  dynamic "schema_column" {
    for_each = var.spec.snowflake.schema_column
    content {
      name      = schema_column.value.name
      type      = schema_column.value.type != "" ? schema_column.value.type : null
      precision = schema_column.value.precision
      scale     = schema_column.value.scale
    }
  }
}

# A SQL Server table.
resource "azurerm_data_factory_dataset_sql_server_table" "main" {
  count = var.spec.sql_server_table != null ? 1 : 0

  name                = var.spec.name
  data_factory_id     = var.spec.data_factory_id
  linked_service_name = var.spec.linked_service_name

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  folder                = var.spec.folder != "" ? var.spec.folder : null

  table_name = var.spec.sql_server_table.table_name != "" ? var.spec.sql_server_table.table_name : null

  dynamic "schema_column" {
    for_each = var.spec.sql_server_table.schema_column
    content {
      name        = schema_column.value.name
      type        = schema_column.value.type != "" ? schema_column.value.type : null
      description = schema_column.value.description != "" ? schema_column.value.description : null
    }
  }
}
