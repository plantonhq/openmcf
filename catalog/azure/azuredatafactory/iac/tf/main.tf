# Create the Azure Data Factory -- the workspace every other Data
# Factory resource lives inside. The repository binding (GitHub or
# Azure DevOps, at most one -- spec CEL mirrors the provider's
# ConflictsWith) travels through a separate configure-repo call the
# provider makes AFTER the factory exists, and REMOVING the block does
# not detach the repository (the provider calls no repo-clear API) --
# documented on the spec fields.
resource "azurerm_data_factory" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Platform default true, sent explicitly (mirrors Azure's own
  # default; the provider maps the bool onto ARM's Enabled/Disabled
  # tokens itself).
  public_network_enabled = coalesce(var.spec.public_network_enabled, true)

  # Enabling is an in-place update (the provider creates the managed
  # network named "default" after the factory); DISABLING it is
  # ForceNew (the provider's one CustomizeDiff on this resource) --
  # documented on the spec field. False and omitted are the same wire
  # shape, so the platform default is sent for an explicit plan.
  managed_virtual_network_enabled = coalesce(var.spec.managed_virtual_network_enabled, false)

  purview_id = var.spec.purview_id != "" ? var.spec.purview_id : null

  # Customer-managed-key encryption composed onto the factory's own
  # inline fields (the provider's standalone CMK resource writes the
  # same encryption object; it demands a VERSIONED key where these
  # inline fields accept versionless too -- the spec prefers
  # versionless so rotation propagates). The unwrap identity is
  # required with the key (spec makes it required in the block,
  # front-loading the provider's create-time CustomizeDiff).
  customer_managed_key_id          = var.spec.customer_managed_key != null ? var.spec.customer_managed_key.key_vault_key_id : null
  customer_managed_key_identity_id = var.spec.customer_managed_key != null ? var.spec.customer_managed_key.user_assigned_identity_id : null

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  dynamic "github_configuration" {
    for_each = var.spec.github_configuration != null ? [var.spec.github_configuration] : []
    content {
      account_name    = github_configuration.value.account_name
      branch_name     = github_configuration.value.branch_name
      git_url         = github_configuration.value.git_url != "" ? github_configuration.value.git_url : null
      repository_name = github_configuration.value.repository_name
      root_folder     = github_configuration.value.root_folder

      # Platform default true, sent explicitly (Azure stores the
      # inverse disablePublish flag; the provider translates).
      publishing_enabled = coalesce(github_configuration.value.publishing_enabled, true)
    }
  }

  dynamic "vsts_configuration" {
    for_each = var.spec.vsts_configuration != null ? [var.spec.vsts_configuration] : []
    content {
      account_name       = vsts_configuration.value.account_name
      branch_name        = vsts_configuration.value.branch_name
      project_name       = vsts_configuration.value.project_name
      repository_name    = vsts_configuration.value.repository_name
      root_folder        = vsts_configuration.value.root_folder
      tenant_id          = vsts_configuration.value.tenant_id
      publishing_enabled = coalesce(vsts_configuration.value.publishing_enabled, true)
    }
  }

  # Workspace-wide parameters (names unique -- spec CEL front-loads
  # the provider's own duplicate check). Array/Object values travel as
  # JSON text; Azure stores the typed value.
  dynamic "global_parameter" {
    for_each = var.spec.global_parameters
    content {
      name  = global_parameter.value.name
      type  = global_parameter.value.type
      value = global_parameter.value.value
    }
  }

  tags = local.final_tags
}

# Composed user-managed-identity credentials -- one provider resource
# per named credential, keyed by name (renames replace only that
# credential), in lockstep with the Pulumi module's per-name
# resources. Linked services reference credentials BY NAME.
resource "azurerm_data_factory_credential_user_managed_identity" "main" {
  for_each = local.user_managed_identity_credentials_by_name

  name            = each.value.name
  data_factory_id = azurerm_data_factory.main.id
  identity_id     = each.value.identity_id
  description     = each.value.description != "" ? each.value.description : null
  annotations     = length(each.value.annotations) > 0 ? each.value.annotations : null
}

# Composed service-principal credentials. The principal's key is
# resolved through a Key Vault LINKED SERVICE by name -- the secret
# itself never travels through this module.
resource "azurerm_data_factory_credential_service_principal" "main" {
  for_each = local.service_principal_credentials_by_name

  name                 = each.value.name
  data_factory_id      = azurerm_data_factory.main.id
  tenant_id            = each.value.tenant_id
  service_principal_id = each.value.service_principal_id
  description          = each.value.description != "" ? each.value.description : null
  annotations          = length(each.value.annotations) > 0 ? each.value.annotations : null

  dynamic "service_principal_key" {
    for_each = each.value.service_principal_key != null ? [each.value.service_principal_key] : []
    content {
      linked_service_name = service_principal_key.value.linked_service_name
      secret_name         = service_principal_key.value.secret_name
      secret_version      = service_principal_key.value.secret_version != "" ? service_principal_key.value.secret_version : null
    }
  }
}

# Composed managed private endpoints -- private egress from the
# factory's managed virtual network (spec CEL requires the network).
# Each endpoint is CREATE-ONLY in the provider (no Update method):
# every field change replaces that endpoint, siblings stay untouched.
# Exactly one arm is set per entry (spec CEL): subresource_name for
# regular ARM targets, fqdns for Private Link Service targets. The
# TARGET side must approve the endpoint's connection before traffic
# flows -- approval happens outside this resource, and the endpoint
# provisions to Succeeded while the connection is still Pending.
resource "azurerm_data_factory_managed_private_endpoint" "main" {
  for_each = local.managed_private_endpoints_by_name

  name               = each.value.name
  data_factory_id    = azurerm_data_factory.main.id
  target_resource_id = each.value.target_resource_id
  subresource_name   = each.value.subresource_name != "" ? each.value.subresource_name : null
  fqdns              = length(each.value.fqdns) > 0 ? each.value.fqdns : null
}
