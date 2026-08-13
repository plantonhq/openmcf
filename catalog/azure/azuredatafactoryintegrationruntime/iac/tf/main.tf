# One kind, 3 provider resources: Azure stores every integration
# runtime flavor in the SAME factory-scoped namespace
# ({factory_id}/integrationRuntimes/{name}), so the spec's variant
# block selects which resource is created. Shared fields (name,
# factory, description) travel identically on every flavor; each
# resource below adds only its variant's own arguments.
#
# Optional arguments are sent only when set: the provider's own
# schema defaults then apply (8 cores / General compute on the azure
# flavor; 1 node / 1 parallel execution / Standard edition /
# LicenseIncluded on SSIS).

# The managed data-flow compute: serverless Spark Azure provisions
# when a data flow runs. virtual_network_enabled requires the
# factory's managed virtual network (Azure rejects the create
# otherwise); the interactive authoring TTL travels on a separate
# enable-interactive-authoring operation the provider performs after
# the runtime is online.
resource "azurerm_data_factory_integration_runtime_azure" "main" {
  count = var.spec.azure != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id
  location        = var.spec.azure.region

  description = var.spec.description != "" ? var.spec.description : null

  # Platform-default TRUE (tear the cluster down after every run) --
  # sent explicitly so the manifest's intent is always on the wire,
  # mirroring the provider's own default.
  cleanup_enabled = coalesce(var.spec.azure.cleanup_enabled, true)

  compute_type     = var.spec.azure.compute_type != "" ? var.spec.azure.compute_type : null
  core_count       = var.spec.azure.core_count != 0 ? var.spec.azure.core_count : null
  time_to_live_min = var.spec.azure.time_to_live_min != 0 ? var.spec.azure.time_to_live_min : null

  interactive_authoring_time_to_live_in_minutes = var.spec.azure.interactive_authoring_time_to_live_in_minutes != 0 ? var.spec.azure.interactive_authoring_time_to_live_in_minutes : null

  virtual_network_enabled = var.spec.azure.virtual_network_enabled
}

