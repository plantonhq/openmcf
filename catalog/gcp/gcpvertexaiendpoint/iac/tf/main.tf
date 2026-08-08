# Enable the Vertex AI API — the control plane that owns endpoints.
# disable_on_destroy is false: tearing down one endpoint must never disable
# the API for everything else in the project (other endpoints keep serving).
resource "google_project_service" "aiplatform_api" {
  project = local.project_id
  service = "aiplatform.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Vertex AI endpoint — the durable serving surface models deploy onto.
# Creating the endpoint is the infrastructure concern; deploying models to
# it is an operational step performed via the Vertex AI API. Location,
# network, CMEK, and the numeric name are all immutable (ForceNew).
resource "google_vertex_ai_endpoint" "this" {
  name         = local.endpoint_name
  display_name = local.display_name
  location     = local.location
  # The provider resolves the Vertex AI API host from `region`
  # (https://{region}-aiplatform.googleapis.com), never from `location` —
  # without it, deploys fail with "Cannot determine region" unless the
  # provider config happens to carry one. The two fields are the same axis
  # for this regional API, so the module pins region to location and the
  # spec keeps a single honest field.
  region  = local.location
  project = local.project_id
  labels  = local.final_labels

  description                = var.spec.description != "" ? var.spec.description : null
  dedicated_endpoint_enabled = var.spec.dedicated_endpoint_enabled ? true : null

  # VPC-peered private networking (requires Private Services Access on the
  # network; mutually exclusive with PSC — enforced pre-deploy by the
  # spec's CEL rule).
  network = var.spec.network != "" ? var.spec.network : null

  # CMEK: the endpoint and all sub-resources encrypted under this key.
  dynamic "encryption_spec" {
    for_each = var.spec.kms_key_name != "" ? [var.spec.kms_key_name] : []
    content {
      kms_key_name = encryption_spec.value
    }
  }

  # Private Service Connect: the endpoint exposed via a service attachment.
  # Secure PSC (IAM authorization on top of network reachability) is not
  # offered: the GA provider does not expose it, and GA is the catalog's
  # parity baseline.
  dynamic "private_service_connect_config" {
    for_each = var.spec.private_service_connect_config != null ? [var.spec.private_service_connect_config] : []
    content {
      enable_private_service_connect = true
      project_allowlist              = length(private_service_connect_config.value.project_allowlist) > 0 ? private_service_connect_config.value.project_allowlist : null
    }
  }

  # Request/response logging: samples online predictions into BigQuery —
  # the raw material for drift monitoring and audit. sampling_rate 0 means
  # "not set" (the API then applies its own default).
  dynamic "predict_request_response_logging_config" {
    for_each = var.spec.request_response_logging_config != null ? [var.spec.request_response_logging_config] : []
    content {
      enabled       = predict_request_response_logging_config.value.enabled
      sampling_rate = predict_request_response_logging_config.value.sampling_rate != 0 ? predict_request_response_logging_config.value.sampling_rate : null

      dynamic "bigquery_destination" {
        for_each = predict_request_response_logging_config.value.bigquery_destination_uri != "" ? [predict_request_response_logging_config.value.bigquery_destination_uri] : []
        content {
          output_uri = bigquery_destination.value
        }
      }
    }
  }

  depends_on = [
    google_project_service.aiplatform_api,
  ]
}
