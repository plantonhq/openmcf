# Enable the Compute Engine API so a fresh project can host the health check.
# disable_on_destroy is false: tearing down one health check must never
# disable the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Compute Engine health check — the probe backend services consult before
# sending traffic to a backend, and managed instance groups consult before
# auto-healing an instance.
#
# One kind, two provider resources: GCP models global and regional health
# checks as separate API collections with an otherwise identical surface, so
# this module creates google_compute_health_check when spec.region is empty
# and google_compute_region_health_check when it is set. Exactly one of the
# two count guards is 1 — never both.
#
# name and project are immutable (ForceNew): changing either destroys and
# recreates the check, briefly breaking every backend service referencing the
# old self_link. All probe knobs (interval, thresholds, protocol settings)
# update in place.
#
# Ports deliberately fall through to the API defaults (http/tcp 80,
# https/http2/ssl 443) when unset — hardcoding them here would silently pin
# behavior the provider may evolve.
resource "google_compute_health_check" "this" {
  count = local.is_regional ? 0 : 1

  name        = local.health_check_name
  project     = local.project_id
  description = var.spec.description

  check_interval_sec  = local.check_interval_sec
  timeout_sec         = local.timeout_sec
  healthy_threshold   = local.healthy_threshold
  unhealthy_threshold = local.unhealthy_threshold

  # Global-only: pin probing to exactly 3 named regions so a regional outage
  # cannot flip global health verdicts (the regional resource has no
  # equivalent field).
  source_regions = length(var.spec.source_regions) > 0 ? var.spec.source_regions : null

  dynamic "http_health_check" {
    for_each = local.http != null ? [local.http] : []
    content {
      host               = http_health_check.value.host
      port               = http_health_check.value.port
      port_name          = http_health_check.value.port_name
      port_specification = http_health_check.value.port_specification
      proxy_header       = http_health_check.value.proxy_header
      request_path       = http_health_check.value.request_path
      response           = http_health_check.value.response
    }
  }

  dynamic "https_health_check" {
    for_each = local.https != null ? [local.https] : []
    content {
      host               = https_health_check.value.host
      port               = https_health_check.value.port
      port_name          = https_health_check.value.port_name
      port_specification = https_health_check.value.port_specification
      proxy_header       = https_health_check.value.proxy_header
      request_path       = https_health_check.value.request_path
      response           = https_health_check.value.response
    }
  }

  dynamic "http2_health_check" {
    for_each = local.http2 != null ? [local.http2] : []
    content {
      host               = http2_health_check.value.host
      port               = http2_health_check.value.port
      port_name          = http2_health_check.value.port_name
      port_specification = http2_health_check.value.port_specification
      proxy_header       = http2_health_check.value.proxy_header
      request_path       = http2_health_check.value.request_path
      response           = http2_health_check.value.response
    }
  }

  dynamic "tcp_health_check" {
    for_each = local.tcp != null ? [local.tcp] : []
    content {
      port               = tcp_health_check.value.port
      port_name          = tcp_health_check.value.port_name
      port_specification = tcp_health_check.value.port_specification
      proxy_header       = tcp_health_check.value.proxy_header
      request            = tcp_health_check.value.request
      response           = tcp_health_check.value.response
    }
  }

  dynamic "ssl_health_check" {
    for_each = local.ssl != null ? [local.ssl] : []
    content {
      port               = ssl_health_check.value.port
      port_name          = ssl_health_check.value.port_name
      port_specification = ssl_health_check.value.port_specification
      proxy_header       = ssl_health_check.value.proxy_header
      request            = ssl_health_check.value.request
      response           = ssl_health_check.value.response
    }
  }

  dynamic "grpc_health_check" {
    for_each = local.grpc != null ? [local.grpc] : []
    content {
      grpc_service_name  = grpc_health_check.value.grpc_service_name
      port               = grpc_health_check.value.port
      port_name          = grpc_health_check.value.port_name
      port_specification = grpc_health_check.value.port_specification
    }
  }

  dynamic "grpc_tls_health_check" {
    for_each = local.grpc_tls != null ? [local.grpc_tls] : []
    content {
      grpc_service_name  = grpc_tls_health_check.value.grpc_service_name
      port               = grpc_tls_health_check.value.port
      port_specification = grpc_tls_health_check.value.port_specification
    }
  }

  # The block is always emitted so disabling logging is an explicit false,
  # not an absent block the API back-fills (which would show as drift).
  log_config {
    enable = var.spec.enable_logging
  }

  depends_on = [google_project_service.compute_api]
}

