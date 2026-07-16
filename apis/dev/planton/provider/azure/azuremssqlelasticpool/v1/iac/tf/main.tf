# The Azure SQL elastic pool: shared compute that member databases
# (AzureMssqlDatabase with sku_name "ElasticPool") draw from instead of
# carrying their own SKU. The server is addressed by name + resource
# group (azurerm's contract), both derived from the spec's server ARM id;
# the region must match the server's (ARM rejects a mismatch).
resource "azurerm_mssql_elasticpool" "main" {
  name                = var.spec.pool_name
  location            = var.spec.region
  resource_group_name = local.resource_group_name
  server_name         = local.server_name

  # The tier and family are derived from the sku name (pure functions of
  # it), so a mismatched combination is unrepresentable.
  sku {
    name     = var.spec.sku_name
    tier     = local.sku_tier
    family   = local.sku_family
    capacity = var.spec.capacity
  }

  # What any ONE member database may consume: min is guaranteed
  # (reserved even while idle), max caps noisy neighbors.
  per_database_settings {
    min_capacity = var.spec.per_database_settings.min_capacity
    max_capacity = var.spec.per_database_settings.max_capacity
  }

  # The pool's total storage cap -- gigabytes XOR bytes (spec-validated);
  # both null lets ARM apply the SKU default.
  max_size_gb    = var.spec.max_size_gb
  max_size_bytes = var.spec.max_size_bytes

  zone_redundant = var.spec.zone_redundant

  # Every database in the pool must share the enclave type, so it lives
  # at the pool level. Changing it is disruptive -- plan accordingly.
  enclave_type = local.enclave_type

  license_type = local.license_type

  # Hyperscale pools only: readable HA replicas per member database.
  high_availability_replica_count = var.spec.high_availability_replica_count

  # Member databases inherit this window (and must not set their own).
  maintenance_configuration_name = var.spec.maintenance_configuration_name

  tags = local.final_tags
}
