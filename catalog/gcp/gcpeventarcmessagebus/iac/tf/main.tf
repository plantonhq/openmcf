# Enable the Eventarc API so a fresh project can host the family.
# disable_on_destroy is false: tearing down one bus must never disable
# Eventarc for everything else in the project.
resource "google_project_service" "eventarc_api" {
  project = local.project_id
  service = "eventarc.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Eventarc Advanced message bus — the central conduit. Satellites are
# wired to it BY RESOURCE REFERENCE (its computed `name` is the full
# resource name), never by string assembly, so the ambient-project case
# renders correctly and dependency order comes free — byte-identical
# results to the Pulumi module.
resource "google_eventarc_message_bus" "this" {
  message_bus_id = local.message_bus_id
  location       = var.spec.location
  project        = local.project_id

  display_name    = var.spec.display_name != "" ? var.spec.display_name : null
  labels          = local.final_labels
  annotations     = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  crypto_key_name = var.spec.crypto_key != "" ? var.spec.crypto_key : null

  dynamic "logging_config" {
    for_each = var.spec.log_severity != "" ? [var.spec.log_severity] : []
    content {
      log_severity = logging_config.value
    }
  }

  deletion_policy = local.deletion_policy

  depends_on = [google_project_service.eventarc_api]
}

# Google API sources — each publishes Google-service events INTO this bus
# (destination auto-wired to the created bus; a source feeding another bus
# belongs to that bus's kind instance).
resource "google_eventarc_google_api_source" "this" {
  for_each = { for source in var.spec.google_api_sources : source.source_id => source }

  google_api_source_id = each.value.source_id
  location             = var.spec.location
  project              = local.project_id

  destination = google_eventarc_message_bus.this.name

  display_name    = each.value.display_name != "" ? each.value.display_name : null
  labels          = merge(each.value.labels, local.final_labels)
  annotations     = length(each.value.annotations) > 0 ? each.value.annotations : null
  crypto_key_name = each.value.crypto_key != "" ? each.value.crypto_key : null

  dynamic "logging_config" {
    for_each = each.value.log_severity != "" ? [each.value.log_severity] : []
    content {
      log_severity = logging_config.value
    }
  }

  deletion_policy = local.deletion_policy
}

# Pipelines — deliver messages OUT of the bus. The API supports exactly
# one destination per pipeline (the provider's own schema note); the spec
# models exactly one and the single-element destinations block restores
# the provider's list shape. The spec's pipeline-level
# output_payload_format also lives INSIDE that destination element (the
# provider's shape).
resource "google_eventarc_pipeline" "this" {
  for_each = { for pipeline in var.spec.pipelines : pipeline.pipeline_id => pipeline }

  pipeline_id = each.value.pipeline_id
  location    = var.spec.location
  project     = local.project_id

  destinations {
    dynamic "http_endpoint" {
      for_each = each.value.destination.http_endpoint != null ? [each.value.destination.http_endpoint] : []
      content {
        uri                      = http_endpoint.value.uri
        message_binding_template = http_endpoint.value.message_binding_template != "" ? http_endpoint.value.message_binding_template : null
      }
    }

    # Required for HTTP endpoints, forbidden otherwise (provider rule,
    # spec-enforced) — the spec carries it inside the arm; the wiring
    # restores the provider's sibling network_config shape.
    dynamic "network_config" {
      for_each = each.value.destination.http_endpoint != null ? [each.value.destination.http_endpoint] : []
      content {
        network_attachment = network_config.value.network_attachment
      }
    }

    topic       = each.value.destination.topic != "" ? each.value.destination.topic : null
    workflow    = each.value.destination.workflow != "" ? each.value.destination.workflow : null
    message_bus = each.value.destination.message_bus != "" ? each.value.destination.message_bus : null

    dynamic "authentication_config" {
      for_each = each.value.authentication != null ? [each.value.authentication] : []
      content {
        dynamic "google_oidc" {
          for_each = authentication_config.value.google_oidc != null ? [authentication_config.value.google_oidc] : []
          content {
            service_account = google_oidc.value.service_account
            audience        = google_oidc.value.audience != "" ? google_oidc.value.audience : null
          }
        }
        dynamic "oauth_token" {
          for_each = authentication_config.value.oauth_token != null ? [authentication_config.value.oauth_token] : []
          content {
            service_account = oauth_token.value.service_account
            scope           = oauth_token.value.scope != "" ? oauth_token.value.scope : null
          }
        }
      }
    }

    dynamic "output_payload_format" {
      for_each = each.value.output_payload_format != null ? [each.value.output_payload_format] : []
      content {
        dynamic "avro" {
          for_each = output_payload_format.value.avro != null ? [output_payload_format.value.avro] : []
          content {
            schema_definition = avro.value.schema_definition != "" ? avro.value.schema_definition : null
          }
        }
        dynamic "json" {
          for_each = output_payload_format.value.json ? [true] : []
          content {}
        }
        dynamic "protobuf" {
          for_each = output_payload_format.value.protobuf != null ? [output_payload_format.value.protobuf] : []
          content {
            schema_definition = protobuf.value.schema_definition != "" ? protobuf.value.schema_definition : null
          }
        }
      }
    }
  }

  dynamic "input_payload_format" {
    for_each = each.value.input_payload_format != null ? [each.value.input_payload_format] : []
    content {
      dynamic "avro" {
        for_each = input_payload_format.value.avro != null ? [input_payload_format.value.avro] : []
        content {
          schema_definition = avro.value.schema_definition != "" ? avro.value.schema_definition : null
        }
      }
      dynamic "json" {
        for_each = input_payload_format.value.json ? [true] : []
        content {}
      }
      dynamic "protobuf" {
        for_each = input_payload_format.value.protobuf != null ? [input_payload_format.value.protobuf] : []
        content {
          schema_definition = protobuf.value.schema_definition != "" ? protobuf.value.schema_definition : null
        }
      }
    }
  }

  # The API allows at most ONE mediation (transformation) per pipeline —
  # the single spec template renders as the single-element list.
  dynamic "mediations" {
    for_each = each.value.mediation_transformation_template != "" ? [each.value.mediation_transformation_template] : []
    content {
      transformation {
        transformation_template = mediations.value
      }
    }
  }

  dynamic "retry_policy" {
    for_each = each.value.retry_policy != null ? [each.value.retry_policy] : []
    content {
      max_attempts    = retry_policy.value.max_attempts != 0 ? retry_policy.value.max_attempts : null
      min_retry_delay = retry_policy.value.min_retry_delay != "" ? retry_policy.value.min_retry_delay : null
      max_retry_delay = retry_policy.value.max_retry_delay != "" ? retry_policy.value.max_retry_delay : null
    }
  }

  display_name    = each.value.display_name != "" ? each.value.display_name : null
  labels          = merge(each.value.labels, local.final_labels)
  annotations     = length(each.value.annotations) > 0 ? each.value.annotations : null
  crypto_key_name = each.value.crypto_key != "" ? each.value.crypto_key : null

  dynamic "logging_config" {
    for_each = each.value.log_severity != "" ? [each.value.log_severity] : []
    content {
      log_severity = logging_config.value
    }
  }

  deletion_policy = local.deletion_policy

  depends_on = [google_project_service.eventarc_api]
}

# Enrollments — the routing table: each selects messages from the bus with
# a CEL expression and delivers them to one of this spec's pipelines
# (wired to the created pipeline's computed full name; the sibling-id
# contract is spec-validated, so the lookup always resolves).
resource "google_eventarc_enrollment" "this" {
  for_each = { for enrollment in var.spec.enrollments : enrollment.enrollment_id => enrollment }

  enrollment_id = each.value.enrollment_id
  location      = var.spec.location
  project       = local.project_id

  cel_match   = each.value.cel_match
  message_bus = google_eventarc_message_bus.this.name
  destination = google_eventarc_pipeline.this[each.value.pipeline].name

  display_name = each.value.display_name != "" ? each.value.display_name : null
  labels       = merge(each.value.labels, local.final_labels)
  annotations  = length(each.value.annotations) > 0 ? each.value.annotations : null

  deletion_policy = local.deletion_policy
}
