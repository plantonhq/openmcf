locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The cloud-side name defaults to metadata.name when the spec leaves
  # health_check_name empty — the same naming basis every kind uses.
  health_check_name = (
    var.spec.health_check_name != null && var.spec.health_check_name != ""
    ? var.spec.health_check_name
    : var.metadata.name
  )

  # Empty region means a GLOBAL health check; a set region selects the
  # regional resource. Both branches in main.tf key off this.
  is_regional = var.spec.region != null && var.spec.region != ""

  # Probe cadence and verdict thresholds. Planton middleware applies the
  # proto defaults (5/5/2/2) before the module runs; the coalesce here is a
  # safety net for direct tfvars invocations so behavior stays identical.
  check_interval_sec  = coalesce(var.spec.check_interval_sec, 5)
  timeout_sec         = coalesce(var.spec.timeout_sec, 5)
  healthy_threshold   = coalesce(var.spec.healthy_threshold, 2)
  unhealthy_threshold = coalesce(var.spec.unhealthy_threshold, 2)

  # Normalize optional strings inside protocol blocks: the tfvars converter
  # emits empty strings for unset proto fields, but the provider treats "" and
  # null differently for some attributes (e.g. an empty port_specification is
  # rejected). Map "" -> null once here so main.tf stays declarative.
  http = var.spec.http == null ? null : {
    host               = try(var.spec.http.host, null) != "" ? var.spec.http.host : null
    port               = try(var.spec.http.port, null) != 0 ? var.spec.http.port : null
    port_name          = try(var.spec.http.port_name, null) != "" ? var.spec.http.port_name : null
    port_specification = try(var.spec.http.port_specification, null) != "" ? var.spec.http.port_specification : null
    proxy_header       = try(var.spec.http.proxy_header, null) != "" ? var.spec.http.proxy_header : null
    request_path       = try(var.spec.http.request_path, null) != "" ? var.spec.http.request_path : null
    response           = try(var.spec.http.response, null) != "" ? var.spec.http.response : null
  }
  https = var.spec.https == null ? null : {
    host               = try(var.spec.https.host, null) != "" ? var.spec.https.host : null
    port               = try(var.spec.https.port, null) != 0 ? var.spec.https.port : null
    port_name          = try(var.spec.https.port_name, null) != "" ? var.spec.https.port_name : null
    port_specification = try(var.spec.https.port_specification, null) != "" ? var.spec.https.port_specification : null
    proxy_header       = try(var.spec.https.proxy_header, null) != "" ? var.spec.https.proxy_header : null
    request_path       = try(var.spec.https.request_path, null) != "" ? var.spec.https.request_path : null
    response           = try(var.spec.https.response, null) != "" ? var.spec.https.response : null
  }
  http2 = var.spec.http2 == null ? null : {
    host               = try(var.spec.http2.host, null) != "" ? var.spec.http2.host : null
    port               = try(var.spec.http2.port, null) != 0 ? var.spec.http2.port : null
    port_name          = try(var.spec.http2.port_name, null) != "" ? var.spec.http2.port_name : null
    port_specification = try(var.spec.http2.port_specification, null) != "" ? var.spec.http2.port_specification : null
    proxy_header       = try(var.spec.http2.proxy_header, null) != "" ? var.spec.http2.proxy_header : null
    request_path       = try(var.spec.http2.request_path, null) != "" ? var.spec.http2.request_path : null
    response           = try(var.spec.http2.response, null) != "" ? var.spec.http2.response : null
  }
  tcp = var.spec.tcp == null ? null : {
    port               = try(var.spec.tcp.port, null) != 0 ? var.spec.tcp.port : null
    port_name          = try(var.spec.tcp.port_name, null) != "" ? var.spec.tcp.port_name : null
    port_specification = try(var.spec.tcp.port_specification, null) != "" ? var.spec.tcp.port_specification : null
    proxy_header       = try(var.spec.tcp.proxy_header, null) != "" ? var.spec.tcp.proxy_header : null
    request            = try(var.spec.tcp.request, null) != "" ? var.spec.tcp.request : null
    response           = try(var.spec.tcp.response, null) != "" ? var.spec.tcp.response : null
  }
  ssl = var.spec.ssl == null ? null : {
    port               = try(var.spec.ssl.port, null) != 0 ? var.spec.ssl.port : null
    port_name          = try(var.spec.ssl.port_name, null) != "" ? var.spec.ssl.port_name : null
    port_specification = try(var.spec.ssl.port_specification, null) != "" ? var.spec.ssl.port_specification : null
    proxy_header       = try(var.spec.ssl.proxy_header, null) != "" ? var.spec.ssl.proxy_header : null
    request            = try(var.spec.ssl.request, null) != "" ? var.spec.ssl.request : null
    response           = try(var.spec.ssl.response, null) != "" ? var.spec.ssl.response : null
  }
  grpc = var.spec.grpc == null ? null : {
    grpc_service_name  = try(var.spec.grpc.grpc_service_name, null) != "" ? var.spec.grpc.grpc_service_name : null
    port               = try(var.spec.grpc.port, null) != 0 ? var.spec.grpc.port : null
    port_name          = try(var.spec.grpc.port_name, null) != "" ? var.spec.grpc.port_name : null
    port_specification = try(var.spec.grpc.port_specification, null) != "" ? var.spec.grpc.port_specification : null
  }
  grpc_tls = var.spec.grpc_tls == null ? null : {
    grpc_service_name  = try(var.spec.grpc_tls.grpc_service_name, null) != "" ? var.spec.grpc_tls.grpc_service_name : null
    port               = try(var.spec.grpc_tls.port, null) != 0 ? var.spec.grpc_tls.port : null
    port_specification = try(var.spec.grpc_tls.port_specification, null) != "" ? var.spec.grpc_tls.port_specification : null
  }
}
