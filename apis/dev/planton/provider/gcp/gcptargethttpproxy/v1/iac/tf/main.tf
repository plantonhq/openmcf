# Enable the Compute Engine API so a fresh project can host the proxy.
# disable_on_destroy is false: tearing down one proxy must never disable the
# API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A global Compute Engine target HTTP proxy — the plaintext-HTTP frontend
# adapter that binds a global forwarding rule (the VIP) to a URL map (the
# routing brain). The proxy is deliberately thin: TLS lives on the target
# HTTPS proxy sibling, routing on the URL map, traffic policy on the backend
# service. The standard production pattern points this proxy at a
# redirect-only URL map (http→https 301) while the HTTPS proxy serves the
# real application.
#
# url_map is the only mutable field — GCP repoints it in place via a
# dedicated setUrlMap call, so a live frontend can move to a new routing
# table with no downtime. Everything else (name, description, keep-alive,
# proxy_bind, project) is immutable (ForceNew) and briefly breaks any
# forwarding rule referencing the old self_link on recreate.
resource "google_compute_target_http_proxy" "this" {
  name        = local.proxy_name
  project     = local.project_id
  description = local.description

  # Arrives resolved to a literal self-link (or plain name, which the
  # provider expands against the project).
  url_map = var.spec.url_map

  # Only honored by EXTERNAL_MANAGED load balancers; null keeps GCP's
  # default (610s) in charge.
  http_keep_alive_timeout_sec = local.http_keep_alive_timeout_sec

  # Traffic Director binding; null lets the API compute its default (false).
  proxy_bind = local.proxy_bind

  depends_on = [google_project_service.compute_api]
}
