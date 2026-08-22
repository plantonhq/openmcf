# DigitalOcean Database Db
#
# Provisions an additional logical database inside a DigitalOcean managed
# database cluster -- the complete digitalocean_database_db resource
# surface. Both arguments are create-only: any change replaces the logical
# database and DROPS its data, so treat renames as migrations. DigitalOcean's
# read is a bare existence check; connection credentials live on the cluster
# and its users, not here.

resource "digitalocean_database_db" "database" {
  cluster_id = var.spec.cluster
  name       = var.spec.database_name
}
