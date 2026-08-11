# Enable the Vertex AI API — the control plane that owns index endpoints.
# disable_on_destroy is false: tearing down one endpoint must never disable
# the API for everything else in the project (other Vertex resources keep
# working).
resource "google_project_service" "aiplatform_api" {
  project = local.project_id
  service = "aiplatform.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Vector Search index endpoint — the serving surface deployed indexes
# answer queries through. This is a DIFFERENT GCP resource from the
# online-prediction google_vertex_ai_endpoint (which serves models). GCP
# assigns the numeric resource ID; display_name is the human handle. Every
# connectivity choice (public / peered network / PSC) is immutable
# (ForceNew); display_name, description, and labels PATCH in place.
resource "google_vertex_ai_index_endpoint" "this" {
  display_name = local.display_name
  region       = local.location
  project      = local.project_id
  labels       = local.final_labels

  description = var.spec.description != "" ? var.spec.description : null

  # Public querying arm: deployed indexes become reachable through
  # public_endpoint_domain_name. False is the provider default, so only
  # true is sent.
  public_endpoint_enabled = var.spec.public_endpoint_enabled ? true : null

  # VPC-peered private querying (requires Private Services Access on the
  # network; mutually exclusive with the other arms — enforced pre-deploy
  # by the spec's CEL rules). local.network carries the API's relative
  # form regardless of whether the spec supplied a self-link.
  network = local.network != "" ? local.network : null

  # Private Service Connect: consumers reach deployed indexes through a
  # service attachment (surfaced on the GcpVertexAiDeployedIndex outputs
  # once an index is deployed).
  dynamic "private_service_connect_config" {
    for_each = var.spec.private_service_connect_config != null ? [var.spec.private_service_connect_config] : []
    content {
      enable_private_service_connect = true
      project_allowlist              = length(private_service_connect_config.value.project_allowlist) > 0 ? private_service_connect_config.value.project_allowlist : null

      # PSC endpoints Vertex AI provisions automatically in consumer
      # projects. The API wants the relative network form; references
      # arrive as self-links and are normalized in locals.tf.
      dynamic "psc_automation_configs" {
        for_each = private_service_connect_config.value.psc_automation_configs
        content {
          network    = replace(psc_automation_configs.value.network, "/^https://www\\.googleapis\\.com/compute/v1//", "")
          project_id = psc_automation_configs.value.project_id
        }
      }
    }
  }

  # CMEK: the key must be in the endpoint's region and the Vertex AI
  # service agent needs cryptoKeyEncrypterDecrypter on it. Omitted means
  # Google-managed encryption. Immutable.
  dynamic "encryption_spec" {
    for_each = var.spec.kms_key_name != "" ? [var.spec.kms_key_name] : []
    content {
      kms_key_name = encryption_spec.value
    }
  }

  # Client-side destroy behavior (DELETE deletes the endpoint and stops
  # every deployment on it; PREVENT refuses; ABANDON drops from state
  # but keeps serving). Empty follows the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [
    google_project_service.aiplatform_api,
  ]
}