# The managed SSIS package runtime: a cluster of Azure-managed VMs.
# Creating it leaves it STOPPED and unbilled -- node-hours bill only
# after the runtime is started (an operational action, not part of
# this definition).
resource "azurerm_data_factory_integration_runtime_azure_ssis" "main" {
  count = var.spec.azure_ssis != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id
  location        = var.spec.azure_ssis.region
  node_size       = var.spec.azure_ssis.node_size

  description = var.spec.description != "" ? var.spec.description : null

  number_of_nodes                  = var.spec.azure_ssis.number_of_nodes != 0 ? var.spec.azure_ssis.number_of_nodes : null
  max_parallel_executions_per_node = var.spec.azure_ssis.max_parallel_executions_per_node != 0 ? var.spec.azure_ssis.max_parallel_executions_per_node : null
  edition                          = var.spec.azure_ssis.edition != "" ? var.spec.azure_ssis.edition : null
  license_type                     = var.spec.azure_ssis.license_type != "" ? var.spec.azure_ssis.license_type : null
  credential_name                  = var.spec.azure_ssis.credential_name != "" ? var.spec.azure_ssis.credential_name : null

  dynamic "catalog_info" {
    for_each = var.spec.azure_ssis.catalog_info != null ? [var.spec.azure_ssis.catalog_info] : []
    content {
      server_endpoint        = catalog_info.value.server_endpoint
      administrator_login    = catalog_info.value.administrator_login != "" ? catalog_info.value.administrator_login : null
      administrator_password = catalog_info.value.administrator_password != "" ? catalog_info.value.administrator_password : null
      pricing_tier           = catalog_info.value.pricing_tier != "" ? catalog_info.value.pricing_tier : null
      elastic_pool_name      = catalog_info.value.elastic_pool_name != "" ? catalog_info.value.elastic_pool_name : null
      dual_standby_pair_name = catalog_info.value.dual_standby_pair_name != "" ? catalog_info.value.dual_standby_pair_name : null
    }
  }

  dynamic "custom_setup_script" {
    for_each = var.spec.azure_ssis.custom_setup_script != null ? [var.spec.azure_ssis.custom_setup_script] : []
    content {
      blob_container_uri = custom_setup_script.value.blob_container_uri
      sas_token          = custom_setup_script.value.sas_token
    }
  }

  dynamic "express_custom_setup" {
    for_each = var.spec.azure_ssis.express_custom_setup != null ? [var.spec.azure_ssis.express_custom_setup] : []
    content {
      environment        = length(express_custom_setup.value.environment) > 0 ? express_custom_setup.value.environment : null
      powershell_version = express_custom_setup.value.powershell_version != "" ? express_custom_setup.value.powershell_version : null

      # When both password/license forms are set, Azure receives the
      # inline value (the provider's own precedence) -- prefer the Key
      # Vault forms.
      dynamic "command_key" {
        for_each = express_custom_setup.value.command_key
        content {
          target_name = command_key.value.target_name
          user_name   = command_key.value.user_name
          password    = command_key.value.password != "" ? command_key.value.password : null

          dynamic "key_vault_password" {
            for_each = command_key.value.key_vault_password != null ? [command_key.value.key_vault_password] : []
            content {
              linked_service_name = key_vault_password.value.linked_service_name
              secret_name         = key_vault_password.value.secret_name
              secret_version      = key_vault_password.value.secret_version != "" ? key_vault_password.value.secret_version : null
              parameters          = length(key_vault_password.value.parameters) > 0 ? key_vault_password.value.parameters : null
            }
          }
        }
      }

      dynamic "component" {
        for_each = express_custom_setup.value.component
        content {
          name    = component.value.name
          license = component.value.license != "" ? component.value.license : null

          dynamic "key_vault_license" {
            for_each = component.value.key_vault_license != null ? [component.value.key_vault_license] : []
            content {
              linked_service_name = key_vault_license.value.linked_service_name
              secret_name         = key_vault_license.value.secret_name
              secret_version      = key_vault_license.value.secret_version != "" ? key_vault_license.value.secret_version : null
              parameters          = length(key_vault_license.value.parameters) > 0 ? key_vault_license.value.parameters : null
            }
          }
        }
      }
    }
  }

  dynamic "express_vnet_integration" {
    for_each = var.spec.azure_ssis.express_vnet_integration != null ? [var.spec.azure_ssis.express_vnet_integration] : []
    content {
      subnet_id = express_vnet_integration.value.subnet_id
    }
  }

  # Exactly one of vnet_id (+ subnet_name) or subnet_id is set (the
  # spec enforces it); subnet_name travels only alongside vnet_id.
  dynamic "vnet_integration" {
    for_each = var.spec.azure_ssis.vnet_integration != null ? [var.spec.azure_ssis.vnet_integration] : []
    content {
      vnet_id     = vnet_integration.value.vnet_id != "" ? vnet_integration.value.vnet_id : null
      subnet_id   = vnet_integration.value.subnet_id != "" ? vnet_integration.value.subnet_id : null
      subnet_name = vnet_integration.value.subnet_name != "" ? vnet_integration.value.subnet_name : null
      public_ips  = length(vnet_integration.value.public_ips) > 0 ? vnet_integration.value.public_ips : null
    }
  }

  dynamic "package_store" {
    for_each = var.spec.azure_ssis.package_store
    content {
      name                = package_store.value.name
      linked_service_name = package_store.value.linked_service_name
    }
  }

  dynamic "copy_compute_scale" {
    for_each = var.spec.azure_ssis.copy_compute_scale != null ? [var.spec.azure_ssis.copy_compute_scale] : []
    content {
      data_integration_unit = copy_compute_scale.value.data_integration_unit != 0 ? copy_compute_scale.value.data_integration_unit : null
      time_to_live          = copy_compute_scale.value.time_to_live != 0 ? copy_compute_scale.value.time_to_live : null
    }
  }

  # number_of_external_nodes is sent correctly but Azure's read API
  # mirrors number_of_pipeline_nodes back for it -- a provider read
  # seam, documented on the spec field.
  dynamic "pipeline_external_compute_scale" {
    for_each = var.spec.azure_ssis.pipeline_external_compute_scale != null ? [var.spec.azure_ssis.pipeline_external_compute_scale] : []
    content {
      number_of_external_nodes = pipeline_external_compute_scale.value.number_of_external_nodes != 0 ? pipeline_external_compute_scale.value.number_of_external_nodes : null
      number_of_pipeline_nodes = pipeline_external_compute_scale.value.number_of_pipeline_nodes != 0 ? pipeline_external_compute_scale.value.number_of_pipeline_nodes : null
      time_to_live             = pipeline_external_compute_scale.value.time_to_live != 0 ? pipeline_external_compute_scale.value.time_to_live : null
    }
  }

  dynamic "proxy" {
    for_each = var.spec.azure_ssis.proxy != null ? [var.spec.azure_ssis.proxy] : []
    content {
      self_hosted_integration_runtime_name = proxy.value.self_hosted_integration_runtime_name
      staging_storage_linked_service_name  = proxy.value.staging_storage_linked_service_name
      path                                 = proxy.value.path != "" ? proxy.value.path : null
    }
  }
}

# The self-hosted agent registration: Azure allocates the
# registration and issues the authorization keys agents join with.
# Creating it is free -- the compute is yours.
resource "azurerm_data_factory_integration_runtime_self_hosted" "main" {
  count = var.spec.self_hosted != null ? 1 : 0

  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description = var.spec.description != "" ? var.spec.description : null

  self_contained_interactive_authoring_enabled = var.spec.self_hosted.self_contained_interactive_authoring_enabled

  # A linked registration (sharing another factory's runtime through
  # RBAC) -- Azure issues no keys for it.
  dynamic "rbac_authorization" {
    for_each = var.spec.self_hosted.rbac_authorization != null ? [var.spec.self_hosted.rbac_authorization] : []
    content {
      resource_id = rbac_authorization.value.resource_id
    }
  }
}