# The regional variant — identical probe surface, addressed under
# regions/<region>/healthChecks. Regional backend services can only
# reference health checks in their own region.
resource "google_compute_region_health_check" "this" {
  count = local.is_regional ? 1 : 0

  name        = local.health_check_name
  project     = local.project_id
  region      = var.spec.region
  description = var.spec.description

  check_interval_sec  = local.check_interval_sec
  timeout_sec         = local.timeout_sec
  healthy_threshold   = local.healthy_threshold
  unhealthy_threshold = local.unhealthy_threshold

  dynamic "http_health_check" {
    for_each = local.http != null ? [local.http] : []
    content {
      host               = http_health_check.value.host
      port               = http_health_check.value.port
      port_name          = http_health_check.value.port_name
      port_specification = http_health_check.value.port_specification
      proxy_header       = http_health_check.value.proxy_header
      request_path       = http_health_check.value.request_path
      response           = http_health_check.value.response
    }
  }

  dynamic "https_health_check" {
    for_each = local.https != null ? [local.https] : []
    content {
      host               = https_health_check.value.host
      port               = https_health_check.value.port
      port_name          = https_health_check.value.port_name
      port_specification = https_health_check.value.port_specification
      proxy_header       = https_health_check.value.proxy_header
      request_path       = https_health_check.value.request_path
      response           = https_health_check.value.response
    }
  }

  dynamic "http2_health_check" {
    for_each = local.http2 != null ? [local.http2] : []
    content {
      host               = http2_health_check.value.host
      port               = http2_health_check.value.port
      port_name          = http2_health_check.value.port_name
      port_specification = http2_health_check.value.port_specification
      proxy_header       = http2_health_check.value.proxy_header
      request_path       = http2_health_check.value.request_path
      response           = http2_health_check.value.response
    }
  }

  dynamic "tcp_health_check" {
    for_each = local.tcp != null ? [local.tcp] : []
    content {
      port               = tcp_health_check.value.port
      port_name          = tcp_health_check.value.port_name
      port_specification = tcp_health_check.value.port_specification
      proxy_header       = tcp_health_check.value.proxy_header
      request            = tcp_health_check.value.request
      response           = tcp_health_check.value.response
    }
  }

  dynamic "ssl_health_check" {
    for_each = local.ssl != null ? [local.ssl] : []
    content {
      port               = ssl_health_check.value.port
      port_name          = ssl_health_check.value.port_name
      port_specification = ssl_health_check.value.port_specification
      proxy_header       = ssl_health_check.value.proxy_header
      request            = ssl_health_check.value.request
      response           = ssl_health_check.value.response
    }
  }

  dynamic "grpc_health_check" {
    for_each = local.grpc != null ? [local.grpc] : []
    content {
      grpc_service_name  = grpc_health_check.value.grpc_service_name
      port               = grpc_health_check.value.port
      port_name          = grpc_health_check.value.port_name
      port_specification = grpc_health_check.value.port_specification
    }
  }

  dynamic "grpc_tls_health_check" {
    for_each = local.grpc_tls != null ? [local.grpc_tls] : []
    content {
      grpc_service_name  = grpc_tls_health_check.value.grpc_service_name
      port               = grpc_tls_health_check.value.port
      port_specification = grpc_tls_health_check.value.port_specification
    }
  }

  log_config {
    enable = var.spec.enable_logging
  }

  depends_on = [google_project_service.compute_api]
}
