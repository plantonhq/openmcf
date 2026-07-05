locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_mssql_server"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over
  # them: user tags deliberately win so an org's governance conventions
  # (cost center, owner) can override the derived values where they
  # collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # The spec's enums arrive as FULL proto value names (the tfvars wire
  # format never strips prefixes); each map below carries the complete
  # verbatim vocabulary for its enum, mapped to ARM's values. A missing
  # entry would silently drop the setting, so the maps are exhaustive by
  # construction.
  connection_policy_map = {
    "DEFAULT"  = "Default"
    "PROXY"    = "Proxy"
    "REDIRECT" = "Redirect"
  }

  identity_type_map = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }

  alert_policy_state_map = {
    "ENABLED"  = "Enabled"
    "DISABLED" = "Disabled"
  }

  # Defender's detector wire vocabulary uses Snake_Pascal casing.
  alert_type_map = {
    "SQL_INJECTION"               = "Sql_Injection"
    "SQL_INJECTION_VULNERABILITY" = "Sql_Injection_Vulnerability"
    "ACCESS_ANOMALY"              = "Access_Anomaly"
    "DATA_EXFILTRATION"           = "Data_Exfiltration"
    "UNSAFE_ACTION"               = "Unsafe_Action"
  }

  # Unset connection_policy maps to null so Azure applies its Default
  # policy without the module materializing a value.
  connection_policy = (
    var.spec.connection_policy == null || var.spec.connection_policy == "" ? null :
    local.connection_policy_map[var.spec.connection_policy]
  )

  # Credentials are sent only when non-empty: an Entra-only server omits
  # them entirely and azurerm rejects empty strings.
  administrator_login = (
    var.spec.administrator_login == null || var.spec.administrator_login == "" ? null :
    var.spec.administrator_login
  )

  administrator_password = (
    var.spec.administrator_password == null || var.spec.administrator_password == "" ? null :
    var.spec.administrator_password
  )

  # The Entra tenant for the administrator grant: the spec's explicit
  # tenant wins, otherwise the deploying credential's tenant -- the
  # correct value for virtually every deployment.
  aad_tenant_id = (
    var.spec.azuread_administrator != null && var.spec.azuread_administrator.tenant_id != null && var.spec.azuread_administrator.tenant_id != ""
    ? var.spec.azuread_administrator.tenant_id
    : data.azurerm_client_config.current.tenant_id
  )
}
