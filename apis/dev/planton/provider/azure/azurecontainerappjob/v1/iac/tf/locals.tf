locals {
  resource_id = var.metadata.id != null ? var.metadata.id : var.metadata.name

  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_container_app_job"
    "resource_name" = var.metadata.name
  }

  org_tag = var.metadata.org != null ? { "organization" = var.metadata.org } : {}
  env_tag = var.metadata.env != null ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over
  # them: user tags deliberately win so an org's governance conventions
  # can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # Probe transport wire values (Azure validates these case-sensitively).
  probe_transport_map = {
    "TCP_SOCKET" = "TCP"
    "HTTP_GET"   = "HTTP"
    "HTTPS_GET"  = "HTTPS"
  }

  # Volume storage-type wire values; absent deploys EmptyDir.
  volume_storage_type_map = {
    "EMPTY_DIR"      = "EmptyDir"
    "AZURE_FILE"     = "AzureFile"
    "NFS_AZURE_FILE" = "NfsAzureFile"
    "SECRET"         = "Secret"
  }

  # Managed-identity type wire values.
  identity_type_map = {
    "SYSTEM_ASSIGNED"          = "SystemAssigned"
    "USER_ASSIGNED"            = "UserAssigned"
    "SYSTEM_AND_USER_ASSIGNED" = "SystemAssigned, UserAssigned"
  }
}
