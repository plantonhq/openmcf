# Enable the Eventarc API so a fresh project can host the trigger.
# disable_on_destroy is false: tearing down one trigger must never disable
# Eventarc for everything else in the project.
resource "google_project_service" "eventarc_api" {
  project = local.project_id
  service = "eventarc.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The provider's default project — read only when the googleChannelConfig
# singleton is armed with an ambient project (its name argument IS the
# fixed singleton path, which needs the resolved project id). The
# count-gate keeps every other configuration credential-free at plan time.
data "google_client_config" "this" {
  count = local.create_google_channel_config && var.spec.project_id == "" ? 1 : 0
}

# Partner (SaaS) channel companion — created when the spec arms
# partner_channel. The channel's one-time activation_token must be handed
# to the partner to complete the handshake; until then the channel stays
# PENDING and delivers nothing. Immutable: name and provider replace the
# channel (a NEW activation token — redo the handshake).
resource "google_eventarc_channel" "this" {
  count = local.create_partner_channel ? 1 : 0

  name     = local.channel_name
  location = var.spec.location
  project  = local.project_id

  third_party_provider = var.spec.partner_channel.third_party_provider
  crypto_key_name      = var.spec.partner_channel.crypto_key != "" ? var.spec.partner_channel.crypto_key : null

  labels = local.final_labels

  deletion_policy = local.deletion_policy

  depends_on = [google_project_service.eventarc_api]
}

# The per-project-per-location googleChannelConfig SINGLETON — the shared
# conduit ALL non-partner triggers in this project+location deliver
# through. Managed here only when google_channel_crypto_key is set; the
# spec comment tells practitioners to manage it from AT MOST ONE trigger
# per project+location. Deleting it is a state-only no-op in the provider
# (the singleton always exists in GCP).
resource "google_eventarc_google_channel_config" "this" {
  count = local.create_google_channel_config ? 1 : 0

  name = "projects/${var.spec.project_id != "" ? var.spec.project_id : data.google_client_config.this[0].project}/locations/${var.spec.location}/googleChannelConfig"

  location = var.spec.location
  project  = local.project_id

  crypto_key_name = var.spec.google_channel_crypto_key

  depends_on = [google_project_service.eventarc_api]
}

# The Eventarc trigger — the routing rule "when THIS event happens, call
# THAT service". The first trigger in a project provisions Eventarc's
# service agent; the first delivery can lag a few minutes behind the apply
# (P4SA propagation).
resource "google_eventarc_trigger" "this" {
  name     = local.trigger_name
  location = var.spec.location
  project  = local.project_id

  dynamic "matching_criteria" {
    for_each = var.spec.matching_criteria
    content {
      attribute = matching_criteria.value.attribute
      value     = matching_criteria.value.value
      operator  = matching_criteria.value.operator != "" ? matching_criteria.value.operator : null
    }
  }

  # Exactly one destination arm (spec-validated — the provider checks only
  # server-side at create time).
  destination {
    dynamic "cloud_run_service" {
      for_each = var.spec.destination.cloud_run_service != null ? [var.spec.destination.cloud_run_service] : []
      content {
        service = cloud_run_service.value.service
        # Optional+Computed in the provider: sent only when set, so GCP
        # resolves the region from the trigger's location otherwise.
        region = cloud_run_service.value.region != "" ? cloud_run_service.value.region : null
        path   = cloud_run_service.value.path != "" ? cloud_run_service.value.path : null
      }
    }

    dynamic "gke" {
      for_each = var.spec.destination.gke != null ? [var.spec.destination.gke] : []
      content {
        cluster   = gke.value.cluster
        location  = gke.value.location
        namespace = gke.value.namespace
        service   = gke.value.service
        path      = gke.value.path != "" ? gke.value.path : null
      }
    }

    workflow = var.spec.destination.workflow != "" ? var.spec.destination.workflow : null

    dynamic "http_endpoint" {
      for_each = var.spec.destination.http_endpoint != null ? [var.spec.destination.http_endpoint] : []
      content {
        uri = http_endpoint.value.uri
      }
    }

    # The provider models the attachment as a sibling network_config block
    # permitted only with HTTP endpoints — the spec carries it inside the
    # arm; the wiring restores the provider shape.
    dynamic "network_config" {
      for_each = var.spec.destination.http_endpoint != null ? [var.spec.destination.http_endpoint] : []
      content {
        network_attachment = network_config.value.network_attachment
      }
    }
  }

  service_account = var.spec.service_account != "" ? var.spec.service_account : null

  # Wired to the created partner channel's full resource name — assembled
  # from the channel's own computed project attribute so the
  # ambient-project case renders correctly.
  channel = (
    local.create_partner_channel
    ? "projects/${google_eventarc_channel.this[0].project}/locations/${var.spec.location}/channels/${google_eventarc_channel.this[0].name}"
    : null
  )

  dynamic "transport" {
    for_each = var.spec.transport_pubsub_topic != "" ? [var.spec.transport_pubsub_topic] : []
    content {
      pubsub {
        topic = transport.value
      }
    }
  }

  event_data_content_type = var.spec.event_data_content_type != "" ? var.spec.event_data_content_type : null

  # Only the value 1 is valid (provider truth) — it DISABLES Eventarc's
  # default retries; Cloud Run destinations only (both spec-validated).
  dynamic "retry_policy" {
    for_each = var.spec.retry_max_attempts != 0 ? [var.spec.retry_max_attempts] : []
    content {
      max_attempts = retry_policy.value
    }
  }

  labels = local.final_labels

  deletion_policy = local.deletion_policy

  depends_on = [
    google_project_service.eventarc_api,
    google_eventarc_google_channel_config.this,
  ]
}
