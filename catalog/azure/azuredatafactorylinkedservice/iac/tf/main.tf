# One kind, 23 provider resources: Azure stores every connection type
# in the SAME factory-scoped linked-service namespace
# ({factory_id}/linkedservices/{name}), so the spec's variant block
# selects which resource is created. Shared fields (name, factory,
# description, annotations, parameters, additional_properties, the
# integration runtime) travel identically on every type; each
# resource below adds only its variant's own arguments.

# Azure Blob Storage -- exactly one of the four connection forms
# (enforced by the spec): connection_string, connection_string_insecure,
# sas_uri, or service_endpoint (+ managed identity / service principal).
resource "azurerm_data_factory_linked_service_azure_blob_storage" "main" {
  count = var.spec.azure_blob_storage != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string          = var.spec.azure_blob_storage.connection_string != "" ? var.spec.azure_blob_storage.connection_string : null
  connection_string_insecure = var.spec.azure_blob_storage.connection_string_insecure != "" ? var.spec.azure_blob_storage.connection_string_insecure : null
  sas_uri                    = var.spec.azure_blob_storage.sas_uri != "" ? var.spec.azure_blob_storage.sas_uri : null
  service_endpoint           = var.spec.azure_blob_storage.service_endpoint != "" ? var.spec.azure_blob_storage.service_endpoint : null

  dynamic "sas_token_linked_key_vault_key" {
    for_each = var.spec.azure_blob_storage.sas_token_linked_key_vault_key != null ? [var.spec.azure_blob_storage.sas_token_linked_key_vault_key] : []
    content {
      linked_service_name = sas_token_linked_key_vault_key.value.linked_service_name
      secret_name         = sas_token_linked_key_vault_key.value.secret_name
    }
  }

  dynamic "service_principal_linked_key_vault_key" {
    for_each = var.spec.azure_blob_storage.service_principal_linked_key_vault_key != null ? [var.spec.azure_blob_storage.service_principal_linked_key_vault_key] : []
    content {
      linked_service_name = service_principal_linked_key_vault_key.value.linked_service_name
      secret_name         = service_principal_linked_key_vault_key.value.secret_name
    }
  }

  storage_kind = var.spec.azure_blob_storage.storage_kind != "" ? var.spec.azure_blob_storage.storage_kind : null

  # Sent only when true -- the provider's conflict check fires on the
  # argument's PRESENCE, so an explicit false alongside a service
  # principal is rejected; unset means false anyway (the provider's
  # default).
  use_managed_identity  = coalesce(var.spec.azure_blob_storage.use_managed_identity, false) ? true : null
  service_principal_id  = var.spec.azure_blob_storage.service_principal_id != "" ? var.spec.azure_blob_storage.service_principal_id : null
  service_principal_key = var.spec.azure_blob_storage.service_principal_key != "" ? var.spec.azure_blob_storage.service_principal_key : null
  tenant_id             = var.spec.azure_blob_storage.tenant_id != "" ? var.spec.azure_blob_storage.tenant_id : null
}

