# Enable the Cloud Monitoring API so a fresh project can host the check.
# disable_on_destroy is false: tearing down one check must never disable
# monitoring for everything else in the project.
resource "google_project_service" "monitoring_api" {
  project = local.project_id
  service = "monitoring.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud Monitoring uptime check — the probe Google runs against the
# configured target from multiple regions on a fixed cadence. A check only
# MEASURES; pair it with a google_monitoring_alert_policy filtering on
# uptime_check_passed and this check's id to actually page.
#
# Exactly one target block (monitored_resource | resource_group |
# synthetic_monitor) and — except for synthetic monitors, which embed their
# own probe logic — exactly one check block (http_check | tcp_check) render;
# the spec's validations enforce the shape before the module runs.
#
# Enum-like strings and durations with API defaults (period, checker_type,
# request_method, content_type) fall through as null when empty so GCP
# applies its own defaults rather than receiving empty strings it would
# reject.
resource "google_monitoring_uptime_check_config" "this" {
  display_name = local.display_name
  timeout      = var.spec.timeout
  project      = local.project_id

  period       = var.spec.period != "" ? var.spec.period : null
  checker_type = var.spec.checker_type != "" ? var.spec.checker_type : null

  selected_regions = length(var.spec.selected_regions) > 0 ? var.spec.selected_regions : null

  log_check_failures = var.spec.log_check_failures ? true : null

  user_labels = local.final_labels

  dynamic "monitored_resource" {
    for_each = var.spec.monitored_resource != null ? [var.spec.monitored_resource] : []
    content {
      type   = monitored_resource.value.type
      labels = monitored_resource.value.labels
    }
  }

  dynamic "resource_group" {
    for_each = var.spec.resource_group != null ? [var.spec.resource_group] : []
    content {
      group_id      = resource_group.value.group_id != "" ? resource_group.value.group_id : null
      resource_type = resource_group.value.resource_type != "" ? resource_group.value.resource_type : null
    }
  }

  dynamic "synthetic_monitor" {
    for_each = var.spec.synthetic_monitor != null ? [var.spec.synthetic_monitor] : []
    content {
      cloud_function_v2 {
        name = synthetic_monitor.value.cloud_function
      }
    }
  }

  dynamic "http_check" {
    for_each = var.spec.http_check != null ? [var.spec.http_check] : []
    content {
      path                = http_check.value.path != "" ? http_check.value.path : null
      port                = http_check.value.port != 0 ? http_check.value.port : null
      request_method      = http_check.value.request_method != "" ? http_check.value.request_method : null
      use_ssl             = http_check.value.use_ssl ? true : null
      validate_ssl        = http_check.value.validate_ssl ? true : null
      body                = http_check.value.body != "" ? http_check.value.body : null
      content_type        = http_check.value.content_type != "" ? http_check.value.content_type : null
      custom_content_type = http_check.value.custom_content_type != "" ? http_check.value.custom_content_type : null
      headers             = length(http_check.value.headers) > 0 ? http_check.value.headers : null
      mask_headers        = http_check.value.mask_headers ? true : null

      dynamic "auth_info" {
        for_each = http_check.value.auth_info != null ? [http_check.value.auth_info] : []
        content {
          username = auth_info.value.username
          password = auth_info.value.password
        }
      }

      dynamic "service_agent_authentication" {
        for_each = http_check.value.service_agent_authentication != null ? [http_check.value.service_agent_authentication] : []
        content {
          type = service_agent_authentication.value.type != "" ? service_agent_authentication.value.type : null
        }
      }

      dynamic "accepted_response_status_codes" {
        for_each = http_check.value.accepted_response_status_codes
        content {
          status_class = accepted_response_status_codes.value.status_class != "" ? accepted_response_status_codes.value.status_class : null
          status_value = accepted_response_status_codes.value.status_value != 0 ? accepted_response_status_codes.value.status_value : null
        }
      }

      dynamic "ping_config" {
        for_each = http_check.value.ping_config != null ? [http_check.value.ping_config] : []
        content {
          pings_count = ping_config.value.pings_count
        }
      }
    }
  }

  dynamic "tcp_check" {
    for_each = var.spec.tcp_check != null ? [var.spec.tcp_check] : []
    content {
      port = tcp_check.value.port

      dynamic "ping_config" {
        for_each = tcp_check.value.ping_config != null ? [tcp_check.value.ping_config] : []
        content {
          pings_count = ping_config.value.pings_count
        }
      }
    }
  }

  dynamic "content_matchers" {
    for_each = var.spec.content_matchers
    content {
      content = content_matchers.value.content
      matcher = content_matchers.value.matcher != "" ? content_matchers.value.matcher : null

      dynamic "json_path_matcher" {
        for_each = content_matchers.value.json_path_matcher != null ? [content_matchers.value.json_path_matcher] : []
        content {
          json_path    = json_path_matcher.value.json_path
          json_matcher = json_path_matcher.value.json_matcher != "" ? json_path_matcher.value.json_matcher : null
        }
      }
    }
  }

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.monitoring_api]
}
