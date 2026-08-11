# Enable the AlloyDB API so a fresh project can host users.
resource "google_project_service" "alloydb_api" {
  project = local.project_id
  service = "alloydb.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A database user on an AlloyDB cluster. Users are first-class nodes: one
# per application with its own credential (ALLOYDB_BUILT_IN) or passwordless
# IAM authentication (ALLOYDB_IAM_USER).
resource "google_alloydb_user" "this" {
  cluster   = var.spec.cluster
  user_id   = var.spec.user_id
  user_type = local.user_type

  password       = local.password
  database_roles = length(var.spec.database_roles) > 0 ? var.spec.database_roles : null

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.alloydb_api]
}