# Azure Databricks -- one authentication method and one cluster choice
# (both exactly-one, enforced by the spec).
resource "azurerm_data_factory_linked_service_azure_databricks" "main" {
  count = var.spec.azure_databricks != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  adb_domain = var.spec.azure_databricks.adb_domain

  msi_workspace_id = var.spec.azure_databricks.msi_workspace_id != "" ? var.spec.azure_databricks.msi_workspace_id : null
  access_token     = var.spec.azure_databricks.access_token != "" ? var.spec.azure_databricks.access_token : null

  dynamic "key_vault_password" {
    for_each = var.spec.azure_databricks.key_vault_password != null ? [var.spec.azure_databricks.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }

  existing_cluster_id = var.spec.azure_databricks.existing_cluster_id != "" ? var.spec.azure_databricks.existing_cluster_id : null

  dynamic "new_cluster_config" {
    for_each = var.spec.azure_databricks.new_cluster_config != null ? [var.spec.azure_databricks.new_cluster_config] : []
    content {
      node_type       = new_cluster_config.value.node_type
      cluster_version = new_cluster_config.value.cluster_version
      # The platform default, sent explicitly (mirrors the provider's
      # own schema default); max 0 means "fixed size" and is omitted.
      min_number_of_workers       = coalesce(new_cluster_config.value.min_number_of_workers, 1)
      max_number_of_workers       = new_cluster_config.value.max_number_of_workers > 0 ? new_cluster_config.value.max_number_of_workers : null
      driver_node_type            = new_cluster_config.value.driver_node_type != "" ? new_cluster_config.value.driver_node_type : null
      log_destination             = new_cluster_config.value.log_destination != "" ? new_cluster_config.value.log_destination : null
      spark_config                = length(new_cluster_config.value.spark_config) > 0 ? new_cluster_config.value.spark_config : null
      spark_environment_variables = length(new_cluster_config.value.spark_environment_variables) > 0 ? new_cluster_config.value.spark_environment_variables : null
      custom_tags                 = length(new_cluster_config.value.custom_tags) > 0 ? new_cluster_config.value.custom_tags : null
      init_scripts                = length(new_cluster_config.value.init_scripts) > 0 ? new_cluster_config.value.init_scripts : null
    }
  }

  dynamic "instance_pool" {
    for_each = var.spec.azure_databricks.instance_pool != null ? [var.spec.azure_databricks.instance_pool] : []
    content {
      instance_pool_id      = instance_pool.value.instance_pool_id
      cluster_version       = instance_pool.value.cluster_version
      min_number_of_workers = coalesce(instance_pool.value.min_number_of_workers, 1)
      max_number_of_workers = instance_pool.value.max_number_of_workers > 0 ? instance_pool.value.max_number_of_workers : null
    }
  }
}

# Azure Files -- connection string, optionally with host-addressed
# access fields.
resource "azurerm_data_factory_linked_service_azure_file_storage" "main" {
  count = var.spec.azure_file_storage != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.azure_file_storage.connection_string
  file_share        = var.spec.azure_file_storage.file_share != "" ? var.spec.azure_file_storage.file_share : null
  host              = var.spec.azure_file_storage.host != "" ? var.spec.azure_file_storage.host : null
  user_id           = var.spec.azure_file_storage.user_id != "" ? var.spec.azure_file_storage.user_id : null
  password          = var.spec.azure_file_storage.password != "" ? var.spec.azure_file_storage.password : null

  dynamic "key_vault_password" {
    for_each = var.spec.azure_file_storage.key_vault_password != null ? [var.spec.azure_file_storage.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }
}

# Azure Function -- the host key inline or held in Key Vault
# (exactly one, enforced by the spec).
resource "azurerm_data_factory_linked_service_azure_function" "main" {
  count = var.spec.azure_function != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  url = var.spec.azure_function.url
  key = var.spec.azure_function.key != "" ? var.spec.azure_function.key : null

  dynamic "key_vault_key" {
    for_each = var.spec.azure_function.key_vault_key != null ? [var.spec.azure_function.key_vault_key] : []
    content {
      linked_service_name = key_vault_key.value.linked_service_name
      secret_name         = key_vault_key.value.secret_name
    }
  }
}

# Azure AI Search -- endpoint URL + admin key.
resource "azurerm_data_factory_linked_service_azure_search" "main" {
  count = var.spec.azure_search != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  url                = var.spec.azure_search.url
  search_service_key = var.spec.azure_search.search_service_key
}

# Azure SQL Database -- connection string (inline or Key-Vault-held)
# plus managed identity / service principal / credential identity.
resource "azurerm_data_factory_linked_service_azure_sql_database" "main" {
  count = var.spec.azure_sql_database != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.azure_sql_database.connection_string != "" ? var.spec.azure_sql_database.connection_string : null

  dynamic "key_vault_connection_string" {
    for_each = var.spec.azure_sql_database.key_vault_connection_string != null ? [var.spec.azure_sql_database.key_vault_connection_string] : []
    content {
      linked_service_name = key_vault_connection_string.value.linked_service_name
      secret_name         = key_vault_connection_string.value.secret_name
    }
  }

  dynamic "key_vault_password" {
    for_each = var.spec.azure_sql_database.key_vault_password != null ? [var.spec.azure_sql_database.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }

  # Sent only when true -- the provider's conflict check fires on the
  # argument's PRESENCE (see the blob storage note above).
  use_managed_identity  = coalesce(var.spec.azure_sql_database.use_managed_identity, false) ? true : null
  service_principal_id  = var.spec.azure_sql_database.service_principal_id != "" ? var.spec.azure_sql_database.service_principal_id : null
  service_principal_key = var.spec.azure_sql_database.service_principal_key != "" ? var.spec.azure_sql_database.service_principal_key : null
  tenant_id             = var.spec.azure_sql_database.tenant_id != "" ? var.spec.azure_sql_database.tenant_id : null
  credential_name       = var.spec.azure_sql_database.credential_name != "" ? var.spec.azure_sql_database.credential_name : null
}

# Azure Table Storage -- connection string only.
resource "azurerm_data_factory_linked_service_azure_table_storage" "main" {
  count = var.spec.azure_table_storage != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.azure_table_storage.connection_string
}

# Azure Cosmos DB (SQL API) -- connection string OR account detail
# (endpoint + key + database), exactly one (enforced by the spec).
resource "azurerm_data_factory_linked_service_cosmosdb" "main" {
  count = var.spec.cosmosdb != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.cosmosdb.connection_string != "" ? var.spec.cosmosdb.connection_string : null
  account_endpoint  = var.spec.cosmosdb.account_endpoint != "" ? var.spec.cosmosdb.account_endpoint : null
  account_key       = var.spec.cosmosdb.account_key != "" ? var.spec.cosmosdb.account_key : null
  database          = var.spec.cosmosdb.database != "" ? var.spec.cosmosdb.database : null
}

# Azure Cosmos DB for MongoDB -- connection string only.
resource "azurerm_data_factory_linked_service_cosmosdb_mongoapi" "main" {
  count = var.spec.cosmosdb_mongoapi != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.cosmosdb_mongoapi.connection_string
  database          = var.spec.cosmosdb_mongoapi.database != "" ? var.spec.cosmosdb_mongoapi.database : null
  # The platform default, sent explicitly (mirrors the provider's own
  # schema default).
  server_version_is_32_or_higher = coalesce(var.spec.cosmosdb_mongoapi.server_version_is_32_or_higher, false)
}

# Any other connector type, as raw type-properties JSON. This is the
# one resource whose integration runtime travels as a block (name +
# per-use parameters) instead of a plain name.
resource "azurerm_data_factory_linked_custom_service" "main" {
  count = var.spec.custom != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description           = var.spec.description != "" ? var.spec.description : null
  annotations           = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null

  type                 = var.spec.custom.type
  type_properties_json = var.spec.custom.type_properties_json

  dynamic "integration_runtime" {
    for_each = var.spec.integration_runtime_name != "" ? [var.spec.integration_runtime_name] : []
    content {
      name       = integration_runtime.value
      parameters = length(var.spec.custom.integration_runtime_parameters) > 0 ? var.spec.custom.integration_runtime_parameters : null
    }
  }
}

# Azure Data Lake Storage Gen2 -- exactly one authentication mode
# (managed identity, account key, or service principal; enforced by
# the spec).
resource "azurerm_data_factory_linked_service_data_lake_storage_gen2" "main" {
  count = var.spec.data_lake_storage_gen2 != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  url = var.spec.data_lake_storage_gen2.url

  # Sent only when true -- the provider's AtLeastOneOf group reads an
  # explicit false as "this mode declared", which would collide with
  # the mode actually chosen.
  use_managed_identity  = coalesce(var.spec.data_lake_storage_gen2.use_managed_identity, false) ? true : null
  storage_account_key   = var.spec.data_lake_storage_gen2.storage_account_key != "" ? var.spec.data_lake_storage_gen2.storage_account_key : null
  service_principal_id  = var.spec.data_lake_storage_gen2.service_principal_id != "" ? var.spec.data_lake_storage_gen2.service_principal_id : null
  service_principal_key = var.spec.data_lake_storage_gen2.service_principal_key != "" ? var.spec.data_lake_storage_gen2.service_principal_key : null
  tenant                = var.spec.data_lake_storage_gen2.tenant != "" ? var.spec.data_lake_storage_gen2.tenant : null
}

# Azure Key Vault -- the vault other linked services' Key-Vault-held
# secrets resolve through (the provider derives the vault's base URI
# from the ARM ID).
resource "azurerm_data_factory_linked_service_key_vault" "main" {
  count = var.spec.key_vault != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  key_vault_id = var.spec.key_vault.key_vault_id
}

# Azure Data Explorer (Kusto) -- exactly one authentication mode
# (managed identity or service principal; enforced by the spec).
resource "azurerm_data_factory_linked_service_kusto" "main" {
  count = var.spec.kusto != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  kusto_endpoint      = var.spec.kusto.kusto_endpoint
  kusto_database_name = var.spec.kusto.kusto_database_name

  # Sent only when true -- the provider's ExactlyOneOf pair reads an
  # explicit false as "this mode declared", which would collide with
  # the service principal actually chosen.
  use_managed_identity  = coalesce(var.spec.kusto.use_managed_identity, false) ? true : null
  service_principal_id  = var.spec.kusto.service_principal_id != "" ? var.spec.kusto.service_principal_id : null
  service_principal_key = var.spec.kusto.service_principal_key != "" ? var.spec.kusto.service_principal_key : null
  tenant                = var.spec.kusto.tenant != "" ? var.spec.kusto.tenant : null
}

# MySQL -- connection string + driver line.
resource "azurerm_data_factory_linked_service_mysql" "main" {
  count = var.spec.mysql != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.mysql.connection_string
  # The platform default, sent explicitly (mirrors the provider's own
  # v5 schema default; the provider rejects V1 on new connections).
  driver_version = coalesce(var.spec.mysql.driver_version, "V2")
}

# OData -- anonymous, or Basic via the credentials block.
resource "azurerm_data_factory_linked_service_odata" "main" {
  count = var.spec.odata != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  url = var.spec.odata.url

  dynamic "basic_authentication" {
    for_each = var.spec.odata.basic_authentication != null ? [var.spec.odata.basic_authentication] : []
    content {
      username = basic_authentication.value.username
      password = basic_authentication.value.password
    }
  }
}

# ODBC -- anonymous, or Basic via the credentials block.
resource "azurerm_data_factory_linked_service_odbc" "main" {
  count = var.spec.odbc != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.odbc.connection_string

  dynamic "basic_authentication" {
    for_each = var.spec.odbc.basic_authentication != null ? [var.spec.odbc.basic_authentication] : []
    content {
      username = basic_authentication.value.username
      password = basic_authentication.value.password
    }
  }
}

# PostgreSQL -- connection string only.
resource "azurerm_data_factory_linked_service_postgresql" "main" {
  count = var.spec.postgresql != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.postgresql.connection_string
}

# SFTP -- credential class follows authentication_type (enforced by
# the spec's pairing rules, mirroring the provider's own).
resource "azurerm_data_factory_linked_service_sftp" "main" {
  count = var.spec.sftp != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  authentication_type = var.spec.sftp.authentication_type
  host                = var.spec.sftp.host
  port                = var.spec.sftp.port
  username            = var.spec.sftp.username

  password = var.spec.sftp.password != "" ? var.spec.sftp.password : null

  dynamic "key_vault_password" {
    for_each = var.spec.sftp.key_vault_password != null ? [var.spec.sftp.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }

  private_key_content_base64 = var.spec.sftp.private_key_content_base64 != "" ? var.spec.sftp.private_key_content_base64 : null

  dynamic "key_vault_private_key_content_base64" {
    for_each = var.spec.sftp.key_vault_private_key_content_base64 != null ? [var.spec.sftp.key_vault_private_key_content_base64] : []
    content {
      linked_service_name = key_vault_private_key_content_base64.value.linked_service_name
      secret_name         = key_vault_private_key_content_base64.value.secret_name
    }
  }

  private_key_path       = var.spec.sftp.private_key_path != "" ? var.spec.sftp.private_key_path : null
  private_key_passphrase = var.spec.sftp.private_key_passphrase != "" ? var.spec.sftp.private_key_passphrase : null

  dynamic "key_vault_private_key_passphrase" {
    for_each = var.spec.sftp.key_vault_private_key_passphrase != null ? [var.spec.sftp.key_vault_private_key_passphrase] : []
    content {
      linked_service_name = key_vault_private_key_passphrase.value.linked_service_name
      secret_name         = key_vault_private_key_passphrase.value.secret_name
    }
  }

  # Pass-through: the provider sends it only when declared (no schema
  # default), so an unset spec field stays unsent.
  skip_host_key_validation = var.spec.sftp.skip_host_key_validation
  host_key_fingerprint     = var.spec.sftp.host_key_fingerprint != "" ? var.spec.sftp.host_key_fingerprint : null
}

# Snowflake -- connection string, password optionally in Key Vault.
resource "azurerm_data_factory_linked_service_snowflake" "main" {
  count = var.spec.snowflake != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.snowflake.connection_string

  dynamic "key_vault_password" {
    for_each = var.spec.snowflake.key_vault_password != null ? [var.spec.snowflake.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }
}

# Azure SQL Managed Instance -- connection string (inline or
# Key-Vault-held) plus an optional whole service principal. The one
# linked service the provider models WITHOUT an additional_properties
# argument (noted on the spec field).
resource "azurerm_data_factory_linked_service_sql_managed_instance" "main" {
  count = var.spec.sql_managed_instance != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.sql_managed_instance.connection_string != "" ? var.spec.sql_managed_instance.connection_string : null

  dynamic "key_vault_connection_string" {
    for_each = var.spec.sql_managed_instance.key_vault_connection_string != null ? [var.spec.sql_managed_instance.key_vault_connection_string] : []
    content {
      linked_service_name = key_vault_connection_string.value.linked_service_name
      secret_name         = key_vault_connection_string.value.secret_name
    }
  }

  dynamic "key_vault_password" {
    for_each = var.spec.sql_managed_instance.key_vault_password != null ? [var.spec.sql_managed_instance.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }

  service_principal_id  = var.spec.sql_managed_instance.service_principal_id != "" ? var.spec.sql_managed_instance.service_principal_id : null
  service_principal_key = var.spec.sql_managed_instance.service_principal_key != "" ? var.spec.sql_managed_instance.service_principal_key : null
  tenant                = var.spec.sql_managed_instance.tenant != "" ? var.spec.sql_managed_instance.tenant : null
}

# SQL Server -- connection string (inline or Key-Vault-held), password
# optionally split out to Key Vault.
resource "azurerm_data_factory_linked_service_sql_server" "main" {
  count = var.spec.sql_server != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.sql_server.connection_string != "" ? var.spec.sql_server.connection_string : null

  dynamic "key_vault_connection_string" {
    for_each = var.spec.sql_server.key_vault_connection_string != null ? [var.spec.sql_server.key_vault_connection_string] : []
    content {
      linked_service_name = key_vault_connection_string.value.linked_service_name
      secret_name         = key_vault_connection_string.value.secret_name
    }
  }

  dynamic "key_vault_password" {
    for_each = var.spec.sql_server.key_vault_password != null ? [var.spec.sql_server.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }

  user_name = var.spec.sql_server.user_name != "" ? var.spec.sql_server.user_name : null
}

# Azure Synapse Analytics -- connection string, password optionally in
# Key Vault.
resource "azurerm_data_factory_linked_service_synapse" "main" {
  count = var.spec.synapse != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  connection_string = var.spec.synapse.connection_string

  dynamic "key_vault_password" {
    for_each = var.spec.synapse.key_vault_password != null ? [var.spec.synapse.key_vault_password] : []
    content {
      linked_service_name = key_vault_password.value.linked_service_name
      secret_name         = key_vault_password.value.secret_name
    }
  }
}

# Web -- Anonymous or Basic (the two forms the provider wires).
resource "azurerm_data_factory_linked_service_web" "main" {
  count = var.spec.web != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description              = var.spec.description != "" ? var.spec.description : null
  annotations              = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  parameters               = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  additional_properties    = length(var.spec.additional_properties) > 0 ? var.spec.additional_properties : null
  integration_runtime_name = var.spec.integration_runtime_name != "" ? var.spec.integration_runtime_name : null

  url                 = var.spec.web.url
  authentication_type = var.spec.web.authentication_type
  username            = var.spec.web.username != "" ? var.spec.web.username : null
  password            = var.spec.web.password != "" ? var.spec.web.password : null
}
