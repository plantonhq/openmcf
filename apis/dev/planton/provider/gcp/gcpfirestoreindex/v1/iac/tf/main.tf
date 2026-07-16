# Enable the Firestore API — composite indexes are managed through the
# Firestore Admin API. disable_on_destroy stays false: tearing down one
# index must never disable the API for everything else in the project.
resource "google_project_service" "firestore_api" {
  project = local.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# Firestore composite index — the prerequisite for multi-field queries.
# Every property is immutable: changing anything replaces the index
# (Firestore rebuilds it in the background; the old index serves queries
# until the new one is ready). Firestore appends __name__ automatically.
resource "google_firestore_index" "this" {
  project     = local.project_id
  database    = local.database
  collection  = var.spec.collection
  query_scope = var.spec.query_scope
  api_scope   = var.spec.api_scope
  density     = var.spec.density != "" ? var.spec.density : null

  dynamic "fields" {
    for_each = var.spec.fields
    content {
      field_path = fields.value.field_path
      order      = fields.value.order != "" ? fields.value.order : null
      array_config = fields.value.array_config != "" ? fields.value.array_config : null

      dynamic "vector_config" {
        for_each = fields.value.vector_config != null ? [fields.value.vector_config] : []
        content {
          dimension = vector_config.value.dimension
          flat {}
        }
      }
    }
  }

  depends_on = [google_project_service.firestore_api]
}
