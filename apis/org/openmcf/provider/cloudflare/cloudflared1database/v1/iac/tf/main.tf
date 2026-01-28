# main.tf

# Create the Cloudflare D1 database
resource "cloudflare_d1_database" "main" {
  account_id = var.spec.account_id
  name       = var.spec.database_name

  # Add optional primary location hint (region) if specified
  primary_location_hint = var.spec.region

  # Add optional read replication configuration if specified
  dynamic "read_replication" {
    for_each = var.spec.read_replication != null ? [var.spec.read_replication] : []
    content {
      mode = read_replication.value.mode
    }
  }
}

