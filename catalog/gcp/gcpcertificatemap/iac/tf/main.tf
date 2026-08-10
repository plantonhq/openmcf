# Enable the Certificate Manager API so a fresh project can host the map.
# disable_on_destroy is false: tearing down one map must never disable
# Certificate Manager for everything else in the project.
resource "google_project_service" "certificatemanager_api" {
  project = local.project_id
  service = "certificatemanager.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The certificate map — the hostname-to-certificate routing table an
# external HTTPS load balancer consults at TLS handshake time. GLOBAL by
# API design (no location argument). Immutable name: changing it replaces
# the map and every entry — detach the map from proxies first.
resource "google_certificate_manager_certificate_map" "this" {
  name    = local.map_name
  project = local.project_id

  description = var.spec.description != "" ? var.spec.description : null
  labels      = local.final_labels

  deletion_policy = local.deletion_policy

  depends_on = [google_project_service.certificatemanager_api]
}

# The entry fan-out: each binds a hostname (or the PRIMARY fallback) to
# 1–15 certificates. Entries are almost entirely IMMUTABLE — hostname,
# matcher, and entry name all replace the entry (a brief window where the
# hostname has no binding; plan changes accordingly). The certificate LIST
# is the mutable part — certificate rotation edits it in place while the
# entry keeps serving.
resource "google_certificate_manager_certificate_map_entry" "this" {
  for_each = { for entry in var.spec.entries : entry.entry_name => entry }

  name    = each.value.entry_name
  map     = google_certificate_manager_certificate_map.this.name
  project = local.project_id

  # Exactly one of hostname / matcher (spec-validated, mirroring the
  # provider's ExactlyOneOf).
  hostname = each.value.hostname != "" ? each.value.hostname : null
  matcher  = each.value.matcher != "" ? each.value.matcher : null

  certificates = each.value.certificates

  description = each.value.description != "" ? each.value.description : null
  labels      = merge(each.value.labels, local.final_labels)

  deletion_policy = local.deletion_policy
}
