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

  # MongoDB-style array indexing; only legal under MONGODB_COMPATIBLE_API
  # (spec CEL enforces the pairing). False is the API default — send null
  # instead so plans stay clean on older indexes.
  multikey = var.spec.multikey ? true : null

  # Uniqueness enforcement across documents.
  unique = var.spec.unique ? true : null

  # Client-side: return once creation is REQUESTED; the background build
  # continues and the index serves queries only when it completes.
  skip_wait = var.spec.skip_wait ? true : null

  # Destroy-time guard: PREVENT fails the destroy; ABANDON unmanages the
  # index without deleting it. Null falls back to the provider default
  # (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  dynamic "fields" {
    for_each = var.spec.fields
    content {
      field_path   = fields.value.field_path
      order        = fields.value.order != "" ? fields.value.order : null
      array_config = fields.value.array_config != "" ? fields.value.array_config : null

      dynamic "vector_config" {
        for_each = fields.value.vector_config != null ? [fields.value.vector_config] : []
        content {
          dimension = vector_config.value.dimension
          # Flat is the only vector index layout GCP offers today; the
          # module pins it so every vector field is expressible without a
          # spec knob that has exactly one legal value.
          flat {}
        }
      }

      # Firestore Enterprise search surface (requires an ENTERPRISE-edition
      # database; text search pairs with api_scope MONGODB_COMPATIBLE_API).
      dynamic "search_config" {
        for_each = fields.value.search_config != null ? [fields.value.search_config] : []
        content {
          dynamic "text_spec" {
            for_each = search_config.value.text_spec != null ? [search_config.value.text_spec] : []
            content {
              dynamic "index_specs" {
                for_each = text_spec.value.index_specs
                content {
                  index_type = index_specs.value.index_type != "" ? index_specs.value.index_type : null
                  match_type = index_specs.value.match_type != "" ? index_specs.value.match_type : null
                }
              }
            }
          }

          dynamic "geo_spec" {
            for_each = search_config.value.geo_spec != null ? [search_config.value.geo_spec] : []
            content {
              geo_json_indexing_disabled = geo_spec.value.geo_json_indexing_disabled
            }
          }
        }
      }
    }
  }

  depends_on = [google_project_service.firestore_api]
}
