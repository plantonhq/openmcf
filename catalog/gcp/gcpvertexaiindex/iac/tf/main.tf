# Enable the Vertex AI API — the control plane that owns indexes.
# disable_on_destroy is false: tearing down one index must never disable
# the API for everything else in the project (other Vertex resources keep
# working).
resource "google_project_service" "aiplatform_api" {
  project = local.project_id
  service = "aiplatform.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Vector Search index — the data structure holding embedding vectors.
# GCP assigns the numeric resource ID; display_name is the human handle.
# Region, index_update_method, and the whole config geometry are immutable
# (ForceNew). Creation is a long-running operation: minutes for an empty
# streaming index, up to hours for a large batch build (provider timeout:
# 180 minutes).
resource "google_vertex_ai_index" "this" {
  display_name = local.display_name
  region       = local.location
  project      = local.project_id
  labels       = local.final_labels

  description = var.spec.description != "" ? var.spec.description : null

  # Empty means "let GCP default" (BATCH_UPDATE). Sent only when set so
  # the plan diff stays honest about who chose the value.
  index_update_method = var.spec.index_update_method != "" ? var.spec.index_update_method : null

  # CMEK: the key must be in the index's region and the Vertex AI
  # service agent needs cryptoKeyEncrypterDecrypter on it. Omitted means
  # Google-managed encryption. Immutable.
  dynamic "encryption_spec" {
    for_each = var.spec.kms_key_name != "" ? [var.spec.kms_key_name] : []
    content {
      kms_key_name = encryption_spec.value
    }
  }

  # Client-side destroy behavior (DELETE deletes the corpus; PREVENT
  # refuses; ABANDON drops from state but keeps the index standing).
  # Empty follows the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  # The provider models the API's Index.metadata blob as this nested
  # block: data location + geometry together. contents_delta_uri is
  # write-only upstream (GCP never reports it back) and a change to it
  # travels in its own single-field PATCH — both quirks are documented in
  # the spec so a drift-looking plan after an out-of-band data load is
  # expected, not alarming.
  metadata {
    contents_delta_uri = local.has_contents ? var.spec.contents_delta_uri : null
    # Only meaningful alongside contents_delta_uri; false is the provider
    # default so sending it unconditionally is harmless and keeps the
    # value honest when contents are present.
    is_complete_overwrite = local.has_contents ? var.spec.is_complete_overwrite : null

    config {
      dimensions = var.spec.config.dimensions

      # Required by the API when tree-AH is used (CEL-enforced pre-deploy);
      # 0 means "not set" for brute-force or GCP-default algorithm.
      approximate_neighbors_count = var.spec.config.approximate_neighbors_count > 0 ? var.spec.config.approximate_neighbors_count : null

      # Computed when omitted: GCP picks a shard size from the data.
      shard_size = var.spec.config.shard_size != "" ? var.spec.config.shard_size : null

      distance_measure_type = var.spec.config.distance_measure_type != "" ? var.spec.config.distance_measure_type : null
      feature_norm_type     = var.spec.config.feature_norm_type != "" ? var.spec.config.feature_norm_type : null

      # At most one algorithm arm exists (CEL-enforced). Omitting both
      # lets GCP default (tree-AH with default tuning).
      dynamic "algorithm_config" {
        for_each = (var.spec.config.tree_ah_config != null || var.spec.config.brute_force_config != null) ? [1] : []
        content {
          dynamic "tree_ah_config" {
            for_each = var.spec.config.tree_ah_config != null ? [var.spec.config.tree_ah_config] : []
            content {
              # 0 means "not set": the provider then applies GCP's
              # defaults (1000 embeddings per leaf, 10% searched).
              leaf_node_embedding_count    = tree_ah_config.value.leaf_node_embedding_count > 0 ? tree_ah_config.value.leaf_node_embedding_count : null
              leaf_nodes_to_search_percent = tree_ah_config.value.leaf_nodes_to_search_percent > 0 ? tree_ah_config.value.leaf_nodes_to_search_percent : null
            }
          }

          dynamic "brute_force_config" {
            for_each = var.spec.config.brute_force_config != null ? [1] : []
            content {}
          }
        }
      }
    }
  }

  depends_on = [
    google_project_service.aiplatform_api,
  ]
}
