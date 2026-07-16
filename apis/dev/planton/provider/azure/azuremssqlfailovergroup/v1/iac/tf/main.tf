# Create the SQL failover group on the primary logical server. The group
# replicates the listed databases to the partner servers and exposes a
# listener endpoint that follows the primary through a failover, so
# applications never change their connection string.
resource "azurerm_mssql_failover_group" "main" {
  name      = var.spec.name
  server_id = var.spec.server_id

  # Every partner server (a different region than the primary). At least one
  # is required.
  dynamic "partner_server" {
    for_each = var.spec.partner_servers
    content {
      id = partner_server.value.server_id
    }
  }

  # The databases to replicate; empty means an empty group (databases added
  # later). Databases must live on the primary server.
  databases = length(var.spec.database_ids) > 0 ? var.spec.database_ids : null

  read_write_endpoint_failover_policy {
    mode = local.failover_mode
    # grace_minutes is required for Automatic and rejected for Manual (the
    # spec CEL guarantees the pairing); send it only for Automatic.
    grace_minutes = local.failover_mode == "Automatic" ? var.spec.read_write_endpoint_failover_policy.grace_minutes : null
  }

  # The provider sends Disabled for the read-only endpoint when this is
  # null, so an unspecified spec deploys identically on both engines.
  readonly_endpoint_failover_policy_enabled = var.spec.readonly_endpoint_failover_policy_enabled

  tags = local.final_tags
}
