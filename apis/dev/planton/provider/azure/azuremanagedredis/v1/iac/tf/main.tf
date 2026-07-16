# The Managed Redis instance: the cluster (compute, load balancer, TLS
# endpoint) plus its default database (the Redis process), which Azure
# maps 1-to-1 -- one resource here, exactly as ARM provisions them.
# Provisioning polls the cluster to its Running state and then creates
# the database; expect tens of minutes end to end (azurerm's own create
# timeout is 45 minutes).
resource "azurerm_managed_redis" "main" {
  name                = var.spec.cluster_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # One value picks the tier family and memory size. Azure validates
  # in-place SKU changes against the live instance at apply time and
  # replaces the instance when the target is not scalable.
  sku_name = local.sku_name

  # HA runs a replica and carries the zone-redundant SLA; disabling it
  # halves cost for dev/test. Fixed at creation.
  high_availability_enabled = var.spec.high_availability_enabled

  # false forces all traffic through Private Link (AzurePrivateEndpoint)
  # -- Managed Redis has no VNet injection or IP firewall.
  public_network_access = local.public_network_access

  # Customer-managed-key encryption. The key id is the VERSIONED Key
  # Vault id (rotation = updating the reference); the same identity must
  # also be attached below -- an ARM pairing enforced at apply time.
  dynamic "customer_managed_key" {
    for_each = var.spec.customer_managed_key != null ? [var.spec.customer_managed_key] : []
    content {
      key_vault_key_id          = customer_managed_key.value.key_vault_key_id
      user_assigned_identity_id = customer_managed_key.value.user_assigned_identity_id
    }
  }

  # The managed identity -- what customer-managed-key encryption
  # authenticates to Key Vault with.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  # The Redis process itself. Required at create (Azure rejects a
  # database-less cluster). The database enums deploy Azure's own
  # defaults explicitly so both engines send identical bodies; changing
  # clustering_policy, geo_replication_group_name, or the module set
  # recreates the DATABASE in place (data loss, brief unavailability)
  # without replacing the cluster.
  default_database {
    # Keyless-first: keys are OFF by default; Entra grants
    # (AzureManagedRedisAccessPolicyAssignment) are how clients connect.
    access_keys_authentication_enabled = var.spec.default_database.access_keys_authentication_enabled

    client_protocol   = local.client_protocol
    clustering_policy = local.clustering_policy
    eviction_policy   = local.eviction_policy

    # Joining a named ACTIVE geo-replication group; membership is linked
    # by AzureManagedRedisGeoReplication.
    geo_replication_group_name = var.spec.default_database.geo_replication_group_name

    # Redis modules (search, JSON, bloom, time series) -- capabilities
    # classic Redis never had.
    dynamic "module" {
      for_each = var.spec.default_database.modules
      content {
        name = module.value.name
        args = module.value.args
      }
    }

    # Setting a frequency ENABLES the matching persistence method. AOF
    # and RDB are mutually exclusive, and both conflict with
    # geo-replication (spec-enforced, mirroring Azure's own contract).
    persistence_append_only_file_backup_frequency = var.spec.default_database.persistence_append_only_file_backup_frequency
    persistence_redis_database_backup_frequency   = var.spec.default_database.persistence_redis_database_backup_frequency
  }

  tags = local.final_tags
}
