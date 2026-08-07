# The SQL (NoSQL) API database -- the namespace containers live in and
# the boundary for shared throughput. Addressed by the (resource group,
# account, name) trio azurerm requires, parsed from the parent account's
# ARM ID. No Azure tags: ARM does not support tags on Cosmos child
# resources, so the platform's identity tags live on the account.
resource "azurerm_cosmosdb_sql_database" "main" {
  name                = var.spec.database_name
  resource_group_name = local.resource_group_name
  account_name        = local.cosmosdb_account_name

  # Shared fixed throughput for the database's containers. Sent only
  # when set: serverless accounts reject provisioned throughput, and
  # unset means containers bring their own.
  throughput = var.spec.throughput

  # Autoscale: Azure scales the shared throughput between 10% and 100%
  # of the ceiling. The spec enforces mutual exclusion with the fixed
  # throughput before the plan ever runs.
  dynamic "autoscale_settings" {
    for_each = var.spec.autoscale_max_throughput != null ? [var.spec.autoscale_max_throughput] : []
    content {
      max_throughput = autoscale_settings.value
    }
  }
}
