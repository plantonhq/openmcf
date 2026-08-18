# DigitalOcean Database Connection Pool
#
# Provisions a PgBouncer connection pool on a DigitalOcean managed
# PostgreSQL cluster -- the complete digitalocean_database_connection_pool
# resource surface.
#
# EVERY argument is create-only: the provider registers no update path for
# pools (DigitalOcean's API could update them in place; mirroring the
# provider is the recorded parity decision), so any change replaces the
# pool and drops its live connections. Plan accordingly in production.

resource "digitalocean_database_connection_pool" "pool" {
  cluster_id = var.spec.cluster
  name       = var.spec.pool_name
  mode       = var.spec.mode
  size       = var.spec.size
  db_name    = var.spec.db_name

  # Omitted user = inbound-user pool (clients bring their own credentials).
  user = local.user
}
